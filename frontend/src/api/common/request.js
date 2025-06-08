export default class Request {
  static _handleResponse(url, method, response, resolve, reject) {
    console.debug(method, url, "received response");

    if (response.status != 200) {
      console.warn(method, url, "error", response.status);
      // handle a few specific error messages
      response.text().then((bodyText) => {
        if (response.status === 500 && bodyText.includes("Proxy erro")) {
          reject(["Internal Server Error: Your request cannot be proxied to the backend. Is the API backend reachable?"]);
        }
        reject(["Error: Your request could not be handled. There is an issue with your request or with the backend."]);
      })
      .catch((error) => {
        console.log("Error, cannot get response body", error);
        reject(["Error: Your request could not be handled. There is an issue with your request or with the backend."]);
      })
      return; // Stop further processing
    }

    response.json()
      .then((data) => {
        console.debug(method, url, "response data:", data);
        // Combined check: data.ok is primary, data.status.ok is secondary
        // This covers cases where 'data.ok' is the sole indicator (like original postRequest)
        // and cases where 'data.status.ok' might also be used (like original getRequest)
        if (data.ok || (data.status && data.status.ok)) {
          resolve(data.data);
        } else {
          console.warn(method, url, "response not ok, due to:", data.msg);
          reject([data.msg]);
        }
      })
      .catch((error) => {
          console.warn(method, url, "cannot parse server response", error);
          reject([error]);
      });
  }

  static getRequest(url) {
    console.debug("GET", url);
    return new Promise((resolve, reject) => {
      fetch(url, {
        method: "GET",
        cache: "no-cache",
      })
        .then((response) => {
          Request._handleResponse(url, "GET", response, resolve, reject);
        })
        .catch((error) => {
          console.error("GET", url, "fetch failed", error);
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
          Request._handleResponse(url, "POST", response, resolve, reject);
        })
        .catch((error) => {
          console.warn("POST", url, "fetch failed", error);
          reject([error]);
        });
    });
  }
}
