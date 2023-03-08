package service

import (
	"context"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type GetMessageResponse StoredMessage

func GetMessageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		repo, ok := request.Context().Value("repo").(MessageRepository)

		if !ok {
			RenderResponse(writer, request, NewInternalServerErr("repo not configured"))
			return
		}

		id := chi.URLParam(request, "id")

		if id == "" {
			RenderResponse(writer, request, NewBadRequestErr("invalid id parameter"))
			return
		}

		msg, err := repo.GetMessage(id)

		if err != nil {
			log.Println(err)
			RenderResponse(writer, request, NewInternalServerErr("repo error"))
			return
		}

		ctx := context.WithValue(request.Context(), "message", &msg)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func GetMessage(writer http.ResponseWriter, request *http.Request) {
	msg, ok := request.Context().Value("message").(*StoredMessage)

	if !ok {
		RenderResponse(writer, request, NewNotFoundErr("no message found with that id"))
		return
	}

	RenderResponse(writer, request, msg)
}
