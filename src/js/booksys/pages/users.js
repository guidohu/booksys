// Check if mobile browser first,
// then continue with the page logic
$(function() { 
    if(BooksysBrowser.isMobile()){
        // // Make it behave like an app
        // BooksysBrowser.setViewportMobile();
        // BooksysBrowser.setManifest();
        // BooksysBrowser.setMetaMobile();
        // // Add mobile style dynamically
        // BooksysBrowser.addMobileCSS();

        $("#body").load("res/users.html", function(){
            loadContent();
        });
    } else {
        // Do not add mobile style in desktop mode
        $("#body").load("res/users.html", function(){
            loadContent();
        });
    }
});

let userGroups;
let displayCurrency;

//var users = new Array();
//var selectedUserIdx = -1;
var userTable;

// Load all content for this page
function loadContent() {
    $.stayInWebApp();

    $.ajax({
        type: "GET",
        url: "api/user.php?action=get_user_groups",
        cache: false,
        success: function(response){
            let json = $.parseJSON(response);
            userGroups = json.user_groups;
            getUsers();
            getUserGroups();
        }
    });

    // register tab switch listener
    // used to fix a bug in header spacing
    $('#tablist').on('shown.bs.tab', function(e){
        $.fn.dataTable.tables( {visible: true, api: true}).columns.adjust();
    });

}

// Displays the loading status
function showStatusModal(){
    $('#status_modal').modal({
        backdrop: 'static',
        keyboard: false,
        show:     true,
    });
    $('#status_footer').hide();
    $('#status_header').html('<h4><span class="glyphicon glyphicon-refresh"></span> Loading</h4>');
    $('#status_body').html('Please wait while data is synced...');
    $('#status_modal').modal('show');
}

// Hides the loading status again
function hideStatusModal(){
    $('#status_modal').modal('hide');
}

// Load all user information
function getUsers(){
    // get the user information from the database
    // get the account information
    $.ajax({
        type: "GET",
        url: "api/user.php?action=get_all_user_detailed",
        cache: false,
        beforeSend: showStatusModal,
        success: function(data){
            var json = $.parseJSON(data);
            users = json;
            updateUsers(users);
        },
        complete: hideStatusModal,
    });
}

// Load all user role information
function getUserGroups(){
    $.ajax({
        type: "GET",
        url: "api/user.php?action=get_user_groups",
        cache: false,
        beforeSend: showStatusModal,
        success: function(data){
            var json = $.parseJSON(data);
            updateUserGroups(json);
        },
        complete: hideStatusModal
    });
}

// Update users table
function updateUsers(data){
    let users          = data.users;
    let table_data     = new Array();

    // update global currency setting
    displayCurrency    = data.currency;

    let i = 0;
    for(let user_id in users){
        table_data[i] = [ 
            users[user_id].id,
            users[user_id].first_name,
            users[user_id].last_name,
            users[user_id].mobile,
            users[user_id].email,
            users[user_id].license,
            (users[user_id].total_payment - users[user_id].total_heat_cost).toFixed(2),
            users[user_id].status,
            users[user_id].locked,
            displayCurrency
        ];
        i++;
    }
    if(table_data.length<1){
        table_data == null;
    }

    // The users table
    if(userTable != null){
        userTable.fnDestroy();
    }

    userTable = $('#user_table').DataTable({
        "data": table_data,
        "columnDefs": [
            { "className": "text-right", "targets": [ 6 ]},
            { "visible": false, "targets": [ 0 ]},
            { "render": licenseRender, "targets": [ 5 ]},
            { "render": currencyRender, "targets": [ 6 ]},
            { "render": lockRender, "targets": [ 8 ]},
            { "render": userRoleRender, "targets": [ 7 ]},
        ],
        "columns": [
            { "title": "id" },
            { "title": "First Name" },
            { "title": "Last Name" },
            { "title": "Mobile" },
            { "title": "E-mail" },
            { "title": "License" },
            { "title": "Balance" },
            { "title": "User Group" },
            { "title": "Locked" },
        ],
        "select": {
            "items":"row",
            "style":"single",
            "className":"selected",
        },
        "order": [[ 1, "asc" ]],
        "searching": false,
        "paging": false,
        "info" : false,
        "retrieve": true,
        "scrollY": "350px",
    });

}

// licenseRender renders the 1|0 boolean into icons
function licenseRender( data, type, row, meta ) {
    if(data == 1){
        return '<i class="fa fa-check"></i>';
    }else{
        return '<i class="fa fa-times"></i>';
    }
}

// currencyRender adds the currency to a number
function currencyRender( data, type, row, meta ) {
    return data + " " + displayCurrency;
}

// lockRender adds the currency to a number
function lockRender( isLocked, type, row, meta ) {
    if(isLocked == 0){
        return '<button class="btn btn-xs btn-success" onClick="lockUser(' + meta.row + ')"><i class="fas fa-unlock-alt"></i></button>';
    }else if(isLocked == 1){
        return '<button class="btn btn-xs btn-danger" onClick="unlockUser(' + meta.row + ')"><i class="fas fa-lock"></i></button>';
    }else{
        return "-";
    }
}

// userRoleRender adds a dropdown to the User Group
function userRoleRender( data, type, row, meta ) {
    var content = content = '<select name="userGroup" onchange="updateUserGroup('+meta.row+', this.options[this.selectedIndex].value)">';
    for(i in userGroups){
        var group = userGroups[i];
        if(group.user_group_id == data){
            content += '<option value=' + data + ' selected>' + group.user_group_name + '</option>';
        }else{
            content += '<option value=' + group.user_group_id + '>' + group.user_group_name + '</option>';
        }
    }
    content += '</select>';
    return content;
}

// Lock a given user
function lockUser(row_index){
    var table = $('#user_table').DataTable();
    var row   = table.row(row_index);
    var rowData = row.data();
    var data = new Object();
    data['user_id'] = rowData[0];
    $.ajax({
        type: "POST",
        url: "api/user.php?action=lock_user",
        cache: false,
        data: JSON.stringify(data),
        beforeSend: showStatusModal,
        success: function(){
            // update the user in the table
            rowData[8] = 1;
            row.invalidate();
        },
        complete: hideStatusModal,
    });
}

// Unlock a given user
function unlockUser(row_index){
    var table    = $('#user_table').DataTable();
    var row      = table.row(row_index);
    var rowData  = row.data();
    var data     = new Object();
    data['user_id'] = rowData[0];
    $.ajax({
        type: "POST",
        url: "api/user.php?action=unlock_user",
        cache: false,
        data: JSON.stringify(data),
        beforeSend: showStatusModal,
        success: function(){
            // update the user in the table
            rowData[8] = 0;
            row.invalidate();
        },
        complete: hideStatusModal,
    });
}

// Update the status of a user
function updateUserGroup(row_index, user_group_id){
    var table   = $('#user_table').DataTable();
    var row     = table.row(row_index);
    var rowData = row.data();
    var data    = new Object();
    data['user_id']   = rowData[0];
    data['status_id'] = user_group_id;
    $.ajax({
        type: "POST",
        url: "api/user.php?action=change_user_group_membership",
        cache: false,
        data: JSON.stringify(data),
        beforeSend: showStatusModal,
        success: function(data){
            //getUsers();
        },
        complete: hideStatusModal,
    });
}

// Update the user groups table
function updateUserGroups(data){
    let userGroups = data.user_groups;
    let table_data = new Array();

    displayCurrency = data.currency;

    let i = 0;
    for(let id in userGroups){
        table_data[i] = [ 
            userGroups[id].user_group_id,
            userGroups[id].user_group_name,
            userGroups[id].user_group_description,
            userGroups[id].user_role_id,
            userGroups[id].user_role_name,
            userGroups[id].user_role_description,
            userGroups[id].price_id,
            userGroups[id].price_min.toFixed(2),
            userGroups[id].price_description,
            displayCurrency
        ];
        i++;
    }
    if(table_data.length<1){
        table_data = null;
    }

    // if table exists update it:
    if($.fn.dataTable.isDataTable('#user_type_table')){
        let userGroupTable = $('#user_type_table').DataTable();
        userGroupTable.clear();
        userGroupTable.rows.add(table_data);
        userGroupTable.draw();
        return;
    }

    // else create a new table from scratch
    let userGroupTable = $('#user_type_table').DataTable( {
        "data": table_data,
        "columnDefs": [
            { "className": "text-left" , "targets": [ 0, 1, 2, 3, 4]},
            { "className": "text-right", "targets": [ 7 ]},
            { "visible":   false  , "targets": [ 0, 2, 3, 4, 6, 8 ]},
            { "render": currencyRender, "targets": [ 7 ]},
        ],
        "columns": [
            { "title": "user_group_id" },
            { "title": "User Group" },
            { "title": "user_group_description"},
            { "title": "user_role_id" },
            { "title": "User Role" },
            { "title": "User Role" },
            { "title": "price_id" },
            { "title": "Price " + displayCurrency + "/min" },
            { "title": "price_description"},
        ],
        "select": {
            "items":"row",
            "style":"single",
            "className":"selected",
        },
        "order": [[ 1, "asc" ]],
        "searching":     false,
        "paging":        false,
        "info":          false,
        "retrieve":      true,
        "scrollY":       "350px"
    });
    userGroupTable.on('select', updateUserGroupNavigation);
    userGroupTable.on('deselect', updateUserGroupNavigation);
}

function updateUserGroupNavigation(){
    var selected = getSelectedUserGroupTableRow();
    if(Object.keys(selected).length < 1){
        $("#btnViewUserGroup").attr("disabled", "disabled");
        $("#btnEditUserGroup").attr("disabled", "disabled");
        $("#btnNewUserGroup").removeAttr("disabled");
        $("#btnDeleteUserGroup").attr("disabled", "disabled");
    }else{
        $("#btnViewUserGroup").removeAttr("disabled");
        $("#btnEditUserGroup").removeAttr("disabled");
        $("#btnNewUserGroup").removeAttr("disabled");
        $("#btnDeleteUserGroup").removeAttr("disabled");
    }
}

function showUserGroupView(){
    userGroup(getSelectedUserGroupTableRow(), false);
}

function editUserGroupView(){
    userGroup(getSelectedUserGroupTableRow(), true);
}

function createUserGroupView(){
    userGroup({}, true);
}

function deleteUserGroupView(){
    // TODO check that un-used and report error
    group = getSelectedUserGroupTableRow();
    $.ajax({
        type: "POST",
        url: "api/user.php?action=delete_user_group",
        cache: false,
        beforeSend: showStatusModal,
        data: JSON.stringify(group),
        success: function(data){
            getUserGroups();
            var table = $('#user_type_table').DataTable();
            table.rows().deselect();
            updateUserGroupNavigation();
        },
        error: function(){
            alert("Cannot delete user group. It might still be used?");
            var table = $('#user_type_table').DataTable();
            table.rows().deselect();
            updateUserGroupNavigation();
        },
        complete: hideModal
    });
}

function getSelectedUserGroupTableRow(){
    var data = {};
    var table = $('#user_type_table').DataTable();
    if(table.rows({ selected: true})){
        var rowData = table.rows( { selected: true } ).data();
        if(rowData.length == 1){   
            data['user_group_id']     = rowData[0][0];
            data['user_group_name']   = rowData[0][1];
            data['user_group_description'] = rowData[0][2];
            data['user_role_id']      = rowData[0][3];
            data['user_role_name']    = rowData[0][4];
            data['user_role_description'] = rowData[0][5]
            data['price_id']          = rowData[0][6];
            data['price_min']         = rowData[0][7];
            data['price_description'] = rowData[0][8];
            data['currency']          = rowData[0][9];
        }
    }
    return data;
}

function userGroup(data, editable){
    // get the existing user roles
    var userRoles;
    $.ajax({
        type: "GET",
        url: "api/user.php?action=get_user_roles",
        cache: false,
        success: function(response){
            userRoles = $.parseJSON(response);
            displayUserGroupApp(userRoles, data, editable);
        }
    });
}

function displayUserGroupApp(userRoles, userGroup, editable){
    console.log(userRoles);
    console.log(userGroup);
    console.log(editable);

    // are the values unset?
    let newGroup = false;
    if(Object.keys(userGroup).length < 1){
        newGroup = true;
    }

    // load the dialog
    $('#modal_content').load("res/view_user_group.vue", function() {
        showModal();

        userGroupApp = new Vue({
            el: '#user_group',
            data: {
                "group": userGroup,
                "newGroup": newGroup,
                "editmode": editable,
                "userRoles": userRoles,
                "currency": displayCurrency,
                "errors": [],
            },
            methods: {
                save: function(event){
                    // validate usergroup
                    this.errors = validateUserGroup(this.group);
                    if(this.errors.length > 0){
                        return;
                    }
                    // save usergroup
                    hideModal();
                    saveUserGroup(this.group);
                },
                cancel: function(event){
                    hideModal();
                },
                userRoleChange: function(event){
                    this.group.user_role_description = userRoles[this.group.user_role_id].user_role_description;
                }
            }
        });
    });
}

function validateUserGroup(group){
    var errors = [];

    if(! group['user_group_name'] || ! group['user_group_name'].match(/^[A-Za-z]+$/)){
        errors.push("User Group Name must be a single word (no special characters allowed)");
    }
    if(! group['user_group_description'] || ! group['user_group_description'].match(/^[A-Za-z\s]+$/)){
        errors.push("User Group Description must be given (no special characters allowed)");
    }
    if(! group['user_role_id']){
        errors.push("User Role must be selected");
    }
    if(! group['price_min'] || ! group['price_min'].toString().match(/^\d+(\.\d{1,2})?$/)){
        errors.push("Price must me a valid number (e.g., 2.50)");
    }
    if(! group['price_description'] || ! group['price_description'].match(/^[A-Za-z\s]+$/)){
        errors.push("Price Description must be given (no special characters allowed)");
    }

    return errors;
}

function saveUserGroup(group){
    $.ajax({
        type: "POST",
        url: "api/user.php?action=save_user_group",
        cache: false,
        beforeSend: showStatusModal,
        data: JSON.stringify(group),
        success: function(data){
            getUserGroups();
        },
        complete: hideModal
    });
}

function showModal(){
    $('#modal').modal({
        backdrop: 'static',
        keyboard: false,
        show:     true,
    });
    $('#modal').modal('show');
}

function hideModal(){
    $('#modal').modal('hide');
    $('.modal-backdrop').remove();
}