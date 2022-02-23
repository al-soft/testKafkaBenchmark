package kafkag

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	kafka "github.com/segmentio/kafka-go"
)

func init() {
	os.Chdir("../../configs")
	testConfigFile := "config.local.yml"
	var cfgFile *string = &testConfigFile

	New(cfgFile)
}

func TestNew(t *testing.T) {
	os.Chdir("../../configs")
	testConfigFile := "config.local.yml"
	var cfgFile *string = &testConfigFile

	krd, kwr := New(cfgFile)

	actual := krd.Config().Topic
	expected := "bcm-seeker-analyzer"
	assert.Equal(t, expected, actual)

	actual = kwr.Topic
	expected = "bcm-analyzer-recorder"
	assert.Equal(t, expected, actual)

}

func TestReadKafkaTopic(t *testing.T) {
	err := ReadKafkaTopic()
	if err != nil {
		fmt.Println(err)
	}
	actual := err

	time.Sleep(10 * time.Second)

	var expected error
	assert.Equal(t, expected, actual)
}

func TestReadMessage(t *testing.T) {
	channel := make(chan kafka.Message)
	actual := ReadMessage(channel)

	msg := <-channel
	fmt.Printf("%s\n", msg.Value)

	time.Sleep(10 * time.Second)

	var expected error
	assert.Equal(t, expected, actual)
}

func TestKafkaPutMessage(t *testing.T) {
	tests := []struct {
		Message  string
		ID       string
		Expected error
	}{
		{"test message check", "qweqvtwert345345wertgwert", nil},
	}

	for _, tt := range tests {
		actual := putMessage2Topic([]byte(tt.Message), tt.ID)
		assert.Equal(t, tt.Expected, actual)
	}
}
