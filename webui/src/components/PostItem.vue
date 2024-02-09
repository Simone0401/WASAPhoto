<script>
import * as base64 from "byte-base64";
import ImageModal from "./ImageModal.vue";

export default {
  name: "PostItem",
  components: {ImageModal},
  props: {
    uid: Number,
    postid: Number,
    likes: Number,
    comments: Array,
    uploadTime: String,
  },
  data: function () {
    return {
      errormsg: null,
      loading: false,
      altText: "Post image",
      imageSrc: null,
      contentType: null,
      userPutLike: null,
      numLikes: this.likes,
      Comments: this.comments,
      numComments: null,
      errno: null,
      likeIconFill: null,
      modalStatus: false,
    }
  },
  methods: {
    async getPostImage() {
      this.loading = true;
      this.errormsg = null;
      try {
        let response = await this.$axios.get("/images/" + this.postid, {
          responseType: 'arraybuffer',
          headers: {
            "Authorization": sessionStorage.userID,
          },
        });

        // Get MIME type from Content-Type header
        this.contentType = response.headers['content-type'];

        const binaryData = new Uint8Array(response.data);
        let img = base64.bytesToBase64(binaryData);

        this.imageSrc = `data:${this.contentType};base64,${img}`;

      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
    },
    getNumberComments() {
      try {
        return this.Comments.length;
      } catch (e) {
        return 0;
      }
    },
    toggleLike() {
      if (!this.userPutLike) {
        this.putLike();
      } else {
       this.deleteLike();
      }
    },
    async getUserLike() {
      this.loading = true;
      this.errormsg = null;
      try {
        await this.$axios.get("/posts/" + this.postid + "/likes/" + sessionStorage.userID, {
          headers: {
            "Authorization": sessionStorage.userID,
          },
        });
        this.userPutLike = true;
      } catch (e) {
        this.errno = e.response.data.errno;
        if (this.errno !== 1) {
          this.errormsg = e.toString();
        } else {
          this.errormsg = null;
          this.userPutLike = false;
        }
      } finally {
        this.fillLikeBtn();
        this.loading = false;
      }
    },
    async putLike() {
      this.loading = true;
      this.errormsg = null;
      try {
        await this.$axios.put("/posts/" + this.postid + "/likes/" + sessionStorage.userID, {}, {
          headers: {
            "Authorization": sessionStorage.userID,
          },
        });
        this.userPutLike = true;
        this.likeIconFill = "#ff3636";
        this.numLikes++;
      } catch (e) {
        this.errormsg = e.toString();
      }
    },
    async deleteLike() {
      this.loading = true;
      this.errormsg = null;
      try {
        await this.$axios.delete("/posts/" + this.postid + "/likes/" + sessionStorage.userID, {
          headers: {
            "Authorization": sessionStorage.userID,
          },
        });
        this.userPutLike = false;
        this.likeIconFill = "#ffffff";
        this.numLikes--;
      } catch (e) {
        this.errormsg = e.toString();
      }
    },
    fillLikeBtn() {
      this.likeIconFill = this.userPutLike ? "#ff3636" : "#ffffff";
    },
    showModal() {
      this.modalStatus = true;
    },
    hideModal() {
      this.modalStatus = false;
    },
    updateLikeStatus(likeStatus, numLikes, comments) {
      this.userPutLike = likeStatus;
      this.numLikes = numLikes;
      this.fillLikeBtn();
      this.Comments = comments;
      this.numComments = this.getNumberComments();
    },
  },
  mounted() {
    this.getPostImage();
    this.getUserLike();
    this.numComments = this.getNumberComments();
  },
}
</script>

<template>
  <div class="d-flex flex-column">
    <ImageModal v-if="modalStatus" v-show="modalStatus" :uid="uid" :postid="postid" :upload-time="uploadTime" :likes="numLikes" :user_put_like="userPutLike" :comments="Comments" @close-modal="hideModal" @toggle-like="updateLikeStatus"></ImageModal>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    <!-- Post image -->
    <div class="d-flex flex-row w-100">
      <img :src="imageSrc" :alt="altText" class="preview" :id="'image-id-' + postid" @click="showModal">
    </div>
    <!-- Post information -->
    <div class="d-flex flex-row w-100 bt-1">
      <div class="fit">
        <svg class="feather align-sub like-icon" :id="'like-icon-' + postid" :style="'fill:' + likeIconFill" @click="toggleLike"><use href="/feather-sprite-v4.29.0.svg#heart"/></svg>
        <span>{{numLikes}}</span>
      </div>
      <div class="fit">
        <svg class="feather align-sub comment-icon" @click="showModal"><use href="/feather-sprite-v4.29.0.svg#message-circle"/></svg>
        <span>{{ numComments }}</span>
      </div>
      <div class="fit">
        <svg class="feather align-sub"><use href="/feather-sprite-v4.29.0.svg#calendar"/></svg>
        <span>{{uploadTime}}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
img {
  float: left;
  width: 20em;
  height: 20em;
  object-fit: contain;
}
.fit {
  width: fit-content;
  padding-right: 10px;
}
.align-sub {
  vertical-align: sub;
}
.fit > span {
  padding-left: 2px;
}
.bt-1 {
  border-top: 1px solid #c2c2c4 !important;
}
.like-icon:hover {
  fill: #ff9090 !important;
}
.comment-icon:hover {
  fill: rgba(55,170,185,0.61);
}
</style>