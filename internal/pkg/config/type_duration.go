package config

import (
	json "encoding/json"
	"fmt"
	"time"
)

// Duration override time.Duration for json marshalling/unmarshalling
type Duration struct {
	time.Duration
}

// UnmarshalJSON unmarshal JSON to Duration
func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' {
		sd := string(b[1 : len(b)-1])
		d.Duration, err = time.ParseDuration(sd)
		return
	}

	var id int64
	id, err = json.Number(string(b)).Int64()
	d.Duration = time.Duration(id)

	return
}

// MarshalJSON marshal Duration to Duration
func (d Duration) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}
