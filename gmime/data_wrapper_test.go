package gmime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDataWrapper(t *testing.T) {
	dw := NewDataWrapper()
	assert.Equal(t, dw.Encoding(), "")
}

func TestNewDataWrapperWithStream(t *testing.T) {

}

func TestDataWrapperStream(t *testing.T) {

}
