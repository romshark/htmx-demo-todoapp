<a href="https://goreportcard.com/report/github.com/romshark/htmx-demo-todoapp">
    <img src="https://goreportcard.com/badge/github.com/romshark/htmx-demo-todoapp" alt="GoReportCard">
</a>
<a href="https://pkg.go.dev/github.com/romshark/htmx-demo-todoapp">
    <img src="https://godoc.org/github.com/romshark/htmx-demo-todoapp?status.svg" alt="GoDoc">
</a>

![todoapp_1](https://github.com/user-attachments/assets/789059b1-84d9-4694-b87c-5f7a32d63e64)

# HTMX Tech-Demo - Todo App

A simple todo application powered by [Go](https://go.dev/), [HTMX](https://htmx.org),
[Templ](https://templ.guide), [TailwindCSS](https://tailwindcss.com/) and
[Alpine.js](https://alpinejs.dev/).

The server can be executed by simply running it in the root directory:

```sh
go run . -host=:8080
```

- **Graceful Degradation**: This app continues to provide limited core functionality
  even when JavaScript is disabled by utilizing
  [303 redirects](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/303),
  [HTML forms](https://www.w3schools.com/html/html_forms.asp)
  and the [`HX-Request`](https://htmx.org/reference/#headers) header.
  JavaScript progressively enhances the UX through partial AJAX reloads
  (full page reloads without JS) using HTMX, allowing
  keyboard shortcuts (press "f" to focus the search part; press
  "n" to focus the "New Todo" input field), and more.
- **Bundled CSS**: [PostCSS](https://postcss.org/) is used to build the CSS bundle.
  [Templiér](https://github.com/romshark/templier) is configured to automatically watch
  all relevant `.css` and `.templ` files, build the bundle and reload the browser tab
  (see [dev mode](#dev-mode)).
- **Bundled JavaScript**: [eslint](https://eslint.org/) and
  [esbuild](https://esbuild.github.io/) are used to build the `dist.js` JavaScript bundle.
  [Templiér](https://github.com/romshark/templier) is configured to automatically watch
  all relevant `.css` and `.templ` files, build the bundle and reload the browser tab
  (see [dev mode](#dev-mode))

## Dev mode

Building the CSS bundle requires [Node.js](https://nodejs.org/) and `npm`
to be installed on your system.

For interactive development mode with automatic live-reload using
[Templiér](https://github.com/romshark/templier), simply run:

```sh
./dev.sh
```
