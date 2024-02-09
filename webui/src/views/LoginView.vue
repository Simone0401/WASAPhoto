<script>
import LoadingSpinner from "../components/LoadingSpinner.vue";
import ErrorMsg from "../components/ErrorMsg.vue";

export default {
    components: {ErrorMsg, LoadingSpinner},
  data: function() {
    return {
      errormsg: null,
      loading: false,
      username: null,
    }
  },
  methods: {
    async login() {
      if (event.key === "Enter" || event.type === "click") {

        this.loading = true;
        this.errormsg = null;

        // Check if username is empty
        if (this.username == null || this.username === "") {
          this.errormsg = "Username cannot be empty";
          this.errorUsername();
          setTimeout(() => {
            this.errormsg = null;
          }, 5000)
        }
        // Check if username is at least 3 char
        else if (this.username.length < 3) {
          this.errormsg = "Username must be at least 3 characters long";
          this.errorUsername();
          setTimeout(() => {
            this.errormsg = null;
          }, 5000)
        }
        // Check if username is at most 20 char
        else if (this.username.length > 20) {
          this.errormsg = "Username must be at most 20 characters long";
          this.errorUsername();
          setTimeout(() => {
            this.errormsg = null;
          }, 5000)
        }
        // Check if username regex is valid ^[a-zA-Z0-9.,!?;:'"\s]+$ (only letters, numbers, spaces and punctuation)
        else if (!this.username.match(/^[a-zA-Z0-9.,!?;:'"\b]+$/)) {
          this.errormsg = "Username can only contain letters, numbers and punctuation";
          this.errorUsername();
          setTimeout(() => {
            this.errormsg = null;
          }, 5000)
        } else {
          try {
            let response = await this.$axios.post("/session", {
              username: this.username,
            });
            sessionStorage.userID = response.data.user_id;
            this.$router.push("/profile/" + sessionStorage.userID);
            this.$emit("logged-in");
          } catch (e) {
            this.errormsg = "Error " + e.response.status + ": " + e.response.data;
            setTimeout(() => {
              this.errormsg = null;
            }, 5000)
          }
        }
        this.loading = false;
      }
    },
    errorUsername() {
      document.getElementById("username").classList.add("alert", "alert-danger");
      document.getElementById("username").focus();
    },
    resetErrorUsername() {
      document.getElementById("username").classList.remove("alert", "alert-danger");
    },
  },
  mounted() {
    sessionStorage.clear();
  },
}
</script>

<template>
  <ErrorMsg v-if="errormsg" :msg="errormsg" style="margin-top: 10px;"></ErrorMsg>
  <div class="login-post" @keydown="login">
    <div class="d-flex flex-column border rounded position-relative">
      <div class="position-relative flex-column" style="margin-bottom: 10px;">
        <div class="overflow-hidden bg-login position-relative">
          <img class="img-fluid rounded-1" src="/sfondo_login.jpg" alt="Login background image" />
        </div>
        <div class="position-absolute text-login w-100 h-100 top-0">
          <h1 class="text-center fw-bold position-relative top-50 title-login">Login</h1>
        </div>
      </div>
      <div class="d-flex flex-wrap flex-column p-1">
        <label for="username">Username:</label>
        <input type="text" id="username" name="username" v-model="username" @keydown="resetErrorUsername">

        <div>
          <label for="remember">Remember me:</label>
          <input class="position-relative " type="checkbox" id="remember" name="remember" style="margin-left: 4px; top: 2px;">
        </div>

        <button type="submit" @click="login">Login</button>
        <LoadingSpinner v-if="loading" :loading="loading"></LoadingSpinner>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-post {
  padding-top: 3em;
  padding-left: 25%;
  padding-right: 25%;
}

.bg-login {
  height: 200px;
  width: 100%;
  background-position: center;
  background-repeat: no-repeat;
  background-size: cover;
  border-top-left-radius: var(--bs-border-radius);
  border-top-right-radius: var(--bs-border-radius);
}

.text-login {
  left: 0;
  color: white;
  font-size: 2em;
  font-weight: bold;
  background-color: #4f5050AF;
  border-top-left-radius: var(--bs-border-radius);
  border-top-right-radius: var(--bs-border-radius);
}

.title-login {
  top: 50%;
  transform: translateY(-50%);
}

label {
  margin-bottom: 5px;
}

input[type="text"] {
  padding: 5px;
  margin-bottom: 10px;
}

input[type="checkbox"] {
  margin-bottom: 10px;
}

input[type="checkbox"]:hover {
  cursor: pointer;
  box-shadow: 0 0 0 0.1rem rgba(166, 200, 255, 0.5);
}

button[type="submit"] {
  padding: 5px 10px;
  margin-top: 10px;
  margin-left: 40%;
  margin-right: 40%;
  margin-bottom: 20px;
  background-color: #4f5050B0;
  color: white;
  border: none;
  cursor: pointer;
  border-radius: 25px;
  transition-duration: 0.4s;
}

button[type="submit"]:hover {
  background-color: #3e8e41;
}
</style>