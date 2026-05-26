package handler

import (
	"html/template"
	"net/http"
	"sports-registration/models"
	"sports-registration/repository"
	"strings"
)

type Handler struct {
	repo *repository.Repository
	tmpl *template.Template
}

func NewHandler(repo *repository.Repository) *Handler {
	h := &Handler{repo: repo}
	h.loadTemplates()
	return h
}

func (h *Handler) loadTemplates() {
	funcMap := template.FuncMap{}
	
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
	))
	
	h.tmpl = tmpl
}

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
