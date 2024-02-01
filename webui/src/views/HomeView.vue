<script>
import ProfileInformation from "../components/ProfileInformation.vue";
import Post from "../components/Post.vue";
import Modal from "../components/ImageModal.vue";
import ImageModal from "../components/ImageModal.vue";
import router from "../router";

export default {
  components: {ImageModal, Modal, Post, ProfileInformation},
	data: function() {
		return {
			errormsg: null,
			loading: false,
			profile: null,
      justifyContent: null,
		}
	},
	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/users/" + sessionStorage.userID + "/profile", {
          headers: {
            "Authorization": sessionStorage.userID,
          },
        });
				this.profile = response.data;
        document.getElementById("log-out").style.display = "block";
        document.getElementById("search").style.display = "block";
			} catch (e) {
				this.errormsg = e.toString();
        this.errormsg += ". Please log in again."
        setTimeout(() => {
          router.push("/login");
        }, 3000);
			}
			this.loading = false;
		},
    moreElement() {
      try {
        if (this.profile.uploaded_post.length > 2) {
          this.justifyContent = "justify-content-between";
        } else {
          this.justifyContent = "justify-content-fixer";
        }
      } catch (e) {
        this.justifyContent = "";
      }
    },
	},
	mounted() {
		this.refresh();
	},
  updated() {
    this.moreElement();
  }
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Home page</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="exportList">
						Export
					</button>
				</div>
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="newItem">
						New
					</button>
				</div>
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
  <div class="d-flex flex-wrap p-3 border-2">
    <ProfileInformation v-if="profile" :user_id="profile.profile_info.user.user_id" :username="profile.profile_info.user.username" :numpost="profile.profile_info.numpost" :numfollowers="profile.profile_info.follower" :numfollowing="profile.profile_info.following"></ProfileInformation>
    <div class="d-flex flex-column w-100 mt-5 custom-border">
      <div class="d-flex flex-row w-100">
        <div class="d-flex flex-row w-100">
          <svg class="feather align-sub post-icon" style="width: 20px; height: 20px;"><use href="/feather-sprite-v4.29.0.svg#grid"/></svg>
        </div>
      </div>
      <div :class="'d-grid grid-post w-100 post-section ' + justifyContent">
          <Post class="img-thumbnail p-2 mx-4 mt-3" v-if="profile" v-for="post in profile.uploaded_post" :uid="post.uid" :postid="post.postid" :likes="post.likes" :uploadTime="post.upload_datetime" :comments="post.comments"></Post>
      </div>
    </div>
  </div>
</template>

<style>
.post-section {
  border-top: 1px solid #dee2e6;
}
.grid-post {
  grid-template-columns: auto auto auto;
}
.post-icon {
  margin-left: 11.2em;
  margin-bottom: 5px;
}
.custom-border::before {
  content: "";
  position: relative;
  top: 26px;
  left: 0;
  width: 31%;
  height: 2px;
  background-color: black;
  z-index: 1;
}
.justify-content-fixer > div:nth-child(2) {
  margin-right: 23.5em !important;
}
.img-thumbnail {
  width: fit-content;
}
.img-thumbnail .preview {
  cursor: pointer;
}
</style>
