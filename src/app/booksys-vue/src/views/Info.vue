<template>
<div class="display">
  <b-card no-body>
    <b-card-header>
      <h3>Infos and News</h3>
    </b-card-header>
    <b-card-body>
      <b-row>
        <div class="text-center" style="min-height: 90;">
        </div>
        <div class="text-center" id="location-address">
          <!-- here comes the address -->
        </div>
        <div class="text-center padding-top-16px">
          <div id="location-map">
            <!-- here comes the embedded map -->
          </div>
        </div>
      </b-row>
    </b-card-body>
    <b-card-footer class="text-right">
      <!-- <div class="bottom"> -->
        <b-button variant="outline-dark" to="/dashboard">
          <b-icon-house></b-icon-house>
          HOME
        </b-button>
      <!-- </div> -->
    </b-card-footer>
  </b-card>
</div>
</template>

