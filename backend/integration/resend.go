package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
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

	subject := "Kode Verifikasi Pendaftaran STMIK Tazkia"
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #194189;">STMIK Tazkia</h2>
        <p>Kode verifikasi email Anda adalah:</p>
        <div style="background-color: #f5f5f5; padding: 20px; text-align: center; margin: 20px 0;">
            <span style="font-size: 32px; font-weight: bold; letter-spacing: 8px; color: #194189;">%s</span>
        </div>
        <p>Kode ini berlaku selama 15 menit.</p>
        <p>Jika Anda tidak meminta kode ini, abaikan email ini.</p>
        <hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
        <p style="font-size: 12px; color: #666;">
            Email ini dikirim secara otomatis. Mohon tidak membalas email ini.
        </p>
    </div>
</body>
</html>
`, otp)

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

	log.Printf("Email OTP sent to %s (resend id: %s)", to, response.ID)
	return nil
}
