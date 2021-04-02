package test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func CmpEqual(t *testing.T, got interface{}, want interface{}) {
	if !cmp.Equal(got, want) {
		t.Fail()
		t.Logf("%v", cmp.Diff(got, want))
	}
}
