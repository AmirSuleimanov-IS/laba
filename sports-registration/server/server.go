package server

import (
	"log"
	"net/http"
	"sports-registration/handler"
	"sports-registration/repository"
)

func StartServer() {
	repo := repository.NewRepository()
	h := handler.NewHandler(repo)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", h.HomeHandler)
	http.HandleFunc("/athletes", h.AthletesHandler)
	http.HandleFunc("/athlete", h.AthleteDetailHandler)
	http.HandleFunc("/events", h.EventsHandler)
	http.HandleFunc("/event", h.EventDetailHandler)
	http.HandleFunc("/team-application", h.TeamApplicationHandler)

	log.Println("🏃 Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
