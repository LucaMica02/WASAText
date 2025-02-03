import { createApp, reactive } from "vue";
import App from "./App.vue";
import router from "./router";
import axios from "./services/axios.js";

import "./assets/dashboard.css";
import "./assets/main.css";
import Conversation from "./components/Conversation.vue";
import User from "./components/User.vue";

const app = createApp(App);
app.config.globalProperties.$axios = axios;

app.component("Conversation", Conversation);
app.component("User", User);

app.use(router);
app.mount("#app");
