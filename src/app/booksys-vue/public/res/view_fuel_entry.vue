<div class="modal fade text-left" id="view_fuel_entry_modal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
	<div class="modal-dialog" role="document">
        <div class="modal-content" id="view_fuel_entry_content">
            <div class="modal-header">
                <h4>
                    Fuel Entry
                </h4>
            </div>
            <form>
                <div class="modal-body" id="view_fuel_entry_dialog">               
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
                            <input v-model="entry.date_text" type="text" class="form-control form-control-white input-sm" placeholder="Date" disabled/>
                        </div>
                        <div class="col-sm-1 col-xs-1"></div>
                    </div>
                    <div class="row row-top-padded">
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-3 col-xs-3">
                            <div>
                                <label class="label-sm" for="">Driver</label>
                            </div>
                        </div>
                        <div class="cold-sm-7 col-xs-7">
                            <input v-model="entry.full_name" type="text" class="form-control form-control-white input-sm" placeholder="Name" disabled>
                        </div>
                        <div class="col-sm-1 col-xs-1"></div>
                    </div>
                    <div class="row row-top-padded">
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-3 col-xs-3">
                            <div>
                                <label class="label-sm" for="">Engine</label>
                            </div>
                        </div>
                        <div class="cold-sm-7 col-xs-7">
                            <div class="input-group">
                                <input v-model="entry.engine_hours" type="text" class="form-control form-control-white input-sm" placeholder="0.0">
                                <span class="input-group-addon">HRS</span>
                            </div>
                        </div>
                        <div class="col-sm-1 col-xs-1"></div>
                    </div>
                    <div class="row row-top-padded">
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-3 col-xs-3">
                            <div>
                                <label class="label-sm" for="">Fuel</label>
                            </div>
                        </div>
                        <div class="cold-sm-7 col-xs-7">
                            <div class="input-group">
                                <input v-model="entry.liters" type="text" class="form-control form-control-white input-sm" placeholder="0.0">
                                <span class="input-group-addon">L</span>
                            </div>
                        </div>
                        <div class="col-sm-1 col-xs-1"></div>
                    </div>
                    <div class="row row-top-padded"  v-if="entry.is_fuel_discount === true">
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-3 col-xs-3">
                            <div>
                                <label class="label-sm" for="">Cost (Gross)</label>
                            </div>
                        </div>
                        <div class="cold-sm-7 col-xs-7">
                            <div class="input-group">
                                <input v-model="entry.cost_brutto" type="text" class="form-control form-control-white input-sm" placeholder="0.0">
                                <span class="input-group-addon">{{currency}}</span>
                            </div>
                        </div>
                        <div class="col-sm-1 col-xs-1"></div>
                    </div>
                    <div class="row row-top-padded"  v-else>
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-3 col-xs-3">
                            <div>
                                <label class="label-sm" for="">Cost</label>
                            </div>
                        </div>
                        <div class="cold-sm-7 col-xs-7">
                            <div class="input-group">
                                <input v-model="entry.cost" type="text" class="form-control form-control-white input-sm" placeholder="0.0">
                                <span class="input-group-addon">{{currency}}</span>
                            </div>
                        </div>
                        <div class="col-sm-1 col-xs-1"></div>
                    </div>
                    <div class="row row-top-padded">
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-3 col-xs-3">
                            <div>
                                <label class="label-sm" for="">Discount</label>
                            </div>
                        </div>
                        <div class="col-sm-7 col-xs-7">
                            <input type="checkbox" v-model="entry.is_fuel_discount" name="isFuelDiscount" id="isFuelDiscount"/>
                        </div>
                    </div>
                    <div class="row row-top-padded" v-if="entry.is_fuel_discount === true">
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-3 col-xs-3">
                            <div>
                                <label class="label-sm" for="">Cost (Net)</label>
                            </div>
                        </div>
                        <div class="cold-sm-7 col-xs-7">
                            <div class="input-group">
                                <input v-model="entry.cost" type="text" class="form-control form-control-white input-sm" placeholder="0.0">
                                <span class="input-group-addon">{{currency}}</span>
                            </div>
                        </div>
                        <div class="col-sm-1 col-xs-1"></div>
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-default" v-on:click="save">
                        <span class="glyphicon glyphicon-save"></span>
                        Save
                    </button>
                    <button type="button" class="btn btn-default" v-on:click="cancel">
                        <span class="glyphicon glyphicon-remove"></span>
                        Cancel
                    </button>
                </div>
            </form>
        </div>
    </div>
</div>

