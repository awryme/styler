# styler
Create css styles / classes in go, serve them with http handler

Use generated classnames in pages

## Example usage with gomponents

<details>
  <summary>main.go</summary>
  
  ```go
package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/awryme/styler"
	"github.com/awryme/styler/dynamicstyler"
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

// create styles from map[string]string
var circle = styler.Default.Props(styler.Props{
	"border-style":     "solid",
	"border-width":     "1px",
	"border-radius":    "100%",
	"width":            "10em",
	"height":           "10em",
	"background-color": "red",
})

var textWhite = styler.Default.Props(styler.Props{
	"color": "white",
})

// create styles from raw string
var flexCenter = styler.Default.Raw(`
display: flex;
justify-content: center;
align-items: center;
`)

// create dynamically updated styles, because why not
var dynStyles = dynamicstyler.New()

// increase font-size on every refresh (every time stylesheet is fetched)
var txtsize uint64 = 10
var incText = dynStyles.Props(func() dynamicstyler.Props {
	size := atomic.AddUint64(&txtsize, 1)
	return dynamicstyler.Props{
		"font-size": fmt.Sprintf("%dpx", size),
	}
})

// basic page with gomponents
func myPage() g.Node {
	return c.HTML5(c.HTML5Props{
		Head: []g.Node{
			h.Link(h.Rel("stylesheet"), h.Href(styler.Default.Url)),
			h.Link(h.Rel("stylesheet"), h.Href(dynStyles.Url)),
		},
		Body: []g.Node{
			// make red circle, center items, white text
			h.Div(c.Classes{circle: true, flexCenter: true, textWhite: true},
				h.Span(h.Class(incText),
					g.Text("text"),
				),
			),
		},
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		myPage().Render(w)
	})
	fmt.Println("static styles at url", styler.Default.Url)
	fmt.Println("dynamic styles at url", dynStyles.Url)
	mux.HandleFunc(styler.Default.Url, styler.Default.Handler())
	mux.HandleFunc(dynStyles.Url, dynStyles.Handler())
	http.ListenAndServe(":9090", mux)
}
```
</details>