package jev

import "encoding/json"

type SystemOneRequest struct {
	Model     string         `json:"model"`
	State     any            `json:"state"`
	Questions map[string]any `json:"questions"`
}

type SystemOneResponse struct {
	Model   string                     `json:"model"`
	Answers map[string]json.RawMessage `json:"answers"`
	Usage   Usage                      `json:"usage"`
}

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type ChoiceQuestion struct {
	Type         string            `json:"type"`
	Instructions string            `json:"instructions,omitempty"`
	Criteria     map[string]string `json:"criteria"`
}

type ChoiceAnswer struct {
	Type          string             `json:"type"`
	Choice        string             `json:"choice"`
	Confidence    float64            `json:"confidence"`
	Probabilities map[string]float64 `json:"probabilities"`
}

type NoulQuestion struct {
	Type         string `json:"type"`
	Instructions string `json:"instructions,omitempty"`
}

type NoulAnswer struct {
	Type string  `json:"type"`
	Noul float64 `json:"noul"`
}
