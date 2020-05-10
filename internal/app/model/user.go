package model

import "golang.org/x/crypto/bcrypt"

// User ...
type User struct {
	ID                int    `json:"id"`
	Email             string `json:"email"`
	Name              string `json:"name"`
	Surname           string `json:"surname"`
	Age               int    `json:"age"`
	Sex               string `json:"sex"`
	Interests         string `json:"interests"`
	City              string `json:"city"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:""`
}

// BeforeCreate ...
func (u *User) BeforeCreate() error {

	if len(u.Password) != 0 {
		encryptedPassword, err := encryptString(u.Password)
		if err != nil {
			return err
		}
		u.EncryptedPassword = encryptedPassword
	}

	return nil
}

// ComparePassword checks if password is correct
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

// Sanitize ...
func (u *User) Sanitize() {
	u.Password = ""
}

func encryptString(s string) (string, error) {

	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
