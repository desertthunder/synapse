package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	sb := strings.Builder{}
	l := NewLogger(&sb, InfoLevel, "[test]")
	msg := "sent from default"

	t.Run("debug level not printed", func(t *testing.T) {
		l.Debug(msg)

		got := sb.String()
		want := ""

		if got != want {
			t.Errorf("got %v but wanted %v", got, want)
		}

		sb.Reset()
	})

	t.Run("debug level printed after updating level", func(t *testing.T) {
		l.SetLevel(DebugLevel)
		l.Debug(msg)

		got := sb.String()
		want := "DEBU"

		if !strings.Contains(got, want) {
			t.Errorf("got %v want %v", got, want)
		}

		sb.Reset()
	})

	for _, tc := range []struct {
		desc string
		fn   func(msg interface{}, keyvals ...interface{})
		want string
	}{
		{"info", l.Info, "INFO"},
		{"warn", l.Warn, "WARN"},
		{"error", l.Error, "ERRO"},
		{"no", l.Print, msg},
	} {
		t.Run(fmt.Sprintf("prints %v level", tc.desc), func(t *testing.T) {
			tc.fn(msg)

			got := sb.String()
			if !strings.Contains(got, tc.want) {
				t.Errorf("got %v want %v", got, tc.want)
			}

			sb.Reset()
		})
	}
}
