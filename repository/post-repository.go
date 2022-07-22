package repository

import (
	"gitlab.zalopay.vn/top/intern/vybnt/gallery-backend/gallery/entity"
	"gitlab.zalopay.vn/top/intern/vybnt/gallery-backend/gallery/pagination"
	"gorm.io/gorm"
	"log"
)

type PostRepository interface {
	InsertPost(post entity.Post) entity.Post
	UpdatePost(post entity.Post) entity.Post
	DeletePost(post entity.Post)
	AllPost(pagination pagination.Pagination) []entity.Post
	FindPostByID(postID uint64) entity.Post
	FindPostByTopicID(topicID uint64) []entity.Post
	GetTrendingPosts(pagination pagination.Pagination) []entity.Post
}

type postConnection struct {
	connection *gorm.DB
}

//NewPostRepository creates an instance BookRepository
func NewPostRepository(databaseConnection *gorm.DB) PostRepository {
	return &postConnection{
		connection: databaseConnection,
	}
}

func (db *postConnection) InsertPost(post entity.Post) entity.Post {
	err := db.connection.Save(&post)
	if err != nil {
		log.Println(err)
	}
	db.connection.Preload("User").Find(&post)
	db.connection.Preload("Topic").Find(&post)
	db.connection.Preload("Comments").Find(&post)
	db.connection.Preload("Likes").Find(&post)
	return post
}

func (db *postConnection) UpdatePost(post entity.Post) entity.Post {
	err := db.connection.Save(&post)
	if err != nil {
		log.Println(err)
	}
	db.connection.Preload("User").Find(&post)
	db.connection.Preload("Topic").Find(&post)
	db.connection.Preload("Comments").Find(&post)
	return post
}

func (db *postConnection) DeletePost(post entity.Post) {
	err := db.connection.Delete(&post)

	if err != nil {
		log.Println(err)
	}
}

func (db *postConnection) FindPostByID(postID uint64) entity.Post {
	var post entity.Post
	db.connection.Preload("User").Find(&post, postID)
	return post
}

func (db *postConnection) AllPost(pagination pagination.Pagination) []entity.Post {
	var posts []entity.Post
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := db.connection.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	result := queryBuilder.Find(&posts)
	if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	return posts
}

func (db *postConnection) FindPostByTopicID(topicID uint64) []entity.Post {
	var posts []entity.Post
	db.connection.Preload("User").Find(&posts, "topic_id = ?", topicID)
	return posts
}

func (db *postConnection) GetTrendingPosts(pagination pagination.Pagination) []entity.Post {
	var posts []entity.Post
	offset := (pagination.Page - 1) * pagination.Limit
	queryBuilder := db.connection.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	result := queryBuilder.Preload("User").Find(&posts, "likes_count > ?", 0)
	if result.Error != nil {
		log.Println(result.Error)
		return nil
	}
	//db.connection.Preload("User").Find(&posts, "likes_count > ?", 0)
	return posts
}
