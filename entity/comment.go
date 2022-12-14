package entity

type Comment struct {
	ID      uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Content string `gorm:"type:text" json:"content"`
	PostID  uint64 `gorm:"not null" json:"post_id"`
	UserID  uint64 `gorm:"not null" json:"-"`
	User    User   `gorm:"foreignkey:UserID;constraint:onUpdate:CASCADE,onDelete:SET NULL" json:"user"`
	Post    Post   `gorm:"foreignkey:PostID;constraint:onUpdate:CASCADE,onDelete:SET NULL" json:"post"`
}
