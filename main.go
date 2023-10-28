package main

import (
  "bytes"
  "fmt"
  "net/smtp"
  "text/template"
  "github.com/gin-gonic/gin"
  "net/http"
)

type EmailRequest struct {
  FromEmail    string `json:"fromEmail"`
  FromPassword string `json:"fromPassword"`
  ToEmail      string `json:"toEmail"`
  Subject string `json:"subject"`
  Body    string `json:"body"`
}




func getEmailRequest(c *gin.Context)  {
  email := EmailRequest{
    FromEmail:    c.Query("email"),
    FromPassword: c.Query("password"),
    ToEmail:      c.Query("to"),
    Subject: c.Query("subject"),
    Body:    c.Query("body"),
  }
  var temp = sendEmail(email);
  // Print the query parameters.
  // fmt.Println("name and age is: ",name, age)
  c.IndentedJSON(http.StatusOK, temp)
}

func main() {
  router := gin.Default()
  router.GET("/sendEmail", getEmailRequest)
  router.Run("localhost:8080")
}

func sendEmail(data EmailRequest) string {
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
  
    t.Execute(&body, struct {
      Name    string
      Message string
    }{
      Name:    "name",
      Message: "message",
    })
  
    // Sending email.
    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
    if err != nil {
      fmt.Println(err)
      return "Email Not Sent!, Error!: " + err.Error();
    }
    fmt.Println("Email Sent!")
    return "Email Sent!";
}