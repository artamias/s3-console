<script>
  import { onMount } from "svelte";
  import { api } from "./lib/api.js";
  import { mediaKind } from "./lib/format.js";
  import { loadDashboard } from "./lib/storage-dashboard.js";
  import { uploadFiles } from "./lib/storage.js";
  import AppHeader from "./lib/components/AppHeader.svelte";
  import DeleteFileDialog from "./lib/components/DeleteFileDialog.svelte";
  import FilePreviewModal from "./lib/components/FilePreviewModal.svelte";
  import FileTable from "./lib/components/FileTable.svelte";
  import Icon from "./lib/components/Icon.svelte";
  import PresignedUrlModal from "./lib/components/PresignedUrlModal.svelte";
  import StorageSummary from "./lib/components/StorageSummary.svelte";
  import UploadPanel from "./lib/components/UploadPanel.svelte";

  let files = [];
  let info = null;
  let listError = "";
  let loading = true;
  let serverOnline = false;
  let uploading = false;
  let uploadProgress = 0;
  let uploadStatus = "";
  let search = "";
  let notice = null;
  let preview = null;
  let urlDialog = null;
  let deleteDialog = null;
  let copiedValue = "";
  let noticeTimer;
  let copyTimer;

  onMount(() => {
    fetchAll();
    return () => {
      clearTimeout(noticeTimer);
      clearTimeout(copyTimer);
    };
  });

  function notify(message, type = "success") {
    clearTimeout(noticeTimer);
    notice = { message, type };
    noticeTimer = setTimeout(
      () => {
        notice = null;
      },
      type === "error" ? 6000 : 4000,
    );
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

  async function refresh() {
    await fetchAll();
    if (!listError) notify("Daftar objek telah diperbarui.");
  }

  async function uploadSelected(selected) {
    if (!selected.length || uploading) return;
    uploading = true;
    uploadProgress = 0;
    const { succeeded, failed } = await uploadFiles(
      selected,
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
      notify(
        `${succeeded} file berhasil diunggah${failed.length ? `, ${failed.length} gagal` : ""}.`,
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
      preview = { ...signed, kind: mediaKind(file.name), file };
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
      notify(`Unduhan ${file.name} dimulai.`);
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
      notify(`${deleteDialog.name} dihapus dari bucket.`);
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
      if (!deleteDialog?.loading) deleteDialog = null;
    }
  }
</script>

<svelte:window onkeydown={closeDialogs} />

<div class="flex min-h-screen flex-col bg-[#F7F8F7] text-[#17211D]">
  <AppHeader online={serverOnline} />

  <main
    class="mx-auto w-full max-w-[1440px] flex-1 px-5 py-9 sm:px-8 sm:py-12 lg:px-12"
  >
    <div class="flex flex-col justify-between gap-5 sm:flex-row sm:items-start">
      <div>
        <nav
          aria-label="Breadcrumb"
          class="mb-3 flex items-center gap-1.5 text-[13px] text-[#6B7770]"
        >
          <span>bucket</span>
          <Icon name="chevron" class="size-3.5" />
          <span class="font-medium text-[#425149]">{info?.bucket || "—"}</span>
        </nav>
        <h1 class="text-[30px] font-semibold leading-tight tracking-[-0.035em]">
          Object Storage
        </h1>
        <p class="mt-2 max-w-2xl text-[15px] leading-6 text-[#59665F]">
          Manage files and assets in your S3 bucket from a single location
        </p>
      </div>
      <button
        onclick={refresh}
        disabled={loading}
        class="inline-flex h-10 shrink-0 items-center justify-center gap-2 self-start rounded-lg border border-[#CBD4CE] bg-white px-4 text-[14px] font-medium text-[#27362F] transition hover:border-[#AEBBB4] hover:bg-[#FAFBFA] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#16745A] disabled:cursor-wait disabled:opacity-60"
      >
        <Icon name="refresh" class="size-4 {loading ? 'animate-spin' : ''}" />
        {loading ? "Memperbarui..." : "Perbarui"}
      </button>
    </div>

    <StorageSummary {info} />
    <UploadPanel
      {uploading}
      progress={uploadProgress}
      status={uploadStatus}
      {notice}
      onUpload={uploadSelected}
    />
    <FileTable
      {files}
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
  </main>
</div>

<FilePreviewModal
  {preview}
  onClose={() => {
    preview = null;
  }}
  onDownload={downloadFile}
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
