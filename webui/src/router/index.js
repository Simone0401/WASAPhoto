import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from "../views/LoginView.vue";
import ProfileView from "../views/ProfileView.vue";
import HomeView from "../views/HomeView.vue";

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/login', component: LoginView},
		{path: '/home', component: HomeView},
		{path: '/profile/:user_id', component: ProfileView},
		{path: '/', redirect: to => '/login'}
	]
})

export default router
