export default class Backend {

  static getStatus() {
    console.log("backend/getStatus called");
    return new Promise((resolve, reject) => {
      fetch('/api/backend.php?action=get_status', {
        method: "GET",
        cache: "no-cache",
      })
      .then(response => {
        response.json()
          .then(data => {
            if(data.ok){
              console.log("backend/getStatus response data:", data);
              resolve(data.data);
            }else{
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