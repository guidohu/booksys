<div id="view_stop_watch">
    <div class="panel panel-default panel-padded">
        <div class="panel-body">
            <div v-if="watch.heatQueue>0" class="row">
                <div class="col-sm-12 col-xs-12">
                    <div class="alert alert-warning">
                        {{watch.heatQueue}} heats waiting to be synced with the server.
                    </div>
                </div>
            </div>
            <div class="row">
                <div class="col-sm-12 col-xs-12">
                        <div class="col-sm-12 col-xs-12 no-padding">
                            <div class="well well-lg text-center text-bigger">
                                {{watch.display}}
                            </div>
                        </div>
                </div>
            </div>
            <div class="row">
                <div class="col-sm-12 col-xs-12">
                    <select v-model="selected" class="form-control text-big line-high" placeholder=".col-sm-10 .col-xs-10" @change="onChange" id="watchRiderSelect" :disabled="watch.running==true">
                        <option value="-1">please select rider</option>
                        <option v-for="user in users" :value="user.id">{{user.first_name}} {{user.last_name}}</option>
                    </select>	
                </div>
                <div v-if="watch.comment != null && watch.comment != ''" class="col-sm-12 col-xs-12">
                    <div class="well well-sm overflow-hidden" v-on:click="addComment">
                        {{watch.comment}}
                    </div>
                </div>
                <div class="col-sm-12 col-xs-12 padding-top-8px">
                    <div class="btn-group btn-group-lg">
                        <button v-if="watch.running==false && watch.isPaused==false" type="button" class="btn btn-default stop-watch-button" v-on:click="back">
                            <span class="glyphicon glyphicon-menu-left"></span>BACK
                        </button>
                        <button v-if="selected>-1 && watch.running==false && watch.isPaused==false" type="button" class="btn btn-default stop-watch-button" v-on:click="start">
                            <span class="glyphicon glyphicon-play"></span>START
                        </button>
                        <button v-if="selected==-1 && watch.running==false && watch.isPaused==false" type="button" class="btn btn-default stop-watch-button" v-on:click="start" disabled="disabled">
                            <span class="glyphicon glyphicon-play"></span>START
                        </button>
                        <button v-if="watch.running==true && watch.isPaused==false" type="button" class="btn btn-default stop-watch-button" v-on:click="pause">
                            <span class="glyphicon glyphicon-pause"></span>PAUSE
                        </button>
                        <button v-if="watch.isPaused==true" type="button" class="btn btn-default stop-watch-button" v-on:click="resume">
                            <span class="glyphicon glyphicon-play"></span>RESUME
                        </button>
                        <button v-if="watch.running==true" type="button" class="btn btn-default stop-watch-button" v-on:click="finish">
                            <span class="glyphicon glyphicon-ok"></span>FINISH
                        </button>
                        <button v-if="selected>-1 || watch.running==true" type="button" class="btn btn-default stop-watch-button dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                            <span class="glyphicon glyphicon-option-vertical"></span> 
                        </button>
                        <ul class="dropdown-menu pull-right">
                            <li v-if="selected>-1 && watch.comment == null || watch.comment == ''">
                                <a href="#" v-on:click="addComment">Add Comment</a>
                            </li>
                            <li v-if="selected>-1 && watch.comment != null && watch.comment != ''">
                                <a href="#" v-on:click="addComment">Edit Comment</a>
                            </li>
                            <li v-if="watch.running==true" role="separator" class="divider"></li>
                            <li v-if="watch.running==true">
                                <a href="#">Cancel Heat</a>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <div class="modal fade" role="dialog" id="addUserModal">
        <div class="modal-dialog" role="document">
            <div class="modal-content">
                <div class="modal-header">
                    <h4 class="modal-title">Heat Comment</h4>
                </div>
                <div class="modal-body">
                    <div class="row">
                        <div class="col-sm-12 col-xs-12">
                            <form id="watch_comment_form">
                                <input v-model="watch.comment" type="text" class="form-control form-control-white input-sm" tabindex="0" placeholder="Add your comment here..." id="watchCommentInput" name="watchCommentInput"></input>
                            </form>
                        </div>
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-danger" data-dismiss="modal" v-on:click="removeComment">Remove</button>
                    <button type="button" class="btn btn-default" data-dismiss="modal" v-on:click="saveComment">Save</button>
                </div>
            </div><!-- /.modal-content -->
        </div><!-- /.modal-dialog -->
    </div><!-- /.modal -->
</div>