import { createRouter, createWebHashHistory } from "vue-router";
import HomeView from "../views/HomeView.vue";
import Login from "../views/Login.vue";
import ProfileView from "../views/ProfileView.vue";
import UsersView from "../views/UsersView.vue";
import GroupInfoView from "../views/GroupInfoView.vue";

const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    { path: "/", component: HomeView },
    { path: "/login", component: Login },
    { path: "/profile", component: ProfileView },
    { path: "/users", component: UsersView },
    { path: "/groupInfo", component: GroupInfoView },
  ],
});

export default router;
