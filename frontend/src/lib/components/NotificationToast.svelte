<script>
  import { Check, CircleAlert, X } from 'lucide-svelte';

  let { toast = null, onClose } = $props();
</script>

{#if toast}
  <div class:error-toast={toast.type === 'error'} class="toast" role="status">
    {#if toast.type === 'error'}<CircleAlert size={17} />{:else}<Check size={17} />{/if}
    <span>{toast.message}</span>
    <button class="toast-close" onclick={onClose} aria-label="Tutup notifikasi"><X size={16} /></button>
  </div>
{/if}
