package main

import "time"

type rpcConfig struct {
	Log struct {
		LogLevel string `yaml:"logLevel" default:"info"`
	} `yaml:"log"`
	PerClient struct {
		Retries int `yaml:"retries" default:"3"`
		Delay   int `yaml:"delay" default:"5"`
		Timeout int `yaml:"timeout" default:"30"`
	} `yaml:"perClient"`
	Clients map[string]string `yaml:"clients" default:"[]"`
}

func (p *rpcConfig) getDelay() time.Duration {
	return time.Second * time.Duration(p.PerClient.Delay)
}

func (p *rpcConfig) getTimeout() time.Duration {
	return time.Second * time.Duration(p.PerClient.Timeout)
}
