// Alarm合并报警邮件，防止在出错高峰，收到大量重复报警邮件，甚至因为邮件过多导致发送失败、接收失败。
package alarm

import (
	"sync"
	"time"
)

type Alarm struct {
	prefix   string
	sender   Sender
	waits    []time.Duration // 报警等待时间
	alarms   map[string]*alarm
	location *time.Location
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

func New(sender Sender, waits []time.Duration, options ...option) *Alarm {
	if len(waits) == 0 {
		waits = []time.Duration{
			time.Second, 10 * time.Second,
			time.Minute, 10 * time.Minute,
			time.Hour,
		}
	}
	a := &Alarm{
		location: time.Now().Location(),
		sender:   sender,
		waits:    waits,
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
		a = &alarm{Alarm: alm}
		alm.alarms[mergeKey] = a
	}
	alm.Unlock()

	a.add(title, content)
}

// Send sent alarm directly without waiting. logger package uses this method.
func (alm *Alarm) Send(title, content string) {
	alm.sender.Send(alm.prefix+title, content, Context{})
}

func (alm *Alarm) send(title, content string, ctx Context) {
	alm.sender.Send(alm.prefix+title, content, ctx)
}

func (alm *Alarm) SetWaits(waits []time.Duration) {
	alm.waits = waits
}
