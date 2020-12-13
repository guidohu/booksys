class BooksysBeacon {

    // Sends periodic beacons to the backend beacon API to 
    // keep the user connected and not logged out (e.g. while using the stop watch)
    constructor(config, callbacks){
        if(config!=null && config.timeout){
            this.timeout = config.timeout;
        }else{
            this.timeout = 10000;
        }
        if(callbacks!=null && callbacks.success!=null){
            this.cbSuccess = callbacks.success;
        }
        if(callbacks!=null && callbacks.redirect!=null){
            this.cbRedirect = callbacks.redirect;
        }
        if(callbacks!=null && callbacks.failure!=null){
            this.cbError = callbacks.failure;
        }

        setTimeout(this.sendBeacon.bind(this), 1000);
    }

    // Send a beacon every X seconds specified by the timeout parameter
    sendBeacon(){
        // send the beacon
        let that = this;
        $.ajax({
            type: "GET",
            url: "api/beacon.php",
            timeout: this.timeout,
            success: function(resp){
                let json = $.parseJSON(resp);
                if(json.ok){
                    that.cbSuccess();
                }else{
                    if(json.msg == "you are not logged in"){
                        that.cbRedirect("/index.html");
                        return;
                    }
                }
                setTimeout(that.sendBeacon.bind(that), that.timeout);
                //beacon = setTimeout("sendBeacon()", that.timeout);
            },
            error: function(data, err_txt, err_thrown){
                that.cbError();
                setTimeout(that.sendBeacon.bind(that), 5000);
                //beacon = setTimeout("sendBeacon()", 5000);
            }
        });
    }
}