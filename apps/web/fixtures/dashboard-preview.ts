import type { ChartFunnelItem } from '@crm/ui-kit'
import type { DashboardCalendarEvent, TeamHeatmapMemberRow } from '~/types/dashboard'

export const DASHBOARD_PREVIEW_WEEKLY_FOLLOWUPS = 12

/** Zone E — Preview fixtures only; never persisted */
export const DASHBOARD_FUNNEL_FIXTURE: ChartFunnelItem[] = [
  { name: 'leads', value: 128 },
  { name: 'qualified', value: 86 },
  { name: 'opportunity', value: 41 },
  { name: 'won', value: 17 },
]

export const DASHBOARD_TEAM_HEATMAP_MEMBERS: TeamHeatmapMemberRow[] = [
  { memberId: 'user-demo-001', memberName: '张三', healthScore: 91, emotionScore: 88 },
  { memberId: 'user-demo-002', memberName: '李四', healthScore: 62, emotionScore: 58 },
  { memberId: 'user-demo-003', memberName: '王磊', healthScore: 45, emotionScore: 52 },
  { memberId: 'user-demo-004', memberName: '赵敏', healthScore: 78, emotionScore: 81 },
  { memberId: 'user-demo-005', memberName: '陈浩', healthScore: 55, emotionScore: 49 },
]

const todayIso = () => new Date().toISOString().slice(0, 10)

export const DASHBOARD_CALENDAR_FIXTURE: DashboardCalendarEvent[] = [
  {
    id: 'demo-meeting-1',
    date: todayIso(),
    time: '15:30',
    title: 'CRM 周会',
    subtitle: 'Google Meet · 团队同步',
  },
  {
    id: 'demo-meeting-2',
    date: todayIso(),
    time: '17:00',
    title: '客户方案演示',
    subtitle: 'Slack · 飞书团队',
  },
]
