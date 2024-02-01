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
    }
  },
  methods: {
    isOwnerProfile() {
      this.isOwner = (sessionStorage.userID == this.user_id);
    }
  },
  mounted() {
    this.isOwnerProfile();
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
                <h3>{{ numfollowers }}</h3>
                <h5>Followers</h5>
              </div>
              <div class="d-flex flex-column w-50">
                <h3>{{ numfollowing }}</h3>
                <h5>Following</h5>
              </div>
            </div>
          </div>
        </div>
        <div class="d-flex flex-row w-100 mt-4" v-if="!isOwner">
          <button type="button" class="btn btn-primary w-25">Follow</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>

</style>