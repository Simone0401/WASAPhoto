<script>
export default {
  props: {
    user_id: Number,
    username: String,
    numpost: Number,
    numfollowers: Number,
    numfollowing: Number,
  },
  data: function () {
    return {
      errormsg: null,
      loading: false,
      profile: null,
      isOwner: false,
      followed: false,
      followers: this.numfollowers,
      following: this.numfollowing,
    }
  },
  methods: {
    isOwnerProfile() {
      this.isOwner = (sessionStorage.userID == this.user_id);
    },
    async follow() {
      this.loading = true;
      this.errormsg = null;
      try {
        await this.$axios.put("/users/" + sessionStorage.userID + "/following/" + this.user_id, {}, {
          headers: {
            "Authorization": sessionStorage.userID,
          },
        });
        this.followed = true;
        this.followers++;
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
    },
    async unfollow() {
      this.loading = true;
      this.errormsg = null;
      try {
        await this.$axios.delete("/users/" + sessionStorage.userID + "/following/" + this.user_id,  {
          headers: {
            "Authorization": sessionStorage.userID,
          },
        });
        this.followed = false;
        this.followers--;
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
    },
    async hasFollowed() {
      try {
        await this.$axios.get("/users/" + sessionStorage.userID + "/following/" + this.user_id, {
          headers: {
            "Authorization": sessionStorage.userID,
          },
        });
        this.followed = true;
      } catch (e) {
        this.followed = false;
        this.errormsg = e.toString();
      }
    },
  },
  mounted() {
    this.isOwnerProfile();
    this.hasFollowed();
  }
}
</script>

<template>
  <div class="d-flex flex-column w-100">
    <div class="d-flex flex-row w-100">
      <div class="d-flex flex-column w-25">
        <img src="https://picsum.photos/200" alt="Profile picture" class="rounded-circle" style="width: 200px; height: 200px;">
      </div>
      <div class="d-flex flex-column w-75">
        <div class="d-flex flex-row w-100">
          <div class="d-flex flex-column w-100">
            <h1 class="fw-bolder" id="username-profile">{{ username }}</h1>
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
</style>