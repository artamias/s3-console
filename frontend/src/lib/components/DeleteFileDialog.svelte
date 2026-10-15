<script>
  import Icon from "./Icon.svelte";

  let { dialog = null, onClose, onDelete } = $props();
</script>

{#if dialog}
  <div
    role="dialog"
    aria-modal="true"
    aria-labelledby="delete-title"
    tabindex="-1"
    onmousedown={(event) => {
      if (event.target === event.currentTarget && !dialog.loading) onClose();
    }}
    class="fixed inset-0 z-50 flex items-center justify-center bg-[#17211D]/45 p-4 sm:p-8"
  >
    <div
      class="w-full max-w-md rounded-xl border border-[#D5DDD8] bg-white p-6 shadow-[0_24px_64px_rgba(23,33,29,0.18)]"
    >
      <span
        class="flex size-10 items-center justify-center rounded-lg bg-[#FCF2F1] text-[#A13A36]"
      >
        <Icon name="trash" class="size-5" />
      </span>
      <h2
        id="delete-title"
        class="mt-4 text-[20px] font-semibold tracking-[-0.02em]"
      >
        Hapus file?
      </h2>
      <p class="mt-1 text-[14px] text-[#59665F]">
        Objek ini akan dihapus permanen dari bucket.
      </p>
      <code
        class="mt-4 block break-all rounded-lg border border-[#DFE5E1] bg-[#F7F8F7] px-3 py-2 font-mono text-[13px] text-[#27362F]"
        >{dialog.name}</code
      >
      <div class="mt-6 flex justify-end gap-2.5">
        <button
          onclick={onClose}
          disabled={dialog.loading}
          class="inline-flex h-9 items-center justify-center rounded-lg border border-[#CBD4CE] bg-white px-4 text-[14px] font-medium text-[#27362F] hover:bg-[#F7F8F7] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#16745A] disabled:opacity-60"
          >Batal</button
        >
        <button
          onclick={onDelete}
          disabled={dialog.loading}
          class="inline-flex h-9 items-center justify-center gap-2 rounded-lg bg-[#A13A36] px-4 text-[14px] font-medium text-white transition hover:bg-[#8A302D] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#A13A36] disabled:opacity-60"
        >
          <Icon
            name={dialog.loading ? "refresh" : "trash"}
            class="size-4 {dialog.loading ? 'animate-spin' : ''}"
          />
          Hapus file
        </button>
      </div>
    </div>
  </div>
{/if}
