package service

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

type GetMessagesResponse []StoredMessage

func (g GetMessagesResponse) Render(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(http.StatusOK)

	return nil
}

func GetMessagesMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		latStr := request.URL.Query().Get("lat")

		if latStr == "" {
			RenderResponse(writer, request, NewBadRequestErr("lat parameter not provided"))
			return
		}

		lat, err := strconv.ParseFloat(latStr, 64)

		if err != nil {
			RenderResponse(writer, request, NewBadRequestErr("invalid lat parameter"))
			return
		}

		longStr := request.URL.Query().Get("long")

		if longStr == "" {
			RenderResponse(writer, request, NewBadRequestErr("long parameter not provided"))
			return
		}

		long, err := strconv.ParseFloat(longStr, 64)

		if err != nil {
			RenderResponse(writer, request, NewBadRequestErr("invalid long parameter"))
			return
		}

		radiusInMetersStr := request.URL.Query().Get("radius")
		radiusInMeters := 100.0

		if radiusInMetersStr != "" {
			radiusInMeters, err = strconv.ParseFloat(radiusInMetersStr, 64)

			if err != nil {
				RenderResponse(writer, request, NewBadRequestErr("invalid radius parameter"))
				return
			}
		}

		limitStr := request.URL.Query().Get("limit")
		limit := 100

		if limitStr != "" {
			limit, err = strconv.Atoi(limitStr)

			if err != nil {
				RenderResponse(writer, request, NewBadRequestErr("invalid limit parameter"))
				return
			}
		}

		afterStr := request.URL.Query().Get("after")
		after := time.UnixMilli(0)

		if afterStr != "" {
			afterInt, err := strconv.ParseInt(afterStr, 10, 64)

			if err != nil {
				RenderResponse(writer, request, NewBadRequestErr("invalid after parameter"))
				return
			}

			after = time.UnixMilli(afterInt)
		}

		repo, ok := request.Context().Value("repo").(MessageRepository)

		if !ok {
			RenderResponse(writer, request, NewInternalServerErr("repo not configured"))
			return
		}

		messages, err := repo.GetMessagesForLocation(Location{
			Long: long,
			Lat:  lat,
		}, radiusInMeters, limit, after)

		if err != nil {
			RenderResponse(writer, request, NewInternalServerErr("repo error"))
			return
		}

		ctx := context.WithValue(request.Context(), "messages", messages)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func GetMessages(writer http.ResponseWriter, request *http.Request) {
	messages, ok := request.Context().Value("messages").([]StoredMessage)

	if !ok {
		RenderResponse(writer, request, NewInternalServerErr("internal error"))
		return
	}

	response := make(GetMessagesResponse, len(messages))
	copy(response, messages)
	RenderResponse(writer, request, response)
}
