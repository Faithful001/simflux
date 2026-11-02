package dto

import "encoding/json"

type WebhookRequest struct {
	/// payment fields
	Reference string         `json:"reference" binding:"required"`
	Status    string         `json:"status" binding:"required"`
	Fields    map[string]any `json:"-"`
}

func (w WebhookRequest) MarshalJSON() ([]byte, error) {
	request := make(map[string]any)

	for key, value := range w.Fields {
		request[key] = value
	}

	request["reference"] = w.Reference
	request["status"] = w.Status

	return json.Marshal(request)

}