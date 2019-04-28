package models

// UserPrivateInfo shows personal player info
type UserPrivateInfo struct {
	ID       int    `json:"-"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Update update all fields
func (confirmed *UserPrivateInfo) Update(another *UserPrivateInfo) {
	updateParameter(&confirmed.Name, another.Name)
	updateParameter(&confirmed.Email, another.Email)
	updateParameter(&confirmed.Password, another.Password)
}

func updateParameter(
	to *string, from string) {
	if *to != from && from != "" {
		*to = from
	}
	return
}
