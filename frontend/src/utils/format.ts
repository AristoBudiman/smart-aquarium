export function formatWib(value: string | null | undefined): string {
  if (!value) return 'Belum ada data'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return 'Belum ada data'
  return new Intl.DateTimeFormat('id-ID', {
    timeZone: 'Asia/Jakarta',
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(date)
}

export function formatNumber(value: number | null | undefined, digits = 1): string {
  if (value === null || value === undefined || !Number.isFinite(value)) return '—'
  return new Intl.NumberFormat('id-ID', { maximumFractionDigits: digits }).format(value)
}
