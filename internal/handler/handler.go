package handler

type Handler struct {
	subscService SubsService
}

func New(subscService SubscService) *Handler {
	return &Handler{
		subsService: SubsService,
	}
}
