package alarm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lovego/email"
	"github.com/lovego/mailer"
)

type Context struct {
	SentCount   int
	MergedCount int
	StartedAt   time.Time
	EndedAt     time.Time
}

func (ctx Context) String() string {
	var strs = []string{fmt.Sprintf("[#%d", ctx.SentCount+1), "", ""}

	if ctx.MergedCount > 1 {
		strs[1] = fmt.Sprintf(" merged:%d", ctx.MergedCount)
	}

	var time string
	if start, end := ctx.StartedAt.Format("15:04:05"), ctx.EndedAt.Format("15:04:05"); start == end {
		time = start
	} else {
		time = fmt.Sprintf("%s~%s", start, end)
	}
	strs[2] = fmt.Sprintf(" %s]", time)
	return strings.Join(strs, "")
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
	title = ctx.String() + title
	err := m.Mailer.Send(&email.Email{
		To:      m.Receivers,
		Subject: title,
		Text:    []byte(content),
	}, time.Minute)
	if err != nil {
		log.Printf("send alarm mail failed: %v", err)
	}
}
