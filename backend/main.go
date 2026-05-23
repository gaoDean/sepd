package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/resend/resend-go/v3"
)

type ResendClient struct {
	client *resend.Client
}

type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

const emailTemplate = `
<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
	<h2 style="color: #000; margin-bottom: 20px;">New Contact Form Submission</h2>
	<div style="background-color: #f9f9f9; padding: 20px; border-radius: 8px; margin-bottom: 20px;">
		<p><strong>Name:</strong> {{.Name}}</p>
		<p><strong>Email:</strong> {{.Email}}</p>
	</div>
	<div style="background-color: #f9f9f9; padding: 20px; border-radius: 8px;">
		<h3 style="margin-top: 0; color: #000;">Message:</h3>
		<p style="white-space: pre-wrap; line-height: 1.6;">{{.Message}}</p>
	</div>
	<hr style="margin: 30px 0; border: none; border-top: 1px solid #ddd;">
	<p style="color: #666; font-size: 14px;">
		This message was sent from the SEPD contact form.
	</p>
</div>`

func (rclient *ResendClient) contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var req ContactRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON payload"})
		return
	}

	if req.Name == "" || req.Email == "" || req.Message == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing required fields"})
		return
	}

	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		log.Println("Template parsing error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var bodyBuffer bytes.Buffer
	if err := tmpl.Execute(&bodyBuffer, req); err != nil {
		log.Println("Template execution error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	params := &resend.SendEmailRequest{
		From:    "contact@emails.besspower.com.au",
		To:      []string{"gaof@sepd.com.au"},
		Html:    bodyBuffer.String(),
		Subject: "New contact form submission from " + req.Name,
	}

	sent, err := rclient.client.Emails.Send(params)
	if err != nil {
		log.Println("Resend API error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to send email"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sent)
}

func main() {
	mux := http.NewServeMux()

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on system environment")
	}

	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		log.Fatal("RESEND_API_KEY environment variable is required")
	}

	rclient := &ResendClient{
		client: resend.NewClient(apiKey),
	}

	mux.HandleFunc("/api/contact", rclient.contactHandler)

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
