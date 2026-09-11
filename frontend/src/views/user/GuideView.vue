<template>
  <AppLayout>
    <main class="guide-page">
      <section class="guide-hero">
        <div class="guide-hero__content">
          <p class="guide-kicker">KQS API GUIDE</p>
          <h1 class="guide-title">矿泉水API 使用教程</h1>
        </div>
        <button
          type="button"
          class="guide-contact"
          aria-haspopup="dialog"
          @click="showSupportQr = true"
        >
          <span class="guide-contact__label">联系管理员</span>
          <strong>扫码添加管理员好友</strong>
          <p>点击打开联系二维码，添加时请备注注册邮箱或用户名，方便核对账号。</p>
          <span class="guide-contact__action">
            打开联系二维码
            <Icon name="externalLink" size="xs" />
          </span>
        </button>
      </section>

      <section class="guide-flow" aria-label="使用流程">
        <article
          v-for="(step, index) in steps"
          :key="step.title"
          class="guide-step"
          :style="{ '--delay': `${index * 70}ms` }"
        >
          <div class="guide-step__icon">
            <Icon :name="step.icon" size="md" />
          </div>
          <span class="guide-step__number">{{ String(index + 1).padStart(2, '0') }}</span>
          <h2>{{ step.title }}</h2>
          <p>{{ step.description }}</p>
          <RouterLink v-if="step.to" :to="step.to" class="guide-link">
            {{ step.linkText }}
            <Icon name="arrowRight" size="xs" />
          </RouterLink>
          <a
            v-else-if="step.href"
            :href="step.href"
            target="_blank"
            rel="noopener noreferrer"
            class="guide-link"
          >
            {{ step.linkText }}
            <Icon name="externalLink" size="xs" />
          </a>
        </article>
      </section>

      <section class="guide-grid">
        <article class="guide-panel guide-panel--ccswitch">
          <div class="guide-ccs-copy">
            <p class="guide-kicker">CCSWITCH FIRST</p>
            <h2>推荐：用 CCSwitch 一键导入</h2>
            <p>
              创建好 Key 后，先下载并打开 CCSwitch，再回到密钥页点“导入到 CCS”。
              导入完成后完全退出并重启 Codex，新对话就会走矿泉水API。
            </p>
            <div class="guide-ccs-actions">
              <a
                class="guide-ccs-action guide-ccs-action--secondary"
                :href="CCSWITCH_DOWNLOAD_URL"
                target="_blank"
                rel="noopener noreferrer"
              >
                <Icon name="download" size="sm" />
                下载 CCSwitch
              </a>
              <RouterLink to="/keys" class="guide-ccs-action guide-ccs-action--primary">
                <Icon name="upload" size="sm" />
                去密钥页导入
              </RouterLink>
            </div>
          </div>

          <ol class="guide-ccs-mini" aria-label="CCSwitch 导入步骤">
            <li v-for="step in ccsMiniSteps" :key="step.title">
              <span>{{ step.number }}</span>
              <strong>{{ step.title }}</strong>
              <p>{{ step.description }}</p>
            </li>
          </ol>
        </article>

        <article class="guide-panel guide-panel--wide guide-panel--manual">
          <p class="guide-kicker">MANUAL FALLBACK</p>
          <h2>手动配置备用方案</h2>
          <p>
            推荐优先用 CCSwitch 导入；如果你的环境导入失败，再手动把 base_url 指向正式域名。
            请求路径由 Codex 自己拼接，不要把
            <code>/api/v1/responses</code> 写进配置。
          </p>
          <pre><code>model_provider = "OpenAI"
model = "gpt-5.5"
review_model = "gpt-5.5"

[model_providers.OpenAI]
name = "OpenAI"
base_url = "https://api.cauai.fun"
wire_api = "responses"
requires_openai_auth = true</code></pre>
          <p class="guide-note">
            环境变量里的 <code>OPENAI_API_KEY</code> 填你自己创建的 Key。不要把完整 Key 发给别人，也不要把它发给管理员。
          </p>
        </article>

        <article class="guide-panel">
          <p class="guide-kicker">SAFE FEEDBACK</p>
          <h2>求助时发什么</h2>
          <ul class="guide-list">
            <li>账号邮箱或用户名</li>
            <li>订单号、套餐名或卡密码后四位</li>
            <li>Key ID，不要发完整 Key</li>
            <li>请求时间和错误截图</li>
          </ul>
        </article>

        <article class="guide-panel">
          <p class="guide-kicker">DO NOT SEND</p>
          <h2>这些不要发</h2>
          <ul class="guide-list guide-list--danger">
            <li>完整 API Key</li>
            <li>完整卡密或验证码</li>
            <li>银行卡、身份证等个人敏感信息</li>
            <li>他人的账号信息或隐私内容</li>
          </ul>
        </article>
      </section>

      <section class="guide-rules">
        <div>
          <p class="guide-kicker">SERVICE RULES</p>
          <h2>重要规则先看清楚</h2>
          <p>
            本站不是 OpenAI 官方服务，也不是任何第三方官方系统。卡密兑换后属于已交付数字权益；
            已兑换成功的卡密原则上不支持自动退款，订单异常可以联系管理员人工核对。
          </p>
        </div>
        <div class="guide-legal-links">
          <RouterLink
            v-for="link in legalLinks"
            :key="link.id"
            :to="link.to"
            class="guide-legal-link"
          >
            <Icon name="document" size="sm" />
            <span>{{ link.title }}</span>
          </RouterLink>
        </div>
      </section>
    </main>

    <Teleport to="body">
      <div
        v-if="showSupportQr"
        class="guide-qr-overlay"
        role="presentation"
        @click.self="showSupportQr = false"
      >
        <section
          class="guide-qr-modal"
          role="dialog"
          aria-modal="true"
          aria-labelledby="guideSupportQrTitle"
        >
          <button
            type="button"
            class="guide-qr-close"
            aria-label="关闭联系二维码"
            @click="showSupportQr = false"
          >
            <Icon name="x" size="sm" />
          </button>
          <p class="guide-kicker">SUPPORT QR</p>
          <h2 id="guideSupportQrTitle">联系管理员</h2>
          <p class="guide-qr-lead">打开手机扫码添加管理员，添加或留言时备注注册邮箱或用户名。</p>
          <a
            class="guide-qr-frame"
            :href="supportQrUrl"
            target="_blank"
            rel="noopener noreferrer"
            aria-label="在新窗口打开联系二维码"
          >
            <img :src="supportQrUrl" alt="管理员联系二维码" />
          </a>
          <div class="guide-qr-tips">
            <span>请勿发送完整 API Key</span>
            <span>请勿发送完整卡密</span>
          </div>
        </section>
      </div>
    </Teleport>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { KQS_LEGAL_LINKS } from '@/constants/legalDocuments'
import { LDXP_SHOP_URL } from '@/constants/onboarding'
import { USER_SUBSCRIPTIONS_VISIBLE } from '@/constants/subscriptions'

type GuideIcon = 'creditCard' | 'gift' | 'key' | 'terminal' | 'upload'

const legalLinks = KQS_LEGAL_LINKS
const showSupportQr = ref(false)
const supportQrUrl = '/support-contact-qr.jpg'
const CCSWITCH_DOWNLOAD_URL = 'https://ccswitch.io/'

const steps: Array<{
  icon: GuideIcon
  title: string
  description: string
  linkText: string
  to?: string
  href?: string
}> = []

if (USER_SUBSCRIPTIONS_VISIBLE) {
  steps.push({
    icon: 'creditCard',
    title: '购买卡密',
    description: '在链动小铺选择套餐，付款后复制卡密。站内不直接处理支付，只负责兑换卡密。',
    linkText: '去买卡密',
    href: LDXP_SHOP_URL
  })
}

steps.push(
  {
    icon: 'gift',
    title: '兑换余额',
    description: '回到站内兑换码页面粘贴卡密，成功后余额会立即加到账户里。',
    linkText: '去兑换',
    to: '/redeem'
  },
  {
    icon: 'key',
    title: '创建 Key',
    description: '进入 API 密钥页面创建自己的 Key。建议每个项目单独建一个 Key，方便看用量和停用。',
    linkText: '查 Key',
    to: '/keys'
  },
  {
    icon: 'upload',
    title: '导入 CCSwitch',
    description: '先下载并打开 CCSwitch，再回密钥页点击“导入到 CCS”，不用手动编辑 config.toml。',
    linkText: '去导入',
    to: '/keys'
  },
  {
    icon: 'terminal',
    title: '重启 Codex',
    description: '导入完成后完全退出并重启 Codex，新开的对话就会使用矿泉水API。',
    linkText: '查看备用配置',
    to: '/keys'
  }
)

const ccsMiniSteps = [
  {
    number: '01',
    title: '下载并打开',
    description: '安装 CCSwitch，保持它处于可导入状态。'
  },
  {
    number: '02',
    title: '导入到 CCS',
    description: '在密钥页点“导入到 CCS”，自动写入 base_url、模型和 Key。'
  },
  {
    number: '03',
    title: '重启 Codex',
    description: '完全退出旧进程后重新打开，开始新的对话测试。'
  }
]
</script>

<style scoped>
.guide-page {
  display: grid;
  gap: 0.75rem;
}

.guide-hero {
  display: grid;
  grid-template-columns: minmax(0, 1.7fr) minmax(280px, 0.8fr);
  gap: 0.75rem;
  align-items: stretch;
}

.guide-hero__content,
.guide-contact,
.guide-step,
.guide-panel,
.guide-rules {
  border: 1px solid rgba(43, 132, 83, 0.38);
  border-radius: 8px;
  background:
    linear-gradient(135deg, rgba(34, 197, 94, 0.07), rgba(5, 29, 24, 0.88) 46%),
    rgba(5, 29, 24, 0.78);
  box-shadow:
    0 24px 68px rgba(0, 0, 0, 0.22),
    inset 0 1px 0 rgba(34, 197, 94, 0.08);
}

.guide-hero__content {
  padding: clamp(1.1rem, 2.4vw, 1.75rem);
}

.guide-kicker {
  margin: 0;
  color: rgba(126, 227, 155, 0.84);
  font-size: 0.74rem;
  font-weight: 850;
  letter-spacing: 0;
  text-transform: uppercase;
}

.guide-title {
  margin: 0.45rem 0 0;
  color: #edf7ee;
  font-size: clamp(2rem, 4.2vw, 3.8rem);
  font-weight: 900;
  letter-spacing: 0;
  line-height: 0.98;
}

.guide-contact {
  appearance: none;
  display: flex;
  min-height: 100%;
  flex-direction: column;
  justify-content: flex-end;
  padding: 1rem;
  font: inherit;
  text-align: left;
  cursor: pointer;
  transition:
    transform 220ms cubic-bezier(0.16, 1, 0.3, 1),
    border-color 220ms ease,
    box-shadow 220ms ease;
}

.guide-contact__label {
  color: rgba(220, 239, 224, 0.62);
  font-size: 0.8rem;
  font-weight: 800;
}

.guide-contact strong {
  margin-top: 0.4rem;
  color: #edf7ee;
  font-size: 1.05rem;
  line-height: 1.6;
}

.guide-contact__action {
  margin-top: 0.95rem;
  display: inline-flex;
  width: fit-content;
  align-items: center;
  gap: 0.35rem;
  border: 1px solid rgba(34, 197, 94, 0.28);
  border-radius: 8px;
  background: rgba(34, 197, 94, 0.12);
  padding: 0.5rem 0.75rem;
  color: #7ee39b;
  font-size: 0.85rem;
  font-weight: 850;
}

.guide-contact:hover {
  transform: translateY(-2px);
  border-color: rgba(34, 197, 94, 0.52);
  box-shadow: 0 28px 72px rgba(0, 0, 0, 0.3);
}

.guide-contact:focus-visible {
  outline: 3px solid rgba(0, 137, 78, 0.22);
  outline-offset: 4px;
}

.guide-contact p,
.guide-panel p,
.guide-rules p,
.guide-note {
  margin-bottom: 0;
  color: rgba(220, 239, 224, 0.68);
  line-height: 1.75;
}

.guide-flow {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 0.75rem;
}

.guide-step {
  position: relative;
  overflow: hidden;
  padding: 0.95rem;
  animation: guide-rise 420ms cubic-bezier(0.16, 1, 0.3, 1) both;
  animation-delay: var(--delay);
}

.guide-step__icon {
  display: flex;
  width: 2.35rem;
  height: 2.35rem;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(34, 197, 94, 0.2);
  border-radius: 8px;
  background: rgba(34, 197, 94, 0.12);
  color: #7ee39b;
}

.guide-step__number {
  position: absolute;
  right: 1rem;
  top: 0.75rem;
  color: rgba(220, 239, 224, 0.12);
  font-size: 2rem;
  font-weight: 900;
}

.guide-step h2,
.guide-panel h2,
.guide-rules h2 {
  margin: 0.8rem 0 0;
  color: #edf7ee;
  font-size: 1.15rem;
  font-weight: 850;
  letter-spacing: 0;
}

.guide-step p {
  margin-top: 0.55rem;
  color: rgba(220, 239, 224, 0.64);
  font-size: 0.9rem;
  line-height: 1.65;
}

.guide-link {
  margin-top: 1rem;
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  color: #7ee39b;
  font-size: 0.9rem;
  font-weight: 850;
  transition: transform 180ms ease, color 180ms ease;
}

.guide-link:hover {
  transform: translateX(2px);
  color: #a7f3bf;
}

.guide-grid {
  display: grid;
  grid-template-columns: 1.2fr 0.8fr;
  gap: 0.75rem;
}

.guide-panel {
  padding: 1rem;
}

.guide-panel--wide {
  grid-row: span 2;
}

.guide-panel--ccswitch {
  position: relative;
  grid-column: 1 / -1;
  display: grid;
  grid-template-columns: minmax(0, 0.9fr) minmax(320px, 1.1fr);
  gap: 0.75rem;
  overflow: hidden;
}

.guide-panel--ccswitch::after {
  content: '';
  position: absolute;
  right: -5rem;
  top: -5rem;
  width: 14rem;
  height: 14rem;
  border-radius: 999px;
  background: rgba(34, 197, 94, 0.08);
  pointer-events: none;
}

.guide-ccs-copy,
.guide-ccs-mini {
  position: relative;
  z-index: 1;
}

.guide-ccs-copy h2 {
  font-size: clamp(1.3rem, 2.5vw, 1.9rem);
}

.guide-ccs-actions {
  margin-top: 1rem;
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.guide-ccs-action {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.45rem;
  border-radius: 8px;
  padding: 0.66rem 0.9rem;
  font-size: 0.9rem;
  font-weight: 850;
  transition:
    transform 180ms cubic-bezier(0.16, 1, 0.3, 1),
    border-color 180ms ease,
    background 180ms ease;
}

.guide-ccs-action:hover {
  transform: translateY(-1px);
}

.guide-ccs-action:active {
  transform: translateY(0) scale(0.98);
}

.guide-ccs-action--primary {
  background: #16a34a;
  color: #fff;
  box-shadow: 0 14px 32px rgba(0, 0, 0, 0.22);
}

.guide-ccs-action--secondary {
  border: 1px solid rgba(43, 132, 83, 0.38);
  background: rgba(3, 22, 18, 0.64);
  color: #edf7ee;
}

.guide-ccs-mini {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.75rem;
}

.guide-ccs-mini li {
  min-height: 9.8rem;
  border-radius: 8px;
  border: 1px solid rgba(43, 132, 83, 0.28);
  background:
    linear-gradient(145deg, rgba(34, 197, 94, 0.08), rgba(3, 22, 18, 0.68)),
    rgba(3, 22, 18, 0.64);
  padding: 0.9rem;
}

.guide-ccs-mini span {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  border: 1px solid rgba(34, 197, 94, 0.2);
  border-radius: 8px;
  background: rgba(34, 197, 94, 0.12);
  color: #7ee39b;
  font-size: 0.78rem;
  font-weight: 900;
}

.guide-ccs-mini strong {
  margin-top: 0.85rem;
  display: block;
  color: #edf7ee;
  font-size: 1rem;
  font-weight: 850;
}

.guide-ccs-mini p {
  margin-top: 0.45rem;
  font-size: 0.84rem;
  line-height: 1.65;
}

.guide-panel pre {
  margin-top: 1rem;
  overflow-x: auto;
  border: 1px solid rgba(43, 132, 83, 0.28);
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.86);
  padding: 1.1rem;
  color: #d7f7df;
  font-size: 0.82rem;
  line-height: 1.7;
}

.guide-panel code {
  border-radius: 6px;
  background: rgba(34, 197, 94, 0.12);
  padding: 0.1rem 0.35rem;
  color: #a7f3bf;
  font-size: 0.86em;
}

.guide-panel pre code {
  background: transparent;
  padding: 0;
  color: inherit;
}

.guide-list {
  margin-top: 1rem;
  display: grid;
  gap: 0.65rem;
  color: rgba(220, 239, 224, 0.72);
  font-size: 0.92rem;
}

.guide-list li {
  border-left: 3px solid rgba(34, 197, 94, 0.3);
  padding-left: 0.75rem;
}

.guide-list--danger li {
  border-left-color: rgba(217, 74, 58, 0.55);
}

.guide-rules {
  display: grid;
  grid-template-columns: minmax(0, 1.05fr) minmax(280px, 0.95fr);
  gap: 0.75rem;
  align-items: start;
  padding: 1rem;
}

.guide-legal-links {
  display: grid;
  gap: 0.65rem;
}

.guide-legal-link {
  display: flex;
  align-items: center;
  gap: 0.7rem;
  border-radius: 8px;
  border: 1px solid rgba(43, 132, 83, 0.28);
  background: rgba(3, 22, 18, 0.64);
  padding: 0.8rem 0.9rem;
  color: #edf7ee;
  font-weight: 800;
  transition: transform 180ms ease, border-color 180ms ease, background 180ms ease;
}

.guide-legal-link:hover {
  transform: translateY(-1px);
  border-color: rgba(34, 197, 94, 0.52);
  background: rgba(34, 197, 94, 0.12);
}

.guide-qr-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: grid;
  place-items: center;
  background:
    radial-gradient(circle at 50% 18%, rgba(0, 137, 78, 0.18), transparent 34%),
    rgba(2, 18, 10, 0.56);
  padding: 1rem;
  backdrop-filter: blur(18px);
  animation: guide-fade-in 160ms ease both;
}

.guide-qr-modal {
  position: relative;
  width: min(92vw, 28rem);
  overflow: hidden;
  border: 1px solid rgba(0, 137, 78, 0.16);
  border-radius: 2rem;
  background:
    radial-gradient(circle at 0% 0%, rgba(238, 195, 73, 0.16), transparent 32%),
    linear-gradient(145deg, rgba(255, 255, 255, 0.98), rgba(239, 249, 243, 0.96));
  box-shadow: 0 34px 100px rgba(0, 20, 10, 0.28);
  padding: clamp(1.25rem, 4vw, 2rem);
  color: #032f1c;
  animation: guide-modal-in 220ms cubic-bezier(0.16, 1, 0.3, 1) both;
}

.guide-qr-modal::after {
  content: '';
  position: absolute;
  right: -4rem;
  bottom: -5rem;
  width: 12rem;
  height: 12rem;
  border-radius: 999px;
  background: rgba(0, 137, 78, 0.1);
  pointer-events: none;
}

.guide-qr-close {
  position: absolute;
  right: 1rem;
  top: 1rem;
  z-index: 1;
  display: grid;
  width: 2.4rem;
  height: 2.4rem;
  place-items: center;
  border-radius: 999px;
  background: rgba(0, 137, 78, 0.08);
  color: #006b31;
  transition: transform 160ms ease, background 160ms ease;
}

.guide-qr-close:hover {
  transform: rotate(8deg);
  background: rgba(0, 137, 78, 0.14);
}

.guide-qr-modal h2 {
  margin-top: 0.85rem;
  color: #032f1c;
  font-size: clamp(1.8rem, 5vw, 2.75rem);
  font-weight: 900;
  letter-spacing: -0.055em;
}

.guide-qr-lead {
  margin-top: 0.65rem;
  color: rgba(3, 47, 28, 0.68);
  line-height: 1.75;
}

.guide-qr-frame {
  position: relative;
  z-index: 1;
  margin-top: 1.25rem;
  display: block;
  overflow: hidden;
  border-radius: 1.5rem;
  border: 1px solid rgba(0, 137, 78, 0.14);
  background: #fff;
  padding: 0.75rem;
  box-shadow: inset 0 0 0 1px rgba(0, 137, 78, 0.06);
}

.guide-qr-frame img {
  display: block;
  width: 100%;
  max-height: min(58vh, 32rem);
  object-fit: contain;
  border-radius: 1rem;
}

.guide-qr-tips {
  position: relative;
  z-index: 1;
  margin-top: 1rem;
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.guide-qr-tips span {
  border-radius: 999px;
  background: rgba(0, 137, 78, 0.09);
  padding: 0.45rem 0.7rem;
  color: rgba(3, 47, 28, 0.7);
  font-size: 0.78rem;
  font-weight: 800;
}

:global(.dark) .guide-hero__content,
:global(.dark) .guide-contact,
:global(.dark) .guide-step,
:global(.dark) .guide-panel,
:global(.dark) .guide-rules {
  border-color: rgba(255, 255, 255, 0.1);
  background:
    linear-gradient(135deg, rgba(10, 38, 25, 0.86), rgba(18, 24, 21, 0.82)),
    radial-gradient(circle at 100% 0%, rgba(238, 195, 73, 0.12), transparent 34%);
  box-shadow: 0 20px 54px rgba(0, 0, 0, 0.22);
}

:global(.dark) .guide-title,
:global(.dark) .guide-step h2,
:global(.dark) .guide-panel h2,
:global(.dark) .guide-rules h2 {
  color: #f8fff9;
}

:global(.dark) .guide-contact p,
:global(.dark) .guide-panel p,
:global(.dark) .guide-rules p,
:global(.dark) .guide-note,
:global(.dark) .guide-step p,
:global(.dark) .guide-list {
  color: rgba(255, 255, 255, 0.68);
}

:global(.dark) .guide-contact strong,
:global(.dark) .guide-link,
:global(.dark) .guide-legal-link {
  color: #bde8cb;
}

:global(.dark) .guide-legal-link {
  border-color: rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.06);
}

:global(.dark) .guide-ccs-mini li {
  border-color: rgba(255, 255, 255, 0.1);
  background:
    linear-gradient(145deg, rgba(255, 255, 255, 0.08), rgba(255, 255, 255, 0.04)),
    radial-gradient(circle at 100% 0%, rgba(0, 137, 78, 0.18), transparent 38%);
}

:global(.dark) .guide-ccs-mini strong {
  color: #f8fff9;
}

:global(.dark) .guide-ccs-action--secondary {
  border-color: rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.07);
  color: #bde8cb;
}

@keyframes guide-rise {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes guide-fade-in {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes guide-modal-in {
  from {
    opacity: 0;
    transform: translateY(14px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

@media (max-width: 1024px) {
  .guide-hero,
  .guide-rules,
  .guide-grid,
  .guide-panel--ccswitch {
    grid-template-columns: 1fr;
  }

  .guide-flow {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .guide-flow {
    grid-template-columns: 1fr;
  }

  .guide-ccs-mini {
    grid-template-columns: 1fr;
  }
}
</style>
