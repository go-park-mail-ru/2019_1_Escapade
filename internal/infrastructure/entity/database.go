package entity

import "time"

type Database struct {
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}
