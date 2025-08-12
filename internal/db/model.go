package db

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"uniqueIndex;size:64;not null"`
	Password  string `gorm:"size:255;not null"` // 存储 bcrypt 哈希
	CreatedAt time.Time
	UpdatedAt time.Time
}
