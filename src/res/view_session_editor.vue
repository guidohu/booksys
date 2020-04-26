<!-- The create session window -->
<div class="modal fade text-left" id="view_session_editor_modal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content" id="status_modal_content">
            <div v-if="session.id==null" class="modal-header">
                <h4>
                    <span class="glyphicon glyphicon-plus"></span>Create Session
                </h4>
            </div>
            <div v-else class="modal-header">
                <h4>
                    <span class="glyphicon glyphicon-pencil"></span>Edit Session
                </h4>
            </div>
            <div class="modal-body">
                <div v-if="errors.length" class="row">
                    <div class="col-sm-1 col-xs-1"></div>
                    <div class="col-sm-10 col-xs-10">
                        <div class="alert alert-danger alert-dismissible" role="alert">
                            <b>Please correct the following error(s):</b>
                            <ul>
                                <li v-for="error in errors">{{ error }}</li>
                            </ul>
                        </div>
                    </div>
                    <div class="col-sm-1 col-xs-1"></div>
                </div>
                <div class="row">
                <form id="edit_session_form">
                    <div class="form-group">
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-10 col-xs-10">
                            <label class="label-sm" for="edit_session_title">Title</label>
                            <input type="text" v-model="session.title" class="form-control form-control-white input-sm" name="edit_session_title" id="edit_session_title" placeholder="Title"></input>
                        </div>
                        <div class="clearfix"></div>
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-10 col-xs-10">
                            <label class="label-sm" for="edit_session_description">Description</label>
                            <textarea rows="3" v-model="session.description" class="form-control form-control-white input-sm" name="edit_session_description" id="edit_session_description" placeholder="More detailed description"></textarea>
                        </div>
                        <div class="clearfix"></div>
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-10 col-xs-10">
                            <label class="label-sm" for="edit_session_start">Start</label>
                            <input type='text' class="form-control form-control-white input-sm" id='startPicker'></input>
                        </div>
                        <div class="clearfix"></div>
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-10 col-xs-10">
                            <label class="label-sm" for="edit_session_end">End</label>
                            <input type='text' class="form-control form-control-white input-sm" id='endPicker'></input>
                        </div>
                        <div class="clearfix"></div>
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-5 col-xs-5">
                            <label class="label-sm" for="edit_session_max_riders">Maximum Riders</label>
                        </div>
                        <div class="clearfix"></div>
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-5 col-xs-5">
                            <input type="number" v-model="session.maxRiders" class="form-control form-control-white input-sm" name="edit_session_max_riders" id="edit_session_max_riders"></input>
                        </div>
                        <div class="clearfix"></div>
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-10 col-xs-10">
                            <label class="label-sm" for="edit_session_type">Type</label>
                        </div>
                        <div class="clearfix"></div>
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-5 col-xs-5">
                            <input type="checkbox" v-model="session.isCourse" id="isCourse"></input>
                        </div>
                    </div>
                </form>
                </div>
            </div>
            <div class="modal-footer" id="edit_session_footer">
                <button type="button" class="btn btn-default" v-on:click="cancel">CANCEL</button>
                <button type="button" class="btn btn-default" v-on:click="save">SAVE</button>
            </div>
        </div>
    </div>
</div>