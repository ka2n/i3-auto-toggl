package i3autotoggl

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

func findConfigFile(byFlag string) (string, error) {
	if byFlag != "" {
		p, err := filepath.Abs(byFlag)
		if err != nil {
			return "", fmt.Errorf("failed to get absolute path of config file: %w", err)
		}
		return p, nil
	} else {
		p, err := xdg.SearchConfigFile("i3-auto-toggl/i3-auto-toggl.yml")
		if err != nil {
			return "", fmt.Errorf("failed to search config file: %w", err)
		}
		return p, nil
	}
}

func watchaFile(ctx context.Context, path string, cb func()) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	dir := filepath.Dir(path)
	fname := filepath.Base(path)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					if filepath.Base(event.Name) == fname {
						cb()
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			case <-ctx.Done():
				return
			}
		}
	}()

	log.Println("Watching config file: ", path)
	err = watcher.Add(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to watch config file: %w", err)
	}
	return watcher, nil
}

func readConfigFile(configPath string) (CompiledConfig, error) {
	var cfg Config
	b, err := os.ReadFile(configPath)
	if err != nil {
		return CompiledConfig{}, fmt.Errorf("failed to read config file: %w", err)
	}
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return CompiledConfig{}, fmt.Errorf("failed to parse config as YAML: %w", err)
	}

	ccfg, err := CompileConfig(cfg)
	if err != nil {
		return CompiledConfig{}, fmt.Errorf("failed to parse config: %w", err)
	}

	return ccfg, nil
}
