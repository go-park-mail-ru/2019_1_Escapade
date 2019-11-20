package models

// Record show best score/time of that 'difficult' lvl.
//easyjson:json
type Record struct {
	Score       int     `json:"score,omitempty" minimum:"0"`
	Time        float64 `json:"time,omitempty"`
	Difficult   int     `json:"difficult,omitempty" minimum:"0"`
	SingleTotal int     `json:"singleTotal,omitempty" minimum:"0"`
	OnlineTotal int     `json:"onlineTotal,omitempty" minimum:"0"`
	SingleWin   int     `json:"singleWin,omitempty" minimum:"0"`
	OnlineWin   int     `json:"onlineWin,omitempty" minimum:"0"`
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
