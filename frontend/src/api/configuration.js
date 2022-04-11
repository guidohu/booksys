import Request from "@/api/common/request.js";

export default class Configuration {
  /**
   * Returns the recaptcha key.
   */
  static getRecaptchaKey() {
    console.log("configuration/getRecaptchaKey called");
    return Request.getRequest('/api/configuration.php?action=get_recaptcha_key');
  }

  /**
   * Returns the path to the logo.
   */
  static getLogoFile() {
    console.log("configuration/getLogoFile called");
    return Request.getRequest('/api/configuration.php?action=get_logo_file');
  }

  static getDbConfig() {
    console.log("configuration/getDbConfig called");
    return Request.getRequest('/api/configuration.php?action=get_db_config');
  }

  static setDbConfig(config) {
    console.log("configuration/setDbConfig called with:", config);
    return new Promise((resolve, reject) => {
      const requestData = {
        db_server: config.host,
        db_name: config.name,
        db_user: config.user,
        db_password: config.password,
      };

      fetch("/api/configuration.php?action=setup_db_config", {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(requestData),
      })
        .then((response) => {
          response
            .json()
            .then((data) => {
              console.log("configuration/setDbConfig response data:", data);
              if (data.ok) {
                resolve();
              } else {
                console.log(
                  "configuration/setDbConfig: Cannot set db config, due to:",
                  data.msg
                );
                reject([data.msg]);
              }
            })
            .catch((error) => {
              console.error(
                "configuration/setDbConfig: Cannot parse server response",
                error
              );
              reject([error]);
            });
        })
        .catch((error) => {
          console.error("configuration/setDbConfig", error);
          reject([error]);
        });
    });
  }

  static getConfiguration(cbSuccess, cbFailure) {
    fetch("/api/configuration.php?action=get_configuration", {
      method: "GET",
      cache: "no-cache",
    })
      .then((response) => {
        return response.json();
      })
      .then((data) => {
        if (data.ok == true) {
          cbSuccess(data.data);
        } else {
          cbFailure("Configuration retrieval not successful");
        }
      })
      .catch((error) => {
        console.log("Configuration/getConfiguration:", error);
        cbFailure(error);
      });
  }

  static setConfiguration(params) {
    return new Promise((resolve, reject) => {
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
        mynautique_user: params.mynautique_user,
        mynautique_boat_id: params.mynautique_boat_id,
        mynautique_fuel_capacity: params.mynautique_fuel_capacity,
      };

      // only set password in case it is really given
      if (params.smtp_password != "hidden") {
        request.smtp_password = params.smtp_password;
      }
      // only set the myNautique password in case it is given
      if (params.mynautique_password != "hidden") {
        request.mynautique_password = params.mynautique_password;
      }

      fetch("/api/configuration.php?action=set_configuration", {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(request),
      })
        .then((response) => {
          response
            .json()
            .then((data) => {
              console.log("change configuration response data:", data);
              if (data.ok) {
                resolve(data.data);
              } else {
                console.log(
                  "configuration/setCustomizationParameters: Cannot store configuration, due to:",
                  data.msg
                );
                reject([data.msg]);
              }
            })
            .catch((error) => {
              console.error(
                "configuration/setCustomizationParameters: Cannot parse server response",
                error
              );
              reject([error]);
            });
        })
        .catch((error) => {
          console.error("configuration/setCustomizationParameters", error);
          reject([error]);
        });
    });
  }

  static setMyNautiqueConfig(config) {
    console.log("configuration/setMyNautiqueConfig called with:", config);
    return new Promise((resolve, reject) => {
      const requestData = {
        mynautique_enabled: config.enabled,
        mynautique_user: config.mynautiqueUser,
        mynautique_password: config.mynautiquePassword,
      };

      fetch("/api/configuration.php?action=setup_mynautique_config", {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(requestData),
      })
        .then((response) => {
          response
            .json()
            .then((data) => {
              console.log(
                "configuration/setMyNautiqueConfig response data:",
                data
              );
              if (data.ok) {
                resolve();
              } else {
                console.log(
                  "configuration/setMyNautiqueConfig: Cannot set db myNautique, due to:",
                  data.msg
                );
                reject([data.msg]);
              }
            })
            .catch((error) => {
              console.error(
                "configuration/setMyNautiqueConfig: Cannot parse server response",
                error
              );
              reject([error]);
            });
        })
        .catch((error) => {
          console.error("configuration/setMyNautiqueConfig", error);
          reject([error]);
        });
    });
  }

  static needsDbUpdate(cbSuccess, cbFailure) {
    fetch("/api/backend.php?action=admin_check_database_update", {
      method: "GET",
      cache: "no-cache",
    })
      .then((response) => {
        return response.json();
      })
      .then((data) => {
        cbSuccess(data);
      })
      .catch((error) => {
        cbFailure(error);
      });
  }

  static getDbVersion() {
    return new Promise((resolve, reject) => {
      fetch("/api/backend.php?action=get_version", {
        method: "GET",
        cache: "no-cache",
      })
        .then((response) => {
          return response.json();
        })
        .then((data) => {
          resolve(data);
        })
        .catch((error) => {
          reject([error]);
        });
    });
  }

  static updateDb(cbSuccess, cbFailure) {
    fetch("/api/backend.php?action=update_database", {
      method: "GET",
      cache: "no-cache",
    })
      .then((response) => {
        return response.json();
      })
      .then((data) => {
        cbSuccess(data);
      })
      .catch((error) => {
        cbFailure(error);
      });
  }
}
