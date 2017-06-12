package errlist

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ErrList_Add(t *testing.T) {
	t.Parallel()

	t.Run("ErrList not nil", func(t *testing.T) {
		var errs1, errs2 ErrList
		errs2.AddFieldErr("2", errors.New("b"))
		errs2.AddFieldErr("3", errors.New("c"))

		errs1.AddFieldErr("1", errors.New("a"))
		errs1.Add(errs2)

		assert.Equal(t, `[1: a, 2: b, 3: c]`, errs1.Error())
	})

	t.Run("ErrList nil", func(t *testing.T) {
		var errs1 ErrList

		errs1.AddFieldErr("1", errors.New("a"))
		errs1.Add(nil)

		assert.Equal(t, `[1: a]`, errs1.Error())
	})

	t.Run("unkown field", func(t *testing.T) {
		var errs1 ErrList

		errs1.AddFieldErr("1", errors.New("a"))
		errs1.Add(errors.New("b"))

		assert.Equal(t, `[1: a, unknown: b]`, errs1.Error())
	})

	t.Run("many errors", func(t *testing.T) {
		var errs ErrList

		for i := 0; i < 100; i++ {
			errs.AddFieldErrf("a", fmt.Sprintf("%d", i))
		}

		assert.Len(t, errs, 100)
	})

	t.Run("nil", func(t *testing.T) {
		var errs ErrList
		errs.AddFieldErr("a", nil)
		assert.Len(t, errs, 0)
		assert.Nil(t, errs)
		assert.Equal(t, `[]`, errs.Error())
	})

	t.Run("normal", func(t *testing.T) {
		var errs ErrList

		errs.AddFieldErr("1", errors.New("a"))
		errs.AddFieldErr("2", errors.New("b"))

		assert.Len(t, errs, 2)
		assert.NotNil(t, errs)
		assert.Equal(t, `[1: a, 2: b]`, errs.Error())
	})
}

func Test_ErrList_Marshal(t *testing.T) {
	var errs ErrList

	errs.AddFieldErrf("x", "bla")
	errs.AddFieldErrf("y", "poop")

	res, err := json.Marshal(errs)
	require.NoError(t, err)
	assert.Equal(t, `[{"field":"x","error":"bla"},{"field":"y","error":"poop"}]`, string(res))
}
