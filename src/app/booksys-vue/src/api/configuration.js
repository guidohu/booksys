export class Configuration {

  static getConfiguration(cbSuccess, cbFailure){
    fetch('/api/configuration.php?action=get_configuration', {
      method: 'GET',
      cache: 'no-cache'
    })
    .then(response => {
      return response.json(); 
    })
    .then(data => {
      if(data.ok == true){
        cbSuccess(data.data);
      }else{
        cbFailure("Configuration retrieval not successful");
      }
    })
    .catch(error => {
      console.log('Configuration/getConfiguration:', error)
      cbFailure(
        error
      )
    })
  }

  static setConfiguration(params){
    return new Promise((resolve, reject) => {
      const request = {
        location_time_zone       : params.location_time_zone,
        location_longitude       : params.location_longitude,
        location_latitude        : params.location_latitude,
        location_map             : params.location_map,
        location_address         : params.location_address,
        currency                 : params.currency,
        payment_account_owner    : params.payment_account_owner,
        payment_account_iban     : params.payment_account_iban,
        payment_account_bic      : params.payment_account_bic,
        payment_account_comment  : params.payment_account_comment,
        smtp_sender              : params.smtp_sender,
        smtp_server              : params.smtp_server,
        smtp_username            : params.smtp_username,
        recaptcha_privatekey     : params.recaptcha_privatekey,
        recaptcha_publickey      : params.recaptcha_publickey,
      };

      // only set password in case it is really given
      if(params.smtp_password != 'hidden'){
        request.smtp_password = params.smtp_password
      }
       
      fetch('/api/configuration.php?action=set_configuration', {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(request)
      })
      .then(response => {
        response.json()
          .then(data => {
            console.log("change configuration response data:", data);
            if(data.ok){
              resolve(data.data);
            }else{
              console.log("configuration/setCustomizationParameters: Cannot store configuration, due to:", data.msg);
              reject([data.msg]);
            }
          })
          .catch(error => {
            console.error("configuration/setCustomizationParameters: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error("configuration/setCustomizationParameters", error);
        reject([error]);
      })
    })
  }

  static needsDbUpdate(cbSuccess, cbFailure){
    fetch('/api/backend.php?action=admin_check_database_update', {
      method: 'GET',
      cache: 'no-cache'
    })
    .then(response => {
      return response.json()
    })
    .then(data => {
      cbSuccess(data)
    })
    .catch(error =>{
      cbFailure(error)
    })
  }

  static getDbVersion(cbSuccess, cbFailure){
    fetch('/api/backend.php?action=get_version', {
      method: 'GET',
      cache: 'no-cache'
    })
    .then(response => {
      return response.json()
    })
    .then(data => {
      cbSuccess(data)
    })
    .catch(error => {
      cbFailure(error)
    })
  }

  static updateDb(cbSuccess, cbFailure){
    fetch('/api/backend.php?action=update_database', {
      method: 'GET',
      cache: 'no-cache'
    })
    .then(response => {
      return response.json()
    })
    .then(data => {
      cbSuccess(data)
    })
    .catch(error => {
      cbFailure(error)
    })
  }
}