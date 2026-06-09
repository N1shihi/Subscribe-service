package handler

type Handler struct {
	subsService SubsService
}

func New(subsService SubsService) *Handler {
	return &Handler{
		subsService: subsService,
	}
}
