package alarm

import (
	"sync"
	"time"
)

// Alarm合并报警邮件，防止在出错高峰，收到大量重复报警邮件，
// 甚至因为邮件过多导致发送失败、接收失败。
type Alarm struct {
	prefix        string
	sender        Sender
	min, inc, max time.Duration // 发送间隔时间最小值，增加值，最大值
	alarms        map[string]*alarm
	sync.Mutex
}

func New(
	prefix string, sender Sender, min, inc, max time.Duration,
) *Alarm {
	return &Alarm{
		prefix: prefix,
		sender: sender,
		min:    min,
		inc:    inc,
		max:    max,
		alarms: make(map[string]*alarm),
	}
}

func (alm *Alarm) Alarm(title, content, mergeKey string) {
	alm.Lock()
	a := alm.alarms[mergeKey]
	if a == nil {
		a = &alarm{Alarm: alm, interval: alm.min}
		alm.alarms[mergeKey] = a
	}
	alm.Unlock()

	a.add(title, content)
}

func (alm *Alarm) send(title, content string) {
	alm.sender.Send(alm.prefix+title, content)
}
