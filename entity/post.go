package entity

// Post struct represents a post in the config.

type Post struct {
	ID          uint64 `gorm:"primary_key:auto_increment" json:"id" form:"id"`
	Title       string `gorm:"index,type:varchar(255)" json:"title" form:"title"`
	Description string `gorm:"type:text" json:"description" form:"description"`
	UserID      uint64 `gorm:"not null" json:"-"`
	User        User   `gorm:"foreignkey:UserID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"user" form:"user"`
	TopicID     uint64 `gorm:"not null" json:"-" form:"topic_id"`
	//Topic       Topic  `gorm:"many2many:post_topic;association_join_table_foreignkey:post_id;foreignkey:id;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"topic"`
	Topic      Topic      `gorm:"foreignkey:TopicID;constraint:onUpdate:CASCADE,onDelete:SET NULL" json:"topic" form:"topic"`
	Comments   *[]Comment `gorm:"many2many:comments;constraint:onUpdate:CASCADE,onDelete:SET NULL" json:"comments,omitempty"`
	Likes      *[]Like    `gorm:"many2many:likes;constraint:onUpdate:CASCADE,onDelete:SET NULL" json:"likes,omitempty"`
	LikesCount int        `json:"likes_count"`
	ImagePath  string     `gorm:"type:varchar(255)" json:"image_path" form:"image_path"`
	//Photo      *multipart.FileHeader `json:"photo,omitempty"`
	//PhotoUrl    string `json:"photo_url"`
	//Photo  *multipart.FileHeader `json:"photo" form:"photo" binding:"required"`
}
