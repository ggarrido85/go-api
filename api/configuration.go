package api

import "time"

type Configuration struct {
	Env                 string
	AppName             string
	Port                string
	AppUrl        string
	BackofficeUrl string
	RequestLoggingLevel string
	TokenLifetimeMinute int
	SegmentWriteKey     string
	DefaultTimeout  time.Duration
}
