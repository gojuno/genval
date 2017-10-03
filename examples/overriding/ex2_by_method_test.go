package overriding

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Request2_Validate(t *testing.T) {
	t.Parallel()

	validRequest := Request2{
		Age:  Age2{Value: 10},
		Some: 10,
	}

	t.Run("valid", func(t *testing.T) {
		assert.NoError(t, validRequest.Validate())
	})

	t.Run("invalid", func(t *testing.T) {
		t.Run("too young, using overridden rule", func(t *testing.T) {
			r := validRequest
			r.Age.Value = 5

			err := r.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[Age: field Age is less than 10]`, err.Error())
		})

		t.Run("too young, using generated type validator", func(t *testing.T) {
			r := validRequest
			r.Age.Value = 5

			err := r.Age.Validate()
			require.NoError(t, err)
		})

		t.Run("too young, using func validator", func(t *testing.T) {
			r := validRequest
			r.Age.Value = 5

			err := r.Age.ValidateMin10()
			require.NotNil(t, err)
			assert.Equal(t, `field Age is less than 10`, err.Error())
		})

		t.Run("too old, using func validator", func(t *testing.T) {
			r := validRequest
			r.Age.Value = 65

			err := r.Age.ValidateMin10()
			require.NotNil(t, err)
			assert.Equal(t, `field Age is more than 64`, err.Error())
		})

		t.Run("too old, using generated validator", func(t *testing.T) {
			r := validRequest
			r.Age.Value = 65

			err := r.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[Age: field Age is more than 64]`, err.Error())
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
