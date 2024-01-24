package domain

import (
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
	"time"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Video struct {
	ID         string    `json:"encoded_vide_folder" valid:"uuid"`
	ResourceID string    `json:"resource_id" valid:"notnull"`
	FilePath   string    `json:"file_path" valid:"notnull"`
	CreatedAt  time.Time `json:"-" valid:"-"`
	Jobs       []*Job    `json:"-" valid:"-"`
}

func NewVideo(resourceID, filePath string) (*Video, error) {
	video := &Video{
		ID:         uuid.NewV4().String(),
		ResourceID: resourceID,
		FilePath:   filePath,
		CreatedAt:  time.Now(),
	}

	if err := video.Validate(); err != nil {
		return nil, err
	}

	return video, nil
}

func (v *Video) Validate() error {
	_, err := govalidator.ValidateStruct(v)
	if err != nil {
		return err
	}

	return nil
}
