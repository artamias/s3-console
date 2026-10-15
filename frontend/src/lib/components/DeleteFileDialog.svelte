<script>
  import { LoaderCircle, Trash2 } from 'lucide-svelte';

  let { dialog = null, onClose, onDelete } = $props();
</script>

{#if dialog}
  <div class="modal-backdrop" role="presentation" onclick={(event) => { if (event.target === event.currentTarget && !dialog.loading) onClose(); }}>
    <div class="modal delete-modal" role="dialog" tabindex="-1" aria-modal="true" aria-label="Konfirmasi hapus file">
      <span class="delete-mark"><Trash2 size={19} /></span>
      <h2>Hapus file?</h2>
      <p>Objek ini akan dihapus permanen dari bucket.</p>
      <code>{dialog.name}</code>
      <div class="modal-actions">
        <button class="button button-secondary" onclick={onClose} disabled={dialog.loading}>Batal</button>
        <button class="button button-delete" onclick={onDelete} disabled={dialog.loading}>{#if dialog.loading}<LoaderCircle size={16} class="spin" />{:else}<Trash2 size={16} />{/if}Hapus file</button>
      </div>
    </div>
  </div>
{/if}
