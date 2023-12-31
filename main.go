package main

import (
	"bytes"
	// "encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"text/template"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type EmailRequest struct {
	FromEmail    string `json:"fromEmail"`
	FromPassword string `json:"fromPassword"`
	ToEmail      string `json:"toEmail"`
	Subject      string `json:"subject"`
	Body         string `json:"body"`
	CampaignId   string `json:"campaignId"`
	TemplateId   string `json:"templateId"`
}

type TriggerResponse struct {
	ClientEmail   string `json:"clientEmail"`
	TargetEmail   string `json:"targetEmail"`
	TemplateId    string `json:"templateId"`
	DateTime      string `json:"dateTime"`
	Status        string `json:"status"`
	CampaignId    string `json:"campaignId"`
	NumberOfOpens int    `json:"numberOfOpens"`
}

func getEmailRequest(c *gin.Context) {
	// c = new(gin.Context)
	var email EmailRequest
	c.BindJSON(&email)
	// Print the query parameters.
	// fmt.Println("name and age is: ",name, age)
	c.IndentedJSON(http.StatusOK, sendEmail(email, true))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	fmt.Println("Starting server...")
	router := gin.Default()
	router.GET("/sendEmail", getEmailRequest)
	router.GET("/emailTrigger", onEmailOpen)
	router.GET("/", func(c *gin.Context) {
		var temp EmailRequest
		c.BindJSON(&temp)

		c.IndentedJSON(http.StatusOK, gin.H{
			"message": temp,
			"err":     err,
		})
	})
	router.Run(
	// For run on specific port
	)
}

func onEmailOpen(c *gin.Context) {

	trigger := TriggerResponse{
		ClientEmail:   c.Query("email"),
		TargetEmail:   c.Query("to"),
		TemplateId:    c.Query("templateId"),
		DateTime:      time.Now().Format(time.RFC1123),
		Status:        c.Query("status"),
		CampaignId:    c.Query("campaignId"),
		NumberOfOpens: 1,
	}
	fmt.Println(trigger)
	fmt.Println("password is: ", os.Getenv("EMAIL_PASSWORD"))
	// add data to mongodb
	// Print the query parameters.
	// onEmailOpen

	response := sendEmail(
		EmailRequest{
			FromEmail:    trigger.ClientEmail,
			FromPassword: os.Getenv("EMAIL_PASSWORD"),
			ToEmail:      trigger.TargetEmail,
			Subject:      "EMAIL_Opened",
			Body:         "templateID:_" + trigger.TemplateId + "_opened_at_:" + time.Now().Format(time.RFC850) + "_____status_:" + trigger.Status + "____campaignId_:" + trigger.CampaignId,
			CampaignId:   trigger.CampaignId,
			TemplateId:   trigger.TemplateId,
		},
		false,
	)
	c.IndentedJSON(http.StatusOK, response)
}

func sendEmail(data EmailRequest, TrackerActive bool) string {
	// Sender data.
	from := data.FromEmail
	password := data.FromPassword

	// Receiver email address.
	to := []string{
		data.ToEmail,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("template.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject:"+data.Subject+" \n%s\n\n", mimeHeaders)))
	baseUrl := "http://localhost:8080"
	if os.Getenv("ENV") == "prod" {
		baseUrl = os.Getenv("BASE_URL")
	}
	if TrackerActive {
		t.Execute(&body, struct {
			Name    string
			Message string
			Tracker string
		}{
			Name:    "name",
			Message: data.Body,
			Tracker: "<img src=" + baseUrl + "/emailTrigger" +
				"?email=" + data.FromEmail +
				"&to=" + data.ToEmail +
				"&subject=" + data.Subject +
				"&body=" + data.Body +
				"&campaignId=" + data.CampaignId +
				"&templateId=" + data.TemplateId +
				"&status=opened" +
				" hidden>",
		})
	} else {
		t.Execute(&body, struct {
			Name    string
			Message string
			Tracker string
		}{
			Name:    "name",
			Message: data.Body,
			Tracker: "",
		})
	}

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return "Email Not Sent!, Error!: " + err.Error()
	}
	fmt.Println("Email Sent!")
	return "Email Sent!"
}
