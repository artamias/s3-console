<script>
  import { fileKind, formatDate, formatSize, isFolder, mediaKind } from "../format.js";
  import Icon from "./Icon.svelte";

  let {
    files = [],
    loading = false,
    online = false,
    error = "",
    search = "",
    onSearch,
    onPreview,
    onDownload,
    onUrl,
    onDelete,
  } = $props();

  let openMenu = $state("");
  let visibleFiles = $derived(
    files.filter((file) =>
      file.name.toLowerCase().includes(search.trim().toLowerCase()),
    ),
  );

  const ROW_ICONS = { folder: "folder", image: "image", video: "video" };

  function rowIcon(file) {
    return ROW_ICONS[mediaKind(file.name)] ?? "file";
  }

  function runAction(action, file) {
    openMenu = "";
    action(file);
  }

  function closeMenu(event) {
    if (!event.target.closest?.("[data-row-menu]")) openMenu = "";
  }
</script>

<svelte:window
  onclick={closeMenu}
  onkeydown={(event) => {
    if (event.key === "Escape") openMenu = "";
  }}
/>

<section class="mt-10" aria-labelledby="files-title">
  <div
    class="mb-4 flex flex-col justify-between gap-4 sm:flex-row sm:items-end"
  >
    <div>
      <h2 id="files-title" class="text-[20px] font-semibold tracking-[-0.02em]">
        List File
      </h2>
    </div>
    <label class="relative block w-full sm:w-72">
      <span class="sr-only">Cari file</span>
      <Icon
        name="search"
        class="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-[#7C8881]"
      />
      <input
        value={search}
        oninput={(event) => onSearch(event.currentTarget.value)}
        placeholder="Cari nama file..."
        class="h-10 w-full rounded-lg border border-[#CBD4CE] bg-white pl-9 pr-9 text-[14px] outline-none transition placeholder:text-[#909A94] focus:border-[#16745A] focus:ring-2 focus:ring-[#16745A]/15"
      />
      {#if search}
        <button
          aria-label="Hapus pencarian"
          onclick={() => onSearch("")}
          class="absolute right-2 top-1/2 flex size-7 -translate-y-1/2 items-center justify-center rounded-md text-[#6B7770] hover:bg-[#F0F2F1] focus-visible:outline-2 focus-visible:outline-[#16745A]"
        >
          <Icon name="x" class="size-4" />
        </button>
      {/if}
    </label>
  </div>

  <div class="overflow-visible rounded-xl border border-[#DFE5E1] bg-white">
    <div
      class="hidden grid-cols-[minmax(280px,2fr)_1fr_120px_180px_76px] items-center border-b border-[#DFE5E1] bg-[#FAFBFA] px-5 py-3 text-[13px] font-medium text-[#59665F] md:grid"
    >
      <div>Name</div>
      <div>Type</div>
      <div class="flex items-center justify-end">Size</div>
      <div class="pl-8">Updated</div>
      <div class="text-right">Action</div>
    </div>

    {#if loading && !files.length}
      <div
        class="flex min-h-60 flex-col items-center justify-center text-[#59665F]"
      >
        <Icon name="refresh" class="mb-3 size-5 animate-spin" />
        <span class="text-[14px]">Memuat objek...</span>
      </div>
    {:else if !online}
      <div
        class="flex min-h-64 flex-col items-center justify-center px-6 text-center"
      >
        <div
          class="flex size-10 items-center justify-center rounded-lg bg-[#FCF2F1] text-[#A13A36]"
        >
          <Icon name="alert" class="size-5" />
        </div>
        <div class="mt-3 text-[15px] font-medium">S3 gagal terhubung</div>
        <p class="mt-1 text-[13px] text-[#6B7770]">
          Pastikan service S3 berjalan.
        </p>
      </div>
    {:else if error}
      <div
        class="flex min-h-64 flex-col items-center justify-center px-6 text-center"
      >
        <div
          class="flex size-10 items-center justify-center rounded-lg bg-[#FCF2F1] text-[#A13A36]"
        >
          <Icon name="alert" class="size-5" />
        </div>
        <div class="mt-3 text-[15px] font-medium">Daftar file gagal dimuat</div>
        <p class="mt-1 text-[13px] text-[#6B7770]">{error}</p>
      </div>
    {:else if !visibleFiles.length}
      <div
        class="flex min-h-64 flex-col items-center justify-center px-6 text-center"
      >
        <div
          class="flex size-10 items-center justify-center rounded-lg bg-[#F0F2F1] text-[#6B7770]"
        >
          <Icon name={search ? "search" : "folder"} class="size-5" />
        </div>
        <div class="mt-3 text-[15px] font-medium">
          {search ? "Tidak ada file yang ditemukan" : "Bucket masih kosong"}
        </div>
        <p class="mt-1 text-[13px] text-[#6B7770]">
          {search
            ? "Coba gunakan kata kunci lain atau hapus pencarian."
            : "Unggah file untuk mulai mengisi bucket."}
        </p>
        {#if search}
          <button
            onclick={() => onSearch("")}
            class="mt-4 text-[14px] font-medium text-[#16745A] hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#16745A]"
            >Hapus pencarian</button
          >
        {/if}
      </div>
    {:else}
      {#each visibleFiles as file (file.name)}
        {@const folder = isFolder(file.name)}
        <div
          class="group relative grid grid-cols-[1fr_auto] gap-x-3 border-b border-[#E8ECE9] px-4 py-4 transition last:border-0 hover:bg-[#FAFBFA] md:grid-cols-[minmax(280px,2fr)_1fr_120px_180px_76px] md:items-center md:px-5 md:py-0 md:min-h-16"
        >
          <div class="flex min-w-0 items-center gap-3">
            <span
              class="flex size-9 shrink-0 items-center justify-center rounded-lg {folder
                ? 'bg-[#E7F2EE] text-[#16745A]'
                : 'bg-[#F0F2F1] text-[#59665F]'}"
            >
              <Icon name={rowIcon(file)} class="size-[18px]" />
            </span>
            <span
              class="truncate text-[15px] font-medium text-[#27362F]"
              title={file.name}>{file.name}</span
            >
          </div>
          <div
            class="col-start-1 mt-2 flex items-center gap-2 pl-12 text-[13px] text-[#6B7770] md:col-auto md:mt-0 md:block md:pl-0 md:text-[15px]"
          >
            <span>{fileKind(file.name)}</span><span class="md:hidden">·</span
            ><span class="md:hidden">{folder ? "—" : formatSize(file.size)}</span>
          </div>
          <div
            class="hidden text-right text-[15px] tabular-nums text-[#425149] md:block"
          >
            {folder ? "—" : formatSize(file.size)}
          </div>
          <div
            class="col-start-1 pl-12 text-[13px] text-[#6B7770] md:col-auto md:pl-8 md:text-[15px] md:text-[#425149]"
          >
            {formatDate(file.last_modified)}
          </div>
          <div
            data-row-menu
            class="relative col-start-2 row-span-3 row-start-1 flex items-center justify-end md:col-auto md:row-auto"
          >
            <button
              aria-label={`Aksi untuk ${file.name}`}
              aria-haspopup="menu"
              aria-expanded={openMenu === file.name}
              onclick={() => {
                openMenu = openMenu === file.name ? "" : file.name;
              }}
              class="flex size-9 items-center justify-center rounded-lg text-[#59665F] transition hover:bg-[#EDF0EE] hover:text-[#17211D] focus-visible:outline-2 focus-visible:outline-[#16745A]"
            >
              <Icon name="more" class="size-[18px]" />
            </button>
            {#if openMenu === file.name}
              <div
                role="menu"
                class="absolute right-0 top-10 z-20 w-44 overflow-hidden rounded-lg border border-[#D8DFDB] bg-white py-1 shadow-[0_8px_24px_rgba(23,33,29,0.12)] md:top-11"
              >
                {#if !folder}
                  <button
                    role="menuitem"
                    onclick={() => runAction(onDownload, file)}
                    class="flex w-full items-center gap-2.5 px-3 py-2 text-left text-[14px] text-[#27362F] hover:bg-[#F3F5F4] focus:bg-[#F3F5F4] focus:outline-none"
                  >
                    <Icon name="download" class="size-4 text-[#6B7770]" /> Unduh
                  </button>
                  <button
                    role="menuitem"
                    onclick={() => runAction(onPreview, file)}
                    class="flex w-full items-center gap-2.5 px-3 py-2 text-left text-[14px] text-[#27362F] hover:bg-[#F3F5F4] focus:bg-[#F3F5F4] focus:outline-none"
                  >
                    <Icon name="eye" class="size-4 text-[#6B7770]" /> Lihat
                  </button>
                  <button
                    role="menuitem"
                    onclick={() => runAction(onUrl, file)}
                    class="flex w-full items-center gap-2.5 px-3 py-2 text-left text-[14px] text-[#27362F] hover:bg-[#F3F5F4] focus:bg-[#F3F5F4] focus:outline-none"
                  >
                    <Icon name="link" class="size-4 text-[#6B7770]" /> Salin URL
                  </button>
                {/if}
                <button
                  role="menuitem"
                  onclick={() => runAction(onDelete, file)}
                  class="flex w-full items-center gap-2.5 px-3 py-2 text-left text-[14px] text-[#A13A36] hover:bg-[#FCF2F1] focus:bg-[#FCF2F1] focus:outline-none"
                >
                  <Icon name="trash" class="size-4" /> Hapus
                </button>
              </div>
            {/if}
          </div>
        </div>
      {/each}
    {/if}
    <div
      class="border-t border-[#DFE5E1] bg-[#FAFBFA] px-4 py-3 text-[13px] text-[#6B7770] md:px-5"
    >
      {files.length} items in this bucket
    </div>
  </div>
</section>
