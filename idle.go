package i3autotoggl

import (
	"context"
	"log"
	"time"

	"github.com/ka2n/go-idle"
)

const (
	IDLE_CHECK_INTERVAL = time.Minute
	IDLE_THRESHOLD      = 5 * time.Minute
)

func DetectIdle(ctx context.Context) (in <-chan TimelineEvent) {
	out := make(chan TimelineEvent)
	go detectIdleLoop(ctx, out)
	return out
}

func detectIdleLoop(ctx context.Context, in chan<- TimelineEvent) {
	ticker := time.NewTicker(IDLE_CHECK_INTERVAL)
	defer ticker.Stop()
	defer close(in)

	for {
		select {
		case <-ticker.C:
			duration, err := idle.Get()
			if err != nil {
				log.Println("failed to get idle duration:", err)
				continue
			}

			if duration > IDLE_THRESHOLD {
				in <- TimelineEvent{
					Type: TimelineEvent_Idle,
					Time: time.Now().Add(-1 * duration), // Idle Since
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
