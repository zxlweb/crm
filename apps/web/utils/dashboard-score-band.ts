/** Score pill colors aligned with dashboard Preview (≥80 / 60–79 / <60) */

export type ScoreBand = 'excellent' | 'watch' | 'alert'

export function scoreBand(value: number): ScoreBand {
  if (value >= 80) return 'excellent'
  if (value >= 60) return 'watch'
  return 'alert'
}

export function scoreBandPillClass(band: ScoreBand): string {
  switch (band) {
    case 'excellent':
      return 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950/50 dark:text-emerald-300'
    case 'watch':
      return 'bg-amber-100 text-amber-800 dark:bg-amber-950/50 dark:text-amber-300'
    case 'alert':
      return 'bg-red-100 text-red-700 dark:bg-red-950/50 dark:text-red-300'
  }
}

export function scoreBandLegendClass(band: ScoreBand): string {
  switch (band) {
    case 'excellent':
      return 'bg-emerald-500'
    case 'watch':
      return 'bg-amber-400'
    case 'alert':
      return 'bg-red-400'
  }
}
