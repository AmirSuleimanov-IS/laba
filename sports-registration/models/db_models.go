package models

// Пользователь
type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	CreatedAt string `json:"created_at"`
}

// Услуга (спортивное событие)
type Service struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Status      string  `json:"status"` // "active", "deleted"
	ImageURL    *string `json:"image_url"`
	ServiceType string  `json:"service_type"`
	EventDate   *string `json:"event_date"`
	Location    string  `json:"location"`
	Price       float64 `json:"price"`
	CreatedAt   string  `json:"created_at"`
}

// Заявка
type Application struct {
	ID           int      `json:"id"`
	UserID       int      `json:"user_id"`
	Status       string   `json:"status"` // "draft", "deleted", "formed", "completed", "rejected"
	CreatedAt    string   `json:"created_at"`
	FormedAt     *string  `json:"formed_at"`
	CompletedAt  *string  `json:"completed_at"`
	ModeratorID  *int     `json:"moderator_id"`
	TeamName     string   `json:"team_name"`
	TotalAmount  float64  `json:"total_amount"`
	Services     []AppService `json:"services,omitempty"`
}

// Связь заявки и услуги (M:M)
type AppService struct {
	ApplicationID int     `json:"application_id"`
	ServiceID     int     `json:"service_id"`
	Quantity      int     `json:"quantity"`
	SortOrder     int     `json:"sort_order"`
	IsPrimary     bool    `json:"is_primary"`
	AddedAt       string  `json:"added_at"`
	Service       Service `json:"service,omitempty"`
}
