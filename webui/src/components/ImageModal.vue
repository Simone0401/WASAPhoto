<script>
import Comment from "./Comment.vue";
import ErrorMsg from "./ErrorMsg.vue";
import LoadingSpinner from "./LoadingSpinner.vue";

export default {
  name: "ImageModal",
  components: {LoadingSpinner, ErrorMsg, Comment},
  props: {
    uid: Number,
    postid: Number,
    comments: Array,
    uploadTime: String,
    likes: Number,
    user_put_like: Boolean,
  },
  data: function () {
    return {
      errormsg: null,
      loading: false,
      altText: "Post image",
      userPutLike: this.user_put_like,
      numLikes: this.likes,
      likeIconFill: null,
      imageSrc: null,
      usernameOwner: null,
      displayMenu: false,
      isOwner: null,
      Comments: this.comments,
      numComments: null,
    }
  },
  methods: {
    async refresh() {
      this.loading = true;
      this.errormsg = null;
      try {
        let response = await this.$axios.get("/posts/" + this.postid, {
          headers: {
            "Authorization": sessionStorage.userID,
          },
        });
        this.numLikes = response.data.post.likes;
        this.Comments = response.data.post.comments;
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
      this.numComments = this.getNumberComments();
      this.$emit('toggle-like', this.userPutLike, this.numLikes, this.Comments);
    },
    displayImage() {
      this.imageSrc = document.getElementById("image-id-" + this.postid).src;
    },
    closeModal() {
      this.$emit("close-modal");
      document.body.style.overflow = "scroll";
    },
    getUsername() {
      this.usernameOwner = document.getElementById("username-profile").innerHTML;
    },
    async removePost() {
      this.loading = true;
      this.errormsg = null;
      try {
        await this.$axios.delete("/users/" + sessionStorage.userID + "/posts/" + this.postid, {
          headers: {
            "Authorization": sessionStorage.userID,
          },
        });
        location.reload();
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
    },
    toggleMenu() {
      this.displayMenu = !this.displayMenu;
    },
    isOwnerPost() {
      this.isOwner = (sessionStorage.userID == this.uid);
    },
    downloadImage() {
      let a = document.createElement("a");
      a.href = this.imageSrc;
      let imageName = "image";
      imageName += (this.imageSrc.charAt(11) === "p") ? ".png" : ".jpg";
      a.download = imageName;
      a.click();
    },
    getNumberComments() {
      try {
        return this.Comments.length;
      } catch (e) {
        return 0;
      }
    },
    async toggleLike() {
      if (!this.userPutLike) {
        await this.putLike();
      } else {
        await this.deleteLike();
      }
      this.$emit('toggle-like', this.userPutLike, this.numLikes, this.Comments);
    },
    async getUserLike() {
      this.loading = true;
      this.errormsg = null;
      try {
        let response = await this.$axios.get("/posts/" + this.postid + "/likes/" + sessionStorage.userID, {
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
        let response = await this.$axios.put("/posts/" + this.postid + "/likes/" + sessionStorage.userID, {}, {
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
      this.loading = false;
    },
    async deleteLike() {
      this.loading = true;
      this.errormsg = null;
      try {
        let response = await this.$axios.delete("/posts/" + this.postid + "/likes/" + sessionStorage.userID, {
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
      this.loading = false;
    },
    fillLikeBtn() {
      this.likeIconFill = this.userPutLike ? "#ff3636" : "#ffffff";
    },
    async addComment() {
      this.loading = true;
      this.errormsg = null;
      let textComment = document.getElementById("commentInput").value;
      try {
        textComment = textComment.trim();
        if (textComment.length === 0 || textComment.length > 255) {
          throw new Error("Comment length must be between 1 and 255 characters");
        } else if (!textComment.match(/^[a-zA-Z0-9.,!?;:'"\s]+$/)) {
          throw new Error("Comment can only contain letters, numbers, punctuation and spaces");
        }

        let response = await this.$axios.post("/posts/" + this.postid + "/comments/", {
          "comment": {
            "uid": Number(sessionStorage.userID),
            "message": textComment,
          }
        }, {
          headers: {
            "Authorization": sessionStorage.userID,
            "Content-Type": "application/json",
          },
        });
        await this.refresh();
        document.getElementById("commentInput").value = "";
        this.$emit('toggle-like', this.userPutLike, this.numLikes, this.Comments);
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
      this.numComments = this.getNumberComments();
      this.scrollCommentsToEnd();
    },
    enterPressed(event) {
      if(event.keyCode === 13) {
        this.addComment();
      }
    },
    inputCommentFocus() {
      document.getElementById("commentInput").focus();
    },
    scrollCommentsToEnd() {
      const commentContainer = this.$refs.commentContainer;
      // Make sure the reference exists and has a scroll method.
      if (commentContainer && commentContainer.scrollTo) {
        // Scroll to the bottom of the container.
        commentContainer.scrollTo({
          top: commentContainer.scrollHeight,
          behavior: 'smooth' // Optional, adds animation
        });
      }
    }
  },
  mounted() {
    this.displayImage();
    this.getUsername();
    this.isOwnerPost();
    this.getUserLike();
    this.numComments = this.getNumberComments();
    document.body.style.overflow = "hidden";
  },
}
</script>

<template>
<div class="modal-structure">
  <div class="modal-background" @click="closeModal"></div>
  <div class="modal-content d-flex flex-row">
    <div class="modal-image h-100 w-75">
      <img :src="imageSrc" :alt="altText" :id="'modal-image-id-' + postid">
    </div>
    <div class="modal-post-data h-100 w-25 flex-column">
      <div class="post-information">
        <div class="d-flex flex-column w-100">
          <div class="d-flex flex-row w-100 justify-content-between">
            <div class="d-flex flex-row w-100 flex-wrap">
              <div class="d-flex flex-column w-100">
                <svg class="feather align-sub post-icon" style="width: 20px; height: 20px;" @click="toggleMenu"><use href="/feather-sprite-v4.29.0.svg#more-horizontal"/></svg>
                <div class="dropdown-tail" v-show="displayMenu"></div>
                <div class="dropdown-menu" v-show="displayMenu">
                  <div class="dropdown-content">
                    <div class="dropdown-item dropdown-item-delete" v-show="isOwner" @click="removePost">
                      <svg class="feather vertical-text-top"><use href="/feather-sprite-v4.29.0.svg#trash-2"/></svg>
                      <span class="text-dropdown delete-dropdown">Elimina</span>
                    </div>
                    <div class="dropdown-item dropdown-item-save" @click="downloadImage">
                      <svg class="feather vertical-text-top"><use href="/feather-sprite-v4.29.0.svg#download"/></svg>
                      <span class="text-dropdown save-dropdown">Salva</span>
                    </div>
                  </div>
                </div>
              </div>
              <img src="https://picsum.photos/200" alt="Profile picture" class="rounded-circle" style="width: 50px; height: 50px;">
              <h2 class="fw-bolder username-owner">{{ usernameOwner }}</h2>
              <div class="d-flex flex-row pt-2 upload-date">
                <svg class="feather delete-dropdown-icon"><use href="/feather-sprite-v4.29.0.svg#calendar"/></svg>
                <h6 class="px-1">{{ uploadTime }}</h6>
              </div>
            </div>
          </div>
        </div>
        <!-- Post information -->
        <div class="d-flex flex-row w-100 pl">
          <div class="fit">
            <svg class="feather align-sub like-icon" :id="'like-icon-' + postid" :style="'fill:' + likeIconFill" @click="toggleLike"><use href="/feather-sprite-v4.29.0.svg#heart"/></svg>
            <span>{{ numLikes }}</span>
          </div>
          <div class="fit">
            <svg class="feather align-sub comment-icon" @click="inputCommentFocus"><use href="/feather-sprite-v4.29.0.svg#message-circle"/></svg>
            <span>{{ numComments }}</span>
          </div>
        </div>
      </div>
      <div class="post-comments h-65 border border-2 border-radius overflow-scroll" ref="commentContainer">
        <Comment v-for="comment in Comments" :comment="comment" @removed-comment="refresh"></Comment>
        <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
        <loading-spinner :loading="loading"></loading-spinner>
        <!-- Background placeholder for no comments post -->
        <div class="h-100 d-flex flex-column align-items-center justify-content-center nocomment-section" v-if="!Comments">
          <div class="nocomment-icon-container">
            <svg class="feather nocomment-icon" style="width: 20px; height: 20px;"><use href="/feather-sprite-v4.29.0.svg#message-circle"/></svg>
          </div>
          <div>
            <p class="fw-bold">No comment yet!</p>
          </div>
        </div>
      </div>
      <!-- input for submitting comment -->
      <div class="d-flex flex-column">
        <input type="text" class="mt-2 form-control mb-2 mr-sm-2" id="commentInput" placeholder="Add your comment..." @keydown="enterPressed">
        <button class="btn btn-dark border-radius" type="submit" @click="addComment">
          Send
          <svg class="feather nocomment-icon" style="position: relative; top: -2px;"><use href="/feather-sprite-v4.29.0.svg#send"/></svg>
        </button>
      </div>
    </div>
  </div>
</div>
</template>

<style scoped>
.modal-structure{
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 1000;
  display: flex;
  justify-content: center;
  align-items: center;
}
.modal-background {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1001;
}
.modal-content {
  top: 1.5em;
  width: 85em;
  height: 50em;
  background-color: white;
  border-radius: 10px;
  z-index: 1002;
}
.modal-image > img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}
.modal-post-data {
  padding-left: 10px;
  padding-right: 10px;
}
.post-icon {
  margin-left: 18.5em;
  margin-top: 5px;
}
.post-icon:hover {
  cursor: pointer;
  filter: drop-shadow(0px 1px 1px rgb(0 0 0 / 0.7));
}
.username-owner {
  padding-top: 6px;
  padding-left: 7px;
}
.upload-date {
  display: block;
  float: right;
  margin-left: 10px;
  margin-top: 5px;
}
.dropdown-menu {
  display: block;
  position: absolute;
  top: 1.4em;
  right: 0.2em;
  z-index: 1000;
  float: left;
  min-width: 8rem;
  padding: .5rem 0;
  margin: .125rem 0 0;
  font-size: 1rem;
  color: #212529;
  text-align: left;
  list-style: none;
  background-color: #fff;
  background-clip: padding-box;
  border: 1px solid rgba(0, 0, 0, .15);
  border-radius: .25rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, .1);
}
.dropdown-tail {
  position: absolute;
  top: 1.2em;
  right: 0.8em;
  width: 0;
  height: 0;
  border-style: solid;
  border-width: 0 8px 8px;
  border-color: transparent transparent #d9d9d9 transparent; /* Colore della coda */
  z-index: 1; /* Piazza la coda sopra il menu */
}
.vertical-text-top {
  vertical-align: text-top;
  margin-top: 1px;
}
.text-dropdown {
  margin-left: 0.3em;
}
.dropdown-item-delete:hover {
  color: red;
  cursor: default;
}
.dropdown-item-save:hover {
  color: #36a200;
  cursor: default;
}
.border-radius {
  border-radius: var(--bs-border-radius);
}
.nocomment-section {
  color: rgba(63,60,62,0.38);
}
.nocomment-icon-container {
  position: relative;
}
.nocomment-icon-container::before {
  content: "";
  position: absolute;
  top: 4px;
  left: 4px;
  width: 100%;
  height: 2px;
  background-color: #bbb9bb;
  transform: rotate(45deg);
  transform-origin: top left;
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
.like-icon:hover {
  fill: #ff9090 !important;
}
.comment-icon:hover {
  fill: rgba(55,170,185,0.61);
}
.pl {
  padding-left: 0.75em;
}
.h-65 {
  height: 65% !important;
}
</style>