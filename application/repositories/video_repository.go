package repositories

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/olukkas/go-encoder/domain"
	uuid "github.com/satori/go.uuid"
)

type VideoRepository interface {
	Insert(video *domain.Video) (*domain.Video, error)
	Find(id string) (*domain.Video, error)
}

type VideoRepositoryDB struct {
	DB *gorm.DB
}

func NewVideoRepositoryDB(DB *gorm.DB) *VideoRepositoryDB {
	return &VideoRepositoryDB{DB: DB}
}

func (v *VideoRepositoryDB) Insert(video *domain.Video) (*domain.Video, error) {
	if video.ID == "" {
		video.ID = uuid.NewV4().String()
	}

	err := v.DB.Create(video).Error
	if err != nil {
		return nil, err
	}

	return video, nil
}

func (v *VideoRepositoryDB) Find(id string) (*domain.Video, error) {
	var video domain.Video
	v.DB.First(&video, "id = ?", id)

	if video.ID == "" {
		return nil, fmt.Errorf("video with id %s does not exists \n", id)
	}

	return &video, nil
}
