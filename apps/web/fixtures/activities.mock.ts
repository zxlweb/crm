import { DEMO_LEAD_ID } from '~/fixtures/leads.mock'

export type ActivityTimelineItem = {
  id: string
  at: string
  event_type: string
  direction?: string
  label: string
  sentiment?: string
  body?: string
}

const now = Date.now()
const at = (daysAgo: number, hour = 10) =>
  new Date(now - daysAgo * 86400000 + hour * 3600000).toISOString()

const DEMO_ACTIVITIES: Record<string, ActivityTimelineItem[]> = {
  [DEMO_LEAD_ID]: [
    { id: 'f1000000-0000-4000-8000-000000000005', at: at(3, 9), event_type: 'call', direction: 'outbound', label: '电话：ROI 材料跟进', sentiment: 'hesitant', body: '跟进 ROI 材料阅读情况。' },
    { id: 'f1000000-0000-4000-8000-000000000004', at: at(7, 16), event_type: 'wechat', direction: 'inbound', label: '微信：竞品对比顾虑', sentiment: 'negative', body: '提及竞品对比，需要补充案例材料。' },
    { id: 'f1000000-0000-4000-8000-000000000003', at: at(14, 10), event_type: 'meeting', direction: 'outbound', label: '会议：价格与交付周期', sentiment: 'hesitant', body: '讨论价格与交付周期。' },
    { id: 'f1000000-0000-4000-8000-000000000002', at: at(21, 15), event_type: 'call', direction: 'inbound', label: '电话：预算范围确认', sentiment: 'neutral', body: '客户确认预算区间。' },
    { id: 'f1000000-0000-4000-8000-000000000001', at: at(28, 11), event_type: 'email', direction: 'outbound', label: '邮件：产品方案发送', sentiment: 'positive', body: '产品方案已发送，附 ROI 测算表。' },
  ],
  'a1000000-0000-4000-8000-000000000002': [
    { id: 'f1000000-0000-4000-8000-000000000011', at: at(1, 14), event_type: 'call', direction: 'inbound', label: '电话：预约演示', sentiment: 'positive' },
    { id: 'f1000000-0000-4000-8000-000000000010', at: at(5, 10), event_type: 'email', direction: 'outbound', label: '邮件：白皮书', sentiment: 'positive' },
  ],
  'a1000000-0000-4000-8000-000000000003': [
    { id: 'act-nb-2', at: at(8, 14), event_type: 'call', direction: 'outbound', label: '电话：需求尚不明确', sentiment: 'hesitant', body: '客户表示仍在评估阶段。' },
    { id: 'act-nb-1', at: at(12, 9), event_type: 'email', direction: 'outbound', label: '展会后首次触达', sentiment: 'neutral', body: '发送展会见面纪要。' },
  ],
  'a1000000-0000-4000-8000-000000000004': [
    { id: 'act-xh-4', at: at(30, 9), event_type: 'note', direction: 'outbound', label: '备注：标记为低意向', sentiment: 'negative', body: '暂无预算，建议转入培育池。' },
    { id: 'act-xh-3', at: at(35, 11), event_type: 'call', direction: 'inbound', label: '电话：表示暂无采购计划', sentiment: 'hesitant', body: '对方礼貌拒绝进一步沟通。' },
    { id: 'act-xh-2', at: at(45, 15), event_type: 'email', direction: 'outbound', label: '邮件：资料包发送', sentiment: 'neutral', body: '发送行业案例与报价区间。' },
    { id: 'act-xh-1', at: at(55, 10), event_type: 'call', direction: 'outbound', label: '冷呼：初步介绍产品', sentiment: 'neutral', body: '首次外呼，接通 3 分钟。' },
  ],
  'a1000000-0000-4000-8000-000000000005': [
    { id: 'act-md-1', at: at(2, 11), event_type: 'meeting', direction: 'outbound', label: '季度业务回顾', sentiment: 'positive' },
    { id: 'act-md-2', at: at(15, 15), event_type: 'call', direction: 'inbound', label: '合同续签意向确认', sentiment: 'positive' },
  ],
}

export function mockActivitiesForLead(leadId: string): ActivityTimelineItem[] {
  return DEMO_ACTIVITIES[leadId] ?? []
}
