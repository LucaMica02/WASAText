import { reactive } from "vue";

export const EventBus = reactive({
  isLoggedIn: localStorage.getItem("isLoggedIn") === "true",
});
