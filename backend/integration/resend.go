package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/idtazkia/stmik-admission-api/templates/email"
)

// ResendClient handles email sending via Resend API
type ResendClient struct {
	APIKey string
	From   string
	client *http.Client
}

// NewResendClient creates a new Resend client
// Returns nil if not configured (APIKey or From is empty)
func NewResendClient(apiKey, from string) *ResendClient {
	if apiKey == "" || from == "" {
		log.Println("Resend not configured: email OTP will not be available")
		return nil
	}
	return &ResendClient{
		APIKey: apiKey,
		From:   from,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// IsConfigured returns true if the client is properly configured
func (c *ResendClient) IsConfigured() bool {
	return c != nil && c.APIKey != "" && c.From != ""
}

type resendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

type resendEmailResponse struct {
	ID string `json:"id"`
}

// SendOTP sends an OTP code via email
func (c *ResendClient) SendOTP(to, otp string) error {
	if c == nil {
		return fmt.Errorf("resend client not configured")
	}

	html, err := email.RenderOTP(email.OTPData{OTP: otp})
	if err != nil {
		return fmt.Errorf("failed to render OTP email: %w", err)
	}

	return c.sendEmail(to, "Kode Verifikasi Pendaftaran STMIK Tazkia", html)
}

// PaymentConfirmationData holds data for payment confirmation email
type PaymentConfirmationData struct {
	CandidateName string
	BillingType   string
	Amount        string
	TransferDate  string
	ApprovedAt    string
}

// SendPaymentConfirmation sends a payment confirmation email
func (c *ResendClient) SendPaymentConfirmation(to string, data PaymentConfirmationData) error {
	if c == nil {
		return fmt.Errorf("resend client not configured")
	}

	html, err := email.RenderPaymentConfirmed(email.PaymentConfirmedData{
		CandidateName: data.CandidateName,
		BillingType:   data.BillingType,
		Amount:        data.Amount,
		TransferDate:  data.TransferDate,
		ApprovedAt:    data.ApprovedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to render payment confirmation email: %w", err)
	}

	return c.sendEmail(to, "Pembayaran Anda Telah Dikonfirmasi - STMIK Tazkia", html)
}

// PaymentRejectionData holds data for payment rejection email
type PaymentRejectionData struct {
	CandidateName string
	BillingType   string
	Amount        string
	TransferDate  string
	Reason        string
}

// SendPaymentRejection sends a payment rejection notification email
func (c *ResendClient) SendPaymentRejection(to string, data PaymentRejectionData) error {
	if c == nil {
		return fmt.Errorf("resend client not configured")
	}

	html, err := email.RenderPaymentRejected(email.PaymentRejectedData{
		CandidateName: data.CandidateName,
		BillingType:   data.BillingType,
		Amount:        data.Amount,
		TransferDate:  data.TransferDate,
		Reason:        data.Reason,
	})
	if err != nil {
		return fmt.Errorf("failed to render payment rejection email: %w", err)
	}

	return c.sendEmail(to, "Pembayaran Memerlukan Perbaikan - STMIK Tazkia", html)
}

// DocumentStatusData holds data for document status email
type DocumentStatusData struct {
	CandidateName string
	DocumentType  string
	Status        string // "approved" or "rejected"
	Reason        string // Only for rejected
}

// SendDocumentApproved sends a document approval notification email
func (c *ResendClient) SendDocumentApproved(to string, data DocumentStatusData) error {
	if c == nil {
		return fmt.Errorf("resend client not configured")
	}

	html, err := email.RenderDocumentApproved(email.DocumentApprovedData{
		CandidateName: data.CandidateName,
		DocumentType:  data.DocumentType,
	})
	if err != nil {
		return fmt.Errorf("failed to render document approval email: %w", err)
	}

	return c.sendEmail(to, "Dokumen Anda Telah Diverifikasi - STMIK Tazkia", html)
}

// SendDocumentRejected sends a document rejection notification email
func (c *ResendClient) SendDocumentRejected(to string, data DocumentStatusData) error {
	if c == nil {
		return fmt.Errorf("resend client not configured")
	}

	html, err := email.RenderDocumentRejected(email.DocumentRejectedData{
		CandidateName: data.CandidateName,
		DocumentType:  data.DocumentType,
		Reason:        data.Reason,
	})
	if err != nil {
		return fmt.Errorf("failed to render document rejection email: %w", err)
	}

	return c.sendEmail(to, "Dokumen Memerlukan Perbaikan - STMIK Tazkia", html)
}

// sendEmail is a helper to send an email
func (c *ResendClient) sendEmail(to, subject, html string) error {
	reqBody := resendEmailRequest{
		From:    c.From,
		To:      []string{to},
		Subject: subject,
		HTML:    html,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal email request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("resend API error (status %d): %s", resp.StatusCode, string(body))
	}

	var response resendEmailResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	log.Printf("Email sent to %s (resend id: %s)", to, response.ID)
	return nil
}
