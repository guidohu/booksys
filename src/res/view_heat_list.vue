<div id="view_heat_list">
    <div v-if="heats.length > 0" class="panel panel-default panel-padded">
        <div class="panel-body">
            <div class="row text-left">
                <div class="col-sm-12 col-xs-12 overflow-hidden">
                    <table class="table table-striped table-hover">
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th>Duration</th>
                                <th>Cost</th>
                                <th v-if="hasComments">Comment</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="heat in heats">
                                <td>{{heat.first_name}} {{heat.last_name}}</td>
                                <td>{{heat.duration_string}}</td>
                                <td>{{heat.cost}} {{currency}}</td>
                                <td v-if="hasComments">{{heat.comment}}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
</div>