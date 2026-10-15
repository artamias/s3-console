export function formatSize(bytes = 0) {
  if (!bytes) return '0 B';
  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  const index = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1);
  return `${(bytes / 1024 ** index).toFixed(index === 0 ? 0 : 1)} ${units[index]}`;
}

export function formatDate(value) {
  if (!value) return '—';
  return new Date(value).toLocaleString('id-ID', {
    day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit',
  });
}

export function mediaKind(name) {
  if (/\.(png|jpe?g|gif|webp|svg|ico|bmp|avif)$/i.test(name)) return 'image';
  if (/\.(mp4|webm|ogg|mov|avi|mkv|m4v)$/i.test(name)) return 'video';
  return 'file';
}

export function fileKind(name) {
  const extension = name.split('.').pop()?.toLowerCase();
  if (['png', 'jpg', 'jpeg', 'gif', 'webp', 'svg', 'avif'].includes(extension)) return 'Image';
  if (['mp4', 'mov', 'avi', 'mkv', 'webm'].includes(extension)) return 'Video';
  if (extension === 'pdf') return 'PDF';
  if (['zip', 'tar', 'gz', 'rar', '7z'].includes(extension)) return 'Archive';
  return extension ? extension.toUpperCase() : 'File';
}
