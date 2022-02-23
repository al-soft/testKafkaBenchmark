package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigFile(t *testing.T) {
	os.Chdir("../../configs")

	testCF, err := New()
	if err != nil {
		fmt.Println(err)
	}
	actualFile := *testCF.File
	expectedFile := "config.local.yml"
	assert.Equal(t, expectedFile, actualFile)
	assert.Equal(t, nil, err)
}
