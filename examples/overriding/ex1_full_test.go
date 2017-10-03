package overriding

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Request1_Validate(t *testing.T) {
	t.Parallel()

	validRequest := Request1{
		Age:  Age1{Value: 10},
		Some: 10,
	}

	t.Run("valid", func(t *testing.T) {
		assert.NoError(t, validRequest.Validate())
	})

	t.Run("invalid", func(t *testing.T) {
		t.Run("too young, use overridden rule", func(t *testing.T) {
			r := validRequest
			r.Age.Value = 5

			err := r.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `field Age is less than 10`, err.Error())
		})

		t.Run("too old, using generated validator", func(t *testing.T) {
			r := validRequest
			r.Age.Value = 65

			err := r.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `field Age is more than 64`, err.Error())
		})

		t.Run("too old, using generated type validator", func(t *testing.T) {
			r := validRequest
			r.Age.Value = 65

			err := r.Age.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[Value: more than 64]`, err.Error())
		})

		t.Run("check min.", func(t *testing.T) {
			r := validRequest
			r.Some = 1

			err := r.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `field Some is less than 3`, err.Error())
		})

		t.Run("check max.", func(t *testing.T) {
			r := validRequest
			r.Some = 65

			err := r.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `field Some is more than 64`, err.Error())
		})

	})
}
