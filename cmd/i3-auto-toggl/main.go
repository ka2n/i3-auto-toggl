package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/fatih/color"
	"github.com/k0kubun/pp"
	i3autotoggl "github.com/ka2n/i3-auto-toggle"
	"github.com/samber/lo"
)

var (
	flagConfig = flag.String("c", "", "config file")
)

func main() {
	flag.Parse()

	if err := mainCLI(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func mainCLI(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	cfgPath, err := findConfigFile(*flagConfig)
	if err != nil {
		return fmt.Errorf("failed to find config file: %w", err)
	}

	// Load initial config
	cfg, err := readConfigFile(cfgPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	i3autotoggl.LoadConfig(cfg)
	log.Println("config loaded")

	// Watch config
	cfgWatch, err := watchaFile(ctx, cfgPath, func() {
		cfg, err := readConfigFile(cfgPath)
		if err != nil {
			log.Println(fmt.Errorf("failed to read config: %w", err))
		}
		i3autotoggl.LoadConfig(cfg)

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

	ws := i3autotoggl.NewWindowEventWatcher()
	defer ws.Close()

	stream := lo.FanIn(0, i3autotoggl.DetectIdle(ctx), ws.Watch())

	var lastEvent i3autotoggl.TimelineEvent
	for {
		select {
		case ev, ok := <-stream:
			if !ok {
				return nil
			}

			var prev i3autotoggl.TimelineEvent
			prev, lastEvent = lastEvent, ev

			if prev.EqualTarget(ev) {
				continue
			}

			switch ev.Type {
			case i3autotoggl.TimelineEvent_Idle:
				fmt.Printf("%s since: %s", color.BlueString("Idle"), ev.Time.Format(time.Kitchen))

			case i3autotoggl.TimelineEvent_Start:
				cfg := i3autotoggl.GetConfig()
				if match := i3autotoggl.Match(cfg, ev); match != nil {
					fmt.Printf("%s, event: %+v\n", color.CyanString(pp.Sprint(match.Task)), ev)
				}
			}

		case <-ctx.Done():
			return nil
		}
	}
}
