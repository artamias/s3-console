<script>
  import { HardDrive } from "lucide-svelte";

  let { online = false } = $props();
</script>

<header class="topbar">
  <div class="topbar-inner">
    <a class="brand" href="/" aria-label="S3 Console beranda">
      <span class="brand-mark"><HardDrive size={19} strokeWidth={2.2} /></span>
      <span><strong>Object Storage</strong><small>S3 Console</small></span>
    </a>
    <div
      class="connection"
      title={online ? "Backend terhubung" : "Backend tidak terhubung"}
    >
      <span class:online class="status-dot"></span>
      <span>{online ? "Backend terhubung" : "Backend offline"}</span>
    </div>
  </div>
</header>
