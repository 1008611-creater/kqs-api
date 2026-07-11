<template>
  <div class="relative" ref="dropdownRef">
    <button
      @click="toggleDropdown"
      :disabled="switching"
      class="flex min-h-11 items-center gap-2 rounded-full px-2 py-1.5 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-100 dark:text-primary-100/80 dark:hover:bg-white/10"
      :title="currentLocale?.name"
    >
      <span class="inline-flex h-7 min-w-8 items-center justify-center rounded-full border border-primary-700/15 bg-primary-50 px-2 text-[11px] font-bold tracking-[0.1em] text-primary-700 dark:border-primary-200/20 dark:bg-primary-100/10 dark:text-primary-100">
        {{ currentLocale?.code.toUpperCase() }}
      </span>
      <span class="hidden md:inline">{{ currentLocale?.name }}</span>
      <Icon
        name="chevronDown"
        size="xs"
        class="text-gray-400 transition-transform duration-200"
        :class="{ 'rotate-180': isOpen }"
      />
    </button>

    <transition name="dropdown">
      <div
        v-if="isOpen"
        class="absolute right-0 z-50 mt-1 w-36 overflow-hidden rounded-xl border border-primary-900/10 bg-white shadow-lg dark:border-primary-100/10 dark:bg-[#102318]"
      >
        <button
          v-for="locale in availableLocales"
          :key="locale.code"
          :disabled="switching"
          @click="selectLocale(locale.code)"
          class="flex min-h-11 w-full items-center gap-2 px-3 py-2 text-sm text-gray-700 transition-colors hover:bg-primary-50 dark:text-primary-50 dark:hover:bg-primary-100/10"
          :class="{
            'bg-primary-50 text-primary-600 dark:bg-primary-900/20 dark:text-primary-400':
              locale.code === currentLocaleCode
          }"
        >
          <span class="inline-flex h-7 min-w-8 items-center justify-center rounded-full border border-primary-700/15 bg-primary-50 px-2 text-[11px] font-bold tracking-[0.1em] text-primary-700 dark:border-primary-200/20 dark:bg-primary-100/10 dark:text-primary-100">
            {{ locale.code.toUpperCase() }}
          </span>
          <span>{{ locale.name }}</span>
          <Icon v-if="locale.code === currentLocaleCode" name="check" size="sm" class="ml-auto text-primary-500" />
        </button>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import { setLocale, availableLocales } from '@/i18n'

const { locale } = useI18n()

const isOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)
const switching = ref(false)

const currentLocaleCode = computed(() => locale.value)
const currentLocale = computed(() => availableLocales.find((l) => l.code === locale.value))

function toggleDropdown() {
  isOpen.value = !isOpen.value
}

async function selectLocale(code: string) {
  if (switching.value || code === currentLocaleCode.value) {
    isOpen.value = false
    return
  }
  switching.value = true
  try {
    await setLocale(code)
    isOpen.value = false
  } finally {
    switching.value = false
  }
}

function handleClickOutside(event: MouseEvent) {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target as Node)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.15s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: scale(0.95) translateY(-4px);
}
</style>
