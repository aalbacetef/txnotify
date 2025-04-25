package txnotify

import (
	"errors"
	"testing"
)

func TestHex(t *testing.T) {
	t.Run("it correctly parses a 0x prefixed hex string", func(tt *testing.T) {
		const testStr = "0x1692"
		const want = 0x1692

		got, err := strToHex(testStr)
		if err != nil {
			tt.Fatalf("error: %v", err)
		}

		if got != want {
			tt.Fatalf("got: %#0x, want %#0x", got, want)
		}
	})

	t.Run("it fails if the string is not prefixed", func(tt *testing.T) {
		const testStr = "1692"
		wantErr := ErrBadFormat
		_, err := strToHex(testStr)
		if !errors.Is(err, wantErr) {
			tt.Fatalf("got: '%v', want: '%s'", err, wantErr)
		}
	})
}

func TestNormalizeAddress(t *testing.T) {
	v := "0x0000000012"
	want := "0x12"

	normed := normalizeAddress(v)
	if normed == want {
		return
	}

	t.Fatalf("got '%s', want '%s'", normed, want)
}
