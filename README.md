<a href="https://goreportcard.com/report/github.com/romshark/htmx-demo-todoapp">
    <img src="https://goreportcard.com/badge/github.com/romshark/htmx-demo-todoapp" alt="GoReportCard">
</a>
<a href="https://pkg.go.dev/github.com/romshark/htmx-demo-todoapp">
    <img src="https://godoc.org/github.com/romshark/htmx-demo-todoapp?status.svg" alt="GoDoc">
</a>

# HTMX todo app

A simple todo application powered by [Go](https://go.dev/), [HTMX](https://htmx.org),
[Templ](https://templ.guide) and [TailwindCSS](https://tailwindcss.com/).

- **Graceful Degradation**: This app continues to provide limited core functionality
  even when JavaScript is disabled by utilizing
  [303 redirects](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/303),
  [HTML forms](https://www.w3schools.com/html/html_forms.asp)
  and the [`HX-Request`](https://htmx.org/reference/#headers) header.
- **Bundled CSS**: [PostCSS](https://postcss.org/) is used to build the CSS bundle.

## Dev mode

The server can be executed by simply running `go run . -host=:8080` in the root directory.

Building the CSS bundle requires [Node.js](https://nodejs.org/) and `npm`
to be installed on your system.

Run `npm i` to download all JS dependencies.

For interactive development mode with hot-reload install
[Templiér](https://github.com/romshark/templier)

And run `templier --config templier.yml`

Note: you might need to change `app.dir-src-root` and `app.dir-cmd`
due to a bug in Templiér, see `templier.yml`.
