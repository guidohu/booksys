export default class Request {
  static getRequest(url) {
    console.debug("GET", url);
    return new Promise((resolve, reject) => {
      fetch(url, {
        method: "GET",
        cache: "no-cache",
      })
        .then((response) => {
          response
            .json()
            .then((data) => {
              console.debug("GET", url, "response data:", data);
              if (data.ok || data.status.ok) {
                resolve(data.data);
              } else {
                console.warn("GET", url, "response not ok, due to:", data.msg);
                reject([data.msg]);
              }
            })
            .catch((error) => {
              console.warn("GET", url, "cannot parse server response", error);
              reject([error]);
            });
        })
        .catch((error) => {
          console.error("GET", url, "request failed", error);
          reject([error]);
        });
    });
  }

  static postRequest(url, payload) {
    console.debug("POST", url, payload);
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
              console.debug("POST", url, "response data:", data);
              if (data.ok) {
                resolve(data.data);
              } else {
                console.warn("POST", url, "response not ok, due to:", data.msg);
                reject([data.msg]);
              }
            })
            .catch((error) => {
              console.warn("POST", url, "cannot parse server response", error);
              reject([error]);
            });
        })
        .catch((error) => {
          console.warn("POST", url, "request failed", error);
          reject([error]);
        });
    });
  }
}
