const API_BASE = (import.meta.env.VITE_API_BASE_URL || '/api/v1').replace(/\/$/, '');

async function request(path, options = {}) {
  let response;
  try {
    response = await fetch(`${API_BASE}${path}`, options);
  } catch {
    throw new Error('Backend tidak dapat dijangkau. Pastikan API Go aktif di port 8080.');
  }

  if (!response.ok) {
    let message = `Request gagal (${response.status})`;
    try {
      const body = await response.json();
      message = body.error || message;
    } catch {
      // Keep the HTTP status message when the response is not JSON.
    }
    throw new Error(message);
  }

  if (response.status === 204) return null;
  return response.json();
}

export const api = {
  getFiles: () => request('/files'),
  getInfo: () => request('/info'),
  getPresignedUrl: (name) => request(`/files/${encodeURIComponent(name)}/presigned`),
  getDownloadUrl: (name) => request(`/files/${encodeURIComponent(name)}/presigned?download=true`),
  deleteFile: (name) => request(`/files/${encodeURIComponent(name)}`, { method: 'DELETE' }),
  getUploadUrl: (file) => request('/files/presigned-upload', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      filename: file.name,
      content_type: file.type || 'application/octet-stream',
    }),
  }),
};

export function putFile(uploadUrl, file, onProgress) {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest();
    xhr.open('PUT', uploadUrl);
    xhr.setRequestHeader('Content-Type', file.type || 'application/octet-stream');
    xhr.upload.onprogress = (event) => {
      if (event.lengthComputable) onProgress(event.loaded / event.total);
    };
    xhr.onload = () => {
      if (xhr.status >= 200 && xhr.status < 300) resolve();
      else reject(new Error(`MinIO menolak upload (${xhr.status}). Periksa CORS pada endpoint MinIO.`));
    };
    xhr.onerror = () => reject(new Error('Upload ke MinIO gagal. Periksa koneksi dan konfigurasi CORS MinIO.'));
    xhr.send(file);
  });
}
