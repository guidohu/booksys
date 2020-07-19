class BooksysLoginCheck {

    constructor(){
        
    }

    // checks whether we are logged in and redirects to the url in case we are not
    check(loggedInCb, loggedOutCb){
        let that = this;
        $.ajax({
            type: "GET",
            url: 'api/login.php?action=isLoggedIn',
            success: function(data){
                var json = $.parseJSON(data);
                if(json.loggedIn){
                    // we are already logged in:
                    if(loggedInCb != null){
                        loggedInCb();
                    }
                }else{
                    // we are not logged in and thus clean-up
                    // session related things
                    that.removeUserInStorage();
                    if(loggedOutCb != null){
                        loggedOutCb();
                    }
                }
            },
            error: function(data, text, errorCode){
                // something went wrong
                console.error(data);
                console.error(text);
                console.error(errorCode);
                console.error("Cannot call isLoggedIn, error returned");
            }
        });
    }

    // remove user information in local storage
    removeUserInStorage(){
        if (typeof(Storage) !== "undefined") {
            localStorage.removeItem('userRole');
            localStorage.removeItem('userRoleName');
        }
    }


    
}