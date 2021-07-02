package requests

type PublishRandomRequest struct {
	Topic string `json:"topic" validate:"required"`
}
