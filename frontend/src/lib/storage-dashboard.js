import { api } from './api.js';

/** Ambil data independen secara paralel tanpa menggagalkan seluruh dashboard. */
export async function loadDashboard() {
  const [filesResult, infoResult] = await Promise.allSettled([
    api.getFiles(),
    api.getInfo(),
  ]);

  return {
    files: filesResult.status === 'fulfilled' ? filesResult.value.data || [] : [],
    info: infoResult.status === 'fulfilled' ? infoResult.value : null,
    listError: filesResult.status === 'rejected'
      ? filesResult.reason.message || 'Gagal memuat daftar file'
      : '',
    online: filesResult.status === 'fulfilled' || infoResult.status === 'fulfilled',
  };
}
