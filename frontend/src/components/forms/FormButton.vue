<template>
  <div class="row text-right mt-2">
    <div class="col-9 offset-3">
      <div class="d-grid gap-2">
        <button
          :class="getClass"
          :type="type ? type : 'button'"
          @click="$emit('click')"
        >
          <slot></slot>
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "FormButton",
  props: ["type", "btnStyle", "btnSize"],
  emits: ["click"],
  computed: {
    getClass: function() {
      const style = (this.btnStyle == null) ? "btn-outline-info" : "btn-outline-" + this.btnStyle;
      const size = (this.btnSize == "small") ? "btn-sm" : (this.btnSize == "large") ? "btn-lg" : ""; 
      return "btn " + style + " " + size;
    }
  }
}
</script>