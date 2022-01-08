package models

type Application struct {
	ID         string `json:"id,omitempty"`
	StackID    string `json:"stack_id,omitempty"`
	Parameters string `json:"parameters,omitempty"`
}
