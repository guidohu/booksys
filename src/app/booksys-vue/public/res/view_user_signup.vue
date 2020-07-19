<div id="view_user_signup_modal">
    <div class="modal fade text-left" id="view_user_signup_created_modal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
        <div class="modal-dialog" role="document">
            <div class="modal-content" id="view_user_signup_content">
                <div class="modal-header">
                    Sign Up Successful
                </div>
                <div class="modal-body" id="view_user_signup_dialog">
                    Your account has been created successfully. We will activate it soon and let you know, when your account is ready for your first login.
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-default" v-on:click="done">
                        <span class="glyphicon glyphicon-ok"></span>
                        Done
                    </button>
                </div>
            </div>
        </div>
    </div>
    <div class="modal fade text-left" id="view_user_signup_form_modal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
        <div class="modal-dialog" role="document">
            <div class="modal-content" id="view_user_signup_content">
                <div class="modal-header">
                    <h4>
                        Sign Up
                    </h4>
                </div>
                <form>
                    <div class="modal-body" id="view_user_signup_dialog">               
                        <div v-if="errors.length" class="row">
                            <div class="alert alert-danger alert-dismissible" role="alert">
                                <b>Please correct the following error(s):</b>
                                <ul>
                                    <li v-for="error in errors">{{ error }}</li>
                                </ul>
                            </div>
                        </div>
                        <div>
                            <div class="alert alert-info">
                                Please fill out the form below to create an account:
                            </div>
                            <!-- username / email -->
                            <div class="row">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Email</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.email" type="text" class="form-control form-control-white input-sm" placeholder=""/>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- password -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Password</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.password" type="password" class="form-control form-control-white input-sm" placeholder="">
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- password confirm -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Confirm</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.passwordConfirm" type="password" class="form-control form-control-white input-sm" placeholder="">
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <hr>
                            <!-- first name -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">First Name</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.firstName" type="text" class="form-control form-control-white input-sm" placeholder="">
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- last name -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Last Name</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.lastName" type="text" class="form-control form-control-white input-sm" placeholder="">
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- street/no -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Street / No</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.address" type="text" class="form-control form-control-white input-sm" placeholder="">
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- zip -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Zip Code</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.zipCode" type="text" class="form-control form-control-white input-sm" placeholder="">
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- place -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Place / City</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.city" type="text" class="form-control form-control-white input-sm" placeholder="">
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- phone -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Phone</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.phone" type="text" class="form-control form-control-white input-sm" placeholder="">
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <hr>
                            <!-- driving license -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-10 col-xs-10">
                                    <label for="license">
                                        <input v-model="user.driving_license" type="checkbox"/>
                                        I have a Driving License for Motorboats
                                    </label>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- no school -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-10 col-xs-10">
                                    <label for="school">
                                        <input v-model="user.accepted_no_school" type="checkbox"/>
                                        This is no Wakeboard School. If you are looking for a School please see www.wakeandski.ch, www.ceccotorenas.ch, ...
                                    </label>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- terms and conditions -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-10 col-xs-10">
                                    <label for="license">
                                        <input v-model="user.accepted_terms_and_conditions" type="checkbox"/>
                                        I know the conditions. Most important: Any activity takes place at your OWN RISK!
                                    </label>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- captcha -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-10 col-xs-10">
                                    <div id="captcha"></div>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                        </div>
                        <div v-else>
                            <div class="alert alert-warning">
                                Something went wrong.
                            </div>
                        </div>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-default" v-on:click="cancel">
                            <span class="glyphicon glyphicon-remove"></span>
                            Cancel
                        </button>
                        <button type="button" class="btn btn-default" v-on:click="save">
                            <span class="glyphicon glyphicon-check"></span>
                            Sign-up
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</div>

