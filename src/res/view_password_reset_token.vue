<div id="view_password_reset_token">
    <div class="modal fade text-left" id="view_password_reset_token_modal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
        <div class="modal-dialog" role="document">
            <div class="modal-content" id="view_user_signup_content">
                <div class="modal-header">
                    <h4>
                        Set New Password
                    </h4>
                </div>
                <form>
                    <div class="modal-body" id="view_user_signup_dialog">               
                        <div v-if="reset.errors.length" class="row">
                            <div class="alert alert-danger alert-dismissible" role="alert">
                                <b>Please correct the following error(s):</b>
                                <ul>
                                    <li v-for="error in reset.errors">{{ error }}</li>
                                </ul>
                            </div>
                        </div>
                        <div>
                            <div class="alert alert-info">
                                Please enter the token that you received by email and a new password of your choice:
                            </div>
                            <!-- email -->
                            <div class="row">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Username</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7 padding-top-4px">
                                    {{ reset.email }}
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- token -->
                            <div class="row padding-top-16px">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Token</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="reset.token" type="text" class="form-control form-control-white input-sm" placeholder="your token"/>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- password -->
                            <div class="row padding-top-16px">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Password</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="reset.password" type="password" class="form-control form-control-white input-sm" placeholder="new password"/>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- password confirmation -->
                            <div class="row">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Confirm</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="reset.password2" type="password" class="form-control form-control-white input-sm" placeholder="repeat password"/>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                        </div>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-default" v-on:click="back">
                            <span class="glyphicon glyphicon-arrow-left"></span>
                            Back
                        </button>
                        <button type="button" class="btn btn-default" v-on:click="save">
                            <span class="glyphicon glyphicon-ok"></span>
                            Reset Password
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</div>

