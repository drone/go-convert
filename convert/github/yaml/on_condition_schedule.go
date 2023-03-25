package yaml

type ScheduleCondition struct {
	Cron []string `yaml:"cron,omitempty"`
}

func (s *ScheduleCondition) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var items []map[string]string
	if err := unmarshal(&items); err != nil {
		var m map[string]string
		if err2 := unmarshal(&m); err2 != nil {
			return err
		}
		if cron, ok := m["cron"]; ok {
			s.Cron = []string{cron}
		}
		return nil
	}
	s.Cron = make([]string, len(items))
	for i, item := range items {
		if cron, ok := item["cron"]; ok {
			s.Cron[i] = cron
		}
	}
	return nil
}
