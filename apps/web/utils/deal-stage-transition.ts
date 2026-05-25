import type { DealStage } from '~/types/deal'

/** Mirrors backend/internal/pkg/crm/deal_stage.go */
const STAGE_ORDER: Record<DealStage, number> = {
  qualification: 0,
  proposal: 1,
  negotiation: 2,
  won: 3,
  lost: 3,
}

export function isDealTerminal(stage: DealStage): boolean {
  return stage === 'won' || stage === 'lost'
}

export function canTransitionDealStage(from: DealStage, to: DealStage): boolean {
  if (from === to) return true
  if (isDealTerminal(from)) return false
  if (to === 'won' || to === 'lost') return from === 'negotiation'
  const fromIdx = STAGE_ORDER[from]
  const toIdx = STAGE_ORDER[to]
  if (toIdx < fromIdx) return true
  return toIdx === fromIdx + 1
}

export function allowedNextStages(from: DealStage): DealStage[] {
  const all: DealStage[] = ['qualification', 'proposal', 'negotiation', 'won', 'lost']
  return all.filter((to) => to !== from && canTransitionDealStage(from, to))
}

export const DEAL_DRAG_MIME = 'application/x-crm-deal-id'

export type DealDragPayload = {
  dealId: string
  fromStage: DealStage
}

export function writeDealDragData(dataTransfer: DataTransfer, payload: DealDragPayload) {
  dataTransfer.effectAllowed = 'move'
  dataTransfer.setData(DEAL_DRAG_MIME, JSON.stringify(payload))
  dataTransfer.setData('text/plain', payload.dealId)
}

export function readDealDragData(dataTransfer: DataTransfer): DealDragPayload | null {
  const raw = dataTransfer.getData(DEAL_DRAG_MIME)
  if (!raw) return null
  try {
    const parsed = JSON.parse(raw) as DealDragPayload
    if (parsed.dealId && parsed.fromStage) return parsed
  } catch {
    /* ignore */
  }
  return null
}
