package domain_test

import (
	"github.com/olukkas/go-encoder/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidationWithEmptyVideo(t *testing.T) {
	_, err := domain.NewVideo("", "")
	require.Error(t, err)
}

func TestVideoIDIsNotUUID(t *testing.T) {
	v, _ := domain.NewVideo("resourceID", "FilePath")
	v.ID = "not uuid"

	err := v.Validate()
	require.Error(t, err)
}

func TestVideWithFieldCorrect(t *testing.T) {
	_, err := domain.NewVideo("resourceID", "FilePath")
	require.Nil(t, err)
}
