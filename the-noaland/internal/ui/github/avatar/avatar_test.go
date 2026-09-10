package avatar

import "testing"

func TestPlaceholderAndMissingImage(t *testing.T) {
	if got := Placeholder(3, 2); got != "   \n   " {
		t.Fatalf("placeholder = %q", got)
	}
	if got := Render(nil, 14, 7); got != "" {
		t.Fatalf("missing image output = %q", got)
	}
}
