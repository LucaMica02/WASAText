import { createApp, reactive } from "vue";
import App from "./App.vue";
import router from "./router";
import axios from "./services/axios.js";

import "./assets/dashboard.css";
import "./assets/main.css";
import SearchBar from "./components/SearchBar.vue";
import Conversation from "./components/Conversation.vue";

const app = createApp(App);
app.config.globalProperties.$axios = axios;

app.component("SearchBar", SearchBar);
app.component("Conversation", Conversation);

app.use(router);
app.mount("#app");
