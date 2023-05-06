package common

import (
	"github.com/dustin/go-humanize"
	"time"
)

var (
	uptime time.Time
)

func StartTime() {
	uptime = time.Now()
}

func GetUptime() string {
	h := humanize.Time(uptime)
	return h
}
