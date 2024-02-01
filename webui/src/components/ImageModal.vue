<script>
export default {
  name: "ImageModal",
  props: {
    uid: Number,
    postid: Number,
    comments: Array,
    uploadTime: String,
    userPutLike: Boolean,
  },
  data: function () {
    return {
      errormsg: null,
      loading: false,
      altText: "Post image",
      userPutLike: this.userPutLike,
      numLikes: this.likes,
      likeIconFill: null,
      imageSrc: null,
      usernameOwner: null,
      displayMenu: false,
      isOwner: null,
    }
  },
  methods: {
    displayImage() {
      this.imageSrc = document.getElementById("image-id-" + this.postid).src;
    },
    closeModal() {
      this.$emit("close-modal");
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
    }
  },
  mounted() {
    this.displayImage();
    this.getUsername();
    this.isOwnerPost();
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
      </div>
      <div class="post-comments">
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
</style>