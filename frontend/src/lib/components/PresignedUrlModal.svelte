<script>
  import { Copy, ExternalLink, X } from 'lucide-svelte';

  let { urlDialog = null, onClose, onCopy, copied = false } = $props();
</script>

{#if urlDialog}
  <div class="modal-backdrop" role="presentation" onclick={(event) => { if (event.target === event.currentTarget) onClose(); }}>
    <div class="modal url-modal" role="dialog" tabindex="-1" aria-modal="true" aria-label="URL presigned file">
      <div class="modal-heading">
        <div><h2>URL presigned</h2><p>{urlDialog.name} · berlaku {urlDialog.expires_in}</p></div>
        <button class="icon-button" onclick={onClose} aria-label="Tutup"><X size={19} /></button>
      </div>
      <code class="url-block">{urlDialog.url}</code>
      <div class="modal-actions">
        <button class="button button-primary" onclick={() => onCopy(urlDialog.url)}><Copy size={15} />{copied ? 'Tersalin' : 'Salin URL'}</button>
        <a class="button button-secondary" href={urlDialog.url} target="_blank" rel="noreferrer">Buka <ExternalLink size={15} /></a>
      </div>
    </div>
  </div>
{/if}
