package overriding

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Request5_Validate(t *testing.T) {
	t.Parallel()

	validRequest := Request5{
		Age:  Age5{Value: 5},
		Some: 10,
	}

	t.Run("valid", func(t *testing.T) {
		assert.NoError(t, validRequest.Validate())
	})

	t.Run("invalid", func(t *testing.T) {
		t.Run("definite value, using private validator", func(t *testing.T) {
			r := validRequest
			r.Age.Value = 10
			r.Some = 10

			err := r.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[unknown: fields Age and Some can't be 10 at the same time]`, err.Error())
		})

		t.Run("too young, using generated type validator", func(t *testing.T) {
			r := validRequest
			r.Age.Value = 2

			err := r.Age.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[Value: less than 3]`, err.Error())
		})

		t.Run("too old, using generated type validator", func(t *testing.T) {
			r := validRequest
			r.Age.Value = 65

			err := r.Age.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[Value: more than 64]`, err.Error())
		})

		t.Run("too young, using generated validator", func(t *testing.T) {
			r := validRequest
			r.Age.Value = 2

			err := r.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[Age.Value: less than 3]`, err.Error())
		})

		t.Run("too old, using generated validator", func(t *testing.T) {
			r := validRequest
			r.Age.Value = 65

			err := r.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[Age.Value: more than 64]`, err.Error())
		})

		t.Run("check min.", func(t *testing.T) {
			r := validRequest
			r.Some = 1

			err := r.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[Some: less than 3]`, err.Error())
		})

		t.Run("check max.", func(t *testing.T) {
			r := validRequest
			r.Some = 65

			err := r.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[Some: more than 64]`, err.Error())
		})

	})
}
