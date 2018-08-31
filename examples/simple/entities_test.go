package simple

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_User_Validate(t *testing.T) {
	t.Parallel()

	title := None

	validUser := User{
		Name:   "Vasa",
		Age:    20,
		Dog:    Dog{Name: "Taksa"},
		Emails: map[int]string{1: "vasa@gojuno.com"},
		Title:  &title,
	}

	t.Run("valid", func(t *testing.T) {
		assert.NoError(t, validUser.Validate())
	})

	t.Run("invalid", func(t *testing.T) {
		t.Run("short name", func(t *testing.T) {
			user := validUser
			user.Name = ""

			err := user.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[Name: shorter than 3 chars]`, err.Error())
		})

		t.Run("too long name", func(t *testing.T) {
			user := validUser
			user.Name = "123456789123456789123456789123456789123456789123456789123456789123456789"

			err := user.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[Name: longer than 64 chars]`, err.Error())
		})

		t.Run("nil title", func(t *testing.T) {
			user := validUser
			user.Title = nil

			err := user.Validate()
			require.NotNil(t, err)
			assert.Equal(t, "[Title: cannot be nil]", err.Error())
		})

		t.Run("bad title", func(t *testing.T) {
			user := validUser
			badTitle := Title("Jedi")
			user.Title = &badTitle

			err := user.Validate()
			require.NotNil(t, err)
			assert.Equal(t, "[Title: invalid value for enum Title: Jedi]", err.Error())
		})

		t.Run("no email", func(t *testing.T) {
			user := validUser
			user.Emails = nil

			err := user.Validate()
			require.NotNil(t, err)
			assert.Equal(t, "[Emails: less items than 1]", err.Error())
		})

		t.Run("bad email key and value", func(t *testing.T) {
			user := validUser
			user.Emails = map[int]string{1111: "abc"}

			err := user.Validate()
			require.NotNil(t, err)
			assert.Equal(t, "[Emails.key[1111]: more than 3, Emails.1111: shorter than 5 chars]", err.Error())
		})

		t.Run("too young", func(t *testing.T) {
			user := validUser
			user.Age = 15

			err := user.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[Age: less than 18]`, err.Error())
		})

		t.Run("too old", func(t *testing.T) {
			user := validUser
			user.Age = 111

			err := user.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[Age: more than 95]`, err.Error())
		})

		t.Run("bad dog", func(t *testing.T) {
			user := validUser
			user.Dog = Dog{}

			err := user.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[Dog.Name: shorter than 1 chars]`, err.Error())
		})
	})
}

func Test_Dog_Validate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		assert.NoError(t, Dog{"Jora"}.Validate())
	})

	t.Run("invalid", func(t *testing.T) {
		t.Run("too short", func(t *testing.T) {
			assert.NotNil(t, Dog{}.Validate())
		})
		t.Run("too big", func(t *testing.T) {
			assert.NotNil(t, Dog{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}.Validate())
		})
	})
}

func Test_Title_Validate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		assert.NoError(t, Title(None).Validate())
		assert.NoError(t, Title(Doctor).Validate())
		assert.NoError(t, Title(Sir).Validate())
		assert.NoError(t, Title(Father).Validate())
	})

	t.Run("invalid", func(t *testing.T) {
		assert.Error(t, Title("unknown").Validate())
	})
}
