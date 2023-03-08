package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader := strings.Split(request.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			RenderResponse(writer, request, NewUnauthorizedErr("unauthorized"))
			return
		}

		config, ok := request.Context().Value("config").(Configuration)

		if !ok {
			RenderResponse(writer, request, NewInternalServerErr("config not found"))
			return
		}

		jwtToken := authHeader[1]
		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			if config.GetTokenSecretKey() != "" {
				return []byte(config.GetTokenSecretKey()), nil
			}

			if config.GetTokenPrivateKey() != nil {
				return config.GetTokenPrivateKey(), nil
			}

			return "", errors.New("unknown signing method")

		})

		if err != nil {
			RenderResponse(writer, request, NewUnauthorizedErr(err.Error()))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if _, ok := claims["username"]; !ok {
				RenderResponse(writer, request, NewUnauthorizedErr("unauthorized"))
				return
			}

			if _, ok := claims["email"]; !ok {
				RenderResponse(writer, request, NewUnauthorizedErr("unauthorized"))
				return
			}

			ctx := context.WithValue(request.Context(), "sender", Sender{
				Id:       claims["sub"].(string),
				Username: claims["username"].(string),
			})

			// Access context values in handlers like this
			// props, _ := r.Context().Value("props").(jwt.MapClaims)
			next.ServeHTTP(writer, request.WithContext(ctx))
		} else {
			RenderResponse(writer, request, NewUnauthorizedErr("unauthorized"))
		}
	})
}
