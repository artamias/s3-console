<script>
  import Icon from "./Icon.svelte";

  let { online = false } = $props();
</script>

<header class="border-b border-[#DFE5E1] bg-white">
  <div
    class="mx-auto flex h-16 max-w-[1440px] items-center justify-between px-5 sm:px-8 lg:px-12"
  >
    <a class="flex items-center gap-3" href="/" aria-label="S3 Console beranda">
      <div
        class="flex size-8 items-center justify-center rounded-lg bg-[#16745A] text-white"
      >
        <Icon name="box" class="size-[18px]" />
      </div>
      <span class="text-[15px] font-semibold tracking-[-0.01em]">S3 Console</span>
    </a>
    <div
      class="flex items-center gap-2.5 rounded-full border border-[#DFE5E1] px-3 py-1.5 text-[13px] font-medium text-[#425149]"
      title={online ? "Backend terhubung" : "Backend tidak terhubung"}
    >
      <span
        class="size-2 rounded-full ring-4 {online
          ? 'bg-[#16745A] ring-[#E8F3EF]'
          : 'bg-[#C34C42] ring-[#FBEAE8]'}"
      ></span>
      <span class="hidden sm:inline"
        >{online ? "S3 terhubung" : "S3 tidak terhubung"}</span
      >
      <span class="sm:hidden">{online ? "Terhubung" : "Offline"}</span>
    </div>
  </div>
</header>
