<div class="modal fade text-left" id="view_database_setup_modal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
	<div class="modal-dialog" role="document">
        <div class="modal-content" id="view_database_setup_content">
            <div class="modal-header">
                <h4>
                    Database Setup
                </h4>
            </div>
            <form>
                <div class="modal-body" id="view_database_setup_dialog">               
                    <div v-if="errors.length" class="row">
                        <div class="alert alert-danger alert-dismissible" role="alert">
                            <b>Please correct the following error(s):</b>
                            <ul>
                                <li v-for="error in errors">{{ error }}</li>
                            </ul>
                        </div>
                    </div>
                    <div class="alert alert-info">
                        Please setup a new MySQL database and user first. The database schema will be applied and all prior data will be deleted if there is any. Enter the database connection details below:
                    </div>
                    <div class="row">
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-3 col-xs-3">
                            <div>
                                <label class="label-sm" for="">Host</label>
                            </div>
                        </div>
                        <div class="cold-sm-7 col-xs-7">
                            <input v-model="dbConfig.host" type="text" class="form-control form-control-white input-sm" placeholder="192.168.1.12:1234"/>
                        </div>
                        <div class="col-sm-1 col-xs-1"></div>
                    </div>
                    <div class="row row-top-padded">
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-3 col-xs-3">
                            <div>
                                <label class="label-sm" for="">Database Name</label>
                            </div>
                        </div>
                        <div class="cold-sm-7 col-xs-7">
                            <input v-model="dbConfig.name" type="text" class="form-control form-control-white input-sm" placeholder="myDB">
                        </div>
                        <div class="col-sm-1 col-xs-1"></div>
                    </div>
                    <div class="row row-top-padded">
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-3 col-xs-3">
                            <div>
                                <label class="label-sm" for="">User</label>
                            </div>
                        </div>
                        <div class="cold-sm-7 col-xs-7">
                            <input v-model="dbConfig.user" type="text" class="form-control form-control-white input-sm" placeholder="user">
                        </div>
                        <div class="col-sm-1 col-xs-1"></div>
                    </div>
                    <div class="row row-top-padded">
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-3 col-xs-3">
                            <div>
                                <label class="label-sm" for="">Password</label>
                            </div>
                        </div>
                        <div class="cold-sm-7 col-xs-7">
                            <input v-model="dbConfig.password" type="password" class="form-control form-control-white input-sm" placeholder="mySecretPassword">
                        </div>
                        <div class="col-sm-1 col-xs-1"></div>
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-default" v-on:click="save">
                        <span class="glyphicon glyphicon-wrench"></span>
                        Setup
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

