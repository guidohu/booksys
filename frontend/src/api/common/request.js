export default class Request {
  static getRequest(url) {
    return new Promise((resolve, reject) => {
      fetch(url, {
        method: "GET",
        cache: "no-cache",
      })
        .then((response) => {
          response
            .json()
            .then((data) => {
              console.debug(url, "response data:", data);
              if (data.ok) {
                resolve(data.data);
              } else {
                console.log(url, "response not ok, due to:", data.msg);
                reject([data.msg]);
              }
            })
            .catch((error) => {
              console.error(url, "cannot parse server response", error);
              reject([error]);
            });
        })
        .catch((error) => {
          console.error(url, "request failed", error);
          reject([error]);
        });
    });
  }

  static postRequest(url, payload) {
    return new Promise((resolve, reject) => {
      fetch(url, {
        method: "POST",
        cache: "no-cache",
        body: JSON.stringify(payload),
      })
        .then((response) => {
          response
            .json()
            .then((data) => {
              console.log(url, "response data:", data);
              if (data.ok) {
                resolve(data.data);
              } else {
                console.log(url, "response not ok, due to:", data.msg);
                reject([data.msg]);
              }
            })
            .catch((error) => {
              console.error(url, "cannot parse server response", error);
              reject([error]);
            });
        })
        .catch((error) => {
          console.error(url, "request failed", error);
          reject([error]);
        });
    });
  }
}
