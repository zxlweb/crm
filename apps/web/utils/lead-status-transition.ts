import type { LeadStatus } from '~/types/lead'

/** Mirrors backend/internal/pkg/crm/lead_status_transition.go */
const NEXT_STATUS: Partial<Record<LeadStatus, LeadStatus[]>> = {
  new: ['contacted'],
  contacted: ['qualified', 'unqualified'],
  qualified: ['unqualified'],
  unqualified: [],
  converted: [],
}

export function canTransitionLeadStatus(from: LeadStatus, to: LeadStatus): boolean {
  if (from === to) return true
  if (to === 'converted') return false
  return NEXT_STATUS[from]?.includes(to) ?? false
}

export function allowedLeadStatusTargets(from: LeadStatus): LeadStatus[] {
  const next = NEXT_STATUS[from] ?? []
  return [from, ...next.filter((s) => s !== from)]
}

export function canConvertLead(status: LeadStatus): boolean {
  return status === 'qualified'
}
