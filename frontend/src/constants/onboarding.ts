export interface LdxpPackage {
  id: string
  priceCny: string
  quotaLabel: string
  quotaUsd?: number
  url: string
  recommended?: boolean
  note?: string
}

export const LDXP_SHOP_URL = 'https://pay.ldxp.cn/shop/B59CCLX7'

export const LDXP_BALANCE_PACKAGES: LdxpPackage[] = [
  {
    id: 'balance-18',
    priceCny: '1.5',
    quotaLabel: '18.8 美金额度',
    quotaUsd: 18.8,
    url: LDXP_SHOP_URL,
    note: '入门体验'
  },
  {
    id: 'balance-38',
    priceCny: '3',
    quotaLabel: '38.8 美金额度',
    quotaUsd: 38.8,
    url: LDXP_SHOP_URL,
    note: '轻量补充'
  },
  {
    id: 'balance-88',
    priceCny: '7',
    quotaLabel: '88.8 美金额度',
    quotaUsd: 88.8,
    url: LDXP_SHOP_URL,
    note: '常用补给'
  },
  {
    id: 'balance-188',
    priceCny: '15',
    quotaLabel: '188.8 美金额度',
    quotaUsd: 188.8,
    url: LDXP_SHOP_URL
  },
  {
    id: 'balance-1088',
    priceCny: '88',
    quotaLabel: '1088.8 美金额度',
    quotaUsd: 1088.8,
    url: LDXP_SHOP_URL
  },
  {
    id: 'balance-8888',
    priceCny: '648',
    quotaLabel: '8888 美金额度',
    quotaUsd: 8888,
    url: LDXP_SHOP_URL,
    note: '大额补给'
  }
]

// Backward-compatible alias used by older balance purchase components.
export const LDXP_PACKAGES = LDXP_BALANCE_PACKAGES

export const LDXP_FALLBACK_SUBSCRIPTION_PLANS = [
  {
    name: '日卡 10 刀',
    price: 1,
    validity_days: 1,
    daily_limit_usd: 10,
    weekly_limit_usd: 10,
    sort_order: 10
  },
  {
    name: '日卡 20 刀',
    price: 1.9,
    validity_days: 1,
    daily_limit_usd: 20,
    weekly_limit_usd: 20,
    sort_order: 20
  },
  {
    name: '日卡 45 刀',
    price: 4,
    validity_days: 1,
    daily_limit_usd: 45,
    weekly_limit_usd: 45,
    sort_order: 30
  },
  {
    name: '日卡 60 刀',
    price: 5.2,
    validity_days: 1,
    daily_limit_usd: 60,
    weekly_limit_usd: 60,
    sort_order: 40
  },
  {
    name: '日卡 90 刀',
    price: 7.5,
    validity_days: 1,
    daily_limit_usd: 90,
    weekly_limit_usd: 90,
    sort_order: 50
  },
  {
    name: '周卡 45 刀',
    price: 25,
    validity_days: 7,
    daily_limit_usd: 45,
    weekly_limit_usd: 315,
    sort_order: 110
  },
  {
    name: '周卡 90 刀',
    price: 48,
    validity_days: 7,
    daily_limit_usd: 90,
    weekly_limit_usd: 630,
    sort_order: 120
  },
  {
    name: '周卡 135 刀',
    price: 70,
    validity_days: 7,
    daily_limit_usd: 135,
    weekly_limit_usd: 945,
    sort_order: 130
  },
  {
    name: '周卡 180 刀',
    price: 90,
    validity_days: 7,
    daily_limit_usd: 180,
    weekly_limit_usd: 1260,
    sort_order: 140
  }
] as const

export const SUBSCRIPTION_DISPLAY_RATE_MULTIPLIER = 2

export const CODEX_API_BASE_URL = 'https://api.cauai.fun'
export const CODEX_RESPONSES_ENDPOINT = `${CODEX_API_BASE_URL}/responses`
export const CODEX_DEFAULT_MODEL = 'gpt-5.5'

export const CODEX_DEFAULT_KEY_LIMITS = {
  quotaUsd: 0,
  rateLimit5hUsd: 0,
  rateLimit1dUsd: 0,
  rateLimit7dUsd: 0
} as const
