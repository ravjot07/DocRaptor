package models

type Doc struct {
	ID       string `gorm:"primaryKey;unique;not null" json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	FilePath string `json:"file_path,omitempty"`
}
