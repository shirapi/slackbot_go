package router

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"slackbot-go/infra/di"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupMux() *chi.Mux {

	r := chi.NewRouter()
	r.Use(RootRecoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(RequestLog)

	talker := di.NewHiyokoTalker()
	r.Get("/", routing(Service{Handler: talker.Talk}))

	return r
}

func routing(next http.Handler) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request Header:", r.Header)
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func RootRecoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				fmt.Println("Recovered from panic", rvr)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func RequestLog(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request URL:", r.URL)
		reqID, ok := r.Context().Value(middleware.RequestIDKey).(string)
		if !ok {
			reqID = "unknown"
		}
		fmt.Println("Request ID:", reqID)
		fmt.Println("Request IP:", r.RemoteAddr)

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err == nil {
			fmt.Println("Request Body:", string(body))
			r.Body = io.NopCloser(bytes.NewBuffer(body))
		} else {
			fmt.Println("Request Body Error:", err)
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
