package alarm

import (
	"sync"
	"time"
)

type alarm struct {
	*Alarm
	sync.Mutex

	title       string
	content     string
	startedAt   time.Time
	endedAt     time.Time
	mergedCount int
	sentCount   int
}

func (a *alarm) add(title, content string) {
	a.Lock()
	defer a.Unlock()

	a.mergedCount += 1
	a.endedAt = time.Now().In(a.Alarm.location)
	if a.mergedCount == 1 {
		a.title, a.content, a.startedAt = title, content, time.Now().In(a.Alarm.location)
		go a.send()
	}
}

func (a *alarm) send() {
	a.Lock()
	sentCount := a.sentCount
	a.Unlock()

	if sentCount >= len(a.Alarm.waits) {
		sentCount = len(a.Alarm.waits) - 1
	}
	time.Sleep(a.Alarm.waits[sentCount])

	a.Lock()
	title, content := a.title, a.content
	ctx := Context{
		SentCount:   a.sentCount,
		MergedCount: a.mergedCount,
		StartedAt:   a.startedAt,
		EndedAt:     a.endedAt,
	}
	a.mergedCount = 0
	a.sentCount++
	a.Unlock()

	a.Alarm.send(title, content, ctx)
}
