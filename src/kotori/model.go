package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"time"
)

const (
	CommentBonus = 50
)

type Admin struct {
	Username string
	Password string
}

type Index struct {
	ID    uint `gorm:"AUTO_INCREMENT"`
	Class string `gorm:"not null"`
	Title string
	Attr  string
}

type User struct {
	ID      uint `gorm:"AUTO_INCREMENT"`
	Name    string
	Email   string `gorm:"not null;unique"`
	Website string
	Rank    int64
	Honor   string
}

type Comment struct {
	ID            uint `gorm:"AUTO_INCREMENT"`
	CommentZoneID uint
	FatherID      uint
	ReplyUserID   uint
	ReplyUser     User
	UserID        uint
	User          User
	Content       string
	Type          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Post struct {
	ID        uint `gorm:"AUTO_INCREMENT"`
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func FindComments(db *gorm.DB, commentZoneID uint, fatherID uint, offsetID uint) (comments []Comment, err error) {
	var order string
	var offset string
	if fatherID != 0 {
		order = "id asc"
		offset = "id > ?"
	} else {
		order = "id desc"
		offset = "id < ?"
	}
	if offsetID == 0 {
		err = db.Where("comment_zone_id = ?", commentZoneID).Where("father_id = ?", fatherID).
			Preload("User").Preload("ReplyUser").Order(order).Limit(10).Find(&comments).Error
	} else {
		err = db.Where("comment_zone_id = ?", commentZoneID).Where("father_id = ?", fatherID).
			Where(offset, offsetID).
			Preload("User").Preload("ReplyUser").Order(order).Limit(10).Find(&comments).Error
	}
	if err != nil {
		err = errors.Wrap(err, "ListComments")
		return
	}
	return
}

func StoreComment(db *gorm.DB, comment Comment) (comment_new Comment, err error) {
	var users []User
	var user_cnt uint
	err = db.Model(&User{}).Where("email = ?", comment.User.Email).Find(&users).Count(&user_cnt).Error
	if err != nil {
		err = errors.Wrap(err, "SaveComment")
		return
	}
	if user_cnt != 0 {
		comment.UserID = users[0].ID
		users[0].Name = comment.User.Name
		users[0].Website = comment.User.Website
		users[0].Rank += CommentBonus
		db.Model(&User{}).Updates(&users[0])
	} else {
		db.Create(&comment.User)
		comment.UserID = comment.User.ID
	}
	err = db.Set("gorm:save_associations", false).Create(&comment).Error
	if err != nil {
		err = errors.Wrap(err, "SaveComment")
		return
	}
	err = db.Where("id = ?", &comment.ID).
		Preload("User").Preload("ReplyUser").First(&comment_new).Error
	return
}

func RemoveComment(db *gorm.DB, id uint) (err error) {
	var comment Comment
	err = db.Model(&Comment{}).Where("id = ?", id).Preload("User").First(&comment).Error
	if err != nil {
		err = errors.Wrap(err, "RemoveComment")
		return
	}
	comment.User.Rank -= CommentBonus
	db.Model(&User{}).Updates(&comment.User)
	db.Delete(&comment)
	return
}

func FindIndexes(db *gorm.DB, class string, order string, offsetID uint) (indexes []Index, err error) {
	var offset string
	if order == "asc" {
		offset = "id > ?"
	} else {
		offset = "id < ?"
	}
	if offsetID == 0 {
		err = db.Where("class = ?", class).Order("id " + order).Limit(20).Find(&indexes).Error
	} else {
		err = db.Where("class = ?", class).Order("id " + order).Limit(20).
			Where(offset, offsetID).Find(&indexes).Error
	}
	if err != nil {
		err = errors.Wrap(err, "ListComments")
		return
	}
	return
}

func StoreIndex(db *gorm.DB, index Index) (index_new Index, err error) {
	err = db.Create(&index).Error
	if err != nil {
		err = errors.Wrap(err, "SaveComment")
		return
	}
	index_new = index
	return
}

func UpdateIndex(db *gorm.DB, index Index) (index_new Index, err error) {
	err = db.Model(&index).Updates(index).Error
	if err != nil {
		err = errors.Wrap(err, "UpdateIndex")
		return
	}
	index_new = index
	return
}

func RemoveIndex(db *gorm.DB, id uint) (err error) {
	err = db.Delete(Index{}, "id = ?", id).Error
	if err != nil {
		err = errors.Wrap(err, "RemoveIndex")
		return
	}
	return
}

func FindPost(db *gorm.DB, id uint) (post Post, err error) {
	err = db.Where("id = ?", id).Find(&post).Error
	if err != nil {
		err = errors.Wrap(err, "FindPost")
		return
	}
	return
}

func StorePost(db *gorm.DB, post Post) (post_new Post, err error) {
	err = db.Create(&post).Error
	if err != nil {
		err = errors.Wrap(err, "StorePost")
		return
	}
	post_new = post
	return
}

func UpdatePost(db *gorm.DB, post Post) (post_new Post, err error) {
	err = db.Model(&post).Updates(post).Error
	if err != nil {
		err = errors.Wrap(err, "UpdatePost")
		return
	}
	post_new = post
	return
}

func RemovePost(db *gorm.DB, id uint) (err error) {
	err = db.Delete(Post{}, "id = ?", id).Error
	if err != nil {
		err = errors.Wrap(err, "RemovePost")
		return
	}
	return
}