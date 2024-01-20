package domain_test

import (
	"github.com/olukkas/go-encoder/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

var defaultVideo, _ = domain.NewVideo("resource", "path")

func TestNewJob(t *testing.T) {
	job, err := domain.NewJob("output_path", domain.JobDownloading, defaultVideo)
	require.Nil(t, err)
	require.NotNil(t, job)
}
