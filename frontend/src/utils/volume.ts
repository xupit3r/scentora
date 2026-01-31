/**
 * Utility functions for volume conversions and formatting
 * Used throughout the application for consistent volume display
 */

/**
 * Standard conversion rate: 1 ml = 20 drops (approximate)
 * This is a common conversion used in perfumery
 */
export const DROPS_PER_ML = 20;

/**
 * Convert milliliters to drops
 */
export function mlToDrops(ml: number): number {
  return Math.round(ml * DROPS_PER_ML);
}

/**
 * Convert drops to milliliters
 */
export function dropsToMl(drops: number): number {
  return Math.round((drops / DROPS_PER_ML) * 100) / 100;
}

/**
 * Format volume for display with appropriate unit
 */
export function formatVolume(ml: number, showDrops = false): string {
  const mlFormatted = ml.toFixed(ml < 1 ? 2 : 1);
  
  if (showDrops) {
    const drops = mlToDrops(ml);
    return `${mlFormatted} ml (${drops} drops)`;
  }
  
  return `${mlFormatted} ml`;
}

/**
 * Get stock level based on volume
 */
export function getStockLevel(volumeMl: number, threshold = 10): 'ok' | 'low' | 'critical' {
  if (volumeMl === 0) return 'critical';
  if (volumeMl < threshold / 2) return 'critical';
  if (volumeMl < threshold) return 'low';
  return 'ok';
}

/**
 * Get stock warning message
 */
export function getStockWarning(volumeMl: number, threshold = 10): string | null {
  const level = getStockLevel(volumeMl, threshold);
  
  if (level === 'critical') {
    return volumeMl === 0 ? 'Out of stock' : 'Critical - Reorder now';
  }
  
  if (level === 'low') {
    return 'Low stock';
  }
  
  return null;
}

/**
 * Format volume range for filters
 */
export function formatVolumeRange(minMl?: number, maxMl?: number): string {
  if (minMl !== undefined && maxMl !== undefined) {
    return `${minMl}-${maxMl} ml`;
  }
  if (minMl !== undefined) {
    return `≥${minMl} ml`;
  }
  if (maxMl !== undefined) {
    return `≤${maxMl} ml`;
  }
  return 'All volumes';
}
