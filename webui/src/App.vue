<script setup>
import { RouterLink, RouterView } from 'vue-router'
import SearchItem from "./components/SearchItem.vue";
</script>
<script>
export default {
  data: function () {
    return {
      errormsg: null,
      loading: false,
      wantSearch: false,
      sessionUid: null,
    }
  },
  methods: {
    logOut() {
      sessionStorage.clear();
      this.hideLogOut();
      this.hideSearch();
      this.hideProfile();
    },
    hideLogOut() {
      let uid = sessionStorage.getItem("userID");
      if (uid === null) {
        document.getElementById("log-out").style.display = "none";
      }
    },
    hideSearch() {
      let uid = sessionStorage.getItem("userID");
      if (uid === null) {
        document.getElementById("search").style.display = "none";
      }
    },
    toggleSearch() {
      this.wantSearch = !this.wantSearch;
      if (this.wantSearch) {
        document.getElementById("search").classList.add("no-active");
      } else {
        document.getElementById("search").classList.remove("no-active");
      }
    },
    hideProfile() {
      let uid = sessionStorage.getItem("userID");
      if (uid === null) {
        document.getElementById("profile").style.display = "none";
      } else {
        this.sessionUid = uid;
      }
    },
    reload() {
      location.reload();
    }
  },
  mounted() {
    this.hideLogOut();
    this.hideSearch();
    this.hideProfile();
  },
}
</script>

<template>

	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WASA Photo</a>
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
	</header>

	<div class="container-fluid">
		<div class="row">
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>General</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink to="/home" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
								Home
							</RouterLink>
						</li>
						<li class="nav-item" id="profile" @click="reload">
							<RouterLink :to='"/profile/" + sessionUid' class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user"/></svg>
								Profile
							</RouterLink>
						</li>
						<li class="nav-item" id="search" @click="toggleSearch">
							<RouterLink to="" class="nav-link search no-active">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#search"/></svg>
								Search
							</RouterLink>
						</li>
            <li class="nav-item" id="log-out">
              <RouterLink to="/" class="nav-link log-out" @click="logOut">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#log-out"/></svg>
                  Log out
              </RouterLink>
            </li>
					</ul>

				</div>
        <SearchItem v-if="wantSearch"></SearchItem>
			</nav>

			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView @logged-in="hideProfile"/>
			</main>
		</div>
	</div>
</template>

<style>
.log-out:hover {
  color: #f8293a !important;
}
.search:hover {
  color: #0d41ff !important;
}
.no-active {
  color: #333 !important;
}
</style>
