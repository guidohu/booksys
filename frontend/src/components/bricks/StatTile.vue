<template>
  <div class="stat-tile" :class="[`tone-${tone}`, { 'stat-tile-hero': hero }]">
    <div class="stat-head">
      <span class="stat-icon"><i class="bi" :class="icon"></i></span>
      <span class="stat-label">{{ label }}</span>
    </div>
    <div class="stat-value" :class="{ 'stat-value-toned': toneValue }">
      {{ value }}<span v-if="unit" class="stat-unit">{{ unit }}</span>
    </div>
    <div v-if="context" class="stat-context">{{ context }}</div>
  </div>
</template>

<script setup>
defineProps({
  label: {
    type: String,
    required: true,
  },
  // Already formatted for display, the tile does no number formatting.
  value: {
    type: [String, Number],
    required: true,
  },
  unit: {
    type: String,
    required: false,
  },
  // Small note below the value, e.g. the period the value covers.
  context: {
    type: String,
    required: false,
  },
  icon: {
    type: String,
    required: false,
    default: "bi-dot",
  },
  // positive = money in, negative = money out, neutral = no direction.
  tone: {
    type: String,
    required: false,
    default: "neutral",
    validator: (v) => ["positive", "negative", "neutral"].includes(v),
  },
  // Opt-in for outcome figures (balance, profit) where the sign is the
  // point. The sign itself stays visible, so the color is never the only
  // channel carrying it.
  toneValue: {
    type: Boolean,
    required: false,
  },
  hero: {
    type: Boolean,
    required: false,
  },
});
</script>

<style scoped>
/* The tone colors only ever reinforce what the icon and the label already
   say, so they are never the sole carrier of meaning. The two that do carry
   direction (positive/negative) are separated far enough to survive colour
   vision deficiency, the neutral one is deliberately grey. */
.stat-tile {
  --tone: #6c757d;
  --tone-bg: rgba(108, 117, 125, 0.12);
  position: relative;
  height: 100%;
  padding: 0.65rem 0.75rem 0.6rem 0.85rem;
  background-color: #ffffff;
  border: 1px solid rgba(11, 11, 11, 0.1);
  border-radius: 0.5rem;
  overflow: hidden;
}

.tone-positive {
  --tone: #198754;
  --tone-bg: rgba(25, 135, 84, 0.12);
}

.tone-negative {
  --tone: #d6432a;
  --tone-bg: rgba(214, 67, 42, 0.12);
}

.tone-neutral {
  --tone: #6c757d;
  --tone-bg: rgba(108, 117, 125, 0.12);
}

/* Accent rail on the leading edge */
.stat-tile::before {
  content: "";
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  width: 3px;
  background-color: var(--tone);
}

.stat-head {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  min-width: 0;
  margin-bottom: 0.3rem;
}

.stat-icon {
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  width: 1.4rem;
  height: 1.4rem;
  border-radius: 0.35rem;
  background-color: var(--tone-bg);
  color: var(--tone);
  font-size: 0.8rem;
  line-height: 1;
}

/* Wraps rather than truncates: the label is the channel that carries the
   meaning, so it must stay readable on the narrowest phones. */
.stat-label {
  font-size: 0.7rem;
  font-weight: 600;
  letter-spacing: 0.04em;
  line-height: 1.2;
  color: #6c757d;
  text-transform: uppercase;
}

/* Proportional figures on purpose: tabular-nums makes a standalone number
   look loose at this size. */
.stat-value {
  color: #212529;
  font-size: 1.4rem;
  font-weight: 600;
  line-height: 1.15;
  white-space: nowrap;
}

.stat-value-toned {
  color: var(--tone);
}

.stat-unit {
  margin-left: 0.2rem;
  color: #6c757d;
  font-size: 0.7rem;
  font-weight: 500;
}

.stat-context {
  margin-top: 0.1rem;
  color: #6c757d;
  font-size: 0.7rem;
  font-weight: 400;
}

/* On mobile the tile joins the single panel material like everything else.
   The tone rails and icon chips keep their colours, which is the whole point
   of the component, they just sit on glass instead of white. */
@media (max-width: 992px) {
  .stat-tile {
    border: 1px solid var(--bk-surface-edge);
    border-radius: var(--bk-r-lg);
    background-color: transparent;
    background-image: var(--bk-surface-fill);
  }

  .stat-label,
  .stat-unit,
  .stat-context {
    color: var(--bk-on-dark-soft);
  }

  .stat-value {
    color: var(--bk-on-dark);
  }

  /* The muted grey tone disappears on the dark canvas. */
  .tone-neutral {
    --tone: var(--bk-on-dark-soft);
    --tone-bg: rgba(255, 255, 255, 0.1);
  }

  .tone-positive {
    --tone: var(--bk-success);
    --tone-bg: rgba(61, 220, 145, 0.16);
  }

  .tone-negative {
    --tone: var(--bk-danger);
    --tone-bg: rgba(255, 122, 133, 0.16);
  }
}

/* Keep long amounts inside the tile on narrow phones */
@media (max-width: 400px) {
  .stat-tile:not(.stat-tile-hero) .stat-value {
    font-size: 1.25rem;
  }
}

/* The one number the view leads with */
.stat-tile-hero {
  padding: 0.8rem 1rem 0.75rem 1.1rem;
}

.stat-tile-hero::before {
  width: 4px;
}

.stat-tile-hero .stat-value {
  font-size: 2.1rem;
  line-height: 1.1;
}

.stat-tile-hero .stat-unit {
  font-size: 0.85rem;
}

@media (min-width: 768px) {
  .stat-tile-hero .stat-value {
    font-size: 2.5rem;
  }
}
</style>
