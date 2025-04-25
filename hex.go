package txnotify

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrBadFormat = errors.New("invalid string format, expected '0x' prefix")

func strToHex(s string) (int, error) {
	if !strings.HasPrefix(s, "0x") {
		return 0, ErrBadFormat
	}

	v, err := strconv.ParseInt(strings.TrimLeft(s, "0x"), 16, 64)
	if err != nil {
		return 0, fmt.Errorf("could not parse hex value: %w", err)
	}

	return int(v), nil
}

func numToStr(v int) string {
	return fmt.Sprintf("%#0x", v)
}

func normalizeAddress(s string) string {
	pfx := "0x"
	addr := strings.TrimPrefix(s, pfx)

	for {
		w := strings.TrimPrefix(addr, "0")
		if addr == w {
			break
		}

		addr = w
	}

	return pfx + addr
}
