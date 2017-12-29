package complicated

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Request1_Validate(t *testing.T) {
	t.Parallel()

	userLastName := "LastN"
	childrenCount := 10
	validUser := User{
		Name:          "TestUserName",
		LastName:      &userLastName,
		Age:           35,
		ChildrenCount: &childrenCount,
		Float:         1.11,
		Dog: Dog{
			Name: "Dog",
		},
		Urls: []string{
			"someURL",
		},
		Dogs: []*Dog{
			{
				Name: "Dog",
			},
		},
		SomeArray: make([]interface{}, 1),
		Dict: map[string]int{
			"test":   1,
			"test_2": 2,
		},
		MapOfMap: map[string]map[int]string{
			"key": {
				1: "value",
			},
		},
		MapOfSlice: map[string][]string{
			"key": {
				"value",
			},
		},
		SliceOfMap: []map[string]int{
			{
				"key": 1,
			},
		},
		SliceOfSliceOfSlice: [][][]string{
			{
				{
					"value",
				},
			},
		},
	}

	t.Run("valid", func(t *testing.T) {
		assert.NoError(t, validUser.Validate())
	})

	t.Run("invalid", func(t *testing.T) {

		t.Run("name: too short", func(t *testing.T) {
			r := validUser
			r.Name = "a"

			err := r.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[Name: shorter than 3 chars]`, err.Error())
		})

		t.Run("last_name: too short", func(t *testing.T) {
			r := validUser
			lastName := ""
			r.LastName = &lastName

			err := r.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[LastName: shorter than 1 chars]`, err.Error())
		})

		t.Run("last_name: empty", func(t *testing.T) {
			r := validUser
			r.LastName = nil

			err := r.Validate()
			require.NoError(t, err)
		})

		t.Run("age: too young", func(t *testing.T) {
			r := validUser
			r.Age = 1

			err := r.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[Age: less than 18]`, err.Error())
		})

		t.Run("children_count: not_null rule", func(t *testing.T) {
			r := validUser
			r.ChildrenCount = nil

			err := r.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[ChildrenCount: cannot be nil]`, err.Error())
		})

		t.Run("children_count: too much", func(t *testing.T) {
			r := validUser
			childrenCount := 100
			r.ChildrenCount = &childrenCount

			err := r.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[ChildrenCount: more than 15]`, err.Error())
		})

		t.Run("dog_pointer: nil is ok", func(t *testing.T) {
			r := validUser
			r.DogPointer = nil

			err := r.Validate()
			require.NoError(t, err)
		})

		t.Run("dog_pointer: custom type validator", func(t *testing.T) {
			r := validUser
			r.DogPointer = &Dog{}

			err := r.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[DogPointer.Name: shorter than 1 chars]`, err.Error())
		})

		t.Run("alias: alias type validator", func(t *testing.T) {
			r := validUser
			r.Alias = DogsMapAlias(map[string]Dog{})

			err := r.Validate()
			require.NoError(t, err)
		})

		t.Run("alias_on_alias: alias of alias type validator", func(t *testing.T) {
			r := validUser
			r.AliasOnAlias = AliasOnDogsMapAlias(map[string]Dog{})

			err := r.Validate()
			require.NoError(t, err)
		})

		t.Run("MapOfMap: invalid value", func(t *testing.T) {
			r := validUser
			r.MapOfMap = map[string]map[int]string{
				"key": {
					1: "v",
				},
			}

			err := r.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[MapOfMap.key.1: shorter than 3 chars]`, err.Error())
		})

		t.Run("SliceOfMap: invalid key", func(t *testing.T) {
			r := validUser
			r.SliceOfMap = []map[string]int{
				{
					"k": 1,
				},
			}

			err := r.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[SliceOfMap.0.key[k]: shorter than 3 chars]`, err.Error())
		})

		t.Run("SliceOfSliceOfSlice: invalid length", func(t *testing.T) {
			r := validUser
			r.SliceOfSliceOfSlice = [][][]string{
				{
					{
						"value",
					},
				},
				{
					{},
				},
			}

			err := r.Validate()
			require.NotNil(t, err)
			assert.Equal(t, `[SliceOfSliceOfSlice.1.0: less items than 1]`, err.Error())
		})
	})
}
