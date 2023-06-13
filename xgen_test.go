package goxgen

import (
	"testing"
)

func TestNewXgen(t *testing.T) {
	xgen := NewXgen()
	if xgen == nil {
		t.Errorf("NewXgen() = nil")
	}
}
