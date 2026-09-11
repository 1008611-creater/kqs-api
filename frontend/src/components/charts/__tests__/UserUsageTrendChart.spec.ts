import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'

import UserUsageTrendChart from '../UserUsageTrendChart.vue'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key,
  }),
}))

vi.mock('vue-chartjs', () => ({
  Line: {
    props: ['data', 'options'],
    template: '<div class="chart-data">{{ JSON.stringify(data) }}</div>',
  },
}))

describe('UserUsageTrendChart', () => {
  it('plots actual cost per user instead of token totals', () => {
    const wrapper = mount(UserUsageTrendChart, {
      props: {
        userTrendData: [
          { date: '2026-08-04 09', user_id: 1, email: 'one@example.com', username: '', requests: 2, tokens: 9000, cost: 0.9, actual_cost: 0.12 },
          { date: '2026-08-04 10', user_id: 1, email: 'one@example.com', username: '', requests: 1, tokens: 100, cost: 0.1, actual_cost: 0.03 },
          { date: '2026-08-04 09', user_id: 2, email: '', username: 'two', requests: 1, tokens: 200, cost: 0.2, actual_cost: 0.2 },
        ],
      },
      global: { stubs: { LoadingSpinner: true } },
    })

    const chartData = JSON.parse(wrapper.find('.chart-data').text())
    expect(chartData.labels).toEqual(['2026-08-04 09', '2026-08-04 10'])
    expect(chartData.datasets).toHaveLength(2)
    expect(chartData.datasets[0].data).toEqual([0.12, 0.03])
    expect(chartData.datasets[1].data).toEqual([0.2, 0])
    expect(chartData.datasets[0].label).toContain('#1')
    expect(chartData.datasets[1].label).toContain('two')
  })

  it('shows an empty state when there is no user usage', () => {
    const wrapper = mount(UserUsageTrendChart, {
      props: { userTrendData: [] },
      global: { stubs: { LoadingSpinner: true } },
    })

    expect(wrapper.text()).toContain('admin.dashboard.noDataAvailable')
  })
})
