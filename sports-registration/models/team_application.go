package models

type TeamApplication struct {
	ID           string       `json:"id"`
	TeamName     string       `json:"team_name"`
	CaptainID    string       `json:"captain_id"`
	Members      []TeamMember `json:"members"`
	TotalMembers int          `json:"total_members"`
	MaxMembers   int          `json:"max_members"`
	Status       string       `json:"status"`
	CreatedAt    string       `json:"created_at"`
}

type TeamMember struct {
	AthleteID string `json:"athlete_id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	EventID   string `json:"event_id"`
	EventType string `json:"event_type"`
	ImageKey  string `json:"image_key"`
}
