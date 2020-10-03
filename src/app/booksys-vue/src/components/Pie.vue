<template>
  <div id="pie">
  </div>
</template>

<script>
import Vue from 'vue'
import BooksysPie from "../libs/pie"

export default Vue.extend({
  name: 'Pie',
  components: {
  },
  mounted() {
    let el = document.getElementById("pie")
    this.pieSessions = BooksysPie.drawPie(el, this.sessionData, this.selectHandler, this.properties);
  },
  methods: {
    selectHandler: function(selectedId) {
      this.$emit("selectHandler", this.pieSessions[selectedId]);
    }
  },
  props: ['sessionData', 'properties']
})
</script>

<style scoped>
  .small-pie {
    height: 7em;
    width: 7em;
  }
</style>
