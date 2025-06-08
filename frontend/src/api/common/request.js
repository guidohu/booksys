export default class Request {
  static getRequest(url) {
    console.debug("GET", url);
    return new Promise((resolve, reject) => {
      fetch(url, {
        method: "GET",
        cache: "no-cache",
      })
        .then((response) => {
          if (response.status != 200) {
            // handle a few specific error messages
            response.text().then((bodyText) => {
              if (response.status === 500 && bodyText.includes("Proxy erro")) {
                reject(["Internal Server Error: Your request cannot be proxied to the backend. Is the API backend reachable?"]);
              }
              reject(["Error: Your request could not be handled. There is an issue with your request or with the backend."]);
            })
            .catch((error) => {
              console.log("Internal Server Error, cannot get response body", error);
              reject(["Error: Your request could not be handled. There is an issue with your request or with the backend."]);
            })
          } else {
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
          }
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
          if (response.status != 200) {
            // handle a few specific error messages
            response.text().then((bodyText) => {
              if (response.status === 500 && bodyText.includes("Proxy erro")) {
                reject(["Internal Server Error: Your request cannot be proxied to the backend. Is the API backend reachable?"]);
              }
              reject(["Error: Your request could not be handled. There is an issue with your request or with the backend."]);
            })
            .catch((error) => {
              console.log("Internal Server Error, cannot get response body", error);
              reject(["Error: Your request could not be handled. There is an issue with your request or with the backend."]);
            })
          } else {
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
          }
        })
        .catch((error) => {
          console.warn("POST", url, "request failed", error);
          reject([error]);
        });
    });
  }
}
