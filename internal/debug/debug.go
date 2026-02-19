package debug

import (
	"fmt"
	"os"
)

var enabled bool

func init() {
	enabled = os.Getenv("NUON_DEBUG") == "true"
}

// Enabled returns whether debug logging is active.
func Enabled() bool {
	return enabled
}

// Log prints a debug message to stderr if NUON_DEBUG=true.
func Log(format string, args ...any) {
	if !enabled {
		return
	}
	fmt.Fprintf(os.Stderr, "[debug] "+format+"\n", args...)
}
