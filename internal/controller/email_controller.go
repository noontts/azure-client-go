package controller

import (
	"azureclient/client/azure"
	"encoding/json"
	"net/http"
)

type EmailController struct {
   Client azure.ISendEmailClient
}

type SendEmailRequest struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	PlainText string `json:"plain_text"`
	HTML      string `json:"html"`
}

func NewEmailController(client azure.ISendEmailClient) *EmailController {
   return &EmailController{Client: client}
}

// POST /send-email
func (c *EmailController) SendEmail(w http.ResponseWriter, r *http.Request) error {
	var req SendEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}
	azureReq := azure.SendEmailRequest{
		Sender:    req.Sender,
		Recipient: req.Recipient,
		Subject:   req.Subject,
		PlainText: req.PlainText,
		HTML:      req.HTML,
	}
	if err := c.Client.SendEmail(r.Context(), azureReq); err != nil {
		return err
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "sent"})
	return nil
}
