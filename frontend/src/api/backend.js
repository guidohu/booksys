export default class Backend {

  static getStatus() {
    return new Promise((resolve, reject) => {
      fetch('/api/backend.php?action=get_status', {
        method: "GET",
        cache: "no-cache",
      })
      .then(response => {
        response.json()
          .then(data => {
            if(data.ok){
              resolve(data.data);
            }else{
              console.error("backend/getStatus: Status is not ok", data);
              reject([data.msg]);
            }
          })
          .catch(error => {
            console.error("backend/getStatus: Cannot parse server response", error);
            reject([error]);
          })
      })
      .catch(error => {
        console.error("backend/getStatus", error);
        reject([error]);
      })
    })
  }
}