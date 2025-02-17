package main

import (
	"fmt"
	"github.com/wneessen/go-mail"
	"log"
)

func sendMail(userEmail string, userToken string) error {
	m := mail.NewMsg()
	m.From("mail@gmail.com")
	m.To(userEmail)
	m.Subject("InterFlow Forum - Подтверждение регистрации \n")
	link := fmt.Sprintf("https://forum.com/confirm?token=%s", userToken)
	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<style>
				body { font-family: Arial, sans-serif; }
				.container { max-width: 600px; margin: auto; padding: 20px; border: 1px solid #ddd; border-radius: 8px; }
				h1 { color: #333; }
				p { font-size: 16px; }
				.button {
					display: inline-block;
					padding: 10px 20px;
					color: white;
					background-color: #007bff;
					text-decoration: none;
					border-radius: 5px;
					font-weight: bold;
				}
				.button:hover { background-color: #0056b3; }
			</style>
		</head>
		<body>
			<div class="container">
				<h1>Здравствуйте!</h1>
				<p>Вы зарегистрировались на нашем форуме. Чтобы подтвердить регистрацию, нажмите на кнопку ниже:</p>
				<p><a class="button" href="%s">Подтвердить регистрацию</a></p>
				<p>Если кнопка не работает, скопируйте и вставьте эту ссылку в браузер:</p>
				<p><a href="%s">%s</a></p>
				<p>С уважением, <br> Администрация форума</p>
			</div>
		</body>
		</html>`, link, link, link)

	m.SetBodyString(mail.TypeTextHTML, htmlBody)

	client, err := mail.NewClient(
		"smtp.gmail.com",
		mail.WithPort(587),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername("mail@gmail.com"),
		mail.WithPassword("your_app_password"),
		mail.WithTLSPolicy(mail.TLSMandatory),
	)
	if err != nil {
		return fmt.Errorf("ошибка подключения к SMTP: %v", err)
	}
	if err := client.DialAndSend(m); err != nil {
		return fmt.Errorf("ошибка отправки письма: %v", err)
	}
	log.Println("Письмо успешно отправлено на", userEmail)
	return nil
}
