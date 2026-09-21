<template>
  <!--
    Teleported to the #dialog container index.html provides. A modal is
    declared inside whichever component owns it, which is normally inside a
    card; left there, it inherits that card's stacking and clipping context
    and can end up positioned against the card instead of the viewport.
  -->
  <Teleport to="#dialog">
    <div
      class="modal fade"
      ref="modalRef"
      :id="name"
      data-bs-target="static"
      aria-hidden="false"
      tabindex="-1"
    >
      <div class="modal-dialog">
        <div class="modal-content">
          <slot></slot>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { Modal } from "bootstrap";
import { onMounted, onUnmounted, watch, ref } from "vue";

const props = defineProps({
  visible: {
    type: Boolean,
    required: true,
  },
  name: {
    type: String,
    required: true,
  },
});

const modalRef = ref(null);
let modal = null;

onMounted(() => {
  if (modalRef.value != null) {
    modal = new Modal(modalRef.value);
    if (props.visible) {
      modal.show();
    }
  } else {
    console.error("Modal element (modalRef) not found on mount!");
  }
});

onUnmounted(() => {
  if (modal == null) {
    return;
  }
  modal.dispose();
});

watch(
  () => props.visible,
  (newValue) => {
    if (newValue) {
      if (modal != null) {
        modal.show();
      }
    } else {
      if (modal != null) {
        modal.hide();
      }
    }
  },
  { immediate: true },
);
</script>
