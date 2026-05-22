/** 24×24 stroke 图标，用于活动类型 / 情绪 / 线索标签等徽标 */

export type TagIconName =
  | 'email'
  | 'call'
  | 'meeting'
  | 'wechat'
  | 'note'
  | 'tag'
  | 'positive'
  | 'neutral'
  | 'hesitant'
  | 'negative'

export type TagIconPath = {
  tag: 'path' | 'circle' | 'rect' | 'line' | 'polyline'
  d?: string
  cx?: number
  cy?: number
  r?: number
  x?: number
  y?: number
  width?: number
  height?: number
  rx?: number
  x1?: number
  y1?: number
  x2?: number
  y2?: number
  points?: string
  /** 默认描边；小圆点等可设为 currentColor */
  fill?: string
}

export const TAG_ICONS: Record<TagIconName, TagIconPath[]> = {
  email: [
    {
      tag: 'path',
      d: 'M4 6h16a2 2 0 012 2v8a2 2 0 01-2 2H4a2 2 0 01-2-2V8a2 2 0 012-2z',
    },
    { tag: 'path', d: 'M4 8l8 5.5L20 8' },
  ],
  call: [
    {
      tag: 'path',
      d: 'M6.5 4.5c.5 2.2 1.8 4.5 3.5 6.2s4 3 6.2 3.5l1.3-1.3a1 1 0 011-.2c1.1.6 2.3 1 3.5 1.2a1 1 0 011 1v3.2a1 1 0 01-1 1C10.2 18.5 5.5 13.8 4 7a1 1 0 011-1h3.2a1 1 0 011 1c.2 1.2.6 2.4 1.2 3.5a1 1 0 01-.2 1L6.5 4.5z',
    },
  ],
  meeting: [
    { tag: 'rect', x: 3, y: 5, width: 18, height: 14, rx: 2 },
    { tag: 'path', d: 'M8 3v4M16 3v4M3 10h18' },
    { tag: 'circle', cx: 9, cy: 14, r: 1.25 },
    { tag: 'circle', cx: 15, cy: 14, r: 1.25 },
  ],
  wechat: [
    {
      tag: 'path',
      d: 'M7 18l-2.5 2.5V17a7 7 0 1112.2-4.9A5.5 5.5 0 017 18z',
    },
    { tag: 'circle', cx: 9.5, cy: 11.5, r: 0.75, fill: 'currentColor' },
    { tag: 'circle', cx: 13.5, cy: 11.5, r: 0.75, fill: 'currentColor' },
  ],
  note: [
    { tag: 'path', d: 'M8 4h8a2 2 0 012 2v14l-4-2-4 2V6a2 2 0 012-2z' },
    { tag: 'path', d: 'M10 9h6M10 12h4' },
  ],
  tag: [
    {
      tag: 'path',
      d: 'M4 12V6a2 2 0 012-2h6l8.5 8.5a1.5 1.5 0 010 2.12l-4.88 4.88a1.5 1.5 0 01-2.12 0L4 12z',
    },
    { tag: 'circle', cx: 9.5, cy: 7.5, r: 1 },
  ],
  /** 情绪表情：圆脸 + 眼 + 嘴 */
  positive: [
    { tag: 'circle', cx: 12, cy: 12, r: 8 },
    { tag: 'circle', cx: 9, cy: 10, r: 1, fill: 'currentColor' },
    { tag: 'circle', cx: 15, cy: 10, r: 1, fill: 'currentColor' },
    { tag: 'path', d: 'M8.5 13.5c1.2 1.8 2.8 2.5 3.5 2.5s2.3-.7 3.5-2.5' },
  ],
  neutral: [
    { tag: 'circle', cx: 12, cy: 12, r: 8 },
    { tag: 'circle', cx: 9, cy: 10, r: 1, fill: 'currentColor' },
    { tag: 'circle', cx: 15, cy: 10, r: 1, fill: 'currentColor' },
    { tag: 'line', x1: 8.5, y1: 14.5, x2: 15.5, y2: 14.5 },
  ],
  hesitant: [
    { tag: 'circle', cx: 12, cy: 12, r: 8 },
    { tag: 'circle', cx: 9, cy: 10, r: 1, fill: 'currentColor' },
    { tag: 'circle', cx: 15, cy: 10, r: 1, fill: 'currentColor' },
    { tag: 'path', d: 'M8.5 14.5c1-.5 2.2-.8 3.5-.8s2.5.3 3.5.8' },
    { tag: 'path', d: 'M9 9.5h1M14 9.5h1' },
  ],
  negative: [
    { tag: 'circle', cx: 12, cy: 12, r: 8 },
    { tag: 'circle', cx: 9, cy: 10, r: 1, fill: 'currentColor' },
    { tag: 'circle', cx: 15, cy: 10, r: 1, fill: 'currentColor' },
    { tag: 'path', d: 'M8.5 16c1.2-1.5 2.8-2 3.5-2s2.3.5 3.5 2' },
  ],
}

/** CRM 活动 event_type → 图标 */
export const ACTIVITY_TYPE_ICON: Record<string, TagIconName> = {
  email: 'email',
  call: 'call',
  meeting: 'meeting',
  wechat: 'wechat',
  note: 'note',
}

export const SENTIMENT_ICON: Record<string, TagIconName> = {
  positive: 'positive',
  neutral: 'neutral',
  hesitant: 'hesitant',
  negative: 'negative',
  unknown: 'neutral',
}

export function resolveActivityTypeIcon(eventType: string): TagIconName {
  return ACTIVITY_TYPE_ICON[eventType] ?? 'note'
}

export function resolveSentimentIcon(sentiment: string): TagIconName {
  return SENTIMENT_ICON[sentiment] ?? 'neutral'
}

/** ECharts path:// 符号，与 UiTagIcon 情绪表情一致 */
export const SENTIMENT_ECHART_SYMBOL: Record<
  'positive' | 'neutral' | 'hesitant' | 'negative',
  string
> = {
  positive:
    'path://M12,4a8,8,0,1,1,0,16a8,8,0,1,1,0,-16M9,10a1,1,0,1,0,0,-0.01M15,10a1,1,0,1,0,0,-0.01M8.5,13.5c1.2,1.8,2.8,2.5,3.5,2.5s2.3,-0.7,3.5,-2.5',
  neutral:
    'path://M12,4a8,8,0,1,1,0,16a8,8,0,1,1,0,-16M9,10a1,1,0,1,0,0,-0.01M15,10a1,1,0,1,0,0,-0.01M8.5,14.5h7',
  hesitant:
    'path://M12,4a8,8,0,1,1,0,16a8,8,0,1,1,0,-16M9,10a1,1,0,1,0,0,-0.01M15,10a1,1,0,1,0,0,-0.01M9,9.5h1M14,9.5h1M8.5,14.5c1,-0.5,2.2,-0.8,3.5,-0.8s2.5,0.3,3.5,0.8',
  negative:
    'path://M12,4a8,8,0,1,1,0,16a8,8,0,1,1,0,-16M9,10a1,1,0,1,0,0,-0.01M15,10a1,1,0,1,0,0,-0.01M8.5,16c1.2,-1.5,2.8,-2,3.5,-2s2.3,0.5,3.5,2',
}

export function sentimentEchartSymbol(sentiment: string): string {
  const key = resolveSentimentIcon(sentiment) as keyof typeof SENTIMENT_ECHART_SYMBOL
  return SENTIMENT_ECHART_SYMBOL[key] ?? SENTIMENT_ECHART_SYMBOL.neutral
}
