import type { PlanDistribution, TenantHealthData, TopTenants } from '~/types/admin-insights'

export function mockTenantHealth(): TenantHealthData {
  return {
    dimensions: ['activity', 'config_completeness', 'audit_risk', 'data_freshness', 'feature_adoption'],
    items: [
      {
        tenant_id: 't-001',
        tenant_name: 'Acme China',
        scores: { activity: 78, config_completeness: 85, audit_risk: 42, data_freshness: 80, feature_adoption: 66 },
        overall_score: 70,
      },
      {
        tenant_id: 't-002',
        tenant_name: 'Beta Corp',
        scores: { activity: 92, config_completeness: 70, audit_risk: 25, data_freshness: 95, feature_adoption: 88 },
        overall_score: 74,
      },
      {
        tenant_id: 't-003',
        tenant_name: 'Gamma Ltd',
        scores: { activity: 35, config_completeness: 40, audit_risk: 80, data_freshness: 30, feature_adoption: 20 },
        overall_score: 41,
      },
    ],
  }
}

export function mockPlanDistribution(): PlanDistribution {
  return {
    items: [
      { plan: 'Enterprise', count: 5 },
      { plan: 'Professional', count: 12 },
      { plan: 'Starter', count: 28 },
      { plan: 'Free Trial', count: 8 },
    ],
  }
}

export function mockTopTenants(): TopTenants {
  return {
    metric: 'activity',
    items: [
      { tenant_id: 't-002', tenant_name: 'Beta Corp', value: 92 },
      { tenant_id: 't-001', tenant_name: 'Acme China', value: 78 },
      { tenant_id: 't-004', tenant_name: 'Delta Inc', value: 65 },
      { tenant_id: 't-005', tenant_name: 'Epsilon Group', value: 58 },
      { tenant_id: 't-003', tenant_name: 'Gamma Ltd', value: 35 },
    ],
  }
}
