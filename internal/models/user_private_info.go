package models

// UserPrivateInfo shows personal player info
type UserPrivateInfo struct {
	ID       int    `json:"-"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Update update all fields
func (confirmed *UserPrivateInfo) Update(another *UserPrivateInfo) {
	another.ID = confirmed.ID
	updateParameter(&another.Name, confirmed.Name)
	updateParameter(&another.Password, confirmed.Password)
}

func updateParameter(
	to *string, from string) {
	if *to != from && *to == "" {
		*to = from
	}
	return
}
