import { describe, expect, it } from 'vitest'
import {
  allowedNextStages,
  canTransitionDealStage,
  isDealTerminal,
} from '~/utils/deal-stage-transition'

describe('deal-stage-transition', () => {
  it('allows forward one-step moves', () => {
    expect(canTransitionDealStage('qualification', 'proposal')).toBe(true)
    expect(canTransitionDealStage('proposal', 'negotiation')).toBe(true)
  })

  it('allows backward moves within open pipeline', () => {
    expect(canTransitionDealStage('negotiation', 'proposal')).toBe(true)
    expect(canTransitionDealStage('proposal', 'qualification')).toBe(true)
  })

  it('allows won/lost only from negotiation', () => {
    expect(canTransitionDealStage('negotiation', 'won')).toBe(true)
    expect(canTransitionDealStage('negotiation', 'lost')).toBe(true)
    expect(canTransitionDealStage('qualification', 'won')).toBe(false)
    expect(canTransitionDealStage('qualification', 'lost')).toBe(false)
  })

  it('blocks terminal stage changes and illegal skips', () => {
    expect(isDealTerminal('won')).toBe(true)
    expect(canTransitionDealStage('won', 'proposal')).toBe(false)
    expect(canTransitionDealStage('qualification', 'negotiation')).toBe(false)
  })

  it('lists allowed next stages for buttons and drag targets', () => {
    expect(allowedNextStages('qualification')).toEqual(['proposal'])
    expect(allowedNextStages('negotiation').sort()).toEqual(
      ['proposal', 'qualification', 'won', 'lost'].sort(),
    )
  })
})
