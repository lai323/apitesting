<template>
  <Login v-if="!userid" />
  <el-container v-if="userid">
    <el-header>
      <el-menu
        :default-active="activeMainMenuIndex"
        class="main-menu"
        mode="horizontal"
        @select="handlePanelSelect"
      >
        <el-menu-item index="env">环境配置</el-menu-item>
        <el-menu-item index="api">接口管理</el-menu-item>
        <el-menu-item index="myTest" :disabled="username == 'admin'"
          >我的测试</el-menu-item
        >
        <el-menu-item index="publicTest">公共测试</el-menu-item>
        <el-menu-item index="user" :disabled="username != 'admin'"
          >用户管理</el-menu-item
        >
        <!-- <el-menu-item index="wiki">Wiki</el-menu-item> -->
        <el-menu-item style="float: right">
          <el-tooltip content="登出" placement="bottom" effect="light">
            <el-button
              icon="el-icon-bicycle"
              size="small"
              @click="logout"
            ></el-button>
          </el-tooltip>
        </el-menu-item>
      </el-menu>
    </el-header>
    <el-main>
      <Env v-if="panelDisabled.env" />
      <Api v-if="panelDisabled.api" />
      <Testing v-if="panelDisabled.myTest" :testingUserId="userid" />
      <Testing v-if="panelDisabled.publicTest" :testingUserId="adminid" />
      <User v-if="panelDisabled.user" />
      <Wiki v-if="panelDisabled.wiki" />
    </el-main>
  </el-container>
</template>

<script>
import Login from "./components/Login.vue";
import User from "./components/User.vue";
import Env from "./components/Env.vue";
import Api from "./components/Api.vue";
import Testing from "./components/Testing.vue";
import Wiki from "./components/Wiki.vue";

export default {
  name: "App",
  components: {
    Login,
    User,
    Env,
    Api,
    Testing,
    Wiki,
  },
  computed: {
    userid() {
      return this.$store.state.userid;
    },
    username() {
      return this.$store.state.username;
    },
    usertoken() {
      return this.$store.state.usertoken;
    },
  },
  data() {
    return {
      panelDisabled: {
        env: false,
        api: true,
        myTest: false,
        publicTest: false,
        user: false,
        wiki: false,
      },
      activeMainMenuIndex: "api",
      adminid: null,
    };
  },
  mounted() {
    this.$tcli.adminid().then((resp) => {
      this.adminid = resp.data.data;
      console.log("adminid", this.adminid);
    });
  },
  methods: {
    handlePanelSelect(key) {
      let panels = Object.entries(this.panelDisabled);
      for (let p of panels) {
        let pname = p[0];
        if (pname == key) {
          this.panelDisabled[pname] = true;
          this.activeMainMenuIndex;
          continue;
        }
        this.panelDisabled[pname] = false;
      }
    },
    logout() {
      this.$store.commit("logout");
      window.location.reload();
    },
  },
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  /* margin-top: 60px; */
}

.main-menu {
  padding-left: 4%;
}
</style>
