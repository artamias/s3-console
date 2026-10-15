# Frontend Svelte

Frontend ini memakai Svelte 5 dan Vite. Struktur yang disarankan:

```text
frontend/
├── index.html
├── package.json
├── vite.config.js
├── svelte.config.js
├── .env.example
└── src/
    ├── App.svelte
    ├── main.js
    ├── style.css
    └── lib/
        ├── api.js
        └── format.js
```

## Menjalankan

```sh
cp .env.example .env
npm install
npm run dev
```

Vite berjalan di `http://localhost:5173` dan meneruskan `/api` serta `/health` ke backend `http://localhost:8080`. Atur `VITE_BACKEND_ORIGIN` bila backend memakai host atau port lain. `VITE_API_BASE_URL` dapat diatur ke URL API absolut untuk deployment tanpa proxy.

## Konfigurasi production

Simpan konfigurasi backend di `backend/.env` berdasarkan `backend/.env.example`. `APP_HOST` dan `APP_PORT` menentukan alamat bind Go. `MINIO_ENDPOINT` dipakai backend dan URL presigned. Saat production, ubah nilainya ke domain MinIO yang dapat dijangkau browser, lalu set `MINIO_USE_SSL=true` bila menggunakan HTTPS.

Simpan konfigurasi frontend di `frontend/.env.production` berdasarkan `frontend/.env.example`, lalu jalankan `npm run build`. Gunakan `VITE_API_BASE_URL=/api/v1` bila frontend dan API berbagi domain di belakang reverse proxy, atau URL lengkap seperti `https://api.example.com/api/v1` bila domainnya terpisah. Nilai `VITE_*` ditanam ke bundle saat build; `VITE_BACKEND_ORIGIN` hanya digunakan proxy Vite saat development.

Origin frontend juga harus diizinkan oleh CORS MinIO untuk request langsung `GET` dan `PUT`. `CORS_ALLOW_ORIGINS` di backend hanya mengatur CORS API Go, bukan CORS MinIO.

Frontend memakai `GET /api/v1/files`, `GET /api/v1/info`, `POST /api/v1/files/presigned-upload`, `GET /api/v1/files/:name/presigned`, dan `DELETE /api/v1/files/:name`. Tambahkan `?download=true` pada endpoint presigned GET untuk menerima URL MinIO dengan `Content-Disposition: attachment`; byte file tetap dikirim langsung oleh MinIO. Berkas upload juga dikirim langsung dengan `PUT` ke MinIO. Karena kedua request menuju origin MinIO, konfigurasi CORS MinIO harus mengizinkan origin frontend (contoh development: `http://localhost:5173`) dan method `GET` serta `PUT`.
