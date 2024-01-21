package slackbot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"slackbot-go/domain"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const (
	HeaderRetryNum    = "X-Slack-Retry-Num"
	HeaderRetryReason = "X-Slack-Retry-Reason"
	HTTPTimeout       = "http_timeout"
)

type Talker struct {
	Client *slack.Client
	Model  domain.AIModel
	AI     domain.Charactor
}

func (t *Talker) Talk(r *http.Request) domain.Response {
	body, err := t.verify(r)
	if err != nil {
		fmt.Println("verify err:", err)
		return t.newErrorResponse(err)
	}

	// https://github.com/slack-go/slack/blob/master/examples/eventsapi
	// https://www.breakfreeall.net/blog/2023/02/chatgpt-slack-bot/
	// https://blog.linkode.co.jp/entry/2020/03/18/100012
	// https://simple-minds-think-alike.moritamorie.com/entry/verify-requests-with-slack-go
	event, err := slackevents.ParseEvent(
		body,
		// slackevents.OptionNoVerifyToken()
		slackevents.OptionVerifyToken(
			&slackevents.TokenComparator{
				VerificationToken: os.Getenv("SLACK_VERIFICATION_TOKEN"),
			},
		),
	)

	if err != nil {
		fmt.Println("event parse err:", err)
		return t.newErrorResponse(err)
	}
	fmt.Printf("event:%+v\n", event)

	if event.Type == slackevents.URLVerification {
		fmt.Println("challenge")
		var res *slackevents.ChallengeResponse
		if err := json.Unmarshal(body, &res); err != nil {
			fmt.Println("challenge err:", err)
			return domain.Response{Challenge: true, Message: err.Error()}
		}
		return domain.Response{Challenge: true, Message: res.Challenge}
	}

	innerEvent := event.InnerEvent
	switch ev := innerEvent.Data.(type) {
	case *slackevents.AppMentionEvent:
		if ev.Channel != os.Getenv("SLACK_CHANNEL_ID") {
			// 指定のチャンネル以外は応答禁止
			fmt.Println("invalid channel:", err)
			err := t.postErrorMessage(ev, domain.CharactorErrInvalidChannel)
			return t.newErrorResponse(err)
		}
		// 時間がかかるのでリアクションをつけておく
		t.Client.AddReaction(t.AI.GetReactionName(), slack.ItemRef{
			Channel:   ev.Channel,
			Timestamp: ev.TimeStamp,
		})
		histories, err := t.GetHistory(ev)
		if err != nil {
			fmt.Println("history err:", err)
			err := t.postErrorMessage(ev, domain.CharactorErrFatal)
			return t.newErrorResponse(err)
		}
		userMessage := t.GetUserText(ev.Text)
		message, err := t.Model.Call(t.AI, histories, userMessage)
		if err != nil {
			fmt.Println("AI error:", err)
			err := t.postErrorMessage(ev, domain.CharactorErrFatal)
			return t.newErrorResponse(err)
		}
		t.postMessage(ev, message)
	}

	return domain.Response{Message: "success"}
}

func (t *Talker) newErrorResponse(err error) domain.Response {
	return domain.Response{Message: err.Error()}
}

func (t *Talker) isRetry(header http.Header) bool {
	return header.Get(HeaderRetryNum) != "" && header.Get(HeaderRetryReason) == HTTPTimeout
}

func (t *Talker) verify(r *http.Request) (json.RawMessage, error) {
	header := r.Header
	fmt.Printf("header:%+v\n", header)

	if t.isRetry(header) {
		// リトライは無視
		fmt.Println("retry.")
		return nil, nil
	}

	verifier, err := slack.NewSecretsVerifier(r.Header, os.Getenv("SLACK_API_SIGNING_SECRET"))
	if err != nil {
		return nil, err
	}

	bodyReader := io.TeeReader(r.Body, &verifier)
	body, err := io.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if err := verifier.Ensure(); err != nil {
		return nil, err
	}
	return json.RawMessage(body), nil
}

func (t *Talker) GetUserText(text string) string {
	// 受信メッセージからメンションを削除
	re := regexp.MustCompile(`^<@[^>]+>\s*`)
	return re.ReplaceAllString(text, "")
}

func (t *Talker) GetHistory(ev *slackevents.AppMentionEvent) ([]domain.History, error) {
	if ev.ThreadTimeStamp == "" {
		return nil, nil // スレッドなしなら履歴なし
	}
	msgs, _, _, err := t.Client.GetConversationReplies(&slack.GetConversationRepliesParameters{
		ChannelID: ev.Channel,
		Timestamp: ev.ThreadTimeStamp,
		Inclusive: false,
	})
	if err != nil {
		return nil, err
	}
	msgs = msgs[:len(msgs)-1] // 最新のメッセージは重複なので消す
	remMention := func(s string) string {
		re := regexp.MustCompile(`^<@[^>]+>\s*`)
		return re.ReplaceAllString(s, "")
	}
	h := make([]domain.History, 0, len(msgs))
	for _, msg := range msgs {
		from := domain.MessageFromUser
		if msg.BotID != "" {
			from = domain.MessageFromAI
		}
		h = append(h, domain.History{
			MessageFrom: from,
			Message:     remMention(msg.Text),
		})
	}
	fmt.Printf("history============%+v\n", h)
	return h, nil
}

func (t *Talker) postMessage(ev *slackevents.AppMentionEvent, message string) error {
	message = fmt.Sprintf("<@%s> %s", ev.User, message)
	_, ts, err := t.Client.PostMessage(
		ev.Channel,
		slack.MsgOptionTS(ev.TimeStamp),
		slack.MsgOptionUser(ev.User),
		slack.MsgOptionText(message, false),
	)
	fmt.Println("post timestamp:", ts)
	return err
}

func (t *Talker) postErrorMessage(ev *slackevents.AppMentionEvent, errType int) error {
	message := t.AI.GetErrorMessage(errType)
	err := t.postMessage(ev, message)
	if err == nil {
		err = errors.New("エラーメッセージ返信 error")
	}
	return err
}
