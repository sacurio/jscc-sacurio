package model

type Message struct {
	ID        uint   `gorm:"type:int;primaryKey;autoIncrement;not null" json:"id"`
	Content   string `gorm:"type:nvarchar(1000);not null" json:"content"`
	Timestamp int64  `gorm:"type:int;not null" json:"timestamp"`
	UserID    uint   `gorm:"type:int;not null"`
	User      User   `gorm:"foreignKey:UserID"`
}

type History struct {
	User     string        `json:"username"`
	Messages []interface{} `json:"messages"`
}
