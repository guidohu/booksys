<div class="modal fade text-left" id="view_database_update_modal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
	<div class="modal-dialog" role="document">
        <div class="modal-content" id="view_database_update_content">
            <div class="modal-header">
                <h4>
                    Database Update
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
                                <label class="label-sm" for="">App Version</label>
                            </div>
                        </div>
                        <div class="cold-sm-7 col-xs-7">
                            v{{ info.app_version }}
                        </div>
                        <div class="col-sm-1 col-xs-1"></div>
                    </div>
                    <div class="row row-top-padded">
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-3 col-xs-3">
                            <div>
                                <label class="label-sm" for="">DB Version</label>
                            </div>
                        </div>
                        <div class="cold-sm-7 col-xs-7" v-if="info.db_version != null">
                            v{{ info.db_version}}
                        </div>
                        <div class="cold-sm-7 col-xs-7" v-else>
                            outdated
                        </div>
                        <div class="col-sm-1 col-xs-1"></div>
                    </div>
                    <div class="row row-top-padded">
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-10 col-xs-10">
                            <div class="alert alert-info" v-if="info.isUpdated == false && info.isUpdating == false && info.ok == true">
                                The database version is outdated and needs to be upgraded. If the upgrade is not done, there might occur data inconsitencies and the app will stop working properly. Thus it is highly recommended to upgrade now.
                            </div>
                            <div class="alert alert-success" v-if="info.isUpdated == false && info.isUpdating == true">
                                Updating database. Please wait, this might take a few minutes.
                            </div>
                            <div class="alert alert-success" v-if="info.isUpdated == true && info.isUpdating == false && info.ok == true">
                                Update successful.
                            </div>
                            <div class="alert alert-danger" v-if="info.isUpdating == false && info.ok == false">
                                Update failed. Please let the administrator know about this, together with the error message below:<br>
                                <strong>Error:</strong> {{ info.message }}
                            </div>
                        </div>
                        <div class="col-sm-1 col-xs-1"></div>
                    </div>
                    <div class="row row-top-padded" v-if="info.queries != null && info.queries.length > 0 && info.ok == false">
                        <div class="row">
                            <div class="col-sm-1 col-xs-1"></div>
                            <div class="col-sm-10 col-xs-10">
                                <div class="panel panel-default">
                                    <div class="panel-body">
                                        <div class="row" v-for="query in info.queries">
                                            <div class="col-sm-12 col-xs-12">
                                                <span class="glyphicon glyphicon-ok" v-if="query.ok == true"></span>
                                                <span class="glyphicon glyphicon-remove" v-if="query.ok == false"></span>
                                                {{ query.query }}
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div class="col-sm-1 col-xs-1"></div>
                        </div>
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-default" v-on:click="update" v-if="info.isUpdated == false && info.isUpdating == false">
                        <span class="glyphicon glyphicon-refresh"></span>
                        Update
                    </button>
                    <button type="button" class="btn btn-default" disabled v-if="info.isUpdated == false && info.isUpdating == true">
                        <span class="glyphicon glyphicon-refresh"></span>
                        Updating
                    </button>
                    <button type="button" class="btn btn-default" v-on:click="cancel" v-if="info.isUpdated == true && info.isUpdating == false">
                        <span class="glyphicon glyphicon-ok"></span>
                        Done
                    </button>
                    <button type="button" class="btn btn-default" v-on:click="cancel" v-if="info.isUpdated == false && info.isUpdating == false">
                        <span class="glyphicon glyphicon-remove"></span>
                        Ignore
                    </button>
                </div>
            </form>
        </div>
    </div>
</div>

