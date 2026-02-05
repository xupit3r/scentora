import { Decimal } from '@prisma/client/runtime/library';

export function toNumber(value: Decimal | null | undefined): number | undefined {
  if (value == null) return undefined;
  return Number(value);
}

export function toNumberRequired(value: Decimal): number {
  return Number(value);
}

/**
 * Transform a Prisma record for API output:
 * - Renames `id` to `_id`
 * - Converts Decimal fields to numbers
 * - Strips undefined/null optional fields (omitempty behavior)
 */
export function toApiRecord(record: Record<string, unknown>): Record<string, unknown> {
  const result: Record<string, unknown> = {};

  for (const [key, value] of Object.entries(record)) {
    if (value === undefined) continue;

    if (key === 'id') {
      result['_id'] = value;
      continue;
    }

    if (value instanceof Decimal) {
      result[key] = Number(value);
      continue;
    }

    if (value === null) {
      // skip null optional fields to match Go's omitempty
      continue;
    }

    result[key] = value;
  }

  return result;
}
