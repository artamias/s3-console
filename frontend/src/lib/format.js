const numberFormats = new Map();

function formatNumber(value, maximumFractionDigits) {
  if (!numberFormats.has(maximumFractionDigits)) {
    numberFormats.set(maximumFractionDigits, new Intl.NumberFormat('id-ID', { maximumFractionDigits }));
  }
  return numberFormats.get(maximumFractionDigits).format(value);
}

// [satuan, jumlah desimal maksimum] per tingkat 1024.
const SIZE_UNITS = [['B', 0], ['KB', 0], ['MB', 1], ['GB', 2], ['TB', 2]];

export function formatSize(bytes = 0) {
  if (!bytes) return '0 B';
  const index = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), SIZE_UNITS.length - 1);
  const [unit, digits] = SIZE_UNITS[index];
  return `${formatNumber(bytes / 1024 ** index, digits)} ${unit}`;
}

function startOfDay(date) {
  return new Date(date.getFullYear(), date.getMonth(), date.getDate());
}

export function formatDate(value, now = new Date()) {
  if (!value) return '—';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return '—';

  const time = date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' });
  const dayDiff = Math.round((startOfDay(now) - startOfDay(date)) / 86_400_000);
  if (dayDiff === 0) return `Hari ini, ${time}`;
  if (dayDiff === 1) return `Kemarin, ${time}`;

  const day = date.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' });
  return `${day}, ${time}`;
}

const IMAGE_EXTENSIONS = ['png', 'jpg', 'jpeg', 'gif', 'webp', 'svg', 'ico', 'bmp', 'avif'];
const VIDEO_EXTENSIONS = ['mp4', 'webm', 'ogg', 'mov', 'avi', 'mkv', 'm4v'];
const ARCHIVE_EXTENSIONS = ['zip', 'tar', 'gz', 'tgz', 'bz2', 'xz', 'rar', '7z'];
// Label khusus bila berbeda dari ekstensi (design: ".gz" tampil sebagai "Arsip GZIP").
const EXTENSION_LABELS = { gz: 'GZIP', bz2: 'BZIP2' };

/** Placeholder "folder" di MinIO adalah object berakhiran "/". */
export function isFolder(name) {
  return name.endsWith('/');
}

function extensionOf(name) {
  const base = name.split('/').pop();
  const dot = base.lastIndexOf('.');
  return dot > 0 ? base.slice(dot + 1).toLowerCase() : '';
}

export function mediaKind(name) {
  if (isFolder(name)) return 'folder';
  const extension = extensionOf(name);
  if (IMAGE_EXTENSIONS.includes(extension)) return 'image';
  if (VIDEO_EXTENSIONS.includes(extension)) return 'video';
  if (extension === 'pdf') return 'pdf';
  return 'file';
}

/** Label kolom "Type", mengikuti design: "PDF", "Gambar PNG", "Arsip GZIP", "Folder". */
export function fileKind(name) {
  if (isFolder(name)) return 'Folder';
  const extension = extensionOf(name);
  if (!extension) return 'File';
  const label = EXTENSION_LABELS[extension] ?? extension.toUpperCase();
  if (IMAGE_EXTENSIONS.includes(extension)) return `Gambar ${label}`;
  if (VIDEO_EXTENSIONS.includes(extension)) return `Video ${label}`;
  if (ARCHIVE_EXTENSIONS.includes(extension)) return `Arsip ${label}`;
  return label;
}
