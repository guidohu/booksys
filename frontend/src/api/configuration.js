import Request from "@/api/common/request.js";

export default class Configuration {
  /**
   * Returns the recaptcha key.
   */
  static getRecaptchaKey() {
    console.log("configuration/getRecaptchaKey called");
    return Request.getRequest('/api/v2/configuration/recaptcha-key');
  }

  /**
   * Returns the path to the logo.
   */
  static getLogoFile() {
    console.log("configuration/getLogoFile called");
    return Request.getRequest('/api/v2/configuration/logo');
  }

  // static getDbConfig() {
  //   console.log("configuration/getDbConfig called");
  //   return Request.getRequest('/api/v2/database/config');
  // }

  static setDbConfig(config) {
    console.log("configuration/setDbConfig called with:", config);
    const requestData = {
      db_server: config.host,
      db_name: config.name,
      db_user: config.user,
      db_password: config.password,
    };
    return Request.postRequest("/api/v2/database/setup", requestData);
  }

  static getConfiguration() {
    return Request.getRequest("/api/v2/configuration/list");
  }

  static setConfiguration(params) {
    const request = {
      logo_file: params.logo_file,
      engine_hour_format: params.engine_hour_format,
      fuel_payment_type: params.fuel_payment_type,
      location_time_zone: params.location_time_zone,
      location_longitude: params.location_longitude,
      location_latitude: params.location_latitude,
      location_map: params.location_map,
      location_address: params.location_address,
      currency: params.currency,
      payment_account_owner: params.payment_account_owner,
      payment_account_iban: params.payment_account_iban,
      payment_account_bic: params.payment_account_bic,
      payment_account_comment: params.payment_account_comment,
      smtp_sender: params.smtp_sender,
      smtp_server: params.smtp_server,
      smtp_username: params.smtp_username,
      recaptcha_privatekey: params.recaptcha_privatekey,
      recaptcha_publickey: params.recaptcha_publickey,
      mynautique_enabled: params.mynautique_enabled,
    };

    if (params.mynautique_enabled === true) {
      request['mynautique_user'] = params.mynautique_user
      request['mynautique_password'] = params.mynautique_password
      request['mynautique_boat_id'] = parseInt(params.mynautique_boat_id, 10)
      request['mynautique_fuel_capacity'] = parseInt(params.mynautique_fuel_capacity, 10)
      request['mynautique_api_key'] = params.mynautique_api_key
    }

    // only set password in case it is really given
    if (params.smtp_password != "hidden") {
      request.smtp_password = params.smtp_password;
    }
    // only set the myNautique password in case it is given
    if (params.mynautique_password != "hidden") {
      request.mynautique_password = params.mynautique_password;
    }

    return Request.postRequest("/api/v2/admin/configuration/set", request);
  }

  static setMyNautiqueConfig(config) {
    console.log("configuration/setMyNautiqueConfig called with:", config);
    const requestData = {
      mynautique_enabled: config.enabled,
      mynautique_user: config.mynautiqueUser,
      mynautique_password: config.mynautiquePassword,
    };

    return Request.postRequest("/api/v2/mynautique/credentials/setup", requestData);
  }

}
