package alarm

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var redundant = 100 * time.Millisecond
var mutex sync.Mutex

type testSenderT []string

func (ts *testSenderT) Send(title, content string, ctx Context) {
	mutex.Lock()
	*ts = append(*ts, fmt.Sprintf("标题%s [Merged: %d, Time: %s-%s] 内容%s",
		title, ctx.Count, inTime(ctx.StartAt), inTime(ctx.EndAt), content))
	mutex.Unlock()
}

func (ts *testSenderT) equal(target []string) bool {
	if len(*ts) != len(target) {
		return false
	}
	m := map[string]bool{}
	for _, str := range *ts {
		m[str] = true
	}
	for _, str := range target {
		if !m[str] {
			return false
		}
	}
	return true
}

func (ts *testSenderT) empty() {
	*ts = []string{}
}

func TestAlarm(t *testing.T) {
	sender := &testSenderT{}
	min := 500 * time.Millisecond
	inc := time.Second
	max := 5 * time.Second
	location := time.Now().Location()
	alarm := New(sender, min, inc, max, SetLocation(location))

	startAt := time.Now().In(location)
	endAt := startAt.Add(min)
	startAtInTime := inTime(startAt)
	endAtInTime := inTime(endAt)

	sender.empty()
	sendAlarms(alarm, map[string]int{`a`: 3, `b`: 4, `c`: 5})
	time.Sleep(min + redundant)
	assertEqual(t, sender, []string{
		fmt.Sprintf("标题a [Merged: 3, Time: %s-%s] 内容a", startAtInTime, endAtInTime),
		fmt.Sprintf(`标题b [Merged: 4, Time: %s-%s] 内容b`, startAtInTime, endAtInTime),
		fmt.Sprintf(`标题c [Merged: 5, Time: %s-%s] 内容c`, startAtInTime, endAtInTime),
	})

	startAt = time.Now().In(location)
	endAt = startAt.Add(min + inc)
	startAtInTime = inTime(startAt)
	endAtInTime = inTime(endAt)

	sender.empty()
	sendAlarms(alarm, map[string]int{`a`: 3, `b`: 4, `c`: 5})
	time.Sleep(inc + min + redundant)
	assertEqual(t, sender, []string{
		fmt.Sprintf("标题a [Merged: 3, Time: %s-%s] 内容a", startAtInTime, endAtInTime),
		fmt.Sprintf(`标题b [Merged: 4, Time: %s-%s] 内容b`, startAtInTime, endAtInTime),
		fmt.Sprintf(`标题c [Merged: 5, Time: %s-%s] 内容c`, startAtInTime, endAtInTime),
	})
}

func sendAlarms(alarm *Alarm, alarms map[string]int) {
	var wg sync.WaitGroup
	for mergeKey, count := range alarms {
		wg.Add(count)
		for i := 0; i < count; i++ {
			go func(mergeKey string) {
				defer wg.Done()
				alarm.Alarm(mergeKey, mergeKey, mergeKey)
			}(mergeKey)
		}
	}
	wg.Wait()
}

func assertEqual(t *testing.T, s *testSenderT, expect []string) {
	if !s.equal(expect) {
		t.Errorf("expect: %q\ngot: %q\n", expect, s)
	}
}
