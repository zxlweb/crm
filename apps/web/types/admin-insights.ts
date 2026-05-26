export type HealthDimension =
  | 'activity'
  | 'config_completeness'
  | 'audit_risk'
  | 'data_freshness'
  | 'feature_adoption'

export interface TenantHealthScores {
  activity: number
  config_completeness: number
  audit_risk: number
  data_freshness: number
  feature_adoption: number
}

export interface TenantHealthItem {
  tenant_id: string
  tenant_name: string
  scores: TenantHealthScores
  overall_score: number
}

export interface TenantHealthData {
  dimensions: HealthDimension[]
  items: TenantHealthItem[]
}

export interface PlanDistributionItem {
  plan: string
  count: number
}

export interface PlanDistribution {
  items: PlanDistributionItem[]
}

export interface TopTenantItem {
  tenant_id: string
  tenant_name: string
  value: number
}

export interface TopTenants {
  items: TopTenantItem[]
  metric: string
}

export interface AdminInsightsQuery {
  from?: string
  to?: string
  metric?: 'activity' | 'revenue' | 'risk'
  limit?: number
}
