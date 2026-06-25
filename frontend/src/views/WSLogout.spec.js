import { mount } from "@vue/test-utils";
import { describe, it, expect, vi } from "vitest";
import { createStore } from "vuex";
import WSLogout from "./WSLogout.vue";

describe("WSLogout.vue", () => {
  it("dispatches login/logout on mount and handles success", async () => {
    const logoutAction = vi.fn().mockResolvedValue();

    const store = createStore({
      modules: {
        login: {
          namespaced: true,
          actions: {
            logout: logoutAction,
          },
        },
      },
    });

    mount(WSLogout, {
      global: {
        plugins: [store],
        stubs: {
          CardModule: true,
          RouterLink: true,
        }
      },
    });

    expect(logoutAction).toHaveBeenCalled();
  });

  it("logs an error if logout fails", async () => {
    const consoleSpy = vi.spyOn(console, "error").mockImplementation(() => {});

    const store = createStore({
      modules: {
        login: {
          namespaced: true,
          actions: {
            logout: () => Promise.reject(new Error("Network Error")),
          },
        },
      },
    });

    mount(WSLogout, {
      global: {
        plugins: [store],
        stubs: {
          CardModule: true,
          RouterLink: true,
        }
      },
    });

    // Wait for the microtask queue to process the rejected promise
    await new Promise(resolve => setTimeout(resolve, 0));

    expect(consoleSpy).toHaveBeenCalledWith("Logout failed.");

    consoleSpy.mockRestore();
  });
});
