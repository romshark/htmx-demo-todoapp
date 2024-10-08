// Tell eslint that Alpine is a global
/* global Alpine */

document.addEventListener("alpine:init", () => {
  Alpine.data("pageIndex", () => ({
    init() {
      this.$refs.formSearch.action = "javascript:void(0)";
      const h = document.addEventListener("keydown", (e) => {
        if (
          ["INPUT", "TEXTAREA"].includes(document.activeElement.tagName) &&
          document.activeElement.type === "text"
        ) {
          return;
        }

        switch (e.key) {
          case "n": {
            if (this.$refs.inputAddNew) {
              this.$refs.inputAddNew.focus();
              e.preventDefault();
            }
            break;
          }
          case "f": {
            if (this.$refs.inputSearch) {
              this.$refs.inputSearch.focus();
              e.preventDefault();
            }
            break;
          }
        }
      });
      this.$destroy = () => {
        document.removeEventListener("keydown", h);
      };
    },
  }));
});

// hxThresholdLoadingIndicator defines how long to wait before
// showing loading indication after an HTMX request was sent.
// If the request finishes fast (response time is below threshold)
// then no loading indication should be applied.
const hxThresholdLoadingIndicator = 150;

document.addEventListener("htmx:beforeRequest", function (event) {
  var target = event.detail.target;
  target.htmxTimeoutId = setTimeout(function () {
    target.classList.add("non-interactable");
  }, hxThresholdLoadingIndicator);
});

document.addEventListener("htmx:afterRequest", function (event) {
  var target = event.detail.target;

  // Clear the timeout if the request finishes before the threshold
  if (target.htmxTimeoutId) {
    clearTimeout(target.htmxTimeoutId);
    target.htmxTimeoutId = null; // Clear the timeout ID
  }

  // Remove the class (just in case it was applied)
  target.classList.remove("non-interactable");
});
