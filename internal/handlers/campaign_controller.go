package handlers

import (
	"net/http"
	"strconv"

	"github.com/Shubham23jha/digital-post-office/internal/models"
	"github.com/Shubham23jha/digital-post-office/internal/service"
	"github.com/gin-gonic/gin"
)

type CampaignHandler struct {
	service service.CampaignService
}

func NewCampaignHandler(s service.CampaignService) *CampaignHandler {
	return &CampaignHandler{service: s}
}

func (h *CampaignHandler) Create(c *gin.Context) {
	var req models.CreateCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	campaign, err := h.service.CreateCampaign(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, campaign)
}

func (h *CampaignHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	campaign, err := h.service.GetCampaign(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign not found"})
		return
	}

	c.JSON(http.StatusOK, campaign)
}
func (h *CampaignHandler) List(c *gin.Context) {
	campaigns, err := h.service.GetAllCampaigns()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, campaigns)
}
