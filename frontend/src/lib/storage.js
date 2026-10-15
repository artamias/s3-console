import { api, putFile } from './api.js';

/**
 * Mengelola alur upload langsung ke MinIO. UI tetap mengatur state dan pesan
 * agar service ini dapat dipakai kembali tanpa bergantung pada Svelte.
 */
export async function uploadFiles(files, onProgress, onStatus) {
  let succeeded = 0;
  const failed = [];

  for (const [index, file] of files.entries()) {
    try {
      onStatus(`Meminta URL · ${file.name}`);
      const signed = await api.getUploadUrl(file);
      onStatus(`Mengunggah · ${file.name}`);
      await putFile(signed.upload_url, file, (fraction) => {
        onProgress(Math.round(((index + fraction) / files.length) * 100));
      });
      succeeded += 1;
    } catch (error) {
      failed.push({ name: file.name, message: error.message });
    }
  }

  return { succeeded, failed };
}
