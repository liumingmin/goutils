package utils

import "testing"

func TestMathMin(t *testing.T) {
	if Min(10, 9) != 9 {
		t.FailNow()
	}

	if Min(-1, -2) != -2 {
		t.FailNow()
	}

	if Min(3.1, 4.02) != 3.1 {
		t.FailNow()
	}
}

func TestMathMax(t *testing.T) {
	if Max(10, 9) != 10 {
		t.FailNow()
	}

	if Max(-1, -2) != -1 {
		t.FailNow()
	}

	if Max(3.1, 4.02) != 4.02 {
		t.FailNow()
	}
}

func TestMathAbs(t *testing.T) {
	if Abs(-1) != 1 {
		t.FailNow()
	}

	if Abs(1) != 1 {
		t.FailNow()
	}
}
