package styler

import (
	"fmt"
	"strings"
)

// Props are a set of css properties, mapped name->prop
// props are not validated
type Props map[string]string

func formatProps(props Props) string {
	res := &strings.Builder{}
	for name, prop := range props {
		res.WriteString(fmt.Sprintf("%s:%s; ", name, prop))
	}
	return res.String()
}
