# HTMX todo app

A simple todo application powered by HTMX demonstrating an approach
to structuring and connecting the HTTP endpoints with "frames",
each of which may have "variants", for each of which there can be
"actions" (methods, essentially).

```html
<button hx-post="/hx/list;var=search-result;act=toggle/123">
```

The above `hx-post` triggers method `toggle` on todo with id `123` on variant `search-result` for the frame `#frame-list`.
`POST /hx/list;var=search-result;act=toggle/123` returns rendered HTML to replace the outerHTML of frame `#frame-list` with.

Probably very similar to https://turbo.hotwired.dev/.
