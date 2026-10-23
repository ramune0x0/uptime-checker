package target

import "testing"

func TestNew(t *testing.T) {
	Target, err := New("example", "https://example.com")
	if err != nil {
		t.Errorf(`New("example, "https://example.com) = %q, %v, want "example", error`, Target, err)
	}
}
