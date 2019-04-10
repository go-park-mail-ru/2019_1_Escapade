package models

// Record show best score/time of that 'difficult' lvl.
type Record struct {
	Score       int `json:"score"`
	Time        int `json:"time"`
	Difficult   int `json:"difficult"`
	SingleTotal int `json:"singleTotal"`
	OnlineTotal int `json:"onlineTotal"`
	SingleWin   int `json:"singleWin"`
	OnlineWin   int `json:"onlineWin"`
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
