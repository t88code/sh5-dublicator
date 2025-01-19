package helper

import (
	"fmt"
)

func GetPath(filename, postfix string) string {
	return fmt.Sprintf("./json/%s_%s.json", filename, postfix)
}
