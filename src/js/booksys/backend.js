class BooksysBackend {
    static getStatus(actions) {
        // get the backend status information
        $.ajax({
            type: "GET",
            cache: false,
            url: "api/backend.php?action=get_status",
            success: function(data){
                var json = $.parseJSON(data);
                if(json.ok){
                    let redirect = null;
                    if(json.redirect != null){
                        redirect = json.redirect;
                    }
                    actions.green(json.message, null);
                }else{
                    let redirect = null;
                    if(json.redirect != null){
                        redirect = json.redirect;
                    }
                    actions.red(json.message, redirect);
                }
            },
            error: function(data, text, errorCode){
                actions.red(text);
            }
        });
    }
}