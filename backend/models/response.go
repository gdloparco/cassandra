package models

// struct use to store information of cards that will be sent to the frontend
type JSONCard struct {
	CardName       string `json:"name"`
	Type           string `json:"type"`
	MeaningUp      string `json:"meaning_up"`
	MeaningReverse string `json:"meaning_rev"`
	Description    string `json:"desc"`
	ImageName      string `json:"image_file_name"`
	Reversed       bool   `json:"reversed"`
}
