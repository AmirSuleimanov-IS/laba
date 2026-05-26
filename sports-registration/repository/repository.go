package repository

import (
	"sports-registration/models"
)

type Repository struct {
	Athletes      map[string]models.Athlete
	Events        map[string]models.Event
	AthleteEvents []models.AthleteEvent
	TeamApplication models.TeamApplication
}

func NewRepository() *Repository {
	repo := &Repository{
		Athletes: make(map[string]models.Athlete),
		Events:   make(map[string]models.Event),
	}

	repo.initAthletes()
	repo.initEvents()
	repo.initTeamApplication()

	return repo
}

func (r *Repository) initAthletes() {
	r.Athletes["1"] = models.Athlete{
		ID:       "1",
		Name:     "Алексей Волков",
		Status:   "Свободен для команд",
		Category: "PRO",
		Age:      29,
		Country:  "Россия",
		Races:    42,
		Wins:     15,
		Rating:   4,
		Bio:      "Алексей — один из самых перспективных ультрамарафонцев страны. Специализируется на технически сложных горных маршрутах с большим набором высоты.",
		Location: "Красная Поляна, Сочи",
		ImageKey: "1.png",
	}

	r.Athletes["2"] = models.Athlete{
		ID:       "2",
		Name:     "Александр Соколов",
		Status:   "Капитан команды",
		Category: "PRO",
		Age:      31,
		Country:  "Россия",
		Races:    38,
		Wins:     12,
		Rating:   6,
		Location: "Москва",
		ImageKey: "2.png",
	}

	r.Athletes["3"] = models.Athlete{
		ID:       "3",
		Name:     "Михаил Волков",
		Status:   "В команде",
		Category: "Amateur",
		Age:      27,
		Country:  "Россия",
		Races:    25,
		Wins:     8,
		Rating:   3,
		Location: "Санкт-Петербург",
		ImageKey: "3.png",
	}

	r.Athletes["4"] = models.Athlete{
		ID:       "4",
		Name:     "Екатерина Морозова",
		Status:   "В команде",
		Category: "PRO",
		Age:      28,
		Country:  "Россия",
		Races:    35,
		Wins:     10,
		Rating:   5,
		Location: "Казань",
		ImageKey: "4.png",
	}

	r.Athletes["5"] = models.Athlete{
		ID:       "5",
		Name:     "Дмитрий Петров",
		Status:   "Свободен для команд",
		Category: "Amateur",
		Age:      32,
		Country:  "Россия",
		Races:    20,
		Wins:     5,
		Rating:   2,
		Location: "Екатеринбург",
		ImageKey: "5.png",
	}
}

func (r *Repository) initEvents() {
	r.Events["1"] = models.Event{
		ID:          "1",
		Type:        "Марафон",
		Name:        "Московский марафон 2026",
		Date:        "2026-05-15",
		Location:    "Москва",
		Price:       2500.00,
		Description: "Ежегодный городской марафон с дистанциями 42км, 21км, 10км и 5км.",
		ImageKey:    "1.png",
		VideoKey:    "video.mp4",
	}

	r.Events["2"] = models.Event{
		ID:          "2",
		Type:        "Триатлон",
		Name:        "Ironman Sochi",
		Date:        "2026-06-20",
		Location:    "Сочи",
		Price:       15000.00,
		Description: "Полный триатлон: плавание 3.8км, велогонка 180км, бег 42км.",
		ImageKey:    "2.png",
		VideoKey:    "video.mp4",
	}

	r.Events["3"] = models.Event{
		ID:          "3",
		Type:        "Трейл",
		Name:        "Elbrus World Race",
		Date:        "2026-07-10",
		Location:    "Приэльбрусье",
		Price:       8000.00,
		Description: "Горный трейл с набором высоты 3500м на дистанции 50км.",
		ImageKey:    "3.png",
		VideoKey:    "video.mp4",
	}

	r.Events["4"] = models.Event{
		ID:          "4",
		Type:        "Велогонка",
		Name:        "Gran Fondo Russia",
		Date:        "2026-08-05",
		Location:    "Сочи",
		Price:       5000.00,
		Description: "Массовая велогонка на дистанциях 120км и 60км.",
		ImageKey:    "4.png",
		VideoKey:    "video.mp4",
	}
}

func (r *Repository) initTeamApplication() {
	r.TeamApplication = models.TeamApplication{
		ID:           "app-001",
		TeamName:     "Kinetic Elite",
		CaptainID:    "2",
		TotalMembers: 3,
		MaxMembers:   5,
		Status:       "active",
		CreatedAt:    "2026-04-01",
		Members: []models.TeamMember{
			{
				AthleteID: "2",
				Name:      "Александр Соколов",
				Role:      "Капитан команды",
				EventID:   "1",
				EventType: "Марафон",
				ImageKey:  "2.png",
			},
			{
				AthleteID: "3",
				Name:      "Михаил Волков",
				Role:      "Участник",
				EventID:   "2",
				EventType: "Триатлон",
				ImageKey:  "3.png",
			},
			{
				AthleteID: "4",
				Name:      "Екатерина Морозова",
				Role:      "Участник",
				EventID:   "3",
				EventType: "Трейл",
				ImageKey:  "4.png",
			},
		},
	}
}

func (r *Repository) GetAllAthletes() []models.Athlete {
	athletes := make([]models.Athlete, 0, len(r.Athletes))
	for _, athlete := range r.Athletes {
		athletes = append(athletes, athlete)
	}
	return athletes
}

func (r *Repository) GetAllEvents() []models.Event {
	events := make([]models.Event, 0, len(r.Events))
	for _, event := range r.Events {
		events = append(events, event)
	}
	return events
}

func (r *Repository) GetAthleteByID(id string) (models.Athlete, bool) {
	athlete, ok := r.Athletes[id]
	return athlete, ok
}

func (r *Repository) GetEventByID(id string) (models.Event, bool) {
	event, ok := r.Events[id]
	return event, ok
}
