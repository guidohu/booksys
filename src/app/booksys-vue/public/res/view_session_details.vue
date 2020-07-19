<div class="panel panel-default" id="view_session_details">
    <div class="panel-heading">
        <div class="row">
            <div class="col-xm-6 col-xs-6 text-left">
                Session Details
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
                                <span class="glyphicon glyphicon-user"></span>
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
            <div class="col-sm-12 col-xs-12">
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
                        <div v-if="typeof session.riders != 'undefined' && session.riders.length > 0" class="row padding-top-8px text-left">
                            <div class="col-sm-12 col-xs-12">
                                Riders <span class="badge">{{ session.riders.length }}</span> / <span class="badge">{{ session.riders_max}}</span>
                            </div>
                            <div class="col-sm-12 col-xs-12">
                                <div v-for="rider in session.riders">
                                    <div class="col-sm-10 col-xs-10">
                                        {{ rider.first_name}} {{ rider.last_name }}
                                    </div>
                                    <div class="col-sm-2 col-xs-2 text-right">
                                        <button class="btn btn-default btn-inline" v-on:click="removeRider($event, rider.id)" :id="rider.id" :value="rider.id" :key="rider.id">
                                            <i class="fa fa-user-minus"></i>
                                        </button>
                                    </div>
                                </div>				
                            </div>
                        </div>
                        <div v-else class="row padding-top-8px text-left">
                            <span class="col-sm-2 col-xs-2">
                                Riders
                            </span>
                            <div class="col-sm-10 col-xs-10">
                                <button class="btn btn-xs btn-success" v-on:click="addRider">
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
                                <button class="btn btn-xs btn-success" v-on:click="createSession">
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