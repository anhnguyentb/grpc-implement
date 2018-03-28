package models

import "time"

type Audit struct {
	Id       int64
	ClientIp string
	ServerIp string
	Message  string
	Tags     []string
	Created  time.Time `sql:"default:now()"`
}