<script>
import ErrorMsg from "./ErrorMsg.vue";

export default {
  name: "CommentItem",
  components: {ErrorMsg},
  props: {
    comment: Object,
  },
  data: function() {
    return {
      errormsg: null,
      loading: false,
      username: null,
      displayMenu: false,
    }
  },
  methods: {
    refresh() {
      this.getUsername();
    },
    async getUsername() {
      this.loading = true;
      try {
        let response = await this.$axios.get("/users/" + this.comment.uid + "/username", {
          headers: {
            "Authorization": sessionStorage.userID,
          },
        });
        this.username = response.data.username;
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
    },
    isOwner() {
      return this.comment.uid == sessionStorage.userID;
    },
    toggleMenu() {
      this.displayMenu = !this.displayMenu;
    },
    async removeComment() {
      this.loading = true;
      this.errormsg = null;
      try {
        await this.$axios.delete("/posts/" + this.comment.postid + "/comments/" + this.comment.id, {
          headers: {
            "Authorization": sessionStorage.userID,
          },
        });
        this.$emit("removed-comment");
      } catch (e) {
        this.errormsg = e.toString();
      }
      this.loading = false;
    },
  },
  mounted() {
    this.refresh();
  },
  updated() {
    this.refresh();
  },
}
</script>

<template>
  <div class="comment d-flex flex-row pt-2 pl-2 pr-2 border-radius">
    <error-msg v-if="errormsg" :msg="errormsg"></error-msg>
    <img src="https://picsum.photos/40/40" alt="User Avatar" class="avatar rounded-circle h-50">
    <div class="comment-content pl-2">
      <span class="username fw-bold">{{ username }}</span>
      <p class="text">{{ comment.message }}</p>
    </div>
    <div class="comment-more" v-if="isOwner()">
      <svg class="feather comment-icon" style="width: 20px; height: 20px;" @click="toggleMenu"><use href="/feather-sprite-v4.29.0.svg#more-horizontal"/></svg>
      <div class="dropdown-tail" v-show="displayMenu"></div>
      <div class="dropdown-menu" v-show="displayMenu">
        <div class="dropdown-content">
          <div class="dropdown-item dropdown-item-delete" @click="removeComment">
            <svg class="feather vertical-text-top"><use href="/feather-sprite-v4.29.0.svg#trash-2"/></svg>
            <span class="text-dropdown delete-dropdown">Elimina</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.pl-2 {
  padding-left: 0.5em;
}
.pr-2 {
  padding-right: 0.5em;
}
.comment-icon {
  float: right;
  margin-left: auto;
}
.comment-icon:hover {
  cursor: pointer;
  filter: drop-shadow(0px 1px 1px rgb(0 0 0 / 0.7));
}
.dropdown-menu {
  display: block;
  position: relative;
  top: -0.6em;
  right: -0.3em;
  z-index: 1000;
  float: right;
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
  position: relative;
  float: right;
  top: -0.5em;
  left: -0.1em;
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
.comment-more {
  width: 2em;
  margin-left: auto;
}
</style>