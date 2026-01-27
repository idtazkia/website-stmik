package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
)

//go:embed *.html
var templateFS embed.FS

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseFS(templateFS, "*.html")
	if err != nil {
		panic(fmt.Sprintf("failed to parse email templates: %v", err))
	}
}

// OTPData holds data for OTP email
type OTPData struct {
	OTP string
}

// PaymentConfirmedData holds data for payment confirmation email
type PaymentConfirmedData struct {
	CandidateName string
	BillingType   string
	Amount        string
	TransferDate  string
	ApprovedAt    string
}

// PaymentRejectedData holds data for payment rejection email
type PaymentRejectedData struct {
	CandidateName string
	BillingType   string
	Amount        string
	TransferDate  string
	Reason        string
}

// DocumentApprovedData holds data for document approval email
type DocumentApprovedData struct {
	CandidateName string
	DocumentType  string
}

// DocumentRejectedData holds data for document rejection email
type DocumentRejectedData struct {
	CandidateName string
	DocumentType  string
	Reason        string
}

// RenderOTP renders the OTP email template
func RenderOTP(data OTPData) (string, error) {
	return render("otp.html", data)
}

// RenderPaymentConfirmed renders the payment confirmation email template
func RenderPaymentConfirmed(data PaymentConfirmedData) (string, error) {
	return render("payment_confirmed.html", data)
}

// RenderPaymentRejected renders the payment rejection email template
func RenderPaymentRejected(data PaymentRejectedData) (string, error) {
	return render("payment_rejected.html", data)
}

// RenderDocumentApproved renders the document approval email template
func RenderDocumentApproved(data DocumentApprovedData) (string, error) {
	return render("document_approved.html", data)
}

// RenderDocumentRejected renders the document rejection email template
func RenderDocumentRejected(data DocumentRejectedData) (string, error) {
	return render("document_rejected.html", data)
}

func render(templateName string, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, templateName, data); err != nil {
		return "", fmt.Errorf("failed to render template %s: %w", templateName, err)
	}
	return buf.String(), nil
}
