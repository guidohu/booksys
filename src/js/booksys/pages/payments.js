var payment = new Array();
var payment_types = new Array();
var expenditure_types = new Array();
var user    = new Array();
var statistics = new Array();
var paymentTable;
var table_data;

// Stay in the WebApp
$(function() {
    $.stayInWebApp();
});

// Check if mobile browser first,
// then continue with checks
$(function() { 
    if(BooksysBrowser.isMobile()){
        // // Make it behave like an app
        BooksysBrowser.setViewportMobile();
        BooksysBrowser.setManifest();
        BooksysBrowser.setMetaMobile();
        // // Add mobile style dynamically
        BooksysBrowser.addMobileCSS();

        $("#body").load("res/accounting_mobile.html", function(){
            loadContent();
        });
    } else {
        $("#body").load("res/accounting.html", function(){
            loadContent();
        });
    }
});

function loadContent(){
    $.validator.setDefaults({
        errorElement: "validate_msg",
    });

    // Get all the data
    var current_year = new Date().getFullYear();
    getCurrency();
    getPaymentInfo(current_year);
    getUsers();
    getStatistics(current_year);
    getPaymentTypes();
    getYears();
    
    // Configure validation
    setValidation();
    
    // Register some handlers
    $("#add_expenditure_type").change(function(){
        updateExpenditureDialog();
    });
    $("#add_payment_type").change(function(){
        updatePaymentDialog();
    });
    $("#selectYear").change(function(){
        changeYear();
    });
}

// Updates the data which can change during visiting this website
function updateData(year){
    getPaymentInfo(year);
    getStatistics(year);
    getUsers();
    getPaymentTypes();
}

// Setup validation for forms
function setValidation(){
    $('#add_payment_form').validate({
        submitHandler: function(form){
            addPayment();
        },
        rules: {
            add_payment_name: {
                required: true,
                optionSelected: true,
            },
            add_payment_amount: {
                required: true,
                number: true,
            },
            add_payment_type: {
                required: true,
                optionSelected: true,
            },
            add_payment_date: {
                required: true,
            },
        },
        messages: {
            add_payment_name: {
                required: "Please select a rider",
            },
            add_payment_amount: {
                required: "No amount specified",
                number:   "This is no valid input",
            },
            add_payment_type: {
                required: "No type selected",
            },
            add_payment_date: {
                required: "No date selected",
            },
        },
    });
    
    $('#add_expenditure_form').validate({
        submitHandler: function(form){
            addExpenditure();
        },
        rules: {
            add_expenditure_user: {
                required: true,
                optionSelected: true,
            },
            add_expenditure_amount: {
                required: true,
                number: true,
            },
            add_expenditure_type: {
                required: true,
                optionSelected: true,
            },
            add_expenditure_date: {
                required: true,
            },
            add_expenditure_comment: {
                required: function(element){
                    if($("#add_expenditure_type").val()==0){
                        return false;
                    }else{
                        return true;
                    }
                },
                notEmpty: function(element){
                    if($("#add_expenditure_type").val()==0){
                        return false;
                    }else{
                        return true;
                    }
                },
            },
            add_expenditure_liters: {
                required: function(element){
                    if($("#add_expenditure_type").val()==0){
                        return true;
                    }else{
                        return false;
                    }
                },
            },
            add_expenditure_engine_hours: {
                required: function(element){
                    if($("#add_expenditure_type").val()==0){
                        return true;
                    }else{
                        return false;
                    }
                },
            },
        },
        messages: {
            add_expenditure_user: {
                required: "Please select the affected person",
                optionSelected: "Please select the affected person",
            },
            add_expenditure_amount: {
                required: "No amount specified",
                number:   "This is no valid input",
            },
            add_expenditure_type: {
                required: "No type selected",
                optionSelected: "Please select a type",
            },
            add_expenditure_date: {
                required: "No date selected",
            },
            add_expenditure_comment: {
                required: "Please shortly describe",
            },
            add_expenditure_liters: {
                required: "Please add the amount of fuel",
            },
            add_expenditure_engine_hours: {
                required: "Please add the engine hours",
            }
        },
    });
}

// Update the dialog based on the selected type
function updateExpenditureDialog(){
    var type = $('#add_expenditure_type').val();
    if(type == 0){
        $("#add_expenditure_engine_hours_line").slideDown(500);
        $("#add_expenditure_liters_line").slideDown(500);
        $("#add_expenditure_comment_line").slideUp(500);
    }else{
        $("#add_expenditure_engine_hours_line").slideUp(500);
        $("#add_expenditure_liters_line").slideUp(500);
        $("#add_expenditure_comment_line").slideDown(500);
    }
}

// Update the payment dialog based on the selected type
function updatePaymentDialog(){
    var type = $('#add_payment_type').val();
    if(type == 4 || type == 6){
        $("#add_payment_comment_line").slideUp(300);
    }else{
        $("#add_payment_comment_line").slideDown(300);
    }
}

function getCurrency(){
    $.ajax({
        type: "GET",
        url: "api/configuration.php?action=get_customization_parameters",
        cache: false,
        success: function(data){
            let json = $.parseJSON(data);
            $("#add_payment_currency").html(json.currency);
            $("#add_expenditure_currency").html(json.currency);
            
        },
        error: function(request, status, error){
            console.log(request);
            console.log(status);
            console.error(error);
        },
    });
}

// Get the payment history
function getPaymentInfo(year){
    var data = new Object();
    data['year']    = year;

    $('#status_modal').modal({
        backdrop: 'static',
        keyboard: false,
        show:     true,
    });
    $('#status_footer').hide();
    $('#status_header').html('<h4><span class="glyphicon glyphicon-refresh"></span> Loading</h4>');
    $('#status_body').html('Please wait while data is retrieved...');
    $('#status_modal').modal('show');
    $.ajax({
        type: "POST",
        url: "api/payment.php?action=get_transactions",
        cache: false,
        data: JSON.stringify(data),
        success: function(data){
            var json = $.parseJSON(data);
            payment = json;
            updatePaymentInfo(payment);
            $('#status_modal').modal('hide');
        },
        error: function(request, status, error){
            alert("Cannot get the transactions.");				
            $('#status_modal').modal('hide');
        },
    });
}

// Get years for which we have payment data
function getYears(){
    $('#status_modal').modal({
        backdrop: 'static',
        keyboard: false,
        show:     true,
    });
    $('#status_footer').hide();
    $('#status_header').html('<h4><span class="glyphicon glyphicon-refresh"></span> Loading</h4>');
    $('#status_body').html('Please wait while data is retrieved...');
    $('#status_modal').modal('show');
    $.ajax({
        type: "GET",
        url: "api/payment.php?action=get_years",
        cache: false,
        success: function(data){
            var json = $.parseJSON(data);
            updateYearSelection(json);
            $('#status_modal').modal('hide');
        },
        error: function(request, status, error){
            alert("Cannot get the transactions.");				
            $('#status_modal').modal('hide');
        },
    });
}

// Get all users from the database
function getUsers(){
    $.ajax({
        type: "GET",
        url: "api/user.php?action=get_all_users",
        cache: false,
        success: function(data){
            var json = $.parseJSON(data);
            user = json;
            updatePaymentInfo(payment);
        }
    });
}

// Get statistics from the server
function getStatistics(year){
    var data = new Object();
    data['year']    = year;
    
    $.ajax({
        type: "POST",
        url: "api/payment.php?action=get_statistics",
        cache: false,
        data: JSON.stringify(data),
        success: function(data){
            var json = $.parseJSON(data);
            statistics = json;
            updateStatistics(statistics);
        }
    });
}

// Get payment types
function getPaymentTypes(){
    $.ajax({
        type: "GET",
        url: "api/payment.php?action=get_expenditure_types",
        cache: true,
        success: function(data){
            var json = $.parseJSON(data);
            expenditure_types = json;
        }
    });
    $.ajax({
        type: "GET",
        url: "api/payment.php?action=get_payment_types",
        cache: true,
        success: function(data){
            var json = $.parseJSON(data);
            payment_types = json;
        }
    });
}

// Generates the year selection drop-down
function updateYearSelection(years){
    var date = new Date();
    
    var content = "";
    var selected = 0;
    for(var i=0; i<years.length; i++){
        content += '<option value="';
        // Select the current year
        if(years[i].year == date.getFullYear()){
            content += years[i].year;
            content += '" selected="selected">';
            selected = 1;
        }else if(years[i].year != 'any'){
            content += years[i].year;
            content += '">';			
        }else{
            content += "any";
            content += '">';
        }
        
        if(years[i].year == "any"){
            content += years[i].year + "</option>";
        }else{
            content += "Year " + years[i].year + "</option>";
        }
    }
    
    // the current year did not show up in the selection
    if(selected == 0){
        content += "<option value=\"" + date.getFullYear() + "\" selected=\"selected\">";
        content += "Year " + date.getFullYear() + "</option>";
    }
    
    // Add the selections
    $('#selectYear').html(content);
}

// Retrieves the data for the given year
function changeYear(){
    var year = $('#selectYear option:selected').val();
    updateData(year);
}

function deleteEntry(idx){		
    if(!confirm("Are you sure you want to delete this record?\n "
        + table_data[idx][2] + " " + table_data[idx][5] + " " + table_data[idx][8] + " " + table_data[idx][6])){
        return;
    }
    
    // TODO call to API to delete from table row with id
    var data = new Object();
    data['table_id'] = table_data[idx][0];
    data['row_id']   = table_data[idx][1];
    
    $("#add_payment_status").html('Deleting Entry');
    $("#add_payment_status").show();
    
    $.ajax({
        type: "POST",
        url: "api/payment.php?action=delete_payment",
        data: JSON.stringify(data),
        success: function(data){
            updateData();				
            $("#add_payment_status").hide();
        },
        error: function(data, text, errorCode){
            $("#add_payment_status").html(data.responseText);		
        }
    });
}

// Update payment info table
function updatePaymentInfo(data){
    // we might call this function with undefined data
    if(data.transactions == null){
        return;
    }

    table_data  = new Array();
    
    for(var i=0; i<data.transactions.length; i++){
        var sign=1;
        if(data.transactions[i].tbl == '0' || data.transactions[i].tbl == '2'){
            sign = -1;
        }
        table_data[i] = [ 
            data.transactions[i].tbl,
            data.transactions[i].id,
            data.transactions[i].timestamp,
            data.transactions[i].type_name,
            data.transactions[i].fn + ' ' + data.transactions[i].ln,
            sprintf("%.2f", sign*data.transactions[i].amount),
            data.transactions[i].comment,
            "<a class='text-center' id='deleteAction' value='"+i+"'><icon class='glyphicon glyphicon-remove text-center'></icon></a>",
            data.currency
        ];
    }
    
    if(table_data.length<1){
        table_data == null;
    }
    
    // The payment history table
    if(paymentTable != null){
        paymentTable.fnDestroy();
    }
    paymentTable = $('#payment_table').dataTable( {
        "aaData": table_data,
        "aoColumns": [
            { "sTitle": "table_id",
            "bVisible": false,
            },
            { "sTitle": "record_id",
            "bVisible": false,
            },
            { "sTitle": "Date",
            "sClass": "left",
            },
            { "sTitle": "Type",
            "sClass": "left",
            },
            { "sTitle": "Name",
            "sClass": "left",
            },
            { "sTitle": data.currency,
            "sClass": "text-right padding-right-16px",
            },
            { "sTitle": "Comment",
            "sClass": "left",
            },
            { "sTitle": "Action",
            "sClass": "center action",
            }
        ],
        "bLengthChange": false,
        "bFilter":       false,
        "bPaginate":     false,
        "bInfo":         false,
        "bRetrieve":     true,
        "sScrollY":      "200px"
    });
    paymentTable.fnSort( [ [2,'desc'] ] );
    
    // Register the delete Action Handler for every button
    $("a[id='deleteAction']").click( function(){
        var idx = $(this).attr('value');
        deleteEntry(idx);
    });
}

// Updates the statistics
function updateStatistics(data){
    $('#total_payments').html(sprintf("%.2f", data.total_payment_selected_year) + ' ' + data.currency);
    $('#total_expenditures').html(sprintf("%.2f", data.total_expenditure_selected_year) + ' ' + data.currency);
    $('#total_used').html(sprintf("%.2f", data.total_used_selected_year) + ' ' + data.currency);
    $('#total_open').html(sprintf("%.2f", data.total_open) + ' ' + data.currency);
    $('#total_balance').html(sprintf("%.2f", data.current_balance) + ' ' + data.currency);
    $('#total_current_profit').html(sprintf("%.2f", data.current_session_profit_selected_year) + ' ' + data.currency);
    $('#total_used').popover({placement: 'top', trigger: 'hover', html: true, title: "Minute Distribution", 
                            content:   "<div class='row row-padded'><div class='col-sm-5 col-xs-5'><label>Admin</label></div><div class='col-sm-7 col-xs-7 text-right'>" + sprintf("%.2f", data.admin_minutes_selected_year) + "min</div></div>"
                                    + "<div class='row row-padded'><div class='col-sm-5 col-xs-5'><label>Member</label></div><div class='col-sm-7 col-xs-7 text-right'>" + sprintf("%.2f", data.member_minutes_selected_year) + "min</div></div>"
                                    + "<div class='row row-padded'><div class='col-sm-5 col-xs-5'><label>Guest</label></div><div class='col-sm-7 col-xs-7 text-right'>" + sprintf("%.2f", data.guest_minutes_selected_year) + "min</div></div>"});
}

// Show the payment dialog
function showPaymentDialog(){
    var name_text = '<option value="" selected="selected"></option>';        
    for(var i=0; i<user.length; i++){
        name_text += '<option value="' + user[i].id + '">' 
                    + user[i].first_name + ' ' + user[i].last_name + '</option>';
    }        
    $('#add_payment_name').html(name_text);
    
    $('#add_payment_amount').val('');

    var date_str = moment().tz(getTimeZone()).format('YYYY-MM-DD');
    $('#add_payment_date').val(date_str);
    $('#add_payment_comment').val('');
    
    // fill the type option box
    var type_text = '<option value="" selected="selected"></option>';
    for(var i=0; i<payment_types.length; i++){
        type_text += '<option value="' + payment_types[i].id + '">' + payment_types[i].name + '</option>';
    }
    $('#add_payment_type').html(type_text); 
    
    $('#add_payment_modal').modal({
        backdrop: 'static',
        keyboard: false,
        show:     true,
    });
    $('#add_payment_modal').modal('show');
}

// Hide the payment dialog
function hidePayment(){
    $('#add_payment_modal').modal('hide');
}

// Show the expenditure dialog
function showExpenditureDialog(){
    $("#add_expenditure_engine_hours_line").hide();
    $("#add_expenditure_liters_line").hide();
    
    // fill the user option box
    var name_text = '<option value="" selected="selected"></option>';        
    for(var i=0; i<user.length; i++){
        name_text += '<option value="' + user[i].id + '">' 
                    + user[i].first_name + ' ' + user[i].last_name + '</option>';
    }        
    $('#add_expenditure_user').html(name_text);

    // reset all values
    $('#add_expenditure_amount').val('');
    var date_str = moment().tz(getTimeZone()).format('YYYY-MM-DD');
    $('#add_expenditure_date').val(date_str);
    $('#add_expenditure_comment').val('');
    $('#add_expenditure_engine_hours').val('');
    $('#add_expenditure_liters').val('');
    $('#add_expenditure_status').hide();
    
    // fill the type option box
    var type_text = '<option value="-1" selected="selected"></option>';
    for(var i=0; i<expenditure_types.length; i++){
        type_text += '<option value="' + expenditure_types[i].id + '">' + expenditure_types[i].name + '</option>';
    }
    $('#add_expenditure_type').html(type_text);    

    // display window
    $('#add_expenditure_modal').modal({
        backdrop: 'static',
        keyboard: false,
        show:     true,
    });
    $('#add_expenditure_modal').modal('show');
}

// Hide the expenditure dialog
function hideExpenditure(){
    $('#add_expenditure_modal').modal('hide');
}

function addPayment(){
    $("#add_payment_add_btn").attr("disabled", true);
    $("#add_payment_status").html('Sending');
    $("#add_payment_status").show();
    
    var data = new Object();
    data['amount']  = $('#add_payment_amount').val();
    data['type_id'] = $('#add_payment_type').val();
    data['date']    = $('#add_payment_date').val();
    data['user_id'] = $('#add_payment_name').val();
    data['comment'] = $('#add_payment_comment').val();
    
    // add payment to database
    $.ajax({
        type: "POST",
        url: "api/payment.php?action=add_payment",
        data: JSON.stringify(data),
        success: function(data){
            hidePayment();
            updateData();
            $("#add_payment_status").hide();
        },
        error: function(data, text, errorCode){
            $("#add_payment_status").html(data.responseText);
            
        },
        complete: function(jqXHR, testStatus){
            $("#add_payment_add_btn").attr("disabled", false);
        },
    });
}

// Writes the expenditure to the database
function addExpenditure(){
    $("#add_expenditure_add_btn").attr("disabled", true);
    $("#add_expenditure_status").html('Sending');
    $("#add_expenditure_status").show();

    var data = new Object();
    data['user_id']    = $('#add_expenditure_user').val();
    data['type_id']    = $('#add_expenditure_type').val();
    data['date']       = $('#add_expenditure_date').val();
    data['amount']     = $('#add_expenditure_amount').val();
    data['comment']    = $('#add_expenditure_comment').val();
    
    // check if it is a fuel expenditure or not
    var post_url = "api/payment.php?action=add_expenditure";
    if(data['type_id'] == 0){
        post_url = "api/boat.php?action=update_fuel";
        data['engine_hours'] = $('#add_expenditure_engine_hours').val();
        data['liters'] = $('#add_expenditure_liters').val();
        data['cost'] = data['amount'];
    }
    
    $.ajax({
        type: "POST",
        url: post_url,
        data: JSON.stringify(data),
        success: function(data){
            hideExpenditure();
            updateData();
            $("#add_expenditure_status").hide();
        },
        error: function(data, text, errorCode){
            $("#add_expenditure_status").html(data.responseText);
            
        },
        complete: function(jqXHR, testStatus){
            $("#add_expenditure_add_btn").attr("disabled", false);
        },
    });
}

// gets the user's timezone or returns the default
function getTimeZone(){
    if (typeof(Storage) !== "undefined") {
        var timezone = localStorage.getItem("timezone");
        if(timezone != null){
            return timezone;
        }
        else{
            // TODO let user set his timezone
            return "Europe/Berlin";
        }
    } else {
        // Sorry! No Web Storage support..
        return "Europe/Berlin";
    }
}

