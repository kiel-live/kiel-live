export function get24hTime(date: string) {
  const d = new Date(date);
  if (Number.isNaN(d.getTime())) {
    return '';
  }

  const pad =(n: number) => n.toString().padStart(2, '0');
  return `${pad(d.getHours())}:${pad(d.getMinutes())}`;
}
