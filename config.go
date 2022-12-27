package i3autotoggl

import (
	"fmt"
	"sync"

	"github.com/dlclark/regexp2"
)

type Config struct {
	Patterns []ConfigPattern `yaml:"patterns"`
}

type ConfigPattern struct {
	MatchTitle    string `yaml:"match"`
	ClassMatch    string `yaml:"match_class"`
	InstanceMatch string `yaml:"match_instance"`

	TaskName   string `yaml:"name"`
	ClientName string `yaml:"client"`
}

type CompiledConfig struct {
	Patterns []CompiledConfigPattern
}

type CompiledConfigPattern struct {
	Task
	Match CompiledConfigPatternMatches
}

type Task struct {
	Name   string
	Client string
}

type CompiledConfigPatternMatches struct {
	Class    *regexp2.Regexp
	Title    *regexp2.Regexp
	Instance *regexp2.Regexp
}

func CompileConfig(cfg Config) (CompiledConfig, error) {
	cc := CompiledConfig{}
	cc.Patterns = make([]CompiledConfigPattern, len(cfg.Patterns))

	for i, p := range cfg.Patterns {
		var match CompiledConfigPatternMatches

		r, err := regexp2.Compile(p.MatchTitle, regexp2.None)
		if err != nil {
			return CompiledConfig{}, fmt.Errorf("failed to compile match pattern: %w", err)
		}
		match.Title = r

		if p.ClassMatch != "" {
			ar, err := regexp2.Compile(p.ClassMatch, regexp2.None)
			if err != nil {
				return CompiledConfig{}, fmt.Errorf("failed to compile class match pattern: %w", err)
			}
			match.Class = ar
		}

		if p.InstanceMatch != "" {
			ar, err := regexp2.Compile(p.InstanceMatch, regexp2.None)
			if err != nil {
				return CompiledConfig{}, fmt.Errorf("failed to compile instance match pattern: %w", err)
			}
			match.Instance = ar
		}

		cc.Patterns[i].Match = match
		cc.Patterns[i].Task.Name = p.TaskName
		cc.Patterns[i].Task.Client = p.ClientName
	}
	return cc, nil
}

var (
	currentConfig   CompiledConfig
	currentConfigMu sync.Mutex
)

func LoadConfig(cfg CompiledConfig) {
	currentConfigMu.Lock()
	defer currentConfigMu.Unlock()
	currentConfig = cfg
}

func GetConfig() CompiledConfig {
	currentConfigMu.Lock()
	defer currentConfigMu.Unlock()
	return currentConfig
}
