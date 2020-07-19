<div class="modal-header">
    <h4>
        User Group Details
    </h4>
</div>
<div class="modal-body" id="user_group">
    <form id="user_group_detail">
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
                    <label class="label-sm" for="">User Group</label>
                </div>
            </div>
            <div class="cold-sm-7 col-xs-7">
                <span v-if="editmode === false">
                    {{ group.user_group_name }}<br>
                    {{ group.user_group_description}}
                </span>
                <div v-if="editmode === true" class="form-group">
                    <label class="label-xs">Name</label>
                    <input v-model="group.user_group_name" type="text" class="form-control form-control-white input-sm" placeholder="User Group Name">
                    <label class="label-xs">Description</label>
                    <input v-model="group.user_group_description" type="text" class="form-control form-control-white input-sm" placeholder="User Group Description">
                </div>
            </div>
            <div class="col-sm-1 col-xs-1"></div>
        </div>
        <div class="row row-top-padded">
            <div class="col-sm-1 col-xs-1"></div>
            <div class="col-sm-3 col-xs-3">
                <div>
                    <label class="label-sm" for="">User Role</label>
                </div>
            </div>
            <div class="cold-sm-7 col-xs-7">
                <span v-if="editmode === false">
                    {{ group.user_role_name }}<br>
                    {{ group.user_role_description}}
                </span>
                <div v-if="editmode === true" class="form-group">
                    <label class="label-xs">Name</label>
                    <select class="form-control form-control-white" v-model="group.user_role_id" v-on:change="userRoleChange">
                        <option v-if="newGroup === true" value="" disabled selected>Please select the User Role</option>
                        <option v-for="role in userRoles" v-bind:value="role.user_role_id">{{role.user_role_name}}</option>
                    </select>
                    <label class="label-xs">Description</label>
                    <input v-model="group.user_role_description" type="text" class="form-control form-control-white input-sm" disabled>
                </div>
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
                <span v-if="editmode === false">
                    {{ group.price_min }}<br>
                    {{ group.price_description}}
                </span>
                <div v-if="editmode === true" class="form-group">
                    <label class="label-xs">Price {{ currency }}/min</label>
                    <div class="input-group">
                        <input v-model="group.price_min" type="text" class="form-control form-control-white input-sm" placeholder="0.00">
                        <span class="input-group-addon">{{ currency }}</span>
                    </div>
                    <label class="label-xs">Description</label>
                    <input v-model="group.price_description" type="text" class="form-control form-control-white input-sm" placeholder="Price Description">    
                </div>
            </div>
            <div class="col-sm-1 col-xs-1"></div>
        </div>
        <div class="modal-footer">
            <button v-if="editmode === true" type="button" class="btn btn-default" v-on:click="save">
                <span class="glyphicon glyphicon-save"></span>
                Save
            </button>
            <button v-if="editmode === true" type="button" class="btn btn-default" v-on:click="cancel">
                <span class="glyphicon glyphicon-remove"></span>
                Cancel
            </button>
            <button v-else-if="editmode === false" type="button" class="btn btn-default" v-on:click="cancel">
                <span class="glyphicon glyphicon-ok"></span>
                OK
            </button>
        </div>
    </form>
</div>

