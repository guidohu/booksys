<template>
  <div id="test">
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
    console.log("props", this.$props)
    console.log("props2", this.sessionData)
    let el = document.getElementById("test")
    console.log("pass data", this.$props.sessionData)
    BooksysPie.drawPie(el, this.sessionData, null, this.properties)
  },
  props: ['sessionData', 'properties']
})
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
  .small-pie {
    height: 7em;
    width: 7em;
  }
</style>
