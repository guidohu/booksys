<template>
  <div :id="pieElementId" class="text-center full-width">
  </div>
</template>

<script>
import BooksysPie from "../libs/pie";
import moment from 'moment-timezone';

export default {
  name: 'Pie',
  mounted() {
    let el = document.getElementById(this.pieElementId);
    const pie = new BooksysPie();
    this.pie = pie;
    this.pieSessions = pie.drawPie(el, this.sessionData, this.selectHandler, this.properties);
  },
  methods: {
    selectHandler: function(selectedId) {
      this.$emit("selectHandler", this.pieSessions[selectedId]);
    },
    repaint: function(){
      let el = document.getElementById(this.pieElementId);
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
        this.pie.resetSelection();
        return;
      }

      if(newSelectedSession.id != null){
        for(let i=0; i<this.pieSessions.length; i++){
          if(this.pieSessions[i].id == newSelectedSession.id){
            this.pie.selectSector(i);
            break;
          }
        }
      }
    },
    properties: function() {
      this.repaint();
    }
  },
  props: ['sessionData', 'selectedSession', 'properties', 'pieId'],
  data() {
    return {
      pieSessions: [],
      pieElementId: "someId"
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
  },
  created() {
    if(this.pieId != null){
      this.pieElementId = "pie" + this.pieId;
    }
  }
}
</script>

<style scoped>
  /* .small-pie {
    height: 7em;
    width: 7em;
  } */

  .full-width {
    width: 100%;
  }
</style>
