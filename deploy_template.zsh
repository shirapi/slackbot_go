#!/bin/zsh

echo "環境変数を指定してください"

# while [[ -z "$ApiKey" ]]; do
#     vared -p "OPENAI_API_KEY": -c ApiKey
# done

while [[ -z "$AiModelId" ]]; do
    vared -p "AI_MODEL_ID (Bedrockのモデルまたは推論プロファイルID)": -c AiModelId
done

while [[ -z "$ChannelId" ]]; do
    vared -p "SLACK_CHANNEL_ID": -c ChannelId
done

while [[ -z "$SigningSecret" ]]; do
    vared -p "SIGNING_SECRET": -c SigningSecret
done

while [[ -z "$OAuthToken" ]]; do
    vared -p "OAUTH_TOKEN": -c OAuthToken
done

echo "=== deploy info ==="
echo " "
# echo "ApiKey: $ApiKey"
echo "AiModelId: $AiModelId"
echo "ChannelId: $ChannelId"
echo "SigningSecret: $SigningSecret"
echo "OAuthToken: $OAuthToken"
echo "==================="
echo " "

yn=""
vared -p "デプロイを開始しますがよろしいですか？(y/n): " yn
case "$yn" in [yY]*) ;; *) echo "abord." ; exit ;; esac

sam build && sam deploy \
    --stack-name slackbot-go \
    --region ap-northeast-1 \
    --parameter-overrides AiModelId=${AiModelId} ChannelId=${ChannelId} SigningSecret=${SigningSecret} OAuthToken=${OAuthToken}
    # --parameter-overrides ApiKey=${ApiKey} ChannelId=${ChannelId} SigningSecret=${SigningSecret} OAuthToken=${OAuthToken}
