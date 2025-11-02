package notify

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strconv"
)

// EmailNotifier sends notifications via SMTP email.
type EmailNotifier struct {
	Service  string // e.g., "QQ", "GMAIL"
	Email    string // sender and receiver
	Password string // SMTP password or app token
	Name     string // display name
}

// NewEmailNotifier creates an EmailNotifier from environment variables.
// It reads:
//   - SMTP_SERVICE
//   - SMTP_EMAIL
//   - SMTP_PASSWORD
//   - SMTP_NAME (optional, defaults to email)
func NewEmailNotifier() (*EmailNotifier, error) {
	service := os.Getenv("SMTP_SERVICE")
	email := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	if service == "" || email == "" || password == "" {
		return nil, fmt.Errorf("missing env vars: SMTP_SERVICE, SMTP_EMAIL, or SMTP_PASSWORD")
	}

	name := os.Getenv("SMTP_NAME")
	if name == "" {
		name = email
	}

	return &EmailNotifier{
		Service:  service,
		Email:    email,
		Password: password,
		Name:     name,
	}, nil
}

func (e *EmailNotifier) getSMTPConfig() (host string, port int, useSSL bool, err error) {
	switch e.Service {
	case "QQ":
		return "smtp.qq.com", 465, true, nil
	case "GMAIL":
		return "smtp.gmail.com", 587, false, nil
	case "163":
		return "smtp.163.com", 465, true, nil
	case "OUTLOOK":
		return "smtp-mail.outlook.com", 587, false, nil
	default:
		return "", 0, false, fmt.Errorf("unsupported SMTP service: %s", e.Service)
	}
}

// Send sends an email to itself (sender == receiver).
func (e *EmailNotifier) Send(subject, body string) error {
	host, port, useSSL, err := e.getSMTPConfig()
	if err != nil {
		return err
	}

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	auth := smtp.PlainAuth("", e.Email, e.Password, host)

	headers := map[string]string{
		"From":         fmt.Sprintf("%s <%s>", e.Name, e.Email),
		"To":           e.Email,
		"Subject":      subject,
		"Content-Type": "text/plain; charset=UTF-8",
	}

	msg := ""
	for k, v := range headers {
		msg += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	msg += "\r\n" + body

	var client *smtp.Client

	if useSSL {
		// SSL/TLS 直连（端口 465）
		conn, err := tls.Dial("tcp", addr, &tls.Config{
			ServerName: host,
		})
		if err != nil {
			return fmt.Errorf("failed to connect via SSL: %w", err)
		}
		defer conn.Close()

		client, err = smtp.NewClient(conn, host)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client over SSL: %w", err)
		}
	} else {
		// 普通连接 + STARTTLS（端口 587）
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to dial: %w", err)
		}
		defer conn.Close()

		client, err = smtp.NewClient(conn, host)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}

		if ok, _ := client.Extension("STARTTLS"); ok {
			if err = client.StartTLS(&tls.Config{ServerName: host}); err != nil {
				return fmt.Errorf("STARTTLS failed: %w", err)
			}
		}
	}

	defer client.Quit()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	if err = client.Mail(e.Email); err != nil {
		return fmt.Errorf("MAIL command failed: %w", err)
	}

	if err = client.Rcpt(e.Email); err != nil {
		return fmt.Errorf("RCPT command failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA command failed: %w", err)
	}

	_, err = w.Write([]byte(msg))
	w.Close()
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}
