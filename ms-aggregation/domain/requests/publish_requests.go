package requests

type PublishRequest struct {
	Topic   string      `json:"topic" validate:"required"`
	Message interface{} `json:"message" validate:"required"`
}
