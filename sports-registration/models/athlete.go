package models

type Athlete struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Category string `json:"category"`
	Age      int    `json:"age"`
	Country  string `json:"country"`
	Races    int    `json:"races"`
	Wins     int    `json:"wins"`
	Rating   int    `json:"rating"`
	Bio      string `json:"bio"`
	Location string `json:"location"`
	ImageKey string `json:"image_key"`
}
