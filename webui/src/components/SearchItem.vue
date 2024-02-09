<script>
export default {
  name: "SearchItem",
  data: function () {
    return {
      errormsg: null,
      loading: false,
      searchResults: [],
    }
  },
  methods: {
    async search() {
      this.loading = true;
      let searchInput = document.getElementById("username-search").value;
      if (searchInput.length > 0) {
        let response = await this.$axios.get("/users/?search=" + searchInput, {
          headers: {
            "Authorization": sessionStorage.userID,
          },
        });
        setTimeout(() => {
          this.loading = false;
          this.searchResults = response.data;
        }, 300);
      } else {
        this.searchResults = null;
        this.loading = false;
      }
    },
    async loadProfile(uid) {
      this.$router.push("/profile/" + uid);
      setTimeout(() => {
        location.reload();
      }, 50);
    }
  },
}
</script>

<template>
  <div class="search-container">
    <input type="text" id="username-search" placeholder="Search username..." @keyup="search">
    <button id="search-btn" @click="search"><svg class="feather search-icon"><use href="/feather-sprite-v4.29.0.svg#search"/></svg></button>
    <div class="mt-3" v-if="loading">
      <LoadingSpinner :loading="loading"/>
    </div>
    <div id="search-results" class="search-results" v-if="searchResults">
      <div class="result-item" v-for="(user, index) in searchResults.users" :key="index" @click="loadProfile(user.user_id)">
        <img src="https://picsum.photos/30/30" alt="User Avatar" class="rounded-circle">
        <span class="username align-middle px-1">{{ user.username }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
body {
  margin: 0;
  padding: 0;
  font-family: Arial, sans-serif;
}

.search-container {
  position: fixed;
  left: 0;
  top: 27em;
  transform: translateY(-50%);
  background-color: #f1f1f1;
  padding: 10px;
  border-top-right-radius: 10px;
  border-bottom-right-radius: 10px;
  box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
}

input[type="text"] {
  padding: 5px;
  border: none;
  border-radius: 5px;
  margin-right: 5px;
}

button {
  padding: 5px 16px;
  border: none;
  border-radius: 5px;
  background-color: #424649;
  color: white;
  cursor: pointer;
}

button:hover {
  background-color: #585c5e;
}

button:focus {
  outline: none;
}
.search-icon {
  transform: translateY(-1px);
}
.search-results {
  display: block;
  position: absolute;
  top: 100%;
  left: 0;
  width: 100%;
  background-color: #ffffff;
  border-bottom-right-radius: 10px;
  border-bottom-left-radius: 10px;
  box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.1);
}
.result-item {
  padding: 5px;
  border-bottom: 1px solid #ccc;
  font-weight: bold;
}
.result-item:last-child {
  border-bottom: none;
}
.result-item:hover {
  background-color: #e0e0e0;
  cursor: pointer;
}
.result-item:last-child:hover {
  background-color: #e0e0e0;
  cursor: pointer;
  border-radius: 0 0 11px 11px;
}

</style>