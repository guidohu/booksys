export const uploadLogo = (file) => {
  return new Promise((resolve, reject) => {
    console.log("Uploading:", file);
    const formData = new FormData();
    formData.append("logo", file);

    fetch("/api/resources.php?action=upload_logo", {
      method: "POST",
      body: formData,
    })
      .then((response) => {
        response
          .json()
          .then((data) => {
            if (data.ok) {
              console.log("resources/uploadLogo: response:", data.msg);
              resolve(data.data.uri);
            } else {
              console.log(
                "resources/uploadLogo: issue while uploading logo:",
                data.msg
              );
              reject([data.msg]);
            }
          })
          .catch((error) => {
            console.log("resources/uploadLogo: Cannot parse JSON:", error);
            reject([error]);
          });
      })
      .catch((error) => {
        reject([error]);
      });
  });
};
