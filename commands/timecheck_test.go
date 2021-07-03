package commands_test

import (
	"testing"

	"github.com/zetamatta/nyagos/commands"
)

type testListT struct {
	r bool
	d []int
}

func TestStampIsValid(t *testing.T) {
	testlist := []testListT{
		{true, []int{2016, 5, 13, 17, 20, 0}},
		{true, []int{2016, 2, 29, 17, 20, 0}},
		{false, []int{2015, 2, 29, 17, 20, 0}},
		{false, []int{2016, 14, 20, 17, 20, 0}},
		{false, []int{2016, 12, 32, 17, 20, 0}},
		{false, []int{2016, 12, 31, 24, 20, 0}},
		{false, []int{2016, 12, 31, 23, 70, 0}},
	}
	for _, p := range testlist {
		d := p.d
		if p.r != commands.StampIsValid(d[0], d[1], d[2], d[3], d[4], d[5]) {
			t.Fatalf("[NG] %d/%d/%d %d:%d:%d\n", d[0], d[1], d[2], d[3], d[4], d[5])
			t.Fail()
		}
	}
}
