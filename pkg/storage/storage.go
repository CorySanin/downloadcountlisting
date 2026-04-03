package storage

import (
	"time"
)

type View struct {
	Prefix    string
	Filename  string
	Timestamp time.Time
	Removed   bool
}
