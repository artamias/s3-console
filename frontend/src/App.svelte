<script>
  import { onMount } from "svelte";
  import { LoaderCircle, RefreshCw } from "lucide-svelte";
  import { api } from "./lib/api.js";
  import { mediaKind } from "./lib/format.js";
  import { loadDashboard } from "./lib/storage-dashboard.js";
  import { uploadFiles } from "./lib/storage.js";
  import AppHeader from "./lib/components/AppHeader.svelte";
  import DeleteFileDialog from "./lib/components/DeleteFileDialog.svelte";
  import FilePreviewModal from "./lib/components/FilePreviewModal.svelte";
  import FileTable from "./lib/components/FileTable.svelte";
  import NotificationToast from "./lib/components/NotificationToast.svelte";
  import PresignedUrlModal from "./lib/components/PresignedUrlModal.svelte";
  import StorageSummary from "./lib/components/StorageSummary.svelte";
  import UploadPanel from "./lib/components/UploadPanel.svelte";

  let files = [];
  let info = null;
  let listError = "";
  let loading = true;
  let serverOnline = false;
  let selectedFiles = [];
  let uploading = false;
  let uploadProgress = 0;
  let uploadStatus = "";
  let search = "";
  let toast = null;
  let preview = null;
  let urlDialog = null;
  let deleteDialog = null;
  let copiedValue = "";
  let toastTimer;
  let copyTimer;

  onMount(() => {
    fetchAll();
    return () => {
      clearTimeout(toastTimer);
      clearTimeout(copyTimer);
    };
  });

  function notify(message, type = "success") {
    clearTimeout(toastTimer);
    toast = { message, type };
    toastTimer = setTimeout(() => {
      toast = null;
    }, 3600);
  }

  async function fetchAll() {
    loading = true;
    const dashboard = await loadDashboard();
    files = dashboard.files;
    info = dashboard.info;
    listError = dashboard.listError;
    serverOnline = dashboard.online;
    loading = false;
    if (dashboard.listError) notify(dashboard.listError, "error");
  }

  function selectFiles(nextFiles) {
    selectedFiles = nextFiles;
    uploadProgress = 0;
    uploadStatus = "";
  }

  function clearSelection() {
    selectedFiles = [];
    uploadProgress = 0;
    uploadStatus = "";
  }

  async function uploadSelected() {
    if (!selectedFiles.length || uploading) return;
    uploading = true;
    uploadProgress = 0;
    const { succeeded, failed } = await uploadFiles(
      selectedFiles,
      (progress) => {
        uploadProgress = progress;
      },
      (status) => {
        uploadStatus = status;
      },
    );
    uploading = false;
    uploadStatus = "";

    if (succeeded) {
      clearSelection();
      notify(
        `${succeeded} file berhasil diunggah ke S3${failed.length ? `, ${failed.length} gagal` : ""}`,
        failed.length ? "error" : "success",
      );
      await fetchAll();
    } else if (failed.length) {
      notify(`${failed[0].name}: ${failed[0].message}`, "error");
    }
  }

  async function previewFile(file) {
    try {
      const signed = await api.getPresignedUrl(file.name);
      preview = { ...signed, kind: mediaKind(file.name), error: false };
    } catch (error) {
      notify(`Gagal membuka file: ${error.message}`, "error");
    }
  }

  async function showUrl(file) {
    try {
      urlDialog = {
        ...(await api.getPresignedUrl(file.name)),
        name: file.name,
      };
    } catch (error) {
      notify(`Gagal membuat URL: ${error.message}`, "error");
    }
  }

  async function downloadFile(file) {
    try {
      const signed = await api.getDownloadUrl(file.name);
      const anchor = document.createElement("a");
      anchor.href = signed.url;
      anchor.download = file.name.split("/").pop();
      anchor.rel = "noreferrer";
      document.body.append(anchor);
      anchor.click();
      anchor.remove();
      notify(`Download ${file.name} dimulai`);
    } catch (error) {
      notify(`Gagal membuat URL download: ${error.message}`, "error");
    }
  }

  function askDelete(file) {
    deleteDialog = { name: file.name, loading: false };
  }

  async function deleteFile() {
    if (!deleteDialog || deleteDialog.loading) return;
    deleteDialog.loading = true;
    try {
      await api.deleteFile(deleteDialog.name);
      notify(`${deleteDialog.name} berhasil dihapus`);
      deleteDialog = null;
      await fetchAll();
    } catch (error) {
      deleteDialog.loading = false;
      notify(`Gagal menghapus file: ${error.message}`, "error");
    }
  }

  async function copyUrl(value) {
    try {
      await navigator.clipboard.writeText(value);
      copiedValue = value;
      clearTimeout(copyTimer);
      copyTimer = setTimeout(() => {
        copiedValue = "";
      }, 1800);
    } catch {
      notify("Clipboard tidak tersedia di browser ini", "error");
    }
  }

  function closeDialogs(event) {
    if (event.key === "Escape") {
      preview = null;
      urlDialog = null;
      deleteDialog = null;
    }
  }
</script>

<svelte:window onkeydown={closeDialogs} />

<div class="app-shell">
  <AppHeader online={serverOnline} />
  <main class="content">
    <div class="page-heading">
      <div>
        <p class="eyebrow">PENYIMPANAN</p>
        <h1>Object Storage</h1>
        <p class="subheading">Kelola file di bucket S3 dari satu dashboard.</p>
      </div>
      <button
        class="button button-secondary refresh-button"
        onclick={fetchAll}
        disabled={loading}
        title="Muat ulang daftar"
      >
        {#if loading}<LoaderCircle size={16} class="spin" />{:else}<RefreshCw
            size={16}
          />{/if}<span>Refresh</span>
      </button>
    </div>

    <StorageSummary {info} />
    <UploadPanel
      {selectedFiles}
      {uploading}
      progress={uploadProgress}
      status={uploadStatus}
      onSelect={selectFiles}
      onClear={clearSelection}
      onUpload={uploadSelected}
    />
    <FileTable
      {files}
      {info}
      {loading}
      online={serverOnline}
      error={listError}
      {search}
      onSearch={(value) => {
        search = value;
      }}
      onPreview={previewFile}
      onDownload={downloadFile}
      onUrl={showUrl}
      onDelete={askDelete}
    />
    <footer class="page-footer">
      Bucket {info?.bucket || "MinIO"} <span>·</span>Masa berlaku URL {info?.presign_expiry ||
        "—"}
    </footer>
  </main>
</div>

<NotificationToast
  {toast}
  onClose={() => {
    toast = null;
  }}
/>
<FilePreviewModal
  {preview}
  onClose={() => {
    preview = null;
  }}
  onCopy={copyUrl}
  onDownload={downloadFile}
  copied={copiedValue === preview?.url}
/>
<PresignedUrlModal
  {urlDialog}
  onClose={() => {
    urlDialog = null;
  }}
  onCopy={copyUrl}
  copied={copiedValue === urlDialog?.url}
/>
<DeleteFileDialog
  dialog={deleteDialog}
  onClose={() => {
    if (!deleteDialog?.loading) deleteDialog = null;
  }}
  onDelete={deleteFile}
/>
