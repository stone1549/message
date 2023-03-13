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
	ClientId string `json:"clientId"`
}

type StoredMessage struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Message
}

func (s StoredMessage) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusOK)

	return nil
}
