package helper_test

import (
	"testing"

	"github.com/willeponken/picoshop/controller"
)

func TestIsFilled(t *testing.T) {
	type isFilledTest struct {
		err    bool
		fields []string
	}

	tests := []isFilledTest{
		isFilledTest{
			err:    true,
			fields: []string{"", "a", "b", "c"},
		},
		isFilledTest{
			err:    false,
			fields: []string{"a", "b", "c"},
		},
	}

	for _, test := range tests {
		result := controller.IsFilled(test.fields...) != nil
		if result != test.err {
			t.Errorf("result != %v, should return %v", test.err, !test.err)
		}
	}
}
