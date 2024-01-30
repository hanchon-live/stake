package database

import "time"

var (
	TimeoutLong     = 30 * 60 * time.Second
	TimeoutShort    = 15 * time.Second
	TimeoutOneBlock = 6 * time.Second
)
