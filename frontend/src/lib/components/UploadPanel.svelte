<script>
  import { File, LoaderCircle, Upload } from "lucide-svelte";
  import { formatSize } from "../format.js";

  let {
    selectedFiles = [],
    uploading = false,
    progress = 0,
    status = "",
    onSelect,
    onClear,
    onUpload,
  } = $props();

  let fileInput = $state();
  let isDragging = $state(false);

  function selectFiles(event) {
    onSelect(Array.from(event.currentTarget.files || []));
    event.currentTarget.value = "";
  }

  function dropFiles(event) {
    event.preventDefault();
    isDragging = false;
    onSelect(Array.from(event.dataTransfer.files || []));
  }

  let totalSize = $derived(
    selectedFiles.reduce((sum, file) => sum + file.size, 0),
  );
</script>

<section class="upload-section" aria-labelledby="upload-title">
  <div class="section-heading">
    <div>
      <h2 id="upload-title">Unggah file</h2>
      <p>File dikirim langsung ke S3 menggunakan URL bertanda tangan.</p>
    </div>
    <span class="direct-label"><span></span>Direct upload</span>
  </div>

  <div
    class:dragging={isDragging}
    class="dropzone"
    role="button"
    tabindex="0"
    aria-label="Pilih file atau seret file ke sini"
    onkeydown={(event) => {
      if (event.key === "Enter" || event.key === " ") fileInput?.click();
    }}
    onclick={() => fileInput?.click()}
    ondragover={(event) => {
      event.preventDefault();
      isDragging = true;
    }}
    ondragleave={() => {
      isDragging = false;
    }}
    ondrop={dropFiles}
  >
    <input
      bind:this={fileInput}
      type="file"
      multiple
      hidden
      onclick={(event) => event.stopPropagation()}
      onchange={selectFiles}
    />
    <span class="drop-icon"><Upload size={21} /></span>
    <div class="drop-copy">
      <strong
        >{selectedFiles.length
          ? `${selectedFiles.length} file dipilih`
          : "Tarik file ke sini"}</strong
      >
      <span
        >{selectedFiles.length
          ? selectedFiles.map((file) => file.name).join(", ")
          : "atau klik untuk memilih dari perangkat"}</span
      >
    </div>
    <span class="browse-label">Pilih file</span>
  </div>

  {#if selectedFiles.length}
    <div class="upload-actions">
      <div class="selected-summary">
        <File size={15} />{selectedFiles.length} file · {formatSize(totalSize)}
      </div>
      <div class="button-row">
        <button
          class="button button-quiet"
          onclick={onClear}
          disabled={uploading}>Batal</button
        >
        <button
          class="button button-primary"
          onclick={onUpload}
          disabled={uploading}
        >
          {#if uploading}<LoaderCircle size={16} class="spin" />{status ||
              "Mengunggah"}{:else}<Upload size={16} />Unggah file{/if}
        </button>
      </div>
    </div>
  {/if}

  {#if uploading}
    <div class="progress-track" aria-label={`Progres upload ${progress}%`}>
      <span style={`width:${progress}%`}></span>
    </div>
  {/if}
</section>
