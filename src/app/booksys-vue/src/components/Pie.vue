<template>
  <div id="pie" class="text-center">
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
    },
    repaint: function(){
      // we have to re-draw the pie upon any change
      let el = document.getElementById("pie");
      el.innerHTML = '';      
      this.pieSessions = BooksysPie.drawPie(el, this.sessionData, this.selectHandler, this.properties);
    }
  },
  watch: {
    sessionData: function() {
      console.log("Change in sessionData, repaint");
      this.repaint();
    },
    selectedSession: function(newSelectedSession) {
      if(newSelectedSession == null){
        BooksysPie.resetSelection();
        return;
      }

      console.log(this.pieSessions);
      if(newSelectedSession.id != null){
        for(let i=0; i<this.pieSessions.length; i++){
          if(this.pieSessions[i].id == newSelectedSession.id){
            BooksysPie.selectSector(i)
          }
        }
      }
    },
    properties: function() {
      console.log("Change in properties, repaint");
      this.repaint();
    }
  },
  props: ['sessionData', 'selectedSession', 'properties'],
  data() {
    return {
      pieSessions: []
    }
  }
})
</script>

<style scoped>
  .small-pie {
    height: 7em;
    width: 7em;
  }
</style>
