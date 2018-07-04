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

func (ts *testSenderT) Send(title, content string) {
	mutex.Lock()
	*ts = append(*ts, fmt.Sprintf("标题%s 内容%s", title, content))
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
	inc := 100 * time.Millisecond
	max := 1 * time.Second
	alarm := New("", sender, min, inc, max)

	sender.empty()
	sendAlarms(alarm, map[string]int{`a`: 3, `b`: 4, `c`: 5})
	time.Sleep(min + redundant)
	assertEqual(t, sender, []string{
		`标题a [Merged: 3] 内容a`, `标题b [Merged: 4] 内容b`, `标题c [Merged: 5] 内容c`,
	})

	sender.empty()
	sendAlarms(alarm, map[string]int{`a`: 3, `b`: 4, `c`: 5})
	time.Sleep(inc + min + redundant)
	assertEqual(t, sender, []string{
		`标题a [Merged: 3] 内容a`, `标题b [Merged: 4] 内容b`, `标题c [Merged: 5] 内容c`,
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
