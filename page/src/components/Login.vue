<template>
  <el-card class="login-card" :body-style="{ padding: '50px' }">
    <h2>请登录</h2>
    <el-form ref="form" label-width="0px">
      <el-form-item>
        <el-input placeholder="name" v-model="name"></el-input>
      </el-form-item>
      <el-form-item>
        <el-input
          placeholder="password"
          v-model="password"
          show-password
        ></el-input>
      </el-form-item>
      <el-button type="primary" @click="login">登陆</el-button>
    </el-form>
  </el-card>
</template>

<script>
import { elmsg } from "../js/common";

export default {
  name: "Login",
  data() {
    return {
      name: "",
      password: "",
    };
  },
  methods: {
    login() {
      if (this.name == "") {
        elmsg.errmsg("请输入用户名");
      }
      if (this.password == "") {
        elmsg.errmsg("请输入密码");
      }
      this.$tcli.login(this.name, this.password).then((resp) => {
        console.log(resp);
        console.log(
          resp.data.data.Id,
          resp.data.data.Name,
          resp.data.data.Token
        );
        this.$store.commit("login", {
          userid: resp.data.data.Id,
          username: resp.data.data.Name,
          usertoken: resp.data.data.Token,
        });
        elmsg.successmsg("登陆成功");
        window.location.reload();
      });
    },
  },
  mounted() {},
};
</script>

<style>
.login-card {
  margin: auto;
  margin-top: 12%;
  width: 30%;
}
</style>