package main

import "time"

type rpcConfig struct {
	Log       *rpcLogConfig
	Clients   map[string]string `yaml:"clients" default:"[]"`
	PerClient *rpcPerClientConfig
}

type rpcLogConfig struct {
	LogLevel string `yaml:"logLevel" default:"info"`
}

type rpcPerClientConfig struct {
	Retries int `yaml:"retries" default:"3"`
	Delay   int `yaml:"delay" default:"5"`
	Timeout int `yaml:"timeout" default:"30"`
}

func (p *rpcPerClientConfig) getDelay() time.Duration {
	return time.Second * time.Duration(p.Delay)
}

func (p *rpcPerClientConfig) getTimeout() time.Duration {
	return time.Second * time.Duration(p.Timeout)
}
