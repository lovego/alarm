package alarm

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lovego/email"
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
	Mailer    *email.Client
}

func (m MailSender) Send(title, content string, ctx Context) {
	if len(m.Receivers) == 0 {
		return
	}
	title = ctx.String() + title
	sendCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := m.Mailer.Send(sendCtx, &email.Message{
		Headers: []email.Header{
			{Name: email.To, Values: m.Receivers},
			{Name: email.Subject, Values: []string{title}},
		},
		Body: []byte(content),
	})
	if err != nil {
		log.Printf("send alarm mail failed: %v", err)
	}
}
