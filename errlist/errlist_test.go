package errlist

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ErrList_Add(t *testing.T) {
	t.Run("ErrList not nil", func(t *testing.T) {
		var errs1, errs2 ErrList
		errs2.Add(errors.New("b"))
		errs2.Add(errors.New("c"))

		errs1.Add(errors.New("a"))
		errs1.Add(&errs2)

		assert.Equal(t, `['a' 'b' 'c' ]`, errs1.Error())
	})

	t.Run("ErrList nil", func(t *testing.T) {
		var errs1, errs2 ErrList

		errs1.Add(errors.New("a"))
		errs1.Add(&errs2)

		assert.Equal(t, `['a' ]`, errs1.Error())
	})

	t.Run("many errors", func(t *testing.T) {
		var errs ErrList

		for i := 0; i < 100; i++ {
			errs.Add(fmt.Errorf("%d", i))
		}

		assert.Len(t, errs, 100)
	})

	t.Run("nil", func(t *testing.T) {
		var errs ErrList
		errs.Add(nil)
		assert.Len(t, errs, 0)
		assert.Nil(t, errs)
		assert.Equal(t, `[]`, errs.Error())
	})

	t.Run("normal", func(t *testing.T) {
		var errs ErrList

		errs.Add(errors.New("a"))
		errs.Add(errors.New("b"))

		assert.Len(t, errs, 2)
		assert.NotNil(t, errs)
		assert.Equal(t, `['a' 'b' ]`, errs.Error())
	})
}
