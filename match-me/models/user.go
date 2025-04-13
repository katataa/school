package models

import "time"

type User struct {
	ID              uint      `gorm:"primaryKey"`
	Name            string    `json:"name"`
	Email           string    `gorm:"unique;not null" json:"email"`
	Password        string    `json:"-"`
	Info            string    `json:"info"`
	Interests       string    `json:"interests"`
	Location        string    `json:"location"`
	Age             int       `json:"age"`
	ProfilePicture  string    `json:"profile_picture"`
	Latitude        float64   `json:"latitude"`
	Longitude       float64   `json:"longitude"`
	PreferredRadius float64   `json:"preferred_radius"`
	Gender          string    `json:"gender"`
	LookingFor      string    `json:"looking_for"`
	Bio             *Bio      `gorm:"foreignKey:UserID"`
	Profile         *Profile  `gorm:"foreignKey:UserID"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Bio struct {
	ID              uint    `gorm:"primaryKey"`
	UserID          uint    `gorm:"not null"`
	Interests       string  `json:"interests"`
	Age             int     `json:"age"`
	Gender          string  `json:"gender"`
	Location        string  `json:"location"`
	PreferredRadius float64 `json:"preferred_radius"`
	Info            string
}

type Profile struct {
	ID       uint   `gorm:"primaryKey"`
	UserID   uint   `gorm:"not null"`
	Headline string `json:"headline"`
}

type Connection struct {
	ID         uint   `gorm:"primaryKey"`
	SenderID   uint   `gorm:"not null"`
	ReceiverID uint   `gorm:"not null"`
	Status     string `gorm:"not null"`
	Sender     User   `gorm:"foreignKey:SenderID"`
	Receiver   User   `gorm:"foreignKey:ReceiverID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type DeclinedUser struct {
	ID             uint `gorm:"primaryKey"`
	UserID         uint `gorm:"not null"`
	DeclinedUserID uint `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Chat struct {
	ID          uint `gorm:"primaryKey"`
	User1ID     uint `gorm:"not null"`
	User2ID     uint `gorm:"not null"`
	User1       User `gorm:"foreignKey:User1ID"`
	User2       User `gorm:"foreignKey:User2ID"`
	UnreadCount uint `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Message struct {
	ID         uint      `gorm:"primaryKey"`
	ChatID     uint      `json:"chat_id"`
	SenderID   uint      `json:"sender_id"`
	ReceiverID uint      `json:"receiver_id"`
	Content    string    `json:"content"`
	Timestamp  time.Time `json:"timestamp"`
	IsRead     bool      `json:"is_read" gorm:"default:false"`
}
