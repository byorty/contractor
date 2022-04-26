package graylog

import "time"

type Config struct {
	Url      string        `yaml:"url"`
	Username string        `yaml:"username"`
	Password string        `yaml:"password"`
	SpecDir  string        `yaml:"spec_dir"`
	Timeout  time.Duration `yaml:"timeout"`
}
