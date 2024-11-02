package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail(to, subject, body string) error {
    from := os.Getenv("SMTP_USER")
    password := os.Getenv("SMTP_PASS")
    smtpHost := os.Getenv("SMTP_HOST")
    smtpPort := os.Getenv("SMTP_PORT")
    
    msg := fmt.Sprintf("From: %s\r\n"+
        "To: %s\r\n"+
        "Subject: %s\r\n"+
        "\r\n"+
        "%s\r\n", from, to, subject, body)
    
    auth := smtp.PlainAuth("", from, password, smtpHost)
    
    return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
}