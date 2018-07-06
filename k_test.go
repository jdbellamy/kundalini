package kundalini_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "gitlab.com/jdbellamy/kundalini"
)

func TestK_WrapsSingleString(t *testing.T) {

	v := "a"

	actual, err := Coil(v).Release()

	expected := []string{v}

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestK_WrapsSliceOfSingleString(t *testing.T) {

	v := []string{"a"}

	actual, err := Coil(v).Release()

	expected := []string{"a"}

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
