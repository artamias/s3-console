<script>
  import Icon from "./Icon.svelte";

  let { urlDialog = null, onClose, onCopy, copied = false } = $props();
</script>

{#if urlDialog}
  <div
    role="dialog"
    aria-modal="true"
    aria-labelledby="url-title"
    tabindex="-1"
    onmousedown={(event) => event.target === event.currentTarget && onClose()}
    class="fixed inset-0 z-50 flex items-center justify-center bg-[#17211D]/45 p-4 sm:p-8"
  >
    <div
      class="w-full max-w-xl rounded-xl border border-[#D5DDD8] bg-white p-6 shadow-[0_24px_64px_rgba(23,33,29,0.18)]"
    >
      <div class="flex items-start justify-between gap-4">
        <div class="min-w-0">
          <div class="mb-1 text-[13px] font-medium text-[#16745A]">
            URL presigned
          </div>
          <h2
            id="url-title"
            class="truncate text-[20px] font-semibold tracking-[-0.02em]"
          >
            {urlDialog.name}
          </h2>
          <p class="mt-1 text-[13px] text-[#6B7770]">
            Berlaku {urlDialog.expires_in}
          </p>
        </div>
        <button
          onclick={onClose}
          aria-label="Tutup"
          class="flex size-9 shrink-0 items-center justify-center rounded-lg text-[#59665F] hover:bg-[#F0F2F1] hover:text-[#17211D] focus-visible:outline-2 focus-visible:outline-[#16745A]"
        >
          <Icon name="x" class="size-[18px]" />
        </button>
      </div>
      <code
        class="mt-4 block max-h-40 overflow-auto break-all rounded-lg border border-[#DFE5E1] bg-[#F7F8F7] px-3 py-2 font-mono text-[13px] text-[#27362F]"
        >{urlDialog.url}</code
      >
      <div class="mt-6 flex justify-end gap-2.5">
        <a
          href={urlDialog.url}
          target="_blank"
          rel="noreferrer"
          class="inline-flex h-9 items-center justify-center gap-2 rounded-lg border border-[#CBD4CE] bg-white px-4 text-[14px] font-medium text-[#27362F] hover:bg-[#F7F8F7] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#16745A]"
        >
          Buka <Icon name="external" class="size-4" />
        </a>
        <button
          onclick={() => onCopy(urlDialog.url)}
          class="inline-flex h-9 items-center justify-center gap-2 rounded-lg bg-[#16745A] px-4 text-[14px] font-medium text-white transition hover:bg-[#125F49] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#16745A]"
        >
          <Icon name={copied ? "check" : "copy"} class="size-4" />
          {copied ? "Tersalin" : "Salin URL"}
        </button>
      </div>
    </div>
  </div>
{/if}
