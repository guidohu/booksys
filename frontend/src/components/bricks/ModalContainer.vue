<template>
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
</template>

<script setup>
import { Modal } from "bootstrap";
import { onMounted, onUnmounted, watch, ref, nextTick } from "vue";

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
  if(modalRef.value != null) {
    modal = new Modal(modalRef.value);
    if (props.visible) {
      modal.show();
    }
  } else {
    console.error('Modal element (modalRef) not found on mount!');
  }
});

onUnmounted(() => {
  modal.hide();
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
      modal.hide();
    }
  }, { immediate: true }
);
</script>
