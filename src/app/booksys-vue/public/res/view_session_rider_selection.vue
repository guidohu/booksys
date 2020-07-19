<div class="modal fade text-left" id="view_session_rider_selection_modal" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <h4>
                    <span class="glyphicon glyphicon-plus"></span>Add Rider
                </h4>
            </div>
            <div class="modal-body">
                <div class="row">
                <form id="user-form">
                    <div class="form-group">
                        <div class="col-sm-1 col-xs-1"></div>
                        <div class="col-sm-10 col-xs-10">
                            <label class="label-sm" for="session_title">Select Riders</label>
                            <div class="input-group">
                                <input v-on:keyup="filter" type="text" class="form-control form-control-white input-sm" id="search" placeholder="Search Name">
                                <span class="input-group-addon glyphicon glyphicon-search"></span>
                            </div>
                            <div v-if="filteredUsers.length>1">
                                <select class="form-control form-control-white" size="8" name="rider_users" multiple="multiple" title="Select riders (multi-select with Ctrl)" id="rider_users">
                                    <option v-for="user in filteredUsers" v-bind:value="user.id">{{user.first_name}} {{user.last_name}}</option>
                                </select>
                            </div>
                            <div v-if="filteredUsers.length==1">
                                <select class="form-control form-control-white" size="8" name="rider_users" multiple="multiple" title="Select riders (multi-select with Ctrl)" id="rider_users">
                                    <option v-for="user in filteredUsers" v-bind:value="user.id" :selected="true">{{user.first_name}} {{user.last_name}}</option>
                                </select>
                            </div>
                            (use ctrl+click to select multiple users)
                        </div>
                        <div class="clearfix"></div>									
                    </div>
                </form>
                </div>
            </div>
            <div class="modal-footer" id="status_footer">
                <button type="button" class="btn btn-default" v-on:click='cancel'>CANCEL</button>
                <button type="button" class="btn btn-default" v-on:click='add'>ADD</button>
            </div>
        </div>
    </div>
</div>