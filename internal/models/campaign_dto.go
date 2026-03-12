package models

type CreateCampaignRequest struct {
	Subject    string   `json:"subject" binding:"required"`
	Body       string   `json:"body" binding:"required"`
	Recipients []string `json:"recipients" binding:"required,min=1"`
}

type CampaignResponse struct {
	ID          uint   `json:"id"`
	Subject     string `json:"subject"`
	Status      string `json:"status"`
	TotalEmails int    `json:"total_emails"`
}
