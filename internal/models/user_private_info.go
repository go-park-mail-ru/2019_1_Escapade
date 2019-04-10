package models

// UserPrivateInfo shows personal player info
type UserPrivateInfo struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Update update all fields
func (upi *UserPrivateInfo) Update(name string, email string, password string) {
	upi.Name = name
	upi.Email = email
	upi.Password = password
}
