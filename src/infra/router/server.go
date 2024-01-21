package router

import (
	"encoding/json"
	"log"
	"net/http"
	"slackbot-go/domain"
)

type Handler func(*http.Request) domain.Response

type Service struct {
	Handler Handler
}

func (s Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := s.Handler(r)
	writeResponse(w, res)
}

func writeResponse(w http.ResponseWriter, res domain.Response) {
	if res.Challenge {
		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte(res.Message)); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	body, err := json.Marshal(res)
	if err != nil {
		body = []byte(`{"message":"予期せぬエラーです"}`)
	}
	w.Write(body)
}
