<script>
  import { formatSize } from "../format.js";
  import Icon from "./Icon.svelte";

  let { info = null } = $props();

  let items = $derived([
    { label: "Bucket", value: info?.bucket || "—", icon: "box" },
    {
      label: "Files",
      value: info ? `${info.files} files` : "—",
      icon: "file",
    },
    {
      label: "Size",
      value: info ? formatSize(info.total_size) : "—",
      icon: "hard-drive",
    },
    { label: "Endpoint", value: info?.endpoint || "—", icon: "server" },
  ]);
</script>

<section
  aria-label="Ringkasan penyimpanan"
  class="mt-9 grid border-y border-[#DFE5E1] sm:grid-cols-2 lg:grid-cols-[1fr_1fr_1fr_1.5fr]"
>
  {#each items as item, index (item.label)}
    <div
      class="flex min-w-0 items-center gap-3 py-5 {index % 2
        ? 'sm:border-l'
        : ''} {index > 0
        ? 'lg:border-l'
        : ''} border-[#DFE5E1] sm:px-5 lg:px-6 first:pl-0"
    >
      <Icon name={item.icon} class="size-[18px] shrink-0 text-[#6B7770]" />
      <div class="min-w-0">
        <div class="text-[13px] text-[#6B7770]">{item.label}</div>
        <div class="mt-0.5 truncate text-[15px] font-medium text-[#27362F]">
          {item.value}
        </div>
      </div>
    </div>
  {/each}
</section>
