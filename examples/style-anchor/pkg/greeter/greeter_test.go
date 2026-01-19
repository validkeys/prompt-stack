package greeter

import "testing"

func TestHello(t *testing.T) {
	want := "hello, Alice"
	if got := Hello("Alice"); got != want {
		t.Fatalf("Hello(Alice) = %q, want %q", got, want)
	}
}
