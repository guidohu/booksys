<template>
  <b-row>
    <b-col>
      <div class="main-title">
        <h3>
          {{ titleName }}
        </h3>
      </div>
    </b-col>
  </b-row>
</template>

<script>
import { BRow, BCol } from "bootstrap-vue";

export default {
  name: "MainTitle",
  props: ["titleName"],
  components: {
    BRow,
    BCol,
  },
};
</script>

<style>
div.main-title {
  font-family: Arial, Helvetica, sans-serif;
  color: #a1a1a1;
  border-bottom: 1px solid #eee;
  margin: 0px 0px 10px;
  padding-top: 10px;
  padding-bottom: 0px;
  padding-left: 15px;
}

div.main-color {
  color: #a1a1a1;
}
</style>
