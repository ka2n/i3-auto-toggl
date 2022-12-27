package i3autotoggl

import "time"

type TimelineEventType int

const (
	TimelineEvent_Start TimelineEventType = iota
	TimelineEvent_Idle
)

type TimelineEvent struct {
	Type TimelineEventType

	Time         time.Time
	Title        string
	ClassName    string
	InstanceName string
}

func (t TimelineEvent) EqualTarget(b TimelineEvent) bool {
	return t.Type == b.Type && t.Title == b.Title && t.ClassName == b.ClassName && t.InstanceName == b.InstanceName
}
