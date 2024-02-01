package database

import "time"

var (
	TimeoutOneDay   = 24 * 60 * time.Second
	Timeout30min    = 30 * 60 * time.Second
	Timeout15sec    = 15 * time.Second
	TimeoutOneBlock = 6 * time.Second
)
