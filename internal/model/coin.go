package model

import "time"

type Coin struct {
	ID          string
	Name        string
	Ticker      string
	Price       float64
	Change24h   float64
	LastUpdated time.Time
}
