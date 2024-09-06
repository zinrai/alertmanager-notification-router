package domain

type Alert struct {
	Subject    string `json:"subject"`
	Body       string `json:"body"`
	Identifier string `json:"identifier"`
	Urgency    string `json:"urgency"`
}
