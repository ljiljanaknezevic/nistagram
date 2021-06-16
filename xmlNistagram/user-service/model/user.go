package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name        string `json:"name"`
	Email       string `gorm:"unique" json:"email"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	PhoneNumber string `json:"phone"`
	Gender      string `json:"gender"`
	Birhtday    string `json:"birthday"`
	Username    string `gorm:"unique" json:"username"`
	Website     string `json:"website"`
	Biography   string `json:"biography"`
	IsPrivate   bool   `json:"isPrivate"`
	IsVerified  bool   `json:"isVerified"`
	Followers   []Follower `gorm:"many2many:user_followers; json:"followers"`
	WaitingFollowers   []WaitingFollower `gorm:"many2many:user_waitingFollowers; json:"waitingFollowers"`
	Following   []Following `gorm:"many2many:user_following; json:"following"`
}

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Follower struct {
	gorm.Model
	Username   string `json:"username"`
}
type WaitingFollower struct {
	gorm.Model
	Username   string `json:"username"`
}
type Following struct {
	gorm.Model
	Username   string `json:"username"`
}
