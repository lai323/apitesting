<template>
  <el-container>
    <el-card class="aside-card-user" shadow="hover">
      <el-aside class="aside-user">
        <div style="display: flex">
          <el-input
            v-model="filterByName"
            @input="filter"
            placeholder="按名称过滤"
            class="aside-filter-user-input"
          ></el-input>
          <el-button
            type="primary"
            icon="el-icon-plus"
            class="aside-create-user-button"
            @click="add"
          ></el-button>
        </div>
        <div class="aside-list-user">
          <div
            v-for="u in users"
            :key="u.ID"
            :class="{
              'aside-list-user-item': true,
              'choiced-user': choicedUsers[u.ID],
            }"
            @click="choiceUser(u.ID)"
          >
            <!-- <el-tooltip :content="item.route" placement="right" effect="light"> -->
            <span>{{ u.Name }}</span>
            <!-- </el-tooltip> -->
          </div>
        </div>
        <el-pagination
          small
          layout="prev, pager, next"
          :page-size="pageSize"
          :total="totalUser"
          :current-page="currentPage"
          @current-change="pageChange"
        >
        </el-pagination>
      </el-aside>
    </el-card>

    <el-main v-if="currentUser">
      <el-card class="user-box-card" shadow="hover">
        <el-form label-width="80px">
          <el-form-item label="用户名">
            <div style="display: flex">
              <el-input
                v-model="currentUser.Name"
                :disabled="currentUser.ID != ''"
              ></el-input>
              <el-button
                v-if="currentUser.ID != ''"
                type="danger"
                icon="el-icon-delete"
                @click="remove(currentUser.ID)"
                style="margin-left: 20px"
              ></el-button>
            </div>
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="currentUser.Password" show-password></el-input>
          </el-form-item>
          <el-button @click="save" type="primary">保存</el-button>
          <el-button @click="resetCurrentUser">取消</el-button>
        </el-form>
      </el-card>
    </el-main>
  </el-container>
</template>

<script>
import { elmsg, debounce } from "../js/common";

export default {
  name: "User",
  data() {
    return {
      pageSize: 9,
      totalUser: 0,
      currentPage: 1,
      filterByName: "",
      currentUser: null,
      users: [],
      choicedUsers: {},
    };
  },
  mounted() {
    this.filter();
  },
  methods: {
    filter: debounce(function (page = 1) {
      this.currentPage = page;
      this.queryUser(this.filterByName, this.currentPage, this.pageSize);
    }, 500),
    pageChange(p) {
      this.filter(p);
    },
    choiceUser(id) {
      for (let u of this.users) {
        this.choicedUsers[u.ID] = false;
        if (u.ID == id) {
          this.currentUser = u;
          this.choicedUsers[u.ID] = true;
        }
      }
    },
    add() {
      if (this.currentUser != null && this.currentUser.ID == "") {
        return;
      }
      this.currentUser = {
        ID: "",
        Name: "",
        Data: [{ key: "BaseUrl", value: "" }],
      };
      this.choicedEnvs = {};
    },
    resetCurrentUser() {
      this.currentUser = null;
    },
    queryUser(name, page, pagesize) {
      this.$tcli.user(name, page, pagesize).then((resp) => {
        let users = [];
        for (let u of resp.data.data.rows) {
          users.push({
            ID: u.ID,
            Name: u.Name,
            Password: u.Password,
          });
        }
        this.totalUser = resp.data.data.total;
        this.users = users;
      });
    },
    save() {
      this.$tcli.saveUser(this.currentUser).then((resp) => {
        console.log(resp);
        elmsg.successmsg("保存成功");
        this.filter(this.currentPage);
      });
    },
    remove(id) {
      this.$confirm("确定删除 ?", "", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(() => {
          this.$tcli.delUser(id).then((resp) => {
            console.log(resp);
            elmsg.successmsg("删除成功");
            this.resetCurrentUser();
            this.filter(this.currentPage);
          });
        })
        .catch(() => {
          this.$message({
            type: "info",
            message: "已取消删除",
          });
        });
    },
  },
};
</script>

<style>
.choiced-user {
  color: #66b1ff !important;
}

.user-box-card {
  width: 70%;
}

.aside-user {
  width: 100% !important;
}
.aside-card-user {
  margin-left: 3%;
  width: 20% !important;
}
.aside-filter-user-input {
  margin-top: 5px;
  height: 40px;
  width: 100%;
}
.aside-list-user {
  margin-top: 10px;
  height: 450px;
}
.aside-list-user-item {
  border: 1px solid #ebeef5;
  padding: 10px 10px;
  font-weight: 700;
  color: #909399;
  background: #fafafa;
  text-align: left;
  line-height: 1.5;
  cursor: pointer;
}
.aside-create-user-button {
  height: 40px;
  margin-top: 5px;
  margin-left: 5px;
}
</style>