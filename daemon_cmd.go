package i3autotoggl

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/fatih/color"
	"github.com/k0kubun/pp"
	"github.com/samber/lo"
)

func StartDaemonCmd(ctx context.Context, configFilePath string) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	cfgPath, err := findConfigFile(configFilePath)
	if err != nil {
		return fmt.Errorf("failed to find config file: %w", err)
	}

	// Load initial config
	cfg, err := readConfigFile(cfgPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	LoadConfig(cfg)
	log.Println("config loaded")

	// Watch config
	cfgWatch, err := watchaFile(ctx, cfgPath, func() {
		cfg, err := readConfigFile(cfgPath)
		if err != nil {
			log.Println(fmt.Errorf("failed to read config: %w", err))
		}
		LoadConfig(cfg)

	})
	if err != nil {
		return fmt.Errorf("failed to watch config file: %w", err)
	}
	defer cfgWatch.Close()

	log.Println("Connecting to i3")
	return WatchWindow(ctx)
}

func WatchWindow(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ws := NewWindowEventWatcher()
	defer ws.Close()

	stream := lo.FanIn(0, DetectIdle(ctx), ws.Watch())

	var lastEvent TimelineEvent
	for {
		select {
		case ev, ok := <-stream:
			if !ok {
				return nil
			}

			var prev TimelineEvent
			prev, lastEvent = lastEvent, ev

			if prev.EqualTarget(ev) {
				continue
			}

			switch ev.Type {
			case TimelineEvent_Idle:
				fmt.Printf("%s since: %s", color.BlueString("Idle"), ev.Time.Format(time.Kitchen))

			case TimelineEvent_Start:
				cfg := GetConfig()
				if match := Match(cfg, ev); match != nil {
					fmt.Printf("%s, event: %+v\n", color.CyanString(pp.Sprint(match.Task)), ev)
				}
			}

		case <-ctx.Done():
			return nil
		}
	}
}
