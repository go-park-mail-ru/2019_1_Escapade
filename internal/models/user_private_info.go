package models

// UserPrivateInfo shows personal player info
//easyjson:json
type UserPrivateInfo struct {
	ID       int    `json:"-"`
	Name     string `json:"name" maxLength:"30" example:"John" `
	Password string `json:"password" minLength:"6" maxLength:"30" example:"easyPassword" `
}

// Update godoc
// update all fields
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

func (user *UserPrivateInfo) GetName() string {
	return user.Name
}

func (user *UserPrivateInfo) GetPassword() string {
	return user.Password
}

func (user *UserPrivateInfo) SetName(name string) {
	user.Name = name
}

func (user *UserPrivateInfo) SetPassword(password string) {
	user.Password = password
}
