import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import i18n, { initI18n } from './i18n'
import { useAppStore } from '@/stores/app'
import { updateFavicon } from '@/utils/branding'
import { isIOSDevice } from '@/utils/device'
import './style.css'

function initIOSViewportZoomFix() {
  // iOS Safari 在输入框字号小于 16px 时聚焦会自动放大页面，且失焦后不会恢复。
  // 限制 maximum-scale 可阻止该行为；iOS 10+ 用户仍可双指手动缩放，不影响可访问性。
  // 仅在 iOS 设备上注入，避免影响 Android Chrome 的手动缩放能力。
  if (!isIOSDevice()) return

  const viewport = document.querySelector('meta[name="viewport"]')
  if (!viewport) return

  const content = viewport.getAttribute('content') || ''
  if (/maximum-scale/i.test(content)) return
  viewport.setAttribute('content', `${content}, maximum-scale=1.0`)
}

function initThemeClass() {
  // Dark mode is temporarily disabled while the visual system is being refined.
  // Keep the implementation reversible by only forcing the runtime theme here.
  document.documentElement.classList.remove('dark')
  localStorage.setItem('theme', 'light')
}

async function bootstrap() {
  // Apply theme class globally before app mount to keep all routes consistent.
  initThemeClass()
  initIOSViewportZoomFix()

  const app = createApp(App)
  const pinia = createPinia()
  app.use(pinia)

  // Hydrate public settings before the first render. AuthLayout and auth
  // views still request a fresh copy, but the shared store must not transition
  // from its defaults while Vue is mounting the initial route.
  const appStore = useAppStore()
  appStore.initFromInjectedConfig()

<<<<<<< HEAD
  // Set document title immediately after config is loaded
  if (appStore.siteName && appStore.siteName !== 'Sub2API') {
    document.title = `${appStore.siteName} - AI API Gateway`
  }
  updateFavicon(appStore.siteLogo)

=======
>>>>>>> 1157e4cfe6271a81c9dc96df84af6c6a3b22f831
  await initI18n()

  app.use(router)
  app.use(i18n)
  app.mount('#app')

  // Client-side navigation can finish after the root is mounted. Mounting
  // first prevents a redirect during router initialization from invalidating
  // the root vnode's anchor nodes.
  await router.isReady()
}

bootstrap()
