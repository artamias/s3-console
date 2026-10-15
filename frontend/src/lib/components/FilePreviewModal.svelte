<script>
  import { Copy, Download, ExternalLink, File, X } from "lucide-svelte";

  let {
    preview = null,
    onClose,
    onCopy,
    onDownload,
    copied = false,
  } = $props();
</script>

{#if preview}
  <div
    class="modal-backdrop"
    role="presentation"
    onclick={(event) => {
      if (event.target === event.currentTarget) onClose();
    }}
  >
    <div
      class="modal preview-modal"
      role="dialog"
      tabindex="-1"
      aria-modal="true"
      aria-label={`Preview ${preview.filename}`}
    >
      <div class="modal-heading">
        <div>
          <h2>{preview.filename}</h2>
          <p>URL S3 · berlaku {preview.expires_in}</p>
        </div>
        <button class="icon-button" onclick={onClose} aria-label="Tutup"
          ><X size={19} /></button
        >
      </div>
      <div class="preview-content">
        {#if preview.kind === "image" && !preview.error}
          <img
            src={preview.url}
            alt={preview.filename}
            onerror={() => {
              preview.error = true;
            }}
          />
        {:else if preview.kind === "video" && !preview.error}
          <!-- svelte-ignore a11y_media_has_caption -->
          <video
            src={preview.url}
            controls
            onerror={() => {
              preview.error = true;
            }}
          ></video>
        {:else}
          <div class="preview-unavailable">
            <File size={30} /><span
              >{preview.error
                ? "File gagal dimuat."
                : "Pratinjau tidak tersedia untuk tipe file ini."}</span
            >
          </div>
        {/if}
      </div>
      <div class="url-line">
        <code>{preview.url}</code><button
          class="button button-secondary"
          onclick={() => onCopy(preview.url)}
          ><Copy size={15} />{copied ? "Tersalin" : "Salin"}</button
        >
      </div>
      <div class="modal-actions">
        <a
          class="button button-primary"
          href={preview.url}
          target="_blank"
          rel="noreferrer">Buka di tab baru <ExternalLink size={15} /></a
        >
        <button
          class="button button-secondary"
          onclick={() => onDownload({ name: preview.filename })}
          ><Download size={15} />Download</button
        >
      </div>
    </div>
  </div>
{/if}
