package dynamicstyler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/oklog/ulid/v2"
)

// Styler stores and serves styles
// classes are created dynamicaly with ULID
type Styler struct {
	styles map[string]func() string

	ID  string
	Url string
}

// New create new styler, sets its url
// Styler.Handler() can be used with any other url too
func New() *Styler {
	id := ulid.Make().String()
	url := fmt.Sprintf("/static/dynstyles/styler-%s", id)
	return &Styler{
		styles: make(map[string]func() string),
		ID:     id,
		Url:    url,
	}
}

// Raw adds styles as raw string
func (st *Styler) Raw(cssFunc func() string) string {
	class := fmt.Sprintf("styler_%s", ulid.Make().String())
	st.styles[class] = cssFunc

	return class
}

// Props add style from props (map[string]string)
func (st *Styler) Props(propsFunc func() Props) string {
	class := fmt.Sprintf("styler_%s", ulid.Make().String())
	st.styles[class] = func() string {
		return formatProps(propsFunc())
	}

	return class
}

// WriteAll writes styles to io.Writer in standard css format
// .classname { ...props... }
func (st *Styler) WriteAll(w io.Writer) error {
	for class, style := range st.styles {
		_, err := fmt.Fprintf(w, ".%s { %s}\n", class, style())
		if err != nil {
			return err
		}
	}
	return nil
}

// Handler creates new http handler
// it recreates styles on every request
// no caching is done
// recommended to be used with Styler.ID or Styler.Url as path
func (st *Styler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		st.WriteAll(w)
	}
}
