package models

type Event struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Name        string  `json:"name"`
	Date        string  `json:"date"`
	Location    string  `json:"location"`
	Price       float64 `json:"price"`
	ImageKey    string  `json:"image_key"`
	VideoKey    string  `json:"video_key"`
	Description string  `json:"description"`
}

type AthleteEvent struct {
	Event   Event
	Athlete Athlete
	PB      string `json:"pb"`
	Stats   string `json:"stats"`
}
