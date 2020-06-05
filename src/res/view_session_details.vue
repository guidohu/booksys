<div class="panel panel-default" id="view_session_details">
    <div class="panel-heading">
        <div class="row">
            <div class="col-xm-6 col-xs-6">
                <h5 class="text-left">
                Session Details
                </h5>
            </div>
            <div v-if="session!=null" class="col-sm-6 col-xs-6 text-right">
                <div class="btn-group dropdown" id="menu_session">
                    <a class="btn btn-small btn-inline dropdown-toggle" data-toggle="dropdown" href="#">
                        <span class="glyphicon glyphicon-edit"></span>
                    </a>
                    <ul class="dropdown-menu text-left pull-right">
                        <!--<li v-if="session.exists==true">
                            <a href="#" onclick="showEditSession()">
                                <span class="glyphicon glyphicon-edit"></span>
                                Edit Session
                            </a>
                        </li>-->
                        <li v-if="session.exists==true">
                            <a href="#" v-on:click="addRider">
                                <span class="glyphicon glyphicon-plus"></span>
                                Add Riders
                            </a>
                        </li>
                        <li v-if="session.isCreate==true">
                            <a href="#" v-on:click="createSession">
                                <span class="glyphicon glyphicon-plus"></span>
                                Create Session
                            </a>
                        </li>
                        <li v-if="session.exists==true">
                            <a href="#" v-on:click="deleteSession">
                                <span class="glyphicon glyphicon-remove"></span>
                                Delete Session
                            </a>
                        </li>
                    </ul>
                </div>
            </div>
        </div>
    </div>
    <div class="panel-body">
        <div class="row">
            <div class="col-sm-1 col-xs-1">			
            </div>
            <div class="col-sm-10 col-xs-10">
                <div v-if="session == null">
                    <div class="row text-center">
                        <div class="alert alert-light">
                            Please select a session or free slot.
                        </div>
                    </div>
                </div>
                <div v-else>
                    <div class="row text-left">
                        <span class="col-sm-2 col-xs-2">
                            Date
                        </span>
                        <div class="col-sm-10 col-xs-10">
                            {{session.date_string}}
                        </div>
                    </div>
                    <div class="row text-left">
                        <span class="col-sm-2 col-xs-2">
                            Time
                        </span>
                        <div class="col-sm-10 col-xs-10">
                            {{session.time_string}}
                        </div>
                    </div>
                    <div v-if="session.isCreate==false">
                        <div v-if="session.title != null" class="row text-left">
                            <span class="col-sm-2 col-xs-2">
                                Title
                            </span>
                            <div class="col-sm-10 col-xs-10">
                                {{session.title}}                            
                            </div>
                        </div>
                        <div v-if="typeof session.riders != 'undefined' && session.riders.length > 0" class="row text-left">
                            <span class="col-sm-2 col-xs-2">
                                Riders
                            </span>
                            <div class="col-sm-10 col-xs-10">
                                <div v-if="session.riders.length==0">
                                    no riders yet
                                </div>
                                <ul class="unstyled" id="session_lst_riders">
                                    <li v-for="rider in session.riders">{{rider.first_name}} {{rider.last_name}}</li>
                                </ul>					
                            </div>
                        </div>
                        <div v-else class="row text-left">
                            <span class="col-sm-2 col-xs-2">
                                Riders
                            </span>
                            <div class="col-sm-10 col-xs-10">
                                <button class="btn btn-sm btn-success" v-on:click="addRider">
                                    <span class="glyphicon glyphicon-user" aria-hidden="true"></span>
                                    Add Riders
                                </button>					
                            </div>
                        </div>
                    </div>
                    <div v-if="session != null && session.isCreate==true">
                        <div class="row text-left">
                            <span class="col-sm-2 col-xs-2">
                            </span>
                            <div class="col-sm-10 col-xs-10">
                                <button class="btn btn-sm btn-success" v-on:click="createSession">
                                    <span class="glyphicon glyphicon-plus" aria-hidden="true"></span>
                                    Create Session
                                </button>                      
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>