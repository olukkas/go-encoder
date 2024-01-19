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
	ID         string    `valid:"uuid"`
	ResourceID string    `valid:"notnull"`
	FilePath   string    `valid:"notnull"`
	CreatedAt  time.Time `valid:"-"`
}

func NewVideo(resourceID, filePath string) (*Video, error) {
	v := &Video{
		ID:         uuid.NewV4().String(),
		ResourceID: resourceID,
		FilePath:   filePath,
		CreatedAt:  time.Now(),
	}

	if err := v.Validate(); err != nil {
		return nil, err
	}

	return v, nil
}

func (v *Video) Validate() error {
	_, err := govalidator.ValidateStruct(v)
	if err != nil {
		return err
	}

	return nil
}
