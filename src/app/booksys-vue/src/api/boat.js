export default class Boat {

  static getEngineHourLog(){
    console.log("api/getEngineHourLog called");
    return new Promise((resolve, reject) => {
      fetch('/api/boat.php?action=get_engine_hours_log', {
        method: "GET",
        cache: "no-cache",
      })
      .then(response => {
        response.json()
          .then(data => {
            console.log("Boat/getEngineHourLog response data:", data);
            if(data.ok){
              resolve(data.data);
            }else{
              console.log("Boat/getEngineHourLog: Cannot retrieve engine hour logs, due to:", data.msg);
              reject([data.msg]);
            }
          })
          .catch(error => {
            console.error("Boat/getEngineHourLog: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error("Boat/getEngineHourLog", error);
        reject([error]);
      })
    })
  }

  static getEngineHourLogLatest(){
    console.log("api/getEngineHourLogLatest called");
    return new Promise((resolve, reject) => {
      fetch('/api/boat.php?action=get_engine_hours_latest', {
        method: "GET",
        cache: "no-cache",
      })
      .then(response => {
        response.json()
          .then(data => {
            console.log("Boat/getEngineHourLogLatest response data:", data);
            if(data.ok){
              resolve(data.data[0]);
            }else{
              console.log("Boat/getEngineHourLogLatest: Cannot retrieve engine hour logs, due to:", data.msg);
              reject([data.msg]);
            }
          })
          .catch(error => {
            console.error("Boat/getEngineHourLogLatest: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error("Boat/getEngineHourLogLatest", error);
        reject([error]);
      })
    })
  }

  /**
   * Adds a new engine hour entry.
   * @param {} engineHourEntry consisting of:
   * {
   *   user_id: ...
   *   engine_hours_before: ...
   *   engine_hours_after: ...
   *   type: ...                 // currently 0 for private, 1 for course
   * }
   */
  static addEngineHours(engineHourEntry){
    console.log("api/addEngineHours called");
    return new Promise((resolve, reject) => {
      fetch('/api/boat.php?action=update_engine_hours', {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(engineHourEntry)
      })
      .then(response => {
        response.json()
          .then(data => {
            console.log("Boat/addEngineHours response data:", data);
            if(data.ok){
              resolve();
            }else{
              console.log("Boat/addEngineHours: Cannot add engine hour entry, due to:", data.msg);
              reject([data.msg]);
            }
          })
          .catch(error => {
            console.error("Boat/addEngineHours: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error("Boat/addEngineHours", error);
        reject([error]);
      })
    })
  }
}