package service

import (
	"net/http"
	"time"
)

type Sender struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type Location struct {
	Long float64 `json:"long"`
	Lat  float64 `json:"lat"`
}

type Message struct {
	Sender   `json:"sender"`
	Content  string `json:"content"`
	Location `json:"location"`
	ClientId string    `json:"clientId"`
	SentAt   time.Time `json:"sentAt"`
}
type StoredMessage struct {
	Id         string    `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	ReceivedAt time.Time `json:"receivedAt"`
	Message
}

func (s StoredMessage) Render(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(http.StatusOK)

	return nil
}
