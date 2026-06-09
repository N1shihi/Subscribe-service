package handler

import (
	"Subscribe-service/internal/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubsService interface {
	CreateSubs(subs models.Subscription) error
	GetSubs(userID uuid.UUID, serviceName string) (models.Subscription, error)
	UpdateSubs(subs models.Subscription) error
	DeleteSubs(userID uuid.UUID, serviceName string) error
	GetAllSubs() ([]models.Subscription, error)
	GetSubsSum(req models.AggregateRequest) (int, error)
}

// CreateSubsHandler godoc
//
//	@Summary		Создать подписку
//	@Description	Создает новую запись о подписке пользователя
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			subscription	body		models.Subscription	true	"Данные подписки"
//	@Success		201				{object}	models.Subscription
//	@Failure		400				{object}	map[string]interface{}
//	@Failure		409				{object}	map[string]interface{}
//	@Failure		500				{object}	map[string]interface{}
//	@Router			/subscriptions [post]
func (h *Handler) CreateSubs(c *gin.Context) {
	var subs models.Subscription

	log.Println("CREATE request")

	if err := c.ShouldBindJSON(&subs); err != nil {
		log.Printf("CREATE bad request %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.subsService.CreateSubs(subs); err != nil {
		log.Printf("CREATE error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusCreated, subs)
}

// GetAllSubs godoc
//
//	@Summary		Получить список подписок
//	@Description	Возвращает список всех подписок
//	@Tags			subscriptions
//	@Produce		json
//	@Success		200	{array}		models.Subscription
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/subscriptions [get]
func (h *Handler) GetAllSubs(c *gin.Context) {
	subs, err := h.subsService.GetAllSubs()
	if err != nil {
		log.Printf("GET ALL error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, subs)
}

// GetSubs godoc
//
//	@Summary		Получить подписку
//	@Description	Получает подписку по пользователю и названию сервиса
//	@Tags			subscriptions
//	@Produce		json
//	@Param			user_id			query		string	true	"UUID пользователя"
//	@Param			service_name	query		string	true	"Название сервиса"
//	@Success		200				{object}	models.Subscription
//	@Failure		400				{object}	map[string]interface{}
//	@Failure		404				{object}	map[string]interface{}
//	@Failure		500				{object}	map[string]interface{}
//	@Router			/subscriptions/item [get]
func (h *Handler) GetSubs(c *gin.Context) {
	userIDStr := c.Query("user_id")
	serviceName := c.Query("service_name")

	if userIDStr == "" || serviceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user_id and service_name required",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	sub, err := h.subsService.GetSubs(userID, serviceName)
	if err != nil {
		log.Printf("GET not found user %s service %s", userIDStr, serviceName)
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// UpdateSubs godoc
//
//	@Summary		Обновить подписку
//	@Description	Обновляет существующую подписку
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			subscription	body		models.Subscription	true	"Обновленная подписка"
//	@Success		200				{object}	models.Subscription
//	@Failure		400				{object}	map[string]interface{}
//	@Failure		404				{object}	map[string]interface{}
//	@Failure		500				{object}	map[string]interface{}
//	@Router			/subscriptions [put]
func (h *Handler) UpdateSubs(c *gin.Context) {
	var subs models.Subscription

	log.Println("UPDATE request")

	if err := c.ShouldBindJSON(&subs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.subsService.UpdateSubs(subs); err != nil {
		log.Printf("UPDATE error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, subs)
}

// DeleteSubs godoc
//
//	@Summary		Удалить подписку
//	@Description	Удаляет подписку по user_id и service_name
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			subscription	body		models.SubscriptionDelete	true	"Удаляемая подписка"
//	@Success		204
//	@Failure		400	{object}	map[string]interface{}
//	@Failure		404	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/subscriptions [delete]
func (h *Handler) DeleteSubs(c *gin.Context) {
	var req struct {
		UserID      string `json:"user_id"`
		ServiceName string `json:"service_name"`
	}

	log.Println("DELETE request")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	if err := h.subsService.DeleteSubs(userID, req.ServiceName); err != nil {
		log.Printf("DELETE error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetSubscSum godoc
//
//	@Summary		Получить суммарную стоимость подписок
//	@Description	Подсчитывает суммарную стоимость подписок за выбранный период с возможной фильтрацией по пользователю и сервису
//	@Tags			subscriptions
//	@Produce		json
//	@Param			user_id			query		string	false	"UUID пользователя"
//	@Param			service_name	query		string	false	"Название сервиса"
//	@Param			from			query		string	false	"Начало периода (MM-YYYY)"
//	@Param			to				query		string	false	"Конец периода (MM-YYYY)"
//	@Success		200				{object}	map[string]int64
//	@Failure		400				{object}	map[string]interface{}
//	@Failure		500				{object}	map[string]interface{}
//	@Router			/subscriptions/aggregate [get]
func (h *Handler) GetSubsSum(c *gin.Context) {
	req := models.AggregateRequest{
		UserID:      c.Query("user_id"),
		ServiceName: c.Query("service_name"),
		StartDate:   c.Query("start_date"),
		EndDate:     c.Query("end_date"),
	}

	log.Println("AGGREGATE request")

	sum, err := h.subsService.GetSubsSum(req)
	if err != nil {
		log.Printf("AGGREGATE error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": sum,
	})
}
