package server

import (
	"log"
	"net/http"
	"sports-registration/handler"
	"sports-registration/repository"
)

func StartServer() {
	repo := repository.NewRepository()
	
	// Пытаемся подключиться к БД, но не прерываем работу если не удалось
	dbRepo, err := repository.NewDatabaseRepository()
	if err != nil {
		log.Printf("⚠️  База данных недоступна: %v", err)
		log.Println("📋 Режим работы: только статические данные (без БД)")
		log.Println("💡 Для подключения БД запустите PostgreSQL и установите переменные окружения")
		
		h := handler.NewHandler(repo, nil)
		setupRoutes(h)
	} else {
		defer dbRepo.Close()
		h := handler.NewHandler(repo, dbRepo)
		setupRoutes(h)
	}

	log.Println("🏃 Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupRoutes(h *handler.Handler) {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", h.HomeHandler)
	http.HandleFunc("/athletes", h.AthletesHandler)
	http.HandleFunc("/athlete", h.AthleteDetailHandler)
	http.HandleFunc("/events", h.EventsHandler)
	http.HandleFunc("/event", h.EventDetailHandler)
	http.HandleFunc("/team-application", h.TeamApplicationHandler)
	
	// Новые маршруты для работы с БД (только если БД подключена)
	if h.IsDBConnected() {
		http.HandleFunc("/services", h.ServicesHandler)
		http.HandleFunc("/application", h.ApplicationHandler)
		http.HandleFunc("/application/add", h.AddToApplicationHandler)
		http.HandleFunc("/application/delete", h.DeleteApplicationHandler)
		log.Println("📋 Услуги доступны на http://localhost:8080/services")
		log.Println("📝 Заявка доступна на http://localhost:8080/application")
	}
}
