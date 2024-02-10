<script>
import PostItem from "../components/PostItem.vue";
import router from "../router";
import LoadingSpinner from "../components/LoadingSpinner.vue";

export default {
  name: "HomeView",
  components: {LoadingSpinner, PostItem},
  data: function() {
    return {
      errormsg: null,
      loading: false,
      stream: null,
      justifyContent: null,
      imageUrl: null,
      image: null,
    }
  },
  methods: {
    async refresh() {
      this.loading = true;
      this.errormsg = null;

      try {
        let response = await this.$axios.get("/users/" + sessionStorage.userID + "/mystream", {
          headers: {
            "Authorization": sessionStorage.userID,
          },
        });
        this.stream = response.data;
        document.getElementById("log-out").style.display = "block";
        document.getElementById("search").style.display = "block";
        document.getElementById("profile").style.display = "block";
      } catch (e) {
        this.errormsg = e.toString();
        let statusCode = e.response.status;
        if (statusCode === 401) {
          this.errormsg += ". Please log in again."
          setTimeout(() => {
            router.push("/login");
          }, 3000);
        }
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
    canUpload() {
      try {
        return this.profile.profile_info.user.user_id == sessionStorage.userID;
      } catch (e) {
        return false;
      }
    },
    uploadImage() {
      document.getElementById("img-uploader").click();
    },
    async getInputValue(event) {
      this.loading = true;
      const files = event.target.files;
      let filename = files[0].name;
      const fileReader = new FileReader();
      fileReader.addEventListener('load', () => {
        this.imageUrl = fileReader.result;
      })
      fileReader.readAsDataURL(files[0]);
      this.image = files[0];
      try {
        let contentType = this.setImageType(filename);
        if (contentType === null) {
          throw new Error("Invalid file type");
        }
        await this.$axios.post("/users/" + sessionStorage.userID + "/posts/", this.image, {
          headers: {
            "Authorization": sessionStorage.userID,
            "Content-Type": contentType,
          },
        });
      } catch (e) {
        this.errormsg = e.toString();
      }
      setTimeout(() => {
        this.loading = false;
        location.reload();
      }, 2000);
    },
    setImageType(filename) {
      let type = filename.split('.').pop();
      if (type === "png") {
        return "image/png";
      } else if (type === "jpg" || type === "jpeg") {
        return "image/jpeg";
      }
      return null;
    },
    reload() {
      location.reload();
    }
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
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">Home page</h1>
      <div class="btn-toolbar mb-2 mb-md-0">
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="reload">
            Refresh
          </button>
        </div>
      </div>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
  </div>
  <div class="d-flex flex-wrap p-3">
    <div class="d-grid grid-stream w-100" v-if="stream">
      <PostItem class="p-2 mx-4 mt-3 img-thumbnail img-fluid" v-for="(post, index) in stream.posts" :key="index" :uid="post.uid" :postid="post.postid" :likes="post.likes" :uploadTime="post.upload_datetime" :comments="post.comments" :ofStream="true"></PostItem>
    </div>
  </div>
  <div class="mt-2 mb-4" v-if="loading">
    <LoadingSpinner :loading="true"></LoadingSpinner>
  </div>
  <div class="btn-add" v-if="canUpload">
    <button class="btn btn-dark rounded-circle" @click="uploadImage">
      <svg class="feather align-sup add-icon">
        <use href="/feather-sprite-v4.29.0.svg#plus-square" class="mt-2"/>
      </svg>
    </button>
    <input type="file" id="img-uploader" accept="image/*" @change="getInputValue" style="display: none;">
  </div>
</template>

<style>
.grid-stream {
  grid-template-columns: auto;
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
.btn-add {
  width: 3em;
  height: 3em;
  position: fixed;
  bottom: 0;
  right: 0;
  margin: 0 2em 1.5em 0;
}
.align-sup {
  vertical-align: super;
}
.add-icon {
  transform: translate(0px, 6px);
  width: 20px !important;
  height: 20px !important;
}
.grid-stream {
  place-items: center;
}
.grid-stream .img-thumbnail .preview {
  width: 40em;
  height: 25em;
  object-fit: contain;
}
.grid-stream .img-thumbnail .bt-1 {
  margin-top: 0.5em !important;
}
</style>
