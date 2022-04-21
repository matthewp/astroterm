package debug

import (
	"fmt"
	"os"
)

func Log(tmpl string, vals ...interface{}) {
	f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	str := fmt.Sprintf(tmpl, vals...)
	if _, err := f.WriteString(str); err != nil {
		fmt.Println(err)
	}
}
