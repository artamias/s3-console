<script>
  import { fileKind, formatDate, formatSize } from "../format.js";
  import Icon from "./Icon.svelte";

  let { preview = null, onClose, onDownload } = $props();

  let loadError = $state(false);

  $effect(() => {
    preview?.url;
    loadError = false;
  });
</script>

{#if preview}
  <div
    role="dialog"
    aria-modal="true"
    aria-labelledby="preview-title"
    tabindex="-1"
    onmousedown={(event) => event.target === event.currentTarget && onClose()}
    class="fixed inset-0 z-50 flex items-center justify-center bg-[#17211D]/45 p-4 sm:p-8"
  >
    <div
      class="flex max-h-[90vh] w-full max-w-4xl flex-col overflow-hidden rounded-xl border border-[#D5DDD8] bg-white shadow-[0_24px_64px_rgba(23,33,29,0.18)]"
    >
      <div
        class="flex items-center justify-between gap-4 border-b border-[#DFE5E1] px-5 py-4 sm:px-6"
      >
        <div class="min-w-0">
          <div class="mb-1 text-[13px] font-medium text-[#16745A]">
            Pratinjau file
          </div>
          <h2
            id="preview-title"
            class="truncate text-[20px] font-semibold tracking-[-0.02em]"
          >
            {preview.filename}
          </h2>
        </div>
        <button
          onclick={onClose}
          aria-label="Tutup pratinjau"
          class="flex size-9 shrink-0 items-center justify-center rounded-lg text-[#59665F] hover:bg-[#F0F2F1] hover:text-[#17211D] focus-visible:outline-2 focus-visible:outline-[#16745A]"
        >
          <Icon name="x" class="size-[18px]" />
        </button>
      </div>

      <div class="min-h-0 flex-1 overflow-auto bg-[#F7F8F7] p-4 sm:p-6">
        {#if preview.kind === "image" && !loadError}
          <div
            class="flex min-h-80 items-center justify-center rounded-lg border border-[#DFE5E1] bg-white p-4"
          >
            <img
              src={preview.url}
              alt={`Pratinjau ${preview.filename}`}
              class="max-h-[58vh] max-w-full object-contain"
              onerror={() => (loadError = true)}
            />
          </div>
        {:else if preview.kind === "video" && !loadError}
          <div
            class="flex min-h-80 items-center justify-center rounded-lg border border-[#DFE5E1] bg-white p-4"
          >
            <!-- svelte-ignore a11y_media_has_caption -->
            <video
              src={preview.url}
              controls
              class="max-h-[58vh] max-w-full"
              onerror={() => (loadError = true)}
            ></video>
          </div>
        {:else if preview.kind === "pdf" && !loadError}
          <iframe
            title={`Pratinjau ${preview.filename}`}
            src={preview.url}
            class="h-[58vh] w-full rounded-lg border border-[#DFE5E1] bg-white"
          ></iframe>
        {:else}
          <div
            class="flex min-h-80 flex-col items-center justify-center rounded-lg border border-[#DFE5E1] bg-white px-6 text-center"
          >
            <span
              class="flex size-12 items-center justify-center rounded-xl bg-[#E7F2EE] text-[#16745A]"
            >
              <Icon name="file" class="size-6" />
            </span>
            <div class="mt-4 text-[15px] font-medium">
              {loadError
                ? "File gagal dimuat"
                : "Pratinjau visual tidak tersedia"}
            </div>
            <p class="mt-1 max-w-sm text-[13px] leading-5 text-[#6B7770]">
              {loadError
                ? "URL file tidak dapat dibuka. Periksa konfigurasi CORS MinIO atau coba unduh file."
                : "Format file ini tidak dapat ditampilkan langsung. Anda tetap dapat melihat informasi atau mengunduh file."}
            </p>
          </div>
        {/if}
      </div>

      <div
        class="flex flex-col justify-between gap-4 border-t border-[#DFE5E1] bg-[#FAFBFA] px-5 py-3.5 sm:flex-row sm:items-center sm:px-6"
      >
        <div
          class="grid flex-1 grid-cols-2 gap-x-8 gap-y-3 text-[13px] sm:grid-cols-[96px_96px_minmax(180px,1fr)]"
        >
          <div class="flex min-w-0 flex-col gap-1">
            <span class="text-[#6B7770]">Tipe</span>
            <span class="truncate font-medium text-[#27362F]"
              >{fileKind(preview.filename)}</span
            >
          </div>
          <div class="flex min-w-0 flex-col gap-1">
            <span class="text-[#6B7770]">Ukuran</span>
            <span class="truncate font-medium text-[#27362F]"
              >{formatSize(preview.file?.size)}</span
            >
          </div>
          <div class="col-span-2 flex min-w-0 flex-col gap-1 sm:col-span-1">
            <span class="text-[#6B7770]">Diubah</span>
            <span class="truncate font-medium text-[#27362F]"
              >{formatDate(preview.file?.last_modified)}</span
            >
          </div>
        </div>
        <button
          onclick={() => onDownload({ name: preview.filename })}
          class="inline-flex h-9 shrink-0 items-center justify-center gap-2 rounded-lg border border-[#CBD4CE] bg-white px-4 text-[14px] font-medium text-[#27362F] hover:bg-[#F7F8F7] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#16745A]"
        >
          <Icon name="download" class="size-4" /> Unduh
        </button>
      </div>
    </div>
  </div>
{/if}
