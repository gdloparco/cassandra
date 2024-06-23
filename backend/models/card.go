package models

// struct used to store information coming from external tarot deck API
type Card struct {
	CardName       string `json:"name"`
	Type           string `json:"type"`
	MeaningUp      string `json:"meaning_up"`
	MeaningReverse string `json:"meaning_rev"`
	Description    string `json:"desc"`
	ShortName      string `json:"name_short"`
	Reversed       bool   `json:"reversed"`
}
