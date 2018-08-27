package alarm

import (
	"fmt"
	"log"
	"time"

	"github.com/lovego/email"
	"github.com/lovego/mailer"
)

type Context struct {
	Count   int
	StartAt time.Time
	EndAt   time.Time
}

type Sender interface {
	Send(title, content string, ctx Context)
}

type MailSender struct {
	Receivers []string
	Mailer    *mailer.Mailer
}

func (m MailSender) Send(title, content string, ctx Context) {
	if len(m.Receivers) == 0 {
		return
	}
	if ctx.Count > 1 {
		title = fmt.Sprintf("%s [merged: %d, time: %s-%s]", title, ctx.Count, inTime(ctx.StartAt), inTime(ctx.EndAt))
	}
	err := m.Mailer.Send(&email.Email{
		To:      m.Receivers,
		Subject: title,
		Text:    []byte(content),
	}, time.Minute)
	if err != nil {
		log.Printf("send alarm mail failed: %v", err)
	}
}

func inTime(t time.Time) string {
	return fmt.Sprintf("%v:%v:%v", t.Hour(), t.Minute(), t.Second())
}
