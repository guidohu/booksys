<template>
  <!--
    One cell of the mobile dashboard's bento grid. Replaces the fixed 120x120
    square button: a tile sizes itself from the grid instead of from a hard
    pixel value, so the same markup works on a 320px phone and a tablet.
  -->
  <router-link :to="to" class="bk-tile" :class="{ 'is-wide': wide }">
    <span class="bk-tile-icon">
      <i class="bi" :class="icon" aria-hidden="true"></i>
    </span>
    <span class="bk-tile-text">
      <span class="bk-tile-label">{{ label }}</span>
      <span v-if="hint" class="bk-tile-hint">{{ hint }}</span>
    </span>
    <i class="bi bi-chevron-right bk-tile-go" aria-hidden="true"></i>
  </router-link>
</template>

<script setup>
defineProps({
  to: {
    type: String,
    required: true,
  },
  label: {
    type: String,
    required: true,
  },
  icon: {
    type: String,
    required: false,
    default: "bi-dot",
  },
  // Second line under the label, for tiles that can say something useful
  // about what is behind them.
  hint: {
    type: String,
    required: false,
  },
  // Spans the full grid width. Used for the one destination a screen leads
  // with, which is what gives the grid its bento rhythm.
  wide: {
    type: Boolean,
    required: false,
  },
});
</script>

<style scoped>
.bk-tile {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
  align-items: flex-start;
  justify-content: flex-start;
  min-height: 108px;
  padding: 0.85rem;
  border: 1px solid var(--bk-surface-edge);
  border-radius: var(--bk-r-xl);
  /* Exactly the material a content .card uses, from the same tokens: a menu
     tile and a card full of data are the same object on the same water. */
  background-image: var(--bk-surface-fill);
  color: var(--bk-on-dark);
  text-decoration: none;
  transition:
    transform 0.16s var(--bk-ease),
    background-color 0.16s var(--bk-ease),
    border-color 0.16s var(--bk-ease);
}

.bk-tile:active {
  transform: scale(0.975);
  border-color: rgba(52, 208, 216, 0.4);
  background-color: rgba(52, 208, 216, 0.08);
}

/* The lead tile lies down instead of stacking, so the row reads as one wide
   destination rather than a stretched square. */
.bk-tile.is-wide {
  flex-direction: row;
  gap: 0.85rem;
  align-items: center;
  min-height: 84px;
  grid-column: 1 / -1;
}

.bk-tile-icon {
  display: inline-flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  min-width: 42px;
  height: 42px;
  padding: 0 0.3rem;
  border-radius: var(--bk-r-md);
  background-color: rgba(52, 208, 216, 0.14);
  color: var(--bk-aqua-bright);
  font-size: 1.35rem;
  line-height: 1;
}

.bk-tile-text {
  display: flex;
  flex: 1 1 auto;
  flex-direction: column;
  gap: 0.1rem;
  min-width: 0;
}

.bk-tile-label {
  font-size: 0.95rem;
  font-weight: 650;
  letter-spacing: -0.01em;
  line-height: 1.2;
}

.bk-tile-hint {
  color: var(--bk-on-dark-soft);
  font-size: 0.72rem;
  font-weight: 400;
  line-height: 1.25;
}

/* On the narrowest phones the bento grid collapses to a single column,
   where a stacked tile is mostly empty space. Every tile takes the compact
   row form there so the dashboard stays one screen rather than four. */
@media (max-width: 340px) {
  .bk-tile {
    flex-direction: row;
    gap: 0.85rem;
    align-items: center;
    min-height: 72px;
  }

  .bk-tile .bk-tile-go {
    display: inline-block;
  }
}

/* Only the wide tile gets the chevron: on a square tile the whole card is so
   obviously the target that the arrow is just noise. */
.bk-tile-go {
  display: none;
  flex: 0 0 auto;
  color: var(--bk-on-dark-faint);
  font-size: 1rem;
}

.bk-tile.is-wide .bk-tile-go {
  display: inline-block;
}
</style>
