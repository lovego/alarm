package alarm

import (
	"log"
	"testing"
	"time"

	"github.com/lovego/mailer"
)

var testMailSender = getTestMailSender()

func getTestMailSender() Sender {
	theMailer, err := mailer.New(
		`mailer://smtp.qq.com:25/?user=小美<xiaomei-go@qq.com>&pass=zjsbosjlhgugechh`,
	)
	if err != nil {
		log.Panic(err)
	}

	return MailSender{
		Receivers: []string{"applejava<applejava@qq.com>"},
		Mailer:    theMailer,
	}
}

func TestMailAlarm(t *testing.T) {
	alarm := New(`alarm测试`, testMailSender, 0, time.Second, 10*time.Second)
	alarm.Alarm(`title`, `content`, `mergeKey`)
	time.Sleep(3 * time.Second) // wait the alarms been sent
}
