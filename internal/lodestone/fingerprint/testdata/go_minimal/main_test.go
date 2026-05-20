package fixture

import "testing"

func TestNewRoot(t *testing.T) {
	if newRoot() == nil {
		t.Fatal("nil")
	}
}
