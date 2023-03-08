package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type addMessageRequest struct {
	Content  string `json:"content"`
	Location `json:"location"`
}

func AddMessageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		repo, ok := request.Context().Value("repo").(MessageRepository)

		if !ok {
			RenderResponse(writer, request, NewInternalServerErr("repo not configured"))
			return
		}

		var amr addMessageRequest
		decoder := json.NewDecoder(request.Body)
		err := decoder.Decode(&amr)

		if err != nil {
			RenderResponse(writer, request, NewBadRequestErr("invalid request body"))
			return
		}

		sender := request.Context().Value("sender").(Sender)

		storedMessage, err := repo.AddMessage(Message{
			Sender:  sender,
			Content: amr.Content,
			Location: Location{
				Long: amr.Long,
				Lat:  amr.Lat,
			},
		})

		if err != nil {
			log.Println(err)
			RenderResponse(writer, request, NewInternalServerErr("repo error"))
			return
		}

		ctx := context.WithValue(request.Context(), "message", &storedMessage)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func AddMessage(writer http.ResponseWriter, request *http.Request) {
	msg, ok := request.Context().Value("message").(*StoredMessage)

	if !ok {
		RenderResponse(writer, request, NewNotFoundErr("unable to store message"))
		return
	}

	RenderResponse(writer, request, msg)
}
