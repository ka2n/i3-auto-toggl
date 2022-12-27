package i3autotoggl

func Match(cfg CompiledConfig, ev TimelineEvent) *CompiledConfigPattern {
	for _, p := range cfg.Patterns {
		if p.Match.Instance != nil {
			if match, _ := p.Match.Instance.MatchString(ev.InstanceName); !match {
				continue
			}
		}
		if p.Match.Class != nil {
			if match, _ := p.Match.Class.MatchString(ev.ClassName); !match {
				continue
			}
		}
		if match, _ := p.Match.Title.MatchString(ev.Title); match {
			return &p
		}
	}
	return nil
}
