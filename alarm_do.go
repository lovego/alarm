package alarm

import (
	"fmt"
	"sync"
	"time"
)

type alarm struct {
	sync.Mutex
	title, content string
	count          int
	interval       time.Duration
	*Alarm
}

func (a *alarm) add(title, content string) {
	a.Lock()
	a.count += 1
	count := a.count
	a.Unlock()

	if count == 1 {
		a.title, a.content = title, content
		go a.send()
	}
}

func (a *alarm) send() {
	time.Sleep(a.interval)
	a.interval += a.Alarm.inc
	if a.interval > a.Alarm.max {
		a.interval = a.Alarm.max
	}
	a.Lock()
	count, title, content := a.count, a.title, a.content
	a.count = 0
	a.Unlock()
	if count > 1 {
		title += fmt.Sprintf(` [Merged: %d]`, count)
	}
	a.Alarm.send(title, content)
}
