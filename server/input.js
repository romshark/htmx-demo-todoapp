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
