package i3autotoggl

import (
	"sync"
	"time"

	"github.com/samber/lo"
	"go.i3wm.org/i3/v4"
)

const (
	WINDOW_FOCUS_THRESHOLD = 1 * time.Second
)

type WindowEventWatcher struct {
	receiver *i3.EventReceiver
}

func NewWindowEventWatcher() (w WindowEventWatcher) {
	w.receiver = i3.Subscribe(i3.WindowEventType)
	return w
}

func (w WindowEventWatcher) Close() {
	w.receiver.Close()
}

func (w WindowEventWatcher) Watch() (in <-chan TimelineEvent) {
	return readEvents(w.receiver)
}

// ReadEvents reads events from i3 and sends them to the out channel.
// It will close the out channel when the i3 connection is closed.
func readEvents(ws *i3.EventReceiver) (in <-chan TimelineEvent) {
	out := make(chan TimelineEvent)
	go readEventsLoop(ws, out)
	return out
}

func readEventsLoop(ws *i3.EventReceiver, in chan<- TimelineEvent) {
	var event i3.Event
	var eventMu sync.Mutex

	defer close(in)

	f := func() {
		eventMu.Lock()
		defer eventMu.Unlock()
		ev := event

		switch t := ev.(type) {
		case *i3.WindowEvent:
			switch t.Change {
			case "title":
				if !t.Container.Focused {
					return
				}
				fallthrough
			case "focus":
				in <- TimelineEvent{
					Type:         TimelineEvent_Start,
					Time:         time.Now(),
					Title:        t.Container.WindowProperties.Title,
					ClassName:    t.Container.WindowProperties.Class,
					InstanceName: t.Container.WindowProperties.Instance,
				}
			}
		default:
		}
	}

	send, cancel := lo.NewDebounce(WINDOW_FOCUS_THRESHOLD, f)
	defer cancel()

	for {
		if !ws.Next() {
			return
		}

		eventMu.Lock()
		event = ws.Event()
		eventMu.Unlock()
		send()
	}
}
