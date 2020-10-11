<template>
  <div id="pie">
  </div>
</template>

<script>
import Vue from 'vue'
import BooksysPie from "../libs/pie"

export default Vue.extend({
  name: 'Pie',
  mounted() {
    let el = document.getElementById("pie")
    this.pieSessions = BooksysPie.drawPie(el, this.sessionData, this.selectHandler, this.properties);
  },
  beforeUpdate() {
    console.log("Pie would update");
  },
  methods: {
    selectHandler: function(selectedId) {
      this.$emit("selectHandler", this.pieSessions[selectedId]);
    }
  },
  watch: {
    sessionData: function() {
      // we have to re-draw the pie upon any change
      let el = document.getElementById("pie");
      el.innerHTML = '';      
      this.pieSessions = BooksysPie.drawPie(el, this.sessionData, this.selectHandler, this.properties);
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
