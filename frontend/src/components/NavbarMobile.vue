<template>
  <!--
    The mobile navigation shell: a slim top bar for context (where am I, how
    do I get back) and a bottom tab bar for travel. Primary destinations moved
    out of the hamburger drawer and into the bottom bar because that is the
    only part of a phone screen a thumb reaches without regripping, and it
    makes the app's structure visible instead of hiding it behind a menu.
  -->
  <div class="bk-nav">
    <header class="bk-topbar" :class="{ 'is-home': isHome }">
      <button
        v-if="!isHome"
        type="button"
        class="bk-topbar-btn"
        aria-label="Go back"
        @click="goBack"
      >
        <i class="bi bi-chevron-left" aria-hidden="true"></i>
      </button>
      <span v-else class="bk-topbar-slot" aria-hidden="true"></span>

      <h1 class="bk-topbar-title">{{ title }}</h1>

      <router-link
        v-if="isHome"
        to="/logout"
        class="bk-topbar-btn"
        aria-label="Log out"
      >
        <i class="bi bi-box-arrow-right" aria-hidden="true"></i>
      </router-link>
      <span v-else class="bk-topbar-slot" aria-hidden="true"></span>
    </header>

    <nav class="bk-tabbar" aria-label="Main navigation">
      <router-link
        v-for="item in items"
        :key="item.to"
        :to="item.to"
        class="bk-tab"
        active-class="is-active"
      >
        <span class="bk-tab-icon">
          <i class="bi" :class="item.icon" aria-hidden="true"></i>
        </span>
        <span class="bk-tab-label">{{ item.label }}</span>
      </router-link>
    </nav>
  </div>
</template>

<script setup>
import { computed } from "vue";
import { useRouter } from "vue-router";

const props = defineProps({
  title: {
    type: String,
    required: true,
  },
  // The dashboards pass "<role>-home"; its presence is what marks a screen as
  // the root of the app, which is the only place the back affordance is wrong.
  role: {
    type: String,
    required: false,
  },
});

const router = useRouter();

const isHome = computed(() => (props.role || "").endsWith("-home"));

// The same set the hamburger drawer carried, minus logout, which moved to the
// top bar: six is one too many for a bottom bar, and logout is not a
// destination you want one thumb-width from the tab you use most.
const items = [
  { to: "/dashboard", label: "Home", icon: "bi-house" },
  { to: "/boat", label: "Boat", icon: "bi-speedometer2" },
  { to: "/ride", label: "Ride", icon: "bi-stopwatch" },
  { to: "/account", label: "Account", icon: "bi-person" },
  { to: "/info", label: "Info", icon: "bi-info-circle" },
];

const goBack = () => {
  // vue-router records the previous entry in history.state; without one the
  // user deep-linked here and "back" would take them out of the app.
  if (window.history.state && window.history.state.back) {
    router.back();
  } else {
    router.push("/dashboard");
  }
};
</script>

<style scoped>
/* ---- Top bar ------------------------------------------------------------ */
.bk-topbar {
  position: fixed;
  top: 0;
  right: 0;
  left: 0;
  z-index: 1030;
  display: flex;
  align-items: center;
  gap: 0.25rem;
  height: calc(var(--bk-topbar-h) + var(--bk-safe-top));
  padding: var(--bk-safe-top) 0.5rem 0;
  border-bottom: 1px solid var(--bk-surface-edge);
  background-color: var(--bk-glass-strong);
  /* The canvas is a gradient, so the bar blurs what scrolls under it rather
     than painting a flat block over it. */
  backdrop-filter: blur(20px) saturate(140%);
  -webkit-backdrop-filter: blur(20px) saturate(140%);
}

.bk-topbar-title {
  flex: 1 1 auto;
  min-width: 0;
  margin: 0;
  color: var(--bk-on-dark);
  font-size: 1rem;
  font-weight: 600;
  letter-spacing: -0.01em;
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* On the root screen there is nothing to go back to, so the title stops
   pretending to be a centred page heading and behaves like a name. */
.bk-topbar.is-home .bk-topbar-title {
  padding-left: 0.35rem;
  font-size: 1.15rem;
  font-weight: 700;
  text-align: left;
}

.bk-topbar-btn,
.bk-topbar-slot {
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  padding: 0;
  border: 0;
  border-radius: var(--bk-r-pill);
  background-color: transparent;
  color: var(--bk-on-dark);
  font-size: 1.2rem;
  line-height: 1;
  text-decoration: none;
  transition:
    background-color 0.15s var(--bk-ease),
    transform 0.15s var(--bk-ease);
}

.bk-topbar-btn:active {
  background-color: rgba(255, 255, 255, 0.1);
  transform: scale(0.92);
}

/* ---- Bottom tab bar ----------------------------------------------------- */
/* Floats clear of the screen edges rather than sitting on them, so the
   content scrolls visibly underneath it and the bar reads as one more panel
   on the water instead of a chrome strip bolted to the bottom. */
.bk-tabbar {
  position: fixed;
  right: 0;
  bottom: 0;
  left: 0;
  z-index: 1030;
  display: flex;
  align-items: stretch;
  height: var(--bk-bottomnav-h);
  margin: 0 var(--bk-bottomnav-inset)
    calc(var(--bk-safe-bottom) + var(--bk-bottomnav-inset));
  border: 1px solid var(--bk-surface-edge);
  border-radius: var(--bk-r-xl);
  background-color: var(--bk-glass-strong);
  background-image: var(--bk-surface-fill);
  box-shadow: var(--bk-shadow-lg);
  backdrop-filter: blur(22px) saturate(140%);
  -webkit-backdrop-filter: blur(22px) saturate(140%);
}

.bk-tab {
  display: flex;
  flex: 1 1 0;
  flex-direction: column;
  gap: 0.15rem;
  align-items: center;
  justify-content: center;
  min-width: 0;
  color: var(--bk-on-dark-faint);
  text-decoration: none;
  transition: color 0.15s var(--bk-ease);
}

.bk-tab-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 46px;
  height: 28px;
  border-radius: var(--bk-r-sm);
  font-size: 1.15rem;
  line-height: 1;
  transition: background-color 0.2s var(--bk-ease);
}

.bk-tab-label {
  font-size: 0.62rem;
  font-weight: 600;
  letter-spacing: 0.02em;
  line-height: 1;
}

/* The active tab is marked three ways over: a filled pill, an aqua icon and
   a brighter label, so it never rests on colour alone. */
.bk-tab.is-active {
  color: var(--bk-aqua-bright);
}

.bk-tab.is-active .bk-tab-icon {
  background-color: rgba(52, 208, 216, 0.16);
}

.bk-tab:active .bk-tab-icon {
  transform: scale(0.9);
}
</style>
