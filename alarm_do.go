package alarm

import (
	"sync"
	"time"
)

type alarm struct {
	sync.Mutex
	title     string
	content   string
	count     int
	interval  time.Duration
	inc       time.Duration
	max       time.Duration
	startedAt time.Time
	*Alarm
}

func (a *alarm) add(title, content string) {
	a.Lock()
	a.count += 1
	count := a.count
	a.Unlock()

	if count == 1 {
		a.title, a.content = title, content
		a.startedAt = time.Now().In(a.Alarm.location)
		go a.send()
	}
}

func (a *alarm) send() {
	time.Sleep(a.interval)

	a.interval += a.inc
	if a.interval > a.max {
		a.interval = a.max
	}

	a.Lock()
	count, title, content, startedAt := a.count, a.title, a.content, a.startedAt
	a.count = 0
	a.Unlock()

	ctx := Context{
		Count:   count,
		StartAt: startedAt,
		EndAt:   time.Now().In(a.Alarm.location),
	}
	a.Alarm.send(title, content, ctx)
}
