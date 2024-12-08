package models

import "gorm.io/gorm"

type Doc struct {
	gorm.Model
	ID       string `gorm:"uniqueIndex;not null" json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	FilePath string `json:"file_path,omitempty"`
}
