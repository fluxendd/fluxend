package email

type EmailInterface interface {
	Send(to, subject, body string) error
}
