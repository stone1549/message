package service

import (
	"github.com/go-chi/render"
	"log"
	"net/http"
)

func RenderResponse(writer http.ResponseWriter, request *http.Request, renderer render.Renderer) {
	err := render.Render(writer, request, renderer)
	if err != nil {
		log.Println(err)
	}
}
