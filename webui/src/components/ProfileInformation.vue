<script>
import ErrorMsg from "./ErrorMsg.vue";

export default {
  components: {ErrorMsg},
  props: {
    user_id: Number,
    username: String,
    numpost: Number,
    numfollowers: Number,
    numfollowing: Number,
  },
  data: function () {
    return {
      errorBanOrFollow: null,
      loading: false,
      isOwner: false,
      followed: false,
      followers: this.numfollowers,
      following: this.numfollowing,
      hasBan: false,
      wantChange: false,
      usernameVar: this.username,
    }
  },
  methods: {
    isOwnerProfile() {
      this.isOwner = (sessionStorage.userID == this.user_id);
    },
    async follow() {
      this.loading = true;
      this.errorBanOrFollow = null;
      try {
        await this.$axios.put("/users/" + sessionStorage.userID + "/following/" + this.user_id, {}, {
          headers: {
            "Authorization": sessionStorage.userID,
          },
        });
        this.followed = true;
        this.followers++;
      } catch (e) {
        this.errorBanOrFollow = e.toString();
      }
      this.loading = false;
    },
    async unfollow() {
      this.loading = true;
      this.errorBanOrFollow = null;
      try {
        await this.$axios.delete("/users/" + sessionStorage.userID + "/following/" + this.user_id,  {
          headers: {
            "Authorization": sessionStorage.userID,
          },
        });
        this.followed = false;
        this.followers--;
      } catch (e) {
        this.errorBanOrFollow = e.toString();
      }
      this.loading = false;
    },
    async hasFollowed() {
      this.errorBanOrFollow = null;
      try {
        await this.$axios.get("/users/" + sessionStorage.userID + "/following/" + this.user_id, {
          headers: {
            "Authorization": sessionStorage.userID,
          },
        });
        this.followed = true;
      } catch (e) {
        this.followed = false;
      }
    },
    async hasBanned() {
      this.loading = true;
      try {
        await this.$axios.get("/users/" + sessionStorage.userID + "/muted/" + this.user_id, {
          headers: {
            "Authorization": sessionStorage.userID,
          },
        });
        this.hasBan = true;
      } catch (e) {
        if (e.response.status === 404) {
          this.hasBan = false;
        }
      }
    },
    async ban() {
      this.loading = true;
      this.errorBanOrFollow = null;
      try {
        await this.$axios.put("/users/" + sessionStorage.userID + "/muted/" + this.user_id, {}, {
          headers: {
            "Authorization": sessionStorage.userID,
          },
        });
        this.hasBan = true;
      } catch (e) {
        this.errorBanOrFollow = e.toString();
      }
      this.loading = false;
    },
    async unban() {
      this.loading = true;
      this.errorBanOrFollow = null;
      try {
        await this.$axios.delete("/users/" + sessionStorage.userID + "/muted/" + this.user_id,  {
          headers: {
            "Authorization": sessionStorage.userID,
          },
        });
        this.hasBan = false;
      } catch (e) {
        this.errorBanOrFollow = e.toString();
      }
      this.loading = false;
    },
    toggleChangeUsername() {
      if (this.wantChange) {
        this.changeUsername();
        this.wantChange = false;
      } else {
        this.wantChange = true;
        setTimeout(() => {
          document.getElementById("input-username").value = this.usernameVar;
          document.getElementById("input-username").focus();
        }, 100);
      }
    },
    async changeUsername() {
      this.loading = true;
      this.errorBanOrFollow = null;
      let newUsername = document.getElementById("input-username").value;
      // Check if username is empty
      if (newUsername == null || newUsername === "") {
        this.errorBanOrFollow = "Username cannot be empty";
        setTimeout(() => {
          this.errorBanOrFollow = null;
        }, 5000)
      }
      // Check if username is at least 3 char
      else if (newUsername.length < 3) {
        this.errorBanOrFollow = "Username must be at least 3 characters long";
        setTimeout(() => {
          this.errorBanOrFollow = null;
        }, 5000);
      }
      // Check if username is at most 20 char
      else if (newUsername.length > 20) {
        this.errorBanOrFollow = "Username must be at most 20 characters long";
        setTimeout(() => {
          this.errorBanOrFollow = null;
        }, 5000);
      }
      // Check if username regex is valid ^[a-zA-Z0-9.,!?;:'"\s]+$ (only letters, numbers, spaces and punctuation)
      else if (!newUsername.match(/^[a-zA-Z0-9.,!?;:'"\b]+$/)) {
        this.errorBanOrFollow= "Username can only contain letters, numbers and punctuation";
        setTimeout(() => {
          this.errorBanOrFollow = null;
        }, 5000);
      } else {
        try {
          await this.$axios.put("/users/" + sessionStorage.userID + "/username", {
            username: newUsername,
          }, {
            headers: {
              "Authorization": sessionStorage.userID,
            },
          });
          this.usernameVar = newUsername;
        } catch (e) {
          if (e.response.status === 409) {
            this.errorBanOrFollow = "Username already taken";
          } else {
            this.errorBanOrFollow = e.toString();
          }
          setTimeout(() => {
            this.errorBanOrFollow = null;
          }, 5000);
        }
      }
      this.loading = false;
    }
  },
  mounted() {
    this.isOwnerProfile();
    if (!this.isOwner) {
      this.hasFollowed();
      this.hasBanned();
    }
  }
}
</script>

<template>
  <ErrorMsg v-if="errorBanOrFollow" :msg="errorBanOrFollow"/>
  <div class="d-flex flex-column w-100">
    <div class="d-flex flex-row w-100">
      <div class="d-flex flex-column w-25">
        <img src="https://picsum.photos/200" alt="Profile picture" class="rounded-circle" style="width: 200px; height: 200px;">
      </div>
      <div class="d-flex flex-column w-75">
        <div class="d-flex flex-row w-100">
          <div class="d-flex flex-column w-100">
            <div class="d-flex flex-row w-100">
              <input type="text" name="" id="input-username" class="input-username" v-if="wantChange">
              <h1 class="fw-bolder" id="username-profile" v-if="!wantChange">{{ usernameVar }}</h1>
              <button type="button" class="btn btn-warning mx-3 h-75 m-auto" v-if="isOwner" @click="toggleChangeUsername">
                <svg class="feather align-sup add-icon">
                  <use href="/feather-sprite-v4.29.0.svg#edit-2"/>
                </svg>
                <span class="text-button">Change username</span>
              </button>
            </div>
            <div class="d-flex flex-row w-100 mt-2">
              <div class="d-flex flex-column w-50">
                <h3>{{ numpost }}</h3>
                <h5>Posts</h5>
              </div>
              <div class="d-flex flex-column w-50">
                <h3>{{ followers }}</h3>
                <h5>Followers</h5>
              </div>
              <div class="d-flex flex-column w-50">
                <h3>{{ following }}</h3>
                <h5>Following</h5>
              </div>
            </div>
          </div>
        </div>
        <div class="d-flex flex-row w-100 mt-4" v-if="!isOwner">
          <button type="button" class="btn btn-primary w-25" v-if="!followed" @click="follow">
            <svg class="feather align-sup add-icon">
              <use href="/feather-sprite-v4.29.0.svg#user-plus"/>
            </svg>
            <span class="text-button">Follow</span>
          </button>
          <button type="button" class="btn btn-secondary w-25" v-if="followed" @click="unfollow">
            <svg class="feather align-sup add-icon">
              <use href="/feather-sprite-v4.29.0.svg#user-minus"/>
            </svg>
            <span class="text-button">Followed</span>
          </button>
          <button type="button" class="btn btn-danger mx-2" v-if="!hasBan" @click="ban">
            <svg class="feather align-sup add-icon">
              <use href="/feather-sprite-v4.29.0.svg#slash"/>
            </svg>
            <span class="text-button">Ban</span>
          </button>
          <button type="button" class="btn btn-secondary mx-2" v-if="hasBan" @click="unban">
            <svg class="feather align-sup add-icon">
              <use href="/feather-sprite-v4.29.0.svg#slash"/>
            </svg>
            <span class="text-button">Unban</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.text-button {
  position: relative;
  vertical-align: text-bottom;
  margin-left: 0.3em;
  top: -1px;
}
#username-profile {
  margin-top: 0.5rem;
}
.input-username {
  border: 1px solid #ced4da;
  border-radius: 0.25rem;
  margin-right: 0.5rem;
  font-weight: bolder;
  font-size: 2.5rem;
}
</style>