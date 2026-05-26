package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"sports-registration/models"
	"sports-registration/repository"
	"strconv"
	"strings"
)

type Handler struct {
	repo     *repository.Repository
	dbRepo   *repository.DatabaseRepository
	tmpl     *template.Template
	currentUserID int // Для демонстрации - фиксированный пользователь ID=1
	dbConnected bool
}

func NewHandler(repo *repository.Repository, dbRepo *repository.DatabaseRepository) *Handler {
	h := &Handler{
		repo:          repo,
		dbRepo:        dbRepo,
		currentUserID: 1, // Фиксируем пользователя для демонстрации
		dbConnected:   dbRepo != nil,
	}
	h.loadTemplates()
	return h
}

// IsDBConnected возвращает true если БД подключена
func (h *Handler) IsDBConnected() bool {
	return h.dbConnected
}

func (h *Handler) loadTemplates() {
	funcMap := template.FuncMap{
		"mul": func(a, b float64) float64 {
			return a * b
		},
	}
	
	tmpl := template.New("").Funcs(funcMap)
	
	// Загружаем layout
	tmpl = template.Must(tmpl.ParseFiles(
		"templates/layout/header.html",
		"templates/layout/footer.html",
	))
	
	// Загружаем страницы
	tmpl = template.Must(tmpl.ParseFiles(
		"templates/athletes.html",
		"templates/events.html",
		"templates/home.html",
		"templates/index.html",
		"templates/team-application.html",
		"templates/athlete-detail.html",
		"templates/event-detail.html",
		"templates/services.html",
		"templates/application.html",
	))
	
	h.tmpl = tmpl
}

// ServicesHandler - GET /services - страница услуг из БД
func (h *Handler) ServicesHandler(w http.ResponseWriter, r *http.Request) {
	services, err := h.dbRepo.GetAllServices()
	if err != nil {
		http.Error(w, "Error loading services: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем текущую заявку пользователя (черновик)
	app, err := h.dbRepo.GetOrCreateDraftApplication(h.currentUserID)
	if err != nil {
		http.Error(w, "Error loading application: "+err.Error(), http.StatusInternalServerError)
		return
	}

	searchQuery := r.URL.Query().Get("search")
	if searchQuery != "" {
		filtered := make([]models.Service, 0)
		for _, s := range services {
			if strings.Contains(strings.ToLower(s.Name), strings.ToLower(searchQuery)) ||
				strings.Contains(strings.ToLower(s.ServiceType), strings.ToLower(searchQuery)) ||
				strings.Contains(strings.ToLower(s.Location), strings.ToLower(searchQuery)) {
				filtered = append(filtered, s)
			}
		}
		services = filtered
	}

	data := map[string]interface{}{
		"Services":       services,
		"SearchQuery":    searchQuery,
		"Application":    app,
		"HasDraft":       app != nil && app.Status == "draft",
		"ServiceCount":   len(app.Services),
		"PageTitle":      "Услуги",
	}

	h.tmpl.ExecuteTemplate(w, "services", data)
}

// ApplicationHandler - GET /application - текущая заявка (корзина)
func (h *Handler) ApplicationHandler(w http.ResponseWriter, r *http.Request) {
	app, err := h.dbRepo.GetOrCreateDraftApplication(h.currentUserID)
	if err != nil {
		http.Error(w, "Error loading application: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if app.Status != "draft" {
		http.Error(w, "Заявка не найдена или уже обработана", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Application":  app,
		"ServiceCount": len(app.Services),
		"PageTitle":    "Текущая заявка",
	}

	h.tmpl.ExecuteTemplate(w, "application", data)
}

// AddToApplicationHandler - POST /application/add - добавление услуги в заявку (через ORM)
func (h *Handler) AddToApplicationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	serviceIDStr := r.FormValue("service_id")
	serviceID, err := strconv.Atoi(serviceIDStr)
	if err != nil {
		http.Error(w, "Invalid service ID", http.StatusBadRequest)
		return
	}

	// Получаем или создаем черновик
	app, err := h.dbRepo.GetOrCreateDraftApplication(h.currentUserID)
	if err != nil {
		http.Error(w, "Error loading application: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if app.Status != "draft" {
		http.Error(w, "Нельзя добавить услугу в обработанную заявку", http.StatusBadRequest)
		return
	}

	// Добавляем услугу через ORM-метод
	err = h.dbRepo.AddServiceToApplication(app.ID, serviceID, 1, len(app.Services)+1, len(app.Services) == 0)
	if err != nil {
		http.Error(w, "Error adding service: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Услуга #%d добавлена в заявку #%d", serviceID, app.ID)

	// Редирект на страницу заявки
	http.Redirect(w, r, "/application", http.StatusSeeOther)
}

// DeleteApplicationHandler - POST /application/delete - удаление заявки (через SQL UPDATE без ORM)
func (h *Handler) DeleteApplicationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	appIDStr := r.FormValue("application_id")
	appID, err := strconv.Atoi(appIDStr)
	if err != nil {
		http.Error(w, "Invalid application ID", http.StatusBadRequest)
		return
	}

	// Удаляем заявку через прямой SQL UPDATE (без ORM)
	err = h.dbRepo.DeleteApplication(appID)
	if err != nil {
		http.Error(w, "Error deleting application: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем JSON ответ для AJAX или редирект
	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Заявка удалена",
		})
		return
	}

	http.Redirect(w, r, "/services", http.StatusSeeOther)
}

// Старые обработчики (для совместимости)
func (h *Handler) AthletesHandler(w http.ResponseWriter, r *http.Request) {
	athletes := h.repo.GetAllAthletes()

	searchQuery := r.URL.Query().Get("search")
	if searchQuery != "" {
		filtered := make([]models.Athlete, 0)
		for _, a := range athletes {
			if strings.Contains(strings.ToLower(a.Name), strings.ToLower(searchQuery)) ||
				strings.Contains(strings.ToLower(a.Category), strings.ToLower(searchQuery)) {
				filtered = append(filtered, a)
			}
		}
		athletes = filtered
	}

	data := map[string]interface{}{
		"Athletes":    athletes,
		"SearchQuery": searchQuery,
		"CartCount":   h.repo.TeamApplication.TotalMembers,
		"PageTitle":   "Атлеты",
	}

	h.tmpl.ExecuteTemplate(w, "athletes", data)
}

func (h *Handler) EventsHandler(w http.ResponseWriter, r *http.Request) {
	events := h.repo.GetAllEvents()

	searchQuery := r.URL.Query().Get("search")
	if searchQuery != "" {
		filtered := make([]models.Event, 0)
		for _, e := range events {
			if strings.Contains(strings.ToLower(e.Name), strings.ToLower(searchQuery)) ||
				strings.Contains(strings.ToLower(e.Type), strings.ToLower(searchQuery)) ||
				strings.Contains(strings.ToLower(e.Location), strings.ToLower(searchQuery)) {
				filtered = append(filtered, e)
			}
		}
		events = filtered
	}

	data := map[string]interface{}{
		"Events":      events,
		"SearchQuery": searchQuery,
		"CartCount":   h.repo.TeamApplication.TotalMembers,
		"PageTitle":   "События",
	}

	h.tmpl.ExecuteTemplate(w, "events", data)
}

func (h *Handler) AthleteDetailHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Athlete ID required", http.StatusBadRequest)
		return
	}

	athlete, ok := h.repo.GetAthleteByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Athlete":   athlete,
		"CartCount": h.repo.TeamApplication.TotalMembers,
		"PageTitle": athlete.Name,
	}

	h.tmpl.ExecuteTemplate(w, "athlete-detail", data)
}

func (h *Handler) EventDetailHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Event ID required", http.StatusBadRequest)
		return
	}

	event, ok := h.repo.GetEventByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{}{
		"Event":     event,
		"CartCount": h.repo.TeamApplication.TotalMembers,
		"PageTitle": event.Name,
	}

	h.tmpl.ExecuteTemplate(w, "event-detail", data)
}

func (h *Handler) TeamApplicationHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Application": h.repo.TeamApplication,
		"CartCount":   h.repo.TeamApplication.TotalMembers,
		"PageTitle":   "Состав команды",
	}

	h.tmpl.ExecuteTemplate(w, "team-application", data)
}

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"CartCount": h.repo.TeamApplication.TotalMembers,
		"PageTitle": "Kinetic - Главная",
	}

	h.tmpl.ExecuteTemplate(w, "index", data)
}
