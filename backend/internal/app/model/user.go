package model

// User represents the user info structure.
type User struct {
	ID       uint   `gorm:"type:int;primaryKey"`
	Username string `gorm:"type:varchar(50);not null;size:50;index"`
	Email    string `gorm:"type:varchar(100);not null;size:100"`
	Pwd      string `gorm:"type:varchar(32);not null;size:32"`
}
