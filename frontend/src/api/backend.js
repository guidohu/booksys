export const getBackendStatus = () => {
  return new Promise((resolve, reject) => {
    fetch("/api/backend.php?action=get_status", {
      method: "GET",
      cache: "no-cache",
    })
      .then((response) => {
        response
          .json()
          .then((data) => {
            if (data.ok) {
              resolve(data.data);
            } else {
              console.error("getBackendStatus: Status is not ok", data);
              reject([data.msg]);
            }
          })
          .catch((error) => {
            console.error(
              "getBackendStatus: Cannot parse server response",
              error
            );
            reject([error]);
          });
      })
      .catch((error) => {
        console.error("getBackendStatus", error);
        reject([error]);
      });
  });
};
