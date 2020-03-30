// JavaScript Document

// RegExp which define allowed string format for different
// user inputs
var alphaNumeric = new RegExp("^[a-z0-9@_]+$","i");
var phone = new RegExp("^\\+?[0-9]{6,15}$");
var email = new RegExp("[a-z0-9\\.\\-\\_]+@[a-z0-9\\.\\-\\_]+","i");
var number = new RegExp("^[0-9]+\\.?[0-9]+?$");
var date = new RegExp("^[0-9]{4}-[0-9]{2}-[0-9]{2}");

// Returns a URL get parameter's value
function getURLParameter(name) {
    return decodeURI(
        (RegExp(name + '=' + '(.+?)(&|$)').exec(location.search)||[,null])[1]
    );
}

// Sets a class of an element with given id
function setClass(id, style){
	$(id).attr('class', style);
}

// Adds a tooltip to an element with given given id
// the element also changes the style and the text 
// can be provided which will be displayed.
function addToolTipOnFocus(id, text, style){
	setClass(id, style);
	$(id).bind('focusin.balloon', 
  	function(){ 
		$(id).showBalloon({contents: text })
  	});
	$(id).bind('focusout.balloon',
  	function(){ 
		$(id).hideBalloon({contents: text })
  	});
}

// Removes the tooltip of the element again and changes
// the style back to the original
function removeToolTipOnFocus(id, style){
  	setClass(id, style);
	$(id).unbind(".balloon");
}

function verifyUsername(id, style_ok, style_error){
	var content = $(id).val();
		
	/* Test username */
	if(id == '#username'){
	    if(! alphaNumeric.test(content) || content == ''){
		  	addToolTipOnFocus(id, 'only alpha numeric characters and _ or @ allowed', style_error);
			return false;
		}else{
			removeToolTipOnFocus(id, style_ok);
			return true;
		}
	}
}

function verifyPassword(id, id2, style_ok, style_error){
	var content = $(id).val();
	/* Test if passwords are identical */
	if(content != $(id2).val() || content == ''){
	  	addToolTipOnFocus(id2, 'passwords not identical or empty', style_error);	
		setClass(id, style_error);
		return false;
	}else{
	  	removeToolTipOnFocus(id2, style_ok);
		setClass(id, style_ok);
		return true;
	}
}

function verifyMobile(id, style_ok, style_error){
	var content = $(id).val();
	/* Test mobile */
    if(! phone.test(content) || content == ''){
		addToolTipOnFocus(id, 'please enter in the format +41791234567', style_error);
		return false;
	}else{
		removeToolTipOnFocus(id, style_ok);
		return true;
	}
}

function verifyEmail(id, style_ok, style_error){
	var content = $(id).val();
	if(! email.test(content) || content == ''){
	  	addToolTipOnFocus(id, 'please enter in the format abc@example.com', style_error);
		return false;
    }else{
		removeToolTipOnFocus(id, style_ok);
		return true;
	}
}

function verifyNumber(id, style_ok, style_error){
    var content = $(id).val();
	if(! number.test(content) || content == ''){
	  	addToolTipOnFocus(id, 'please enter a value (e.g. 99.95)', style_error);
		return false;
    }else{
		removeToolTipOnFocus(id, style_ok);
		return true;
	}
}

function verifyDate(id, style_ok, style_error){
    var content = $(id).val();
	if(! date.test(content) || content == ''){
	  	addToolTipOnFocus(id, 'please enter a date (e.g. 2013-10-19)', style_error);
		return false;
    }else{
		removeToolTipOnFocus(id, style_ok);
		return true;
	}
}

function verifyEmpty(id, style_ok, style_error){
	var content = $(id).val();
	if(content == ''){
	  	addToolTipOnFocus(id, 'maybe you forgot to enter it', style_error);
		return false;
    }else{
		removeToolTipOnFocus(id, style_ok);
		return true;
	}	
}

function goToURL(url){
    window.location.href = url;
}
