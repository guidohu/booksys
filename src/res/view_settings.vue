<div id="view_settings_modal">
    <div class="modal fade text-left" id="view_settings_created_modal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
        <div class="modal-dialog" role="document">
            <div class="modal-content" id="view_settings_signup_content">
                <div class="modal-header">
                    Sign Up Successful
                </div>
                <div class="modal-body" id="view_settings_dialog">
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
    <div class="modal fade text-left" id="view_settings_form_modal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel" aria-hidden="true">
        <div class="modal-dialog" role="document">
            <div class="modal-content" id="view_settings_content">
                <div class="modal-header">
                    <h4>
                        Settings
                    </h4>
                </div>
                <form>
                    <div class="modal-body" id="view_settings_dialog">               
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
                                Timezone and location settings. Users will see the times in the timezone defined by these settings, also sunrise and sunset will be calculated based on these settings.
                            </div>
                            <!-- Timezone -->
                            <div class="row">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Timezone</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.email" type="text" class="form-control form-control-white input-sm" placeholder="Europe/Berlin"/>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- Longitude -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Longitude</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.email" type="text" class="form-control form-control-white input-sm" placeholder="8.542939"/>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- Latitude -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Latitude</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.email" type="text" class="form-control form-control-white input-sm" placeholder="47.367658"/>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- Location Map -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Map (Google Maps IFrame)</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.email" type="text" class="form-control form-control-white input-sm" placeholder=""/>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- Location Address -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Address</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.email" type="text" class="form-control form-control-white input-sm" placeholder="Street 1, 8000 Zurich"/>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <div class="alert alert-info margin-top-16px">
                                Payment Information. This information is displayed to users in case they want to top-up their balance.
                            </div>
                            <!-- Currency -->
                            <div class="row">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Currency</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.email" type="text" class="form-control form-control-white input-sm" placeholder="CHF"/>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- Payment Account Owner -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Payment Account Owner</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.email" type="text" class="form-control form-control-white input-sm" placeholder="John Doe"/>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- Payment Account IBAN -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">IBAN</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.email" type="text" class="form-control form-control-white input-sm" placeholder=""/>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- Payment Account BIC -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">BIC</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.email" type="text" class="form-control form-control-white input-sm" placeholder=""/>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- Payment Account Comment -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Comment</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.email" type="text" class="form-control form-control-white input-sm" placeholder="Please transfer the money via wire transfer and indicate your name in the payment reference"/>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <div class="alert alert-info margin-top-16px">
                                Notification Settings. The application will send emails to users, therefore it needs to have access to an email account for doing so.
                            </div>
                            <!-- SMTP Sender -->
                            <div class="row">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Sender Address</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.email" type="text" class="form-control form-control-white input-sm" placeholder="foo@bar.com"/>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- SMTP Server -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">SMTP Server</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.email" type="text" class="form-control form-control-white input-sm" placeholder="smtp.myserver.com:589"/>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- SMTP Username -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Username</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.email" type="text" class="form-control form-control-white input-sm" placeholder="foobar"/>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- SMTP Password -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Password</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.email" type="text" class="form-control form-control-white input-sm" placeholder="mypassword"/>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <div class="alert alert-info margin-top-16px">
                                ReCAPTCHA protects the application and its users from SPAMers. Thus it is recommended to use ReCAPTCHA.
                            </div>
                            <!-- Recaptcha Private Key -->
                            <div class="row">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Private Key</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.email" type="text" class="form-control form-control-white input-sm" placeholder=""/>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                            <!-- Recaptcha Public Key -->
                            <div class="row row-top-padded">
                                <div class="col-sm-1 col-xs-1"></div>
                                <div class="col-sm-3 col-xs-3">
                                    <div>
                                        <label class="label-sm" for="">Public Key</label>
                                    </div>
                                </div>
                                <div class="cold-sm-7 col-xs-7">
                                    <input v-model="user.email" type="text" class="form-control form-control-white input-sm" placeholder=""/>
                                </div>
                                <div class="col-sm-1 col-xs-1"></div>
                            </div>
                        </div>
                    </div>
                    <div class="modal-footer">
                        <button type="button" class="btn btn-default" v-on:click="cancel">
                            <span class="glyphicon glyphicon-remove"></span>
                            Back
                        </button>
                        <button type="button" class="btn btn-default" v-on:click="save">
                            <span class="glyphicon glyphicon-check"></span>
                            Save
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</div>

