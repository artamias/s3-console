<script>
  import { iconPaths } from "../icon-paths.js";

  let { name, class: className = "" } = $props();
</script>

<svg
  aria-hidden="true"
  fill="none"
  viewBox="0 0 24 24"
  stroke="currentColor"
  stroke-width="1.8"
  stroke-linecap="round"
  stroke-linejoin="round"
  class={className}
>
  {@html iconPaths[name]}
</svg>
