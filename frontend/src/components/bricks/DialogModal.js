import { createApp, h } from "vue";
import ConfirmDialog from "./dialogTemplates/ConfirmDialog.vue";
import InfoDialog from "./dialogTemplates/InfoDialog.vue";

const defaultConfirmParams = {
  title: "",
  message: "",
  textConfirm: "OK",
  textDeny: "Cancel",
  id: "#dialog",
};

export function confirm(params) {
  const p = { ...defaultConfirmParams, ...params };

  return new Promise((resolve, reject) => {
    const dialog = createApp({
      methods: {
        closeHandler: function (fn, arg) {
          dialog.unmount();
          fn(arg);
        },
      },
      render() {
        return h(ConfirmDialog, {
          title: p.title,
          message: p.message,
          textConfirm: p.textConfirm,
          textDeny: p.textDeny,
          onConfirm: (response) => {
            this.closeHandler(resolve, response);
          },
          onCancel: () => {
            this.closeHandler(reject, new Error("canceled"));
          },
        });
      },
    });
    dialog.mount(p.id);
  });
}

const defaultInfoParams = {
  title: "",
  message: "",
  textConfirm: "OK",
  id: "#dialog",
};

export function info(params) {
  const p = { ...defaultInfoParams, ...params };

  return new Promise((resolve) => {
    const dialog = createApp({
      methods: {
        closeHandler: function (fn, arg) {
          dialog.unmount();
          fn(arg);
        },
      },
      render() {
        return h(InfoDialog, {
          title: p.title,
          message: p.message,
          textConfirm: p.textConfirm,
          onConfirm: (response) => {
            this.closeHandler(resolve, response);
          },
        });
      },
    });
    dialog.mount(p.id);
  });
}
