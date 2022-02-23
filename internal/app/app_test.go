package app

import (
	"os"
	"testing"
)

func TestRun(t *testing.T) {
	os.Chdir("../../configs")
	Run()
}
