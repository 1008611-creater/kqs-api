<template>
  <span :class="['brand-logo', `brand-logo--${variant}`]">
    <span class="brand-logo__plate">
      <img v-if="src" :src="src" :alt="alt" class="brand-logo__image" />
    </span>
  </span>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  src: string
  alt?: string
  variant?: 'sidebar' | 'home' | 'auth' | 'compact'
}>(), {
  alt: 'Logo',
  variant: 'compact',
})
</script>

<style scoped>
.brand-logo {
  --brand-logo-radius: 1rem;
  position: relative;
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border: 1px solid rgba(0, 93, 45, 0.18);
  background:
    linear-gradient(145deg, rgba(255, 255, 255, 0.98), rgba(239, 248, 235, 0.92)),
    linear-gradient(145deg, rgba(0, 135, 60, 0.16), rgba(214, 169, 44, 0.16));
  box-shadow:
    0 16px 34px rgba(2, 43, 18, 0.16),
    0 2px 5px rgba(2, 43, 18, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.86);
  isolation: isolate;
}

.brand-logo::before {
  content: '';
  position: absolute;
  inset: 1px;
  border-radius: inherit;
  border: 1px solid rgba(255, 255, 255, 0.72);
  pointer-events: none;
}

.brand-logo::after {
  content: '';
  position: absolute;
  top: 0;
  left: 12%;
  right: 12%;
  height: 42%;
  border-radius: 999px;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.62), rgba(255, 255, 255, 0));
  pointer-events: none;
}

.brand-logo__plate {
  position: relative;
  z-index: 1;
  display: flex;
  width: calc(100% - 8px);
  height: calc(100% - 8px);
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border-radius: calc(var(--brand-logo-radius) - 0.22rem);
  background: linear-gradient(160deg, rgba(255, 255, 255, 0.98), rgba(250, 253, 246, 0.92));
  box-shadow:
    inset 0 0 0 1px rgba(0, 93, 45, 0.08),
    inset 0 -10px 18px rgba(0, 93, 45, 0.035);
}

.brand-logo__image {
  display: block;
  max-width: 84%;
  max-height: 84%;
  object-fit: contain;
  filter:
    drop-shadow(0 4px 8px rgba(0, 92, 44, 0.14))
    drop-shadow(0 1px 0 rgba(255, 255, 255, 0.92));
}

.brand-logo--sidebar,
.brand-logo--compact {
  --brand-logo-radius: 1rem;
  width: 3rem;
  height: 3.5rem;
  border-radius: 1rem;
}

.brand-logo--home {
  --brand-logo-radius: 1.25rem;
  width: 4rem;
  height: 5rem;
  border-radius: 1.25rem;
}

.brand-logo--auth {
  --brand-logo-radius: 1.35rem;
  width: 5rem;
  height: 6rem;
  border-radius: 1.35rem;
}

@media (min-width: 640px) {
  .brand-logo--home {
    width: 5rem;
    height: 6rem;
  }
}

:global(.dark) .brand-logo {
  border-color: rgba(255, 255, 255, 0.13);
  background:
    linear-gradient(145deg, rgba(255, 255, 255, 0.95), rgba(230, 244, 225, 0.9)),
    linear-gradient(145deg, rgba(0, 135, 60, 0.18), rgba(214, 169, 44, 0.12));
  box-shadow:
    0 18px 34px rgba(0, 0, 0, 0.32),
    0 0 0 1px rgba(255, 255, 255, 0.06),
    inset 0 1px 0 rgba(255, 255, 255, 0.78);
}
</style>
