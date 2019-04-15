package models

// Record show best score/time of that 'difficult' lvl.
type Record struct {
	Score       int     `json:"score,omitempty"`
	Time        float64 `json:"time,omitempty"`
	Difficult   int     `json:"difficult,omitempty"`
	SingleTotal int     `json:"singleTotal,omitempty"`
	OnlineTotal int     `json:"onlineTotal,omitempty"`
	SingleWin   int     `json:"singleWin,omitempty"`
	OnlineWin   int     `json:"onlineWin,omitempty"`
}

func zeroOrOne(value *int) {
	if *value > 0 {
		*value = 1
	} else {
		*value = 0
	}
}

// Fix set fields of game amount to 0 or 1
func (record *Record) Fix() {
	zeroOrOne(&record.SingleTotal)
	zeroOrOne(&record.OnlineTotal)
	zeroOrOne(&record.SingleWin)
	zeroOrOne(&record.OnlineWin)
}
