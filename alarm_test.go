package alarm

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

const redundant = 100 * time.Millisecond

type testSenderT struct {
	msgs []string
	sync.Mutex
}

func (ts *testSenderT) Send(title, content string, ctx Context) {
	t := time.Now().UTC().Round(24 * time.Hour)
	ctx.StartedAt, ctx.EndedAt = t, t
	ts.Lock()
	ts.msgs = append(ts.msgs, fmt.Sprintf("%s 标题%s 内容%s", ctx.String(), title, content))
	ts.Unlock()
}

func (ts *testSenderT) printMsgs() {
	sort.Strings(ts.msgs)
	for _, msg := range ts.msgs {
		fmt.Println(msg)
	}
}

func (ts *testSenderT) empty() {
	ts.msgs = []string{}
}

func ExampleAlarm_test() {
	sender := &testSenderT{}
	waits := []time.Duration{500 * time.Millisecond}
	alarm := New(sender, waits)

	sender.empty()
	sendAlarms(alarm, map[string]int{`a`: 1, `b`: 2, `c`: 3})
	time.Sleep(waits[0] + redundant)
	sender.printMsgs()

	sender.empty()
	sendAlarms(alarm, map[string]int{`a`: 1, `b`: 2, `c`: 3})
	time.Sleep(waits[0] + redundant)
	sender.printMsgs()
	// Output:
	// [#1 0:0:0] 标题a 内容a
	// [#1 merged:2 0:0:0] 标题b 内容b
	// [#1 merged:3 0:0:0] 标题c 内容c
	// [#2 0:0:0] 标题a 内容a
	// [#2 merged:2 0:0:0] 标题b 内容b
	// [#2 merged:3 0:0:0] 标题c 内容c
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
