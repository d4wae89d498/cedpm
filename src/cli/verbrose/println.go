package verbrose

import (
	"fmt"
)

var Enabled bool

func Println(args ...any) {
	if !Enabled {
		return
	}
	fmt.Printf("[[ VERBROSE ]] ")
	fmt.Println(args...)
}


func Printf(str string, args ...any) {
	if !Enabled {
		return
	}
	fmt.Printf("[[ VERBROSE ]] ")
	fmt.Printf(str, args...)
}
