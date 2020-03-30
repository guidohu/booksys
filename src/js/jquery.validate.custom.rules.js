// Custom rules for jquery validator
// - Written by ghungerbuehler@gmail.com
// License: Do whatever you want with this code

var alphaNumericEmail = new RegExp("^[a-z0-9@_\.-]+$","i");
var mobileNumber      = new RegExp("^[\\+0-9]{6,15}$", "i");

// Alphanumeric test (Char set of Emails)
$.validator.addMethod('alphaNumericEmail', function(value, element){
	if(this.optional(element)){
		return true;
	}else{
	    this.hideErrors();
		if(alphaNumericEmail.test(value)){
			return true;
		}else{
			return false;
		}
	}
}, "Only alphanumeric values a-z and 0-9 including [@-._] are allowed");

$.validator.addMethod('mobileNumber', function(value, element){
	if(this.optional(element)){
		return true;
	}else{
	    this.hideErrors();
		if(mobileNumber.test(value)){
			return true;
		}else{
			return false;
		}
	}
}, "Please enter in the form +41791234567");

$.validator.addMethod('optionSelected', function(value, element){
	if(this.optional(element)){
		return true;
	}else{
		this.hideErrors();
		if($('#'+element.id).val() != -1 && element.value != ''){
			return true;
		}else{
			return false;
		}
	}
}, "Please select an option");

$.validator.addMethod('notEmpty', function(value, element){
	if(this.optional(element)){
		return true;
	}else{
		this.hideErrors();
		if(element.value != ''){
			return true;
		}else{
			return false;
		}
	}
}, "Please add some content");