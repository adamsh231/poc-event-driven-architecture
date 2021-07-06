package requests

type PublishRequest struct {
	Topic     string      `json:"topic" validate:"required"`
	Partition int32       `json:"partition"`
	Message   interface{} `json:"message" validate:"required"`
}
