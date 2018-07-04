package alarm

import (
	"time"

	"github.com/lovego/email"
	"github.com/lovego/mailer"
)

type Sender interface {
	Send(title, content string)
}

type MailSender struct {
	Receivers []string
	Mailer    *mailer.Mailer
}

func (m MailSender) Send(title, content string) {
	if len(m.Receivers) == 0 {
		return
	}
	m.Mailer.Send(&email.Email{
		To:      m.Receivers,
		Subject: title,
		Text:    []byte(content),
	}, time.Minute)
}
