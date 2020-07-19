<div class="modal fade text-left" id="view_heat_entry_modal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
	<div class="modal-dialog" role="document">
        <div class="modal-content" id="view_heat_entry_content">
            <div class="modal-header">
                <button type="button" v-on:click="cancel" class="close"><span aria-hidden="true">&times;</span></button>
                <h4>
                    Heat Entry
                </h4>
            </div>
            <form>
                <div class="modal-body" id="view_heat_entry_dialog">               
                    <div v-if="errors.length" class="row">
                        <div class="alert alert-danger alert-dismissible" role="alert">
                            <b>Please correct the following error(s):</b>
                            <ul>
                                <li v-for="error in errors">{{ error }}</li>
                            </ul>
                        </div>
                    </div>
                    <div class="row">
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-3 col-xs-3">
                            <div>
                                <label class="label-sm" for="">Date</label>
                            </div>
                        </div>
                        <div class="cold-sm-7 col-xs-7">
                            <input v-model="heat.date_text" type="text" class="form-control form-control-white input-sm" placeholder="Date" disabled/>
                        </div>
                        <div class="col-sm-1 col-xs-1"></div>
                    </div>
                    <div class="row row-top-padded">
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-3 col-xs-3">
                            <div>
                                <label class="label-sm" for="">Rider</label>
                            </div>
                        </div>
                        <div class="cold-sm-7 col-xs-7">
                            <input v-model="heat.full_name" type="text" class="form-control form-control-white input-sm" placeholder="Name" disabled>
                        </div>
                        <div class="col-sm-1 col-xs-1"></div>
                    </div>
                    <div class="row row-top-padded">
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-3 col-xs-3">
                            <div>
                                <label class="label-sm" for="">Price</label>
                            </div>
                        </div>
                        <div class="cold-sm-7 col-xs-7">
                            <div class="input-group">
                                <input v-model="heat.price_per_min_text" type="text" class="form-control form-control-white input-sm" placeholder="0.00" disabled>
                                <span class="input-group-addon">CHF/min</span>
                            </div>
                        </div>
                        <div class="col-sm-1 col-xs-1"></div>
                    </div>
                    <div class="row row-top-padded">
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-3 col-xs-3">
                            <div>
                                <label class="label-sm" for="">Duration</label>
                            </div>
                        </div>
                        <div class="cold-sm-7 col-xs-7">
                            <div class="input-group">
                                <input v-model="heat.duration_text" v-on:input="durationChange()" type="text" class="form-control form-control-white input-sm" placeholder="0:00">
                                <span class="input-group-addon">mm:ss</span>
                            </div>
                        </div>
                        <div class="col-sm-1 col-xs-1"></div>
                    </div>
                    <div class="row row-top-padded">
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-3 col-xs-3">
                            <div>
                                <label class="label-sm" for="">Cost</label>
                            </div>
                        </div>
                        <div class="cold-sm-7 col-xs-7">
                            <div class="input-group">
                                <input v-model="heat.cost_text" type="text" class="form-control form-control-white input-sm" placeholder="0.00" disabled>
                                <span class="input-group-addon">{{currency}}</span>
                            </div>
                        </div>
                        <div class="col-sm-1 col-xs-1"></div>
                    </div>
                </div>
                <div class="modal-footer">
                    <div class="row row-top-padded">
                        <div class="col-sm-1 hidden-xs"></div>
                        <div class="col-sm-3 col-xs-3 text-left">
                            <button type="button" class="btn btn-danger" v-on:click="remove">
                                <span class="glyphicon glyphicon-trash"></span>
                                Delete
                            </button>
                        </div>
                        <div class="col-sm-7 col-xs-9">
                            <button v-if="showSave==true" type="button" class="btn btn-info" v-on:click="save">
                                <span class="glyphicon glyphicon-save"></span>
                                Save
                            </button>
                            <button type="button" class="btn btn-default" v-on:click="cancel">
                                <span class="glyphicon glyphicon-remove"></span>
                                Cancel
                            </button>
                        </div>
                    </div>
                </div>
            </form>
        </div>
    </div>
</div>

