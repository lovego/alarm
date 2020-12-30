// Alarm合并报警邮件，防止在出错高峰，收到大量重复报警邮件，甚至因为邮件过多导致发送失败、接收失败。
package alarm

import (
	"sync"
	"time"
)

type Alarm struct {
	prefix        string
	sender        Sender
	min, inc, max time.Duration // 发送间隔时间最小值，增加值，最大值
	alarms        map[string]*alarm
	location      *time.Location
	sync.Mutex
}

type option func(*Alarm)

// SetPrefix set the prefix of email title
func SetPrefix(pre string) option {
	return func(a *Alarm) {
		a.prefix = pre + ` `
	}
}

// SetLocation set the location of the `StartAt` and `EndAt` of `Context`
func SetLocation(location *time.Location) option {
	return func(a *Alarm) {
		a.location = location
	}
}

func New(sender Sender, min, inc, max time.Duration, options ...option) *Alarm {
	a := &Alarm{
		location: time.Now().Location(),
		sender:   sender,
		min:      min,
		inc:      inc,
		max:      max,
		alarms:   make(map[string]*alarm),
	}
	for _, op := range options {
		op(a)
	}
	return a
}

func (alm *Alarm) Alarm(title, content, mergeKey string) {
	alm.Lock()
	a := alm.alarms[mergeKey]
	if a == nil {
		a = &alarm{Alarm: alm, interval: alm.min, inc: alm.inc, max: alm.max}
		alm.alarms[mergeKey] = a
	}
	alm.Unlock()

	a.add(title, content)
}

// logger will use
func (alm *Alarm) Send(title, content string) {
	alm.sender.Send(alm.prefix+title, content, Context{})
}

func (alm *Alarm) send(title, content string, ctx Context) {
	alm.sender.Send(alm.prefix+title, content, ctx)
}

func (alm *Alarm) SetDuration(min, inc, max time.Duration) {
	alm.min, alm.inc, alm.max = min, inc, max
}
