package notifier

import (
	"Puff/internal/config"
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

type DomainNotification struct {
	Domain        string
	IsFinalNotice bool
	Status        string
}

func SendNotification(notifications []DomainNotification, cfg *config.Config) error {
	auth := smtp.PlainAuth("", cfg.SMTPUsername, cfg.SMTPPassword, cfg.SMTPServer)

	to := []string{cfg.RecipientEmail}
	subject := "域名状态变更提醒"

	body := generateEmailBody(notifications)

	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", cfg.SMTPUsername, cfg.RecipientEmail, subject, body))

	err := smtp.SendMail(fmt.Sprintf("%s:%d", cfg.SMTPServer, cfg.SMTPPort), auth, cfg.SMTPUsername, to, msg)
	if err != nil {
		return fmt.Errorf("发送邮件失败: %v\n服务器: %s:%d\n发件人: %s\n收件人: %s",
			err, cfg.SMTPServer, cfg.SMTPPort, cfg.SMTPUsername, cfg.RecipientEmail)
	}

	return nil
}

func generateEmailBody(notifications []DomainNotification) string {
	var body strings.Builder

	body.WriteString(`
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { width: 100%; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #4CAF50; color: white; padding: 10px; text-align: center; }
        .content { padding: 20px; background-color: #f9f9f9; }
        .footer { text-align: center; font-size: 0.8em; color: #777; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>域名状态变更提醒</h1>
        </div>
        <div class="content">
            <p>尊敬的用户，</p>
            <p>以下域名的状态发生了变化：</p>
            <ul>
    `)

	for _, n := range notifications {
		body.WriteString(fmt.Sprintf("<li>%s: %s", n.Domain, n.Status))
		if n.IsFinalNotice {
			body.WriteString(" (最终通知)")
		}
		body.WriteString("</li>")
	}

	body.WriteString(fmt.Sprintf(`
            </ul>
            <p>如果您对这些域名感兴趣，请尽快采取相应的行动。</p>
            <p>检测时间：%s</p>
        </div>
        <div class="footer">
            <p>此邮件由 Puff 自动发送，请勿直接回复。</p>
        </div>
    </div>
</body>
</html>
    `, time.Now().Format("2006年01月02日 15:04:05")))

	return body.String()
}
