package internal

import (
	"os"
	"fmt"
)

var DebugEnabled bool

func Debug(format string, a ...any) {
	if DebugEnabled {
		_, _ = os.Stdout.WriteString("[cedpm] " + fmt.Sprintf(format, a...) + "\n")
	}
}
