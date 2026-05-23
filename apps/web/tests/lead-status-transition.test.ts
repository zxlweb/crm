import { describe, expect, it } from 'vitest'
import {
  allowedLeadStatusTargets,
  canConvertLead,
  canTransitionLeadStatus,
} from '~/utils/lead-status-transition'

describe('lead-status-transition', () => {
  it('allows valid forward transitions', () => {
    expect(canTransitionLeadStatus('new', 'contacted')).toBe(true)
    expect(canTransitionLeadStatus('contacted', 'qualified')).toBe(true)
    expect(canTransitionLeadStatus('qualified', 'unqualified')).toBe(true)
  })

  it('blocks converted and illegal jumps', () => {
    expect(canTransitionLeadStatus('qualified', 'converted')).toBe(false)
    expect(canTransitionLeadStatus('unqualified', 'contacted')).toBe(false)
    expect(canTransitionLeadStatus('new', 'qualified')).toBe(false)
  })

  it('lists allowed targets including current', () => {
    expect(allowedLeadStatusTargets('contacted')).toEqual(['contacted', 'qualified', 'unqualified'])
    expect(allowedLeadStatusTargets('unqualified')).toEqual(['unqualified'])
  })

  it('convert only from qualified', () => {
    expect(canConvertLead('qualified')).toBe(true)
    expect(canConvertLead('new')).toBe(false)
  })
})
