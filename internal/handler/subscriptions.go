package handler

import (
	"Subscribe-service/internal/models"
	"net/http"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type SubsService interface {
	CreateSubs(team models.Subscription) error
	GetSubs(userId uuid.UUID) (models.Subscription, error)
	DeleteSubs(subs models.SubscriptionDelete) error
}

func (h *Handler) CreateSubsHandler(c *gin.Context) {
	var subs models.Subscription
	if err := c.ShouldBindJSON(&subs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	if err := h.subsService.CreateSubs(subs); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "internal server error",
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"team": subs,
	})
}

func (h *Handler) GetSubs(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "team_name required",
			},
		})
		return
	}

	subscriptions, err := h.subsService.GetSubs(userId)

	// поменять проверки
	if err != nil {
		if err.Error() == "subs not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "NOT_FOUND",
					"message": "subs not found",
				},
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "internal server error",
				},
			})
		}
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}
func (h *Handler) DeleteSubs(c *gin.Context) {
	var subs models.SubscriptionDelete
	if err := c.ShouldBindJSON(&subs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": err.Error(),
			},
		})
		return
	}

	err := h.subsService.DeleteSubs(subs)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.Status(200)
}
