<template>
  <div :id="date" class="text-center full-width">
  </div>
</template>

<script>
import Vue from 'vue';
import BooksysPie from "../libs/pie";
import moment from 'moment-timezone';

export default Vue.extend({
  name: 'Pie',
  mounted() {
    let el = document.getElementById(this.date);
    const pie = new BooksysPie();
    this.pie = pie;
    this.pieSessions = pie.drawPie(el, this.sessionData, this.selectHandler, this.properties);
  },
  methods: {
    selectHandler: function(selectedId) {
      this.$emit("selectHandler", this.pieSessions[selectedId]);
    },
    repaint: function(){
      let el = document.getElementById(this.date);
      el.innerHTML = '';      
      this.pieSessions = this.pie.drawPie(el, this.sessionData, this.selectHandler, this.properties);
    }
  },
  watch: {
    sessionData: function() {
      // repaint the pie upon any change
      this.repaint();
    },
    selectedSession: function(newSelectedSession) {
      if(newSelectedSession == null){
        BooksysPie.resetSelection();
        return;
      }

      if(newSelectedSession.id != null){
        for(let i=0; i<this.pieSessions.length; i++){
          if(this.pieSessions[i].id == newSelectedSession.id){
            BooksysPie.selectSector(i);
            break;
          }
        }
      }
    },
    properties: function() {
      this.repaint();
    }
  },
  props: ['sessionData', 'selectedSession', 'properties'],
  data() {
    return {
      pieSessions: []
    }
  },
  computed: {
    date: function(){
      if(this.sessionData != null){
        return moment(this.sessionData.window_start).format("YYYY-MM-DD");
      }else{
        return "unknown";
      }
    }
  }
})
</script>

<style scoped>
  .small-pie {
    height: 7em;
    width: 7em;
  }

  .full-width {
    width: 100%;
  }
</style>
