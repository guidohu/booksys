export default class Log {
  /**
   * Get all log entries from the backend
   */
  static getLogs(){
    console.log('Log/getLogs called');
    return new Promise((resolve, reject) => {
      fetch('/api/log.php?action=get_logs', {
        method: "GET",
        cache: "no-cache"
      })
      .then(response => {
        response.json()
          .then(data => {
            console.log("get logs response data:", data);
            if(data.ok){
              resolve(data.data);
            }else{
              console.log("log/getLogs: Cannot get logs, due to:", data.msg);
              reject([data.msg]);
            }
          })
          .catch(error => {
            console.error("log/getLogs: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error("log/getLogs", error);
        reject([error]);
      })
    })
  }
}