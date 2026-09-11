import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import i18n, { initI18n } from './i18n'
import { useAppStore } from '@/stores/app'
import './style.css'

function initThemeClass() {
  // Dark mode is temporarily disabled while the visual system is being refined.
  // Keep the implementation reversible by only forcing the runtime theme here.
  document.documentElement.classList.remove('dark')
  localStorage.setItem('theme', 'light')
}

async function bootstrap() {
  // Apply theme class globally before app mount to keep all routes consistent.
  initThemeClass()

  const app = createApp(App)
  const pinia = createPinia()
  app.use(pinia)

  // Hydrate public settings before the first render. AuthLayout and auth
  // views still request a fresh copy, but the shared store must not transition
  // from its defaults while Vue is mounting the initial route.
  const appStore = useAppStore()
  appStore.initFromInjectedConfig()

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
