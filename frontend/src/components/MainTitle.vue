<template>
  <div class="row">
    <div class="col-12">
      <div class="main-title">
        <h3>
          {{ titleName }}
        </h3>
      </div>
    </div>
  </div>
</template>

<script>

export default {
  name: "MainTitle",
  props: ["titleName"],
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
