package handler

import (
	"Subscribe-service/internal/models"
	"net/http"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type SubsService interface {
	CreateSubs(team models.Subscription) error
	GetSubs(subsId uuid.UUID) (models.Subscription, error)
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

func (h *Handler) GetTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "team_name required",
			},
		})
		return
	}

	team, err := h.teamService.GetTeam(teamName)
	if err != nil {
		if err.Error() == "team not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "NOT_FOUND",
					"message": "team not found",
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

	c.JSON(http.StatusOK, team)
}
