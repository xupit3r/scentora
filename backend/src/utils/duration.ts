/**
 * Parse duration strings like "15m", "7d", "24h" into milliseconds.
 */
export function parseDuration(s: string): number {
  if (s.length < 2) throw new Error(`Invalid duration format: ${s}`);

  const value = parseInt(s.slice(0, -1), 10);
  const unit = s.slice(-1);

  if (isNaN(value)) throw new Error(`Invalid duration value: ${s}`);

  switch (unit) {
    case 's':
      return value * 1000;
    case 'm':
      return value * 60 * 1000;
    case 'h':
      return value * 60 * 60 * 1000;
    case 'd':
      return value * 24 * 60 * 60 * 1000;
    default:
      throw new Error(`Unknown duration unit: ${unit}`);
  }
}
