package greeting

import "fmt"

func Greet(somebody string) string {
	return fmt.Sprintf("hello %s", somebody)
}
