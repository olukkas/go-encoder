package utils_test

import (
	"github.com/olukkas/go-encoder/framework/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsJson(t *testing.T) {
	data := `
	{
		"id": "b5b53262-6a03-447c-9070-7fb3f0f7596a",
		"file_path": "file2.mp4",
		"status": "pending"
	} `

	isJson := utils.IsJson(data)
	require.True(t, isJson)

	isJson = utils.IsJson("not a json")
	require.False(t, isJson)
}
