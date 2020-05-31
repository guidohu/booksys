// Check if mobile browser first,
// then continue with the page logic
$(function() { 
    if(BooksysBrowser.isMobile()){
        // Make it behave like an app
        BooksysBrowser.setViewportMobile();
        BooksysBrowser.setManifest();
        BooksysBrowser.setMetaMobile();
        // Add mobile style dynamically
        BooksysBrowser.addMobileCSS();

        $("#body").load("res/boat_mobile.html", function(){
            loadContent();
        });
    } else {
        // Add mobile style also in desktop mode
        BooksysBrowser.addMobileCSS();
        $("#body").load("res/boat_mobile.html", function(){
            loadContent();
        });
    }
});

// Load dynamic content
function loadContent(){
    $.stayInWebApp();
    
    // Set default values
    Date.prototype.toDateInputValue = (function() {
        var local = new Date(this);
        local.setMinutes(this.getMinutes() - this.getTimezoneOffset());
        return local.toJSON().slice(0,10);
    });
    $('#maintenance_date').val(new Date().toDateInputValue());

    $('#boatTab a').click(function (e) {
        e.preventDefault()
        $(this).tab('show')
    });
    
    // Select first tab
    $('#boatTab a:first').tab('show');
    
    // Initialize Bootstrap Switch
    $("[name='boatUsageType']").bootstrapSwitch({
        onText: "Course",
        onColor: "info",
        offText: "Private",
        animate: true,
        size: 'normal',
        onSwitchChange: function(event, state){ isBoatUsageCourse = state }
    });
    
    // Retrieve updates from server
    updateBoatLog();
    updateLatest();
    updateFuelLog();
    updateMaintenanceLog();

    // Get customization settings
    updateCustomization();
    
    // validation rules for engine_hours form
    $('#engine_hours_form').validate({
        submitHandler: function(form){
                    updateEngineHours();
                },            
        rules: {
                    // engine hours before need are required
                    engine_hours_before: {
                                            required:true,
                                            minlength:1,
                                            maxlength:10,
                                            number:true,
                                        },
                    // engine hours after is only needed if engine hours before is greyed out
                    engine_hours_after: {
                        required: {
                            depends: function(element){
                                if($("#engine_hours_btn").html() == 'Start'){
                                    return false;
                                }else{
                                    return true;
                                }
                            },
                        },
                        minlength:1,
                        maxlength:10,
                        number:true,
                    },
                },
        messages: {
                    engine_hours_after: "Please enter the engine hours",
                    engine_hours_before: "Please enter the engine hours",
                },
    });
    
    // validation rules for fuel form
    $('#fuel_form').validate({
        submitHandler: function(form){
                    updateFuel();
                },            
        rules: {
                    // engine hours are required
                    fuel_engine_hours: {
                                            required:true,
                                            minlength:1,
                                            maxlength:10,
                                            number:true,
                                        },
                    // fuel liters are required
                    fuel_liters: {
                                            required:true,
                                            minlength:1,
                                            maxlength:10,
                                            number:true,
                                        },
                    // fuel cost is required
                    fuel_cost: {
                                            required:true,
                                            minlength:1,
                                            maxlength:10,
                                            number:true,
                                        },
                },
        messages: {
                    fuel_engine_hours: "Please enter the engine hours",
                    fuel_liters: "Please enter the amount of fuel",
                    fuel_cost: "Please enter the cost for the fuel",
                },
    });
    
    // validation rules for maintenance form
    $('#maintenance_form').validate({
        submitHandler: function(form){
                    updateMaintenance();
                },            
        rules: {
                    // engine hours are required
                    maintenance_engine_hours: {
                                            required:true,
                                            minlength:1,
                                            maxlength:10,
                                            number:true,
                                        },
                    // description required
                    maintenance_description: {
                                            required:true,
                                            minlength:1,
                                        },
                },
        messages: {
                    maintenance_engine_hours: "Please enter the engine hours",
                    maintenance_description: "Please describe what has been done",
                },
    });
}

var isBoatUsageCourse = 0;

$.validator.setDefaults({
    invalidHandler: function(event, validator){
        var errors = validator.numberOfInvalids();
        if(errors){							
            $('#panel-warning').html("Please check your input again.");
            $('#panel-warning').show();
        }else{
            $('#panel-warning').fadeOut();							
        }	
    },
    errorElement: "validate_msg",
    success: function(label){
        $('#panel-warning').fadeOut(1000);
    },
});

/* Update the Boat Logbook */
function updateBoatLog(){
    $.ajax({
        type:    "GET",
        url:     "api/boat.php?action=get_engine_hours_log",
        cache:   false,
        success: function(data){
            json = $.parseJSON(data);

            // if the API tells us to redirect
            if(typeof json.redirect != "undefined"){
                window.location.replace(json.redirect);
                return;
            }
            displayBoatLog(json);
        },
        error: function(request, status, error){
        },
        complete: function(){
        },
    });
}

/* Update the Fuel Logbook */
function updateFuelLog(){
    $.ajax({
        type: 	"GET",
        url:	"api/boat.php?action=get_fuel_log",
        cache:   false,
        success: function(data){
            json = $.parseJSON(data);
            displayFuelLog(json);
            displayFuelChart(json);
        }
    });
}

function displayFuelChart(data){
    let fuel_series = new Object();
    let last_engine_hours = 0;
    for( let i=data.length-1; i>=0; i--){
        if(last_engine_hours == 0){
            // very first value, we cannot give an average
            last_engine_hours = data[i].engine_hours;
            continue;
        }
        if(last_engine_hours > data[i].engine_hours){
            // the engine hours decreased -> new boat?
            // re-initialize
            last_engine_hours = data[i].engine_hours
            continue;
        }
        
        // get the engine hour difference
        let engine_hours = data[i].engine_hours - last_engine_hours;
        if(engine_hours == 0){
            last_engine_hours = data[i].engine_hours;
            continue;
        }
        let fuel_consumption = data[i].liters / engine_hours;
        last_engine_hours = data[i].engine_hours;
        
        // get the year
        let date = new Date(data[i].timestamp*1000);
        if (fuel_series[date.getFullYear()]){
            fuel_series[date.getFullYear()].push([Date.UTC(1970, date.getMonth(), date.getDate()), fuel_consumption]);
        }else{
            fuel_series[date.getFullYear()] = [[Date.UTC(1970, date.getMonth(), date.getDate()), fuel_consumption]];
        }			
    }

    // convert to charts.js format
    let datasets = [];
    let counter = 0;
    let size = Object.keys(fuel_series).length;
    for ( let year in fuel_series ) {
        counter++;
        let dataset = {
            label:		year,
            data: 		[],
            showLine: 	true,
            fill: 		false,
            hidden:		true,
        };

        // display the current year
        if(counter == size){
            dataset.hidden = false;
        }

        // only add the last 3 years
        if(counter > size-3){
            fuel_series[year].forEach(element => {
                dataset.data.push(
                    {
                        x: moment(element[0], 'x'),
                        y: element[1]
                    }
                );
            });

            datasets.push(dataset);
        }
    }

    let ctx = document.getElementById("fuel_chart");
    let myLineChart = new Chart(ctx, {
        type: 'scatter',
        data: {
            datasets: datasets
        },
        options: {
            title: {
                display: true,
                text: 'Average Fuel Consumption',
                fontStyle: 'none'
            },
            legend: {
                position: 'bottom'
            },
            tooltips: {
                enabled: true,
                callbacks: {
                    title: function(tooltipItem, data) {
                        return "Average Fuel Consumption";
                    },
                    label: function(tooltipItem, data) {
                        time = data.datasets[tooltipItem.datasetIndex].data[tooltipItem.index].x.format("DD. MMM");
                        liters = data.datasets[tooltipItem.datasetIndex].data[tooltipItem.index].y
                        liters = Math.round(liters * 10)/10;
                        return time + " - " + liters + " L";
                    },
                    afterLabel: function(tooltipItem, data) {
                        return '';
                    }
                },
            },
            scales: {
                xAxes: [{
                    type: 'time',
                    distribution: 'series'
                }],
                yAxes: [{
                    ticks: {
                        min: 0,
                        max: 40,
                        suggestedMin: 0,
                        suggestedMax: 40
                    },
                    scaleLabel: {
                        display: true,
                        labelString: 'liter/hour'
                    }
                }]
            },
            // responsive: true,
            maintainAspectRatio: false,
        }
    });
}

function updateMaintenanceLog(){
    $.ajax({
        type:	"GET",
        url:	"api/boat.php?action=get_maintenance_log",
        cache:	false,
        success:	function(data){
            json = $.parseJSON(data);
            displayMaintenanceLog(json);
        }
    });
}

function updateCustomization(){
    $.ajax({
        type:	"GET",
        url:	"api/configuration.php?action=get_customization_parameters",
        cache:	false,
        success:	function(data){
            json = $.parseJSON(data);
            $("#fuel_cost_currency").html(json.currency.toLowerCase());
        }
    });
}

/* Creates a dateString as we use it everywhere */
function dateToString(date){
    var string = '';
    // Date
    string += date.getFullYear();
    string += '.';
    if(date.getMonth()< 9){
        string+= '0';
    }
    string += date.getMonth()+1;
    string += '.';
    if(date.getDate() < 10){
        string += 0;
    }
    string += date.getDate();
    
    // Time
    string += ' ';
    if(date.getHours() < 10){
        string += 0;
    }
    string += date.getHours();
    string += ':';
    if(date.getMinutes() < 10){
        string += 0;
    }
    string += date.getMinutes();
    return string;		
}

/* Display the Boat Logbook */
function displayBoatLog(data){
    var table = document.getElementById("logbook_content");
    var totalEntries = Object.keys(data).length;

    // remove old entries
    while(table.rows.length > 0){
        table.deleteRow(0);
    }

    // fill the table with all the entries we got
    for(var i=0; i<totalEntries; i++){
        var date = new Date(data[i].time*1000);
        var after = '-';
        var diff  = '-';

        if(data[i].after_hours != null){
            after = data[i].after_hours;
            diff  = data[i].delta_hours;
        }

        var row = table.insertRow(i);
        row.id = data[i].id;
        var dateCell   = row.insertCell(0);
        var userCell   = row.insertCell(1);
        var beforeCell = row.insertCell(2);
        var afterCell  = row.insertCell(3);
        var diffCell   = row.insertCell(4);

        dateCell.innerHTML   = dateToString(date);
        userCell.innerHTML   = data[i].user_first_name;
        beforeCell.innerHTML = data[i].before_hours;
        afterCell.innerHTML  = after;
        diffCell.innerHTML   = diff;

        if(data[i].type_name == 'course'){
            dateCell.setAttribute("class", "highlight");
            userCell.setAttribute("class", "highlight");
            beforeCell.setAttribute("class", "highlight");
            afterCell.setAttribute("class", "highlight");
            diffCell.setAttribute("class", "highlight");
        }

        // add click action
        row.addEventListener("click", function(event){ 
            displayBoatLogEntry(event);
        });
    }
}

/* Display a single Boat Log Entry */
function displayBoatLogEntry(event){
    // cell that was clicked on
    var cell = event.target;
    // get ID of the element
    var id = cell.parentElement.id;

    BooksysViewEngineHoursEntry.displayEntry('user_dialog_modal', id, function(){
        BooksysViewEngineHoursEntry.destroyView('user_dialog_modal');
        updateBoatLog();
    });
}

/* Display the Fuel Logbook */
function displayFuelLog(data){
    var table = document.getElementById("fuel_log_content");
    var totalEntries = Object.keys(data).length;

    // remove old entries
    while(table.rows.length > 0){
        table.deleteRow(0);
    }

    // fill the table with all the entries we got
    // --> sorting: latest has highest array index
    for(var i=totalEntries-1; i>=0; i--){
        console.log(data[i]);
        var date = new Date(data[i].timestamp*1000);
        var row = table.insertRow(0);
        row.id = data[i].id;

        // highlight rows with discount
        if(data[i].cost_brutto != "undefined" && data[i].cost_brutto != null){
            row.classList.add("highlight");
        }

        var dateCell = row.insertCell(0);
        var userCell = row.insertCell(1);
        var engineHoursCell = row.insertCell(2);
        var litersCell = row.insertCell(3);
        var costCell = row.insertCell(4);

        dateCell.innerHTML = dateToString(date);
        userCell.innerHTML = data[i].user_first_name;
        engineHoursCell.innerHTML = data[i].engine_hours;
        litersCell.innerHTML = data[i].liters;
        costCell.innerHTML = data[i].cost;

        // add click action
        row.addEventListener("click", function(event){ 
            displayFuelEntry(event)
        });
    }
}

// displays a specific fuel entry
function displayFuelEntry(event){
    // cell that was clicked on
    var cell = event.target;
    // get the ID of the row
    var id = cell.parentElement.id;

    // get the entry that needs to be displayed
    //var entry = data[idx];
    BooksysViewFuelEntry.displayFuelEntry('user_dialog_modal', id, function(){ 
        BooksysViewFuelEntry.destroyView('user_dialog_modal');
        updateFuelLog();
    });
}

/* Display the Maintenance Logbook */
function displayMaintenanceLog(data){
    var maintlog = '';
    for(var i=0; i<Object.keys(data).length; i++){
        var date = new Date(data[i].timestamp*1000);
        maintlog += '<tr>';
        maintlog += '<th>'+dateToString(date)+'</td>';
        maintlog += '<th>'+data[i].engine_hours+'</td>';
        maintlog += '<th>'+data[i].user_first_name+'</td>';
        var descr = data[i].description.replace(/</g, "&lt;").replace(/>/g, "&gt;");
        maintlog += '<th>'+descr+'</td>';
        maintlog += '</tr>';
    }
    $('#maintenance_log_content').html(maintlog);
}

/* Get the latest engine hour entry from the database */
function updateLatest(){
    $.ajax({
        type:    "GET",
        url:     "api/boat.php?action=get_engine_hours_latest",
        cache:   false,
        success: function(data){
            json = $.parseJSON(data);
            displayLatestLog(json);
        }
    });
}

/* Displays a boat_log record which is not yet complete */
function displayLatestLog(data){
    if(data.length == 0 || data[0].id == null){
        // The very first session starts
        updateMyUser(["engine_hours_user", "fuel_user", "maintenance_user"]);
        $("#engine_hours_btn").html("Start");
        $('#engine_hours_before').prop('disabled', false);
        $("#engine_hours_before").val('');
        $("#engine_hours_after").val('');
        $("#engine_hours_after_row").hide();
        $("#engine_hours_type_row").show();
        return;
    }else if((data[0].delta_hours && data[0].after_hours)
    || data[0].id == ''){
        // A new session starts		   
        updateMyUser(["engine_hours_user", "fuel_user", "maintenance_user"]);
        $("#engine_hours_btn").html("Start");
        $('#engine_hours_before').prop('disabled', false);
        $("#engine_hours_before").val(data[0].after_hours);
        $("#engine_hours_after").val('');
        $("#engine_hours_after_row").hide();
        $("#engine_hours_type_row").show();
        return;
    }else{		
        // An old session still needs to be completed
        updateMyUser(["fuel_user", "maintenance_user"]);
        $('#engine_hours_user').html("<option value='"+data[0].user_id+"'>"
                        +data[0].user_first_name+" "+data[0].user_last_name
                        +"</option>");
        $('#engine_hours_user').prop('disabled', true);
        $('#engine_hours_before').val(data[0].before_hours);
        $('#engine_hours_before').prop('disabled', true);
        $('#engine_hours_after').val('');
        $('#engine_hours_after_row').show();
        $('#engine_hours_type_row').hide();
        $("#engine_hours_btn").html("Finish");
    }
}

/* Updates the user dropbox with the current user */
function updateMyUser(ids){
    $.ajax({
        type:    "GET",
        url:     "api/user.php?action=get_my_user",
        success: function(data){
            let user = $.parseJSON(data);
            for(var i=0; i<ids.length; i++){
                displayMyUser(user, ids[i]);
            }
        }
    });	
}

function displayMyUser(user, id){
    var content = '<option value="'+user.id+'">'
                    +user.first_name+' '+user.last_name
                    +'</option>';
    $('#'+id).html(content);
    $('#'+id).prop('disabled', true);
}

/* Update the engine hours */
function updateEngineHours(){
    var data = new Object();
    data['user_id'] = $('#engine_hours_user').val();
    data['engine_hours_before'] = $('#engine_hours_before').val();
    data['engine_hours_after'] = $('#engine_hours_after').val();
    data['type'] = isBoatUsageCourse; // course or private
    
    $('#engine_hours_btn').html("Sending...");

    $.ajax({
        type: 	"POST",
        url:    "api/boat.php?action=update_engine_hours",
        cache:   false,
        data: 	JSON.stringify(data),
        success: function(data){
            $('#panel-warning').hide();
            updateBoatLog();
            updateLatest();
        },
        error:  function(data, text, errorCode){
            $('#panel-warning').html("Error"+ data.status + ": " + data.responseText);
            $('#panel-warning').show();
        },
        complete: function(jqXHR, testStatus){
            if('#engine_hours_before:disabled'){
                $('#engine_hours_btn').html("Finish");
            }else{
                $('#engine_hours_btn').html("Start");
            }
        }
    });
}

/* Send the fuel form to the server */
function updateFuel(){
    var data = new Object();
    data['user_id']      = $('#fuel_user').val();
    data['engine_hours'] = $('#fuel_engine_hours').val();
    data['liters']       = $('#fuel_liters').val();
    data['cost']     = $('#fuel_cost').val();
    
    $('#fuel_btn').html('Sending...');
    
    $.ajax({
        type:	"POST",
        url:	"api/boat.php?action=update_fuel",
        cache:   false,
        data:	JSON.stringify(data),
        success: function(data){
            $('#panel-warning').hide();
            updateFuelLog();
            $('#fuel_engine_hours').val('');
            $('#fuel_liters').val('');
            $('#fuel_cost').val('');
        },
        error: function(data, text, errorCode){
            $('#panel-warning').html("Error"+ data.status + ": " + data.responseText);
            $('#panel-warning').show();
        },
        complete: function(jqXHR, testStatus){
            $('#fuel_btn').html('Add');
        }
    });
}

function updateMaintenance(){		
    var data = new Object();
    data['user_id']      = $('#maintenance_user').val();
    data['engine_hours'] = $('#maintenance_engine_hours').val();
    data['description']       = $('#maintenance_description').val();
    
    $('#maintenance_btn').html('Sending...');
    
    $.ajax({
        type:	"POST",
        url:	"api/boat.php?action=update_maintenance_log",
        cache:   false,
        data:	JSON.stringify(data),
        success: function(data){
            $('#panel-warning').hide();
            updateMaintenanceLog();
            $('#maintenance_engine_hours').val('');
            $('#maintenance_description').val('');
        },
        error: function(data, text, errorCode){
            $('#panel-warning').html("Error"+ data.status + ": " + data.responseText);
            $('#panel-warning').show();
        },
        complete: function(jqXHR, testStatus){
            $('#maintenance_btn').html('Add');
        }
    });
}
