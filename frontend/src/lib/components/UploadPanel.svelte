<script>
  import Icon from "./Icon.svelte";

  let { uploading = false, progress = 0, status = "", notice = null, onUpload } =
    $props();

  let fileInput = $state();
  let dragging = $state(false);

  function pick(files) {
    const selected = Array.from(files || []);
    if (selected.length) onUpload(selected);
  }

  function handleDrop(event) {
    event.preventDefault();
    dragging = false;
    if (!uploading) pick(event.dataTransfer.files);
  }

  function handleDragLeave(event) {
    if (!event.currentTarget.contains(event.relatedTarget)) dragging = false;
  }
</script>

<section class="mt-8">
  <div
    role="group"
    aria-label="Area unggah file"
    ondragenter={(event) => {
      event.preventDefault();
      dragging = true;
    }}
    ondragover={(event) => event.preventDefault()}
    ondragleave={handleDragLeave}
    ondrop={handleDrop}
    class="flex min-h-28 flex-col items-center justify-center rounded-xl border border-dashed px-6 py-5 text-center transition sm:flex-row sm:justify-between sm:text-left {dragging
      ? 'border-[#16745A] bg-[#F1F8F5]'
      : 'border-[#BFC9C3] bg-white'}"
  >
    <div class="flex min-w-0 flex-col items-center gap-3 sm:flex-row">
      <div
        class="flex size-10 shrink-0 items-center justify-center rounded-lg {dragging
          ? 'bg-[#DDEFE8] text-[#16745A]'
          : 'bg-[#F1F3F2] text-[#59665F]'}"
      >
        <Icon name="upload" class="size-5" />
      </div>
      <div class="min-w-0">
        <div class="text-[15px] font-medium">
          {#if uploading}
            {status || "Mengunggah..."}
          {:else if dragging}
            Lepaskan file untuk mengunggah
          {:else}
            Upload file
          {/if}
        </div>
        <div class="mt-1 text-[13px] text-[#6B7770]">
          {#if uploading}
            {progress}% selesai
          {:else}
            Drag and drop files here, or select them from your device
          {/if}
        </div>
      </div>
    </div>
    <input
      bind:this={fileInput}
      type="file"
      multiple
      class="sr-only"
      tabindex="-1"
      onchange={(event) => {
        pick(event.currentTarget.files);
        event.currentTarget.value = "";
      }}
    />
    <button
      onclick={() => fileInput?.click()}
      disabled={uploading}
      class="mt-4 inline-flex h-9 shrink-0 items-center justify-center rounded-lg bg-[#16745A] px-4 text-[14px] font-medium text-white transition hover:bg-[#125F49] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#16745A] disabled:opacity-60 sm:mt-0"
    >
      Pilih file
    </button>
  </div>

  {#if uploading}
    <div
      class="mt-3 h-1.5 overflow-hidden rounded-full bg-[#E8ECE9]"
      role="progressbar"
      aria-label="Progres upload"
      aria-valuemin="0"
      aria-valuemax="100"
      aria-valuenow={progress}
    >
      <div
        class="h-full rounded-full bg-[#16745A] transition-[width]"
        style:width="{progress}%"
      ></div>
    </div>
  {/if}

  {#if notice}
    <div
      role={notice.type === "error" ? "alert" : "status"}
      class="mt-3 flex items-center gap-2 text-[13px] font-medium {notice.type ===
      'error'
        ? 'text-[#A13A36]'
        : 'text-[#116149]'}"
    >
      <span
        class="flex size-5 shrink-0 items-center justify-center rounded-full {notice.type ===
        'error'
          ? 'bg-[#FCF2F1]'
          : 'bg-[#E1F1EB]'}"
      >
        <Icon
          name={notice.type === "error" ? "alert" : "check"}
          class="size-3.5"
        />
      </span>
      {notice.message}
    </div>
  {/if}
</section>
