package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log/slog"
	"net/http"
	"os"
)

type Response struct {
	Code         int
	Error        error
	ResponseBody any
}

func returnResponse(w http.ResponseWriter, encoder json.Encoder, code int, body any, err error) {
	w.WriteHeader(code)
	encoder.Encode(Response{
		Code:         code,
		Error:        err,
		ResponseBody: body,
	})
}

type MiddlewareFunc func(next http.Handler) http.Handler

func VerifyJWT(next http.Handler) http.Handler {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	var role string
	fn := func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header.Get("Token"), func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				logger.Debug("Request discarded: auth failed")
				returnResponse(w, *encoder, http.StatusUnauthorized, nil, fmt.Errorf("unauthorized"))
				return
			}
			if token.Valid {
				claims, ok := token.Claims.(jwt.MapClaims)
				if ok {
					role = claims["role"].(string)
					r = r.WithContext(context.WithValue(context.Background(), "role", role))
					logger.Debug(fmt.Sprintf("User with claims %s authenticated", role))
					next.ServeHTTP(w, r)
				} else {
					logger.Debug("Request discarded: auth failed: required claim absent")
					w.WriteHeader(http.StatusUnauthorized)
					returnResponse(w, *encoder, http.StatusUnauthorized, nil, fmt.Errorf("unauthorized"))
					return
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				logger.Debug("Request discarded: auth failed: invalid token")
				returnResponse(w, *encoder, http.StatusUnauthorized, nil, fmt.Errorf("unauthorized"))
				return
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			logger.Debug("Request discarded: auth failed: token absent")
			returnResponse(w, *encoder, http.StatusUnauthorized, nil, fmt.Errorf("unauthorized"))
			return
		}
	}
	return http.HandlerFunc(fn)
}
