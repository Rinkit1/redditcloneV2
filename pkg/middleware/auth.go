package middleware

import (
	"context"
	"my/redditclone/pkg/handlers"
	"my/redditclone/pkg/session"
	"net/http"
	"strings"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		temp := r.Header.Get("authorization")
		if temp == "" {
			next.ServeHTTP(w, r)
			return
		}
		token := strings.Split(temp, " ")
		user, err := session.Check(token[1])
		if err != nil {
			text := "unknown error"
			switch {
			case err == session.ErrBadSignMethod:
				text = "bad sign method"
			case err == session.ErrNoPayload:
				text = "no payload"
			case err == session.ErrBadToken:
				text = "bad token"
			}
			handlers.SendMessage(w, http.StatusInternalServerError, text)
			return
		}
		ctx := context.WithValue(r.Context(), session.SessionKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
	/*token := strings.Split(r.Header.Get("authorization"), " ")
	 */
}
