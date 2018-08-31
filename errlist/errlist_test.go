package errlist

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_List_Add(t *testing.T) {
	t.Parallel()

	t.Run("err is errlist", func(t *testing.T) {
		var errs1, errs2 List
		errs2.AddField("2", errors.New("b"))
		errs2.AddField("3", errors.New("c"))

		errs1.AddField("1", errors.New("a"))
		errs1.AddField("errs2", errs2)

		assert.Equal(t, `[1: a, errs2.2: b, errs2.3: c]`, errs1.Error())
	})

	t.Run("List not nil", func(t *testing.T) {
		var errs1, errs2 List
		errs2.AddField("2", errors.New("b"))
		errs2.AddField("3", errors.New("c"))

		errs1.AddField("1", errors.New("a"))
		errs1.Add(errs2)

		assert.Equal(t, `[1: a, 2: b, 3: c]`, errs1.Error())
	})

	t.Run("List nil", func(t *testing.T) {
		var errs1 List

		errs1.AddField("1", errors.New("a"))
		errs1.Add(nil)

		assert.Equal(t, `[1: a]`, errs1.Error())
	})

	t.Run("unkown field", func(t *testing.T) {
		var errs1 List

		errs1.AddField("1", errors.New("a"))
		errs1.Add(errors.New("b"))

		assert.Equal(t, `[1: a, b]`, errs1.Error())
	})

	t.Run("many errors", func(t *testing.T) {
		var errs List

		for i := 0; i < 100; i++ {
			errs.AddFieldf("a", fmt.Sprintf("%d", i))
		}

		assert.Len(t, errs, 100)
	})

	t.Run("nil", func(t *testing.T) {
		var errs List
		errs.AddField("a", nil)
		assert.Len(t, errs, 0)
		assert.Nil(t, errs)
		assert.Equal(t, `[]`, errs.Error())
	})

	t.Run("normal", func(t *testing.T) {
		var errs List

		errs.AddField("1", errors.New("a"))
		errs.AddField("2", errors.New("b"))

		assert.Len(t, errs, 2)
		assert.NotNil(t, errs)
		assert.Equal(t, `[1: a, 2: b]`, errs.Error())
	})
}

func Test_List_Marshal(t *testing.T) {
	var errs List

	errs.AddFieldf("x", "bla")
	errs.AddFieldf("y", "poop")

	res, err := json.Marshal(errs)
	require.NoError(t, err)
	assert.Equal(t, `[{"field":"x","error":"bla"},{"field":"y","error":"poop"}]`, string(res))
}
