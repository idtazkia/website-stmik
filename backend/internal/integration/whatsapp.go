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

// WhatsAppClient handles WhatsApp messaging via custom API
type WhatsAppClient struct {
	APIURL   string
	APIToken string
	client   *http.Client
}

// NewWhatsAppClient creates a new WhatsApp client
// Returns nil if not configured (APIURL or APIToken is empty)
func NewWhatsAppClient(apiURL, apiToken string) *WhatsAppClient {
	if apiURL == "" || apiToken == "" {
		log.Println("WhatsApp not configured: phone OTP will not be available")
		return nil
	}
	return &WhatsAppClient{
		APIURL:   apiURL,
		APIToken: apiToken,
		client:   &http.Client{Timeout: 30 * time.Second},
	}
}

// IsConfigured returns true if the client is properly configured
func (c *WhatsAppClient) IsConfigured() bool {
	return c != nil && c.APIURL != "" && c.APIToken != ""
}

type whatsAppRequest struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

// SendOTP sends an OTP code via WhatsApp
func (c *WhatsAppClient) SendOTP(phone, otp string) error {
	if c == nil {
		return fmt.Errorf("whatsapp client not configured")
	}

	message := fmt.Sprintf(`*STMIK Tazkia*

Kode verifikasi pendaftaran Anda adalah:

*%s*

Kode ini berlaku selama 15 menit.

Jika Anda tidak meminta kode ini, abaikan pesan ini.`, otp)

	reqBody := whatsAppRequest{
		Phone:   phone,
		Message: message,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal whatsapp request: %w", err)
	}

	req, err := http.NewRequest("POST", c.APIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIToken)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send whatsapp message: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("whatsapp API error (status %d): %s", resp.StatusCode, string(body))
	}

	log.Printf("WhatsApp OTP sent to %s", phone)
	return nil
}
