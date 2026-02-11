package model

import (
	"time"
)

type Announcement struct {
	AnnouncementID string    `json:"id" gorm:"column:announcement_id;type:STRING"`
	StockCode      string    `json:"stock_code" gorm:"column:stock_code;type:SYMBOL"`
	Title          string    `json:"title" gorm:"column:title;type:STRING"`
	ContentURL     string    `json:"content_url" gorm:"column:content_url;type:STRING"`
	PublishDate    time.Time `json:"publish_date" gorm:"column:publish_date;type:TIMESTAMP"`
	Exchange       string    `json:"exchange" gorm:"column:exchange;type:STRING"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at;type:TIMESTAMP"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"column:updated_at;type:TIMESTAMP"`
}

func (Announcement) TableName() string {
	return "announcements"
}
