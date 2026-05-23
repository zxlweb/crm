export type MockUserProfile = {
  id: string
  name: string
  email: string
  role: string
}

/** Demo 负责人 — 与 leads.mock owner_id 对齐 */
export const MOCK_USER_PROFILES: Record<string, MockUserProfile> = {
  'user-demo-001': {
    id: 'user-demo-001',
    name: '张明',
    email: 'zhang.ming@demo.local',
    role: 'Sales',
  },
  'user-demo-002': {
    id: 'user-demo-002',
    name: '李悦',
    email: 'li.yue@demo.local',
    role: 'Sales',
  },
}
