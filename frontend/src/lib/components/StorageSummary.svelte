<script>
  import { Archive, FolderOpen, HardDrive, Link } from "lucide-svelte";
  import { formatSize } from "../format.js";

  let { info = null } = $props();
</script>

<section class="summary-row" aria-label="Ringkasan penyimpanan">
  <div class="summary-item">
    <span class="summary-icon teal"><HardDrive size={18} /></span>
    <div>
      <span class="summary-label">Bucket</span><strong
        >{info?.bucket || "—"}</strong
      >
    </div>
  </div>
  <div class="summary-item">
    <span class="summary-icon orange"><FolderOpen size={18} /></span>
    <div>
      <span class="summary-label">Jumlah file</span><strong
        >{info?.files ?? "—"}</strong
      >
    </div>
  </div>
  <div class="summary-item">
    <span class="summary-icon green"><Archive size={18} /></span>
    <div>
      <span class="summary-label">Total ukuran</span><strong
        >{info ? formatSize(info.total_size) : "—"}</strong
      >
    </div>
  </div>
  <div class="summary-item endpoint-item">
    <span class="summary-icon blue"><Link size={18} /></span>
    <div>
      <span class="summary-label">Endpoint S3</span><strong
        >{info?.endpoint || "—"}</strong
      >
    </div>
  </div>
</section>
