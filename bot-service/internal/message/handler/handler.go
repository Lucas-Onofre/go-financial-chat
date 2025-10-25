package handler

type Handler struct {
	service any
}

func New(service any) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandleMessage(message string) string {
	return "Handled message: " + message
}
