package objects

import (
	_ "encoding/json"
)

type User struct {
	Id        int    `json:"id" gorm:"primary_key; index"`
	FirstName string `json:"first_name" gorm:"not null"`
	LastName  string `json:"last_name"`
	Login     string `json:"login" gorm:"not null; unique"`
	Password  string `json:"password" gorm:"not null"`
	Email     string `json:"email" gorm:"not null; unique"`
	UserType  string `json:"user_type" gorm:"not null" sql:"DEFAULT: 'user'"`
}

func (User) TableName() string {
	return "user"
}

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponse struct {
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
	Role        string `json:"role"`
}

type UserCreateRequest struct {
	Profile struct {
		Firstname   string `json:"firstName,omitempty"`
		Lastname    string `json:"lastName,omitempty"`
		Email       string `json:"email,omitempty"`
		Login       string `json:"login,omitempty"`
		Mobilephone string `json:"mobilePhone,omitempty"`
		UserType    string `json:"userType,omitempty"`
	} `json:"profile,omitempty"`
	Credentials struct {
		Password struct {
			Value string `json:"value"`
		} `json:"password"`
	} `json:"credentials"`
	GroupIds []string `json:"groupIds,omitempty"`
}
