(()=>{document.addEventListener("alpine:init",()=>{Alpine.data("pageIndex",()=>({init(){this.$refs.formSearch.action="javascript:void(0)";let t=document.addEventListener("keydown",e=>{if(!(["INPUT","TEXTAREA"].includes(document.activeElement.tagName)&&document.activeElement.type==="text"))switch(e.key){case"n":{this.$refs.inputAddNew&&(this.$refs.inputAddNew.focus(),e.preventDefault());break}case"f":{this.$refs.inputSearch&&(this.$refs.inputSearch.focus(),e.preventDefault());break}}});this.$destroy=()=>{document.removeEventListener("keydown",t)}}}))});})();
