<script>
  import {
    CircleAlert,
    Download,
    ExternalLink,
    File,
    FileImage,
    FileVideo,
    FolderOpen,
    LoaderCircle,
    MoreHorizontal,
    Search,
    Trash2,
    Link,
  } from "lucide-svelte";
  import { fileKind, formatDate, formatSize, mediaKind } from "../format.js";

  let {
    files = [],
    info = null,
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

  function runAction(action, file) {
    openMenu = "";
    action(file);
  }
</script>

<section class="files-section" aria-labelledby="files-title">
  <div class="files-toolbar">
    <div class="files-title-wrap">
      <h2 id="files-title">
        Daftar file <span class="count">{files.length}</span>
      </h2>
      <span class="bucket-context"
        >Bucket <strong>{info?.bucket || "—"}</strong></span
      >
    </div>
    <label class="search-box">
      <Search size={16} />
      <input
        value={search}
        oninput={(event) => onSearch(event.currentTarget.value)}
        type="search"
        placeholder="Cari nama file"
        aria-label="Cari nama file"
      />
    </label>
  </div>

  {#if loading && !files.length}
    <div class="state-panel">
      <LoaderCircle size={22} class="spin" /><span
        >Memuat file dari backend...</span
      >
    </div>
  {:else if !online}
    <div class="state-panel error-state">
      <CircleAlert size={22} /><strong>S3 gagal terhubung</strong><span
        >Pastikan service S3 berjalan.</span
      >
    </div>
  {:else if error}
    <div class="state-panel error-state">
      <CircleAlert size={22} /><strong>Daftar file gagal dimuat</strong><span
        >{error}</span
      >
    </div>
  {:else if !visibleFiles.length}
    <div class="state-panel">
      <FolderOpen size={24} /><strong
        >{search ? "File tidak ditemukan" : "Bucket masih kosong"}</strong
      ><span
        >{search
          ? "Coba kata pencarian lain."
          : "Unggah file untuk mulai mengisi bucket."}</span
      >
    </div>
  {:else}
    <div class="table-scroll">
      <table>
        <thead
          ><tr
            ><th>Nama</th><th>Tipe</th><th>Ukuran</th><th>Terakhir diubah</th
            ><th><span class="sr-only">Aksi</span></th></tr
          ></thead
        >
        <tbody>
          {#each visibleFiles as file (file.name)}
            <tr>
              <td class="name-cell">
                <span class="file-icon">
                  {#if mediaKind(file.name) === "image"}<FileImage
                      size={17}
                    />{:else if mediaKind(file.name) === "video"}<FileVideo
                      size={17}
                    />{:else}<File size={17} />{/if}
                </span>
                <span class="name-copy"
                  ><strong title={file.name}>{file.name}</strong></span
                >
              </td>
              <td class="type-cell">{fileKind(file.name)}</td>
              <td class="mono-cell">{formatSize(file.size)}</td>
              <td class="date-cell">{formatDate(file.last_modified)}</td>
              <td class="action-cell">
                <div class="row-actions">
                  <button
                    class="icon-button"
                    onclick={() => onPreview(file)}
                    title="Lihat file"
                    aria-label={`Lihat ${file.name}`}
                    ><ExternalLink size={16} /></button
                  >
                  <div class="action-menu-wrap">
                    <button
                      class="icon-button"
                      class:active={openMenu === file.name}
                      onclick={() => {
                        openMenu = openMenu === file.name ? "" : file.name;
                      }}
                      title="Aksi lainnya"
                      aria-label={`Aksi lainnya untuk ${file.name}`}
                      aria-expanded={openMenu === file.name}
                      ><MoreHorizontal size={17} /></button
                    >
                    {#if openMenu === file.name}
                      <div class="action-menu" role="menu">
                        <button
                          role="menuitem"
                          onclick={() => runAction(onDownload, file)}
                          ><Download size={15} />Download</button
                        >
                        <button
                          role="menuitem"
                          onclick={() => runAction(onUrl, file)}
                          ><Link size={15} />Salin URL</button
                        >
                        <button
                          class="menu-danger"
                          role="menuitem"
                          onclick={() => runAction(onDelete, file)}
                          ><Trash2 size={15} />Hapus</button
                        >
                      </div>
                    {/if}
                  </div>
                </div>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
  <div class="table-footer">
    Menampilkan {visibleFiles.length} dari {files.length} file
  </div>
</section>
