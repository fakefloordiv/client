package keyvalue

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStorage(t *testing.T) {
	testValues := func(t *testing.T, storage *Storage) {
		for _, tc := range []struct {
			Key    string
			Values []string
		}{
			{
				Key:    "Hello",
				Values: []string{"world"},
			},
			{
				Key:    "Some",
				Values: []string{"multiple", "values"},
			},
			{
				Key:    "sOME",
				Values: []string{"multiple", "values"},
			},
		} {
			value, found := storage.Get(tc.Key)
			require.True(t, found)
			require.Equal(t, tc.Values[0], value)

			values := storage.Values(tc.Key)
			require.Equal(t, tc.Values, values)
		}
	}

	t.Run("Value with manual filling", func(t *testing.T) {
		storage := New()
		storage.Add("Hello", "world")
		storage.Add("Some", "multiple")
		storage.Add("Some", "values")
		testValues(t, storage)
	})

	t.Run("Value with map instantiation", func(t *testing.T) {
		storage := NewFromMap(map[string][]string{
			"Hello": {"world"},
			"Some":  {"multiple", "values"},
		})
		testValues(t, storage)
	})

	t.Run("Has", func(t *testing.T) {
		storage := New()
		storage.Add("Hello", "world")
		require.True(t, storage.Has("Hello"))
		require.True(t, storage.Has("hELLO"))
		require.False(t, storage.Has("random"))
	})

	t.Run("Keys", func(t *testing.T) {
		storage := New()
		storage.Add("Hello", "world")
		storage.Add("sOME", "multiple")
		storage.Add("Some", "values")
		storage.Add("hELLO", "nether")
		require.Equal(t, []string{"Hello", "sOME"}, storage.Keys())
	})
}
