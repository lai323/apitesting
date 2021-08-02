<template>
  <el-container>
    <el-card class="aside-testing-card" shadow="hover">
      <el-aside class="aside-testing">
        <div style="display: flex">
          <el-input
            v-model="filterByName"
            @input="filter()"
            placeholder="按名称过滤"
            class="aside-testing-filter-input"
          ></el-input>
          <el-button
            v-if="userid == testingUserId"
            type="primary"
            icon="el-icon-plus"
            class="aside-testing-create-button"
            @click="add"
          ></el-button>
        </div>
        <el-select
          filterable
          clearable
          placeholder="按项目过滤"
          v-model="filterByProject"
          @change="filter()"
          class="aside-api-filter-input"
        >
          <el-option v-for="p in projects" :key="p" :label="p" :value="p">
          </el-option>
        </el-select>

        <div class="aside-testing-list">
          <div
            v-for="t in testings"
            :key="t.ID"
            :class="{
              'aside-testing-list-item': true,
              'choiced-testing': choicedTestings[t.ID],
            }"
            @click="choiceTesting(t.ID)"
          >
            <span>{{ t.Name }}</span>
          </div>
        </div>
        <el-pagination
          small
          layout="prev, pager, next"
          :page-size="pageSize"
          :total="totalTesting"
          :current-page="currentPage"
          @current-change="pageChange"
        >
        </el-pagination>
      </el-aside>
    </el-card>

    <el-main v-if="currentTesting">
      <el-card shadow="hover">
        <div style="display: flex">
          <el-select
            v-model="currentTesting.Project"
            allow-create
            filterable
            placeholder="项目"
            style="width: 40%"
          >
            <el-option v-for="p in projects" :key="p" :label="p" :value="p">
            </el-option>
          </el-select>
          <el-select
            v-model="selectEnvId"
            filterable
            placeholder="环境"
            style="width: 50%"
          >
            <el-option
              v-for="(env, id) in envs"
              :key="id"
              :label="env.Name"
              :value="id"
            >
            </el-option>
          </el-select>
          <el-input v-model="currentTesting.Name" placeholder="名称"></el-input>
          <el-button
            type="primary"
            @click="openAddTestingItem"
            style="margin-left: 8%"
            >添加</el-button
          >
          <el-button
            type="primary"
            @click="test"
            style="margin-left: 8px"
            :loading="testrunning"
            >测试</el-button
          >
          <el-button
            v-if="userid == testingUserId"
            type="primary"
            @click="save"
            style="margin-left: 8px"
            :loading="saveing"
            >保存</el-button
          >
          <el-button
            v-if="currentTesting.ID && userid == testingUserId"
            type="danger"
            icon="el-icon-delete"
            @click="remove(currentTesting.ID)"
            style="margin-left: 8px"
          ></el-button>
        </div>

        <el-tabs
          v-model="activeTestingItem"
          type="border-card"
          style="margin-top: 10px"
        >
          <template v-for="t in currentTesting.Data" :key="t.id">
            <el-tab-pane :name="t.id">
              <template #label>
                <span
                  >{{ t.name }}
                  <i
                    v-if="testingIsAllpass[t.id] === true"
                    class="el-icon-check"
                    style="color: #4ade80"
                  ></i>
                  <i
                    class="el-icon-close"
                    v-if="testingIsAllpass[t.id] === false"
                    style="color: red"
                  ></i>
                </span>
              </template>
              <TestingItem
                :ref="'t-' + t.id"
                :env="currentEnv"
                :testing="t"
                :showSend="false"
                @remove="removeTestingItem"
              />
            </el-tab-pane>
          </template>
        </el-tabs>
      </el-card>
    </el-main>
  </el-container>

  <el-dialog title="添加测试请求" v-model="addTesting.visible">
    <el-form label-width="100px">
      <el-form-item label="从接口添加">
        <el-select
          style="width: 80%"
          placeholder="输入路由来搜索接口"
          v-model="addTesting.route"
          clearable
          filterable
          remote
          :remote-method="addTestingSearch"
          :loading="addTesting.searchLoading"
        >
          <el-option
            v-for="i in addTesting.apis"
            :key="i.ID"
            :label="i.Project + ' : ' + i.Name + ' : ' + i.Route"
            :value="i.ID"
          >
          </el-option>
        </el-select>
      </el-form-item>
    </el-form>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="addTesting.visible = false">取 消</el-button>
        <el-button type="primary" @click="addTestingItem">确 定</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script>
import { elmsg, debounce, sleep } from "../js/common";
import { genTestFromApi, genEmptyTest } from "../js/openapi";
import TestingItem from "./TestingItem.vue";

export default {
  name: "Testing",
  components: {
    TestingItem,
  },
  props: ["testingUserId"],
  computed: {
    selectEnvId: {
      get() {
        if (this.envs[this.currentTesting.EnvId] !== undefined) {
          return String(this.currentTesting.EnvId);
        }
        return "";
      },
      set(envid) {
        this.currentTesting.EnvId = envid;
      },
    },
    currentEnv() {
      if (this.selectEnvId != "") {
        return this.envs[this.selectEnvId];
      }
      return null;
    },
    userid() {
      return this.$store.state.userid;
    },
  },
  data() {
    return {
      envs: {},
      pageSize: 9,
      totalTesting: 0,
      currentPage: 1,
      filterByName: "",
      filterByProject: null,
      currentTesting: null,
      testings: [],
      choicedTestings: {},
      projects: [],
      activeTestingItem: null,
      testingIsAllpass: {},
      saveing: false,
      testrunning: false,
      addTesting: {
        route: null,
        project: "",
        visible: false,
        searchLoading: false,
        apis: [],
      },
    };
  },
  watch: {
    testingUserId() {
      this.filter();
    },
    currentTesting: {
      handler: debounce(function () {
        this.save(null, true);
      }, 650),
      deep: true,
    },
  },
  mounted() {
    this.initProjects();
    this.initEnv();
    this.filter();
  },
  methods: {
    filter: debounce(function (page = 1) {
      if (this.testingUserId == null) {
        return;
      }
      this.currentPage = page;
      this.queryTesting(
        this.filterByName,
        this.filterByProject,
        this.currentPage,
        this.pageSize
      );
    }, 500),
    initProjects() {
      this.$tcli.testingProjects().then((resp) => {
        this.projects = resp.data.data;
        console.log("projects", this.projects);
      });
    },
    initEnv() {
      this.$tcli.allenv().then((resp) => {
        let envs = {};
        for (let env of resp.data.data.rows) {
          let data = JSON.parse(env.Data);
          let baseurl = "";
          for (let v of data) {
            if (v.key == "BaseUrl") {
              baseurl = v.value;
            }
          }
          envs[env.ID] = {
            ID: env.ID,
            Name: env.Name,
            Data: data,
            BaseUrl: baseurl,
          };
        }
        this.envs = envs;
        console.log("envs", this.envs);
      });
    },
    pageChange(p) {
      this.filter(p);
    },
    choiceTesting(id) {
      for (let t of this.testings) {
        this.choicedTestings[t.ID] = false;
        if (t.ID == id) {
          this.currentTesting = t;
          this.choicedTestings[t.ID] = true;
          if (this.currentTesting.Data.length > 0) {
            this.activeTestingItem = this.currentTesting.Data[0].id;
            break;
          }
          this.activeTestingItem = null;
          break;
        }
      }
    },
    add() {
      if (this.currentTesting != null && this.currentTesting.ID == "") {
        return;
      }
      this.currentTesting = {
        ID: null,
        EnvId: null,
        ApiId: null,
        UserId: null,
        Name: null,
        Project: null,
        Data: [],
      };
      this.choicedTestings = {};
    },
    resetCurrentTesting() {
      this.currentTesting = null;
    },
    queryTesting(name, project, page, pagesize) {
      this.$tcli
        .testing(this.testingUserId, name, project, page, pagesize)
        .then((resp) => {
          console.log(resp);
          let testings = [];
          for (let t of resp.data.data.rows) {
            testings.push({
              ID: t.ID,
              EnvId: t.EnvId,
              ApiId: t.ApiId,
              UserId: t.UserId,
              Name: t.Name,
              Project: t.Project,
              Data: JSON.parse(t.Data),
            });
          }
          this.totalTesting = resp.data.data.total;
          this.testings = testings;
        });
    },
    save(_, auto) {
      if (this.userid != this.testingUserId) {
        if (auto) {
          return;
        }
        elmsg.errmsg("你不能修改公共测试");
        return;
      }
      if (
        (this.currentTesting == null || this.currentTesting.ID == null) &&
        auto
      ) {
        return;
      }
      this.saveing = true;
      let t = {
        ID: this.currentTesting.ID,
        EnvId: this.selectEnvId,
        ApiId: null,
        UserId: this.userid,
        Name: this.currentTesting.Name,
        Project: this.currentTesting.Project,
        Data: this.currentTesting.Data,
      };
      this.$tcli
        .saveTesting(t)
        .then((resp) => {
          console.log(resp);
          if (!auto) {
            elmsg.successmsg("保存成功");
            this.filter(this.currentPage);
            this.cleantest();
          }
          this.saveing = false;
          console.log(this.currentTesting);
          if (this.currentTesting.ID == null) {
            this.resetCurrentTesting();
          }
        })
        .catch((error) => {
          console.log(error);
          elmsg.errmsg("保存失败");
          this.saveing = false;
        });
    },
    async test() {
      if (!this.checktest()) {
        return;
      }
      this.cleantest();

      this.testrunning = true;
      let allpass = true;
      let vars = {};

      console.log(this.currentEnv);
      for (let v of this.currentEnv.Data) {
        vars[v.key] = v.value;
      }
      try {
        for (let t of this.currentTesting.Data) {
          let pass = await this.$refs["t-" + t.id].send(vars);
          this.testingIsAllpass[t.id] = pass;
          if (!pass) {
            allpass = false;
            console.log(this.testingIsAllpass);
            break;
          }
          await sleep(500)
        }
      } catch (error) {
        console.log(error);
        elmsg.errmsg("测试运行异常");
        this.testrunning = false;
        return;
      }

      if (allpass) {
        elmsg.successmsg("🙏🙏🙏 所有测试都成功了");
      } else {
        elmsg.warnmsg("👨‍🍳👨‍🍳👨‍🍳 YOU SUCK ! 测试未通过 !");
      }
      this.testrunning = false;
    },
    cleantest() {
      for (let t of this.currentTesting.Data) {
        this.$refs["t-" + t.id].cleanRet();
      }
      this.testingIsAllpass = {};
    },
    checktest() {
      if (this.selectEnvId == "") {
        elmsg.errmsg("未选择环境");
        return false;
      }
      for (let t of this.currentTesting.Data) {
        if (t.method == undefined || t.method == "") {
          elmsg.errmsg(`${t.name} 未选择方法`);
          return false;
        }
        if (t.route == undefined || t.route == "") {
          elmsg.errmsg(`${t.name} 未填写路由`);
          return false;
        }
      }
      return true;
    },

    addTestingSearch(route) {
      if (route == "") {
        return;
      }
      this.addTesting.searchLoading = true;
      this.$tcli
        .api(route, this.addTesting.project, 0, 0)
        .then((resp) => {
          console.log(resp);
          let apis = [];
          for (let api of resp.data.data.rows) {
            apis.push({
              ID: api.ID,
              Project: api.Project,
              Name: api.Name,
              Route: api.Route,
              Define: api.Define,
            });
          }
          this.addTesting.apis = apis;
          this.addTesting.searchLoading = false;
          console.log(this.addTesting.apis);
        })
        .catch((error) => {
          console.log(error);
          elmsg.errmsg("搜索接口失败");
          this.addTesting.searchLoading = false;
        });
    },
    openAddTestingItem() {
      this.addTesting.project = "";
      this.addTesting.route = null;
      this.addTesting.visible = true;
      this.addTesting.apis = [];
      this.addTesting.searchLoading = false;
    },
    addTestingItem() {
      this.addTesting.visible = false;
      let ti = null;
      if (this.addTesting.route == null) {
        ti = genEmptyTest();
      } else {
        let api = null;
        for (let i of this.addTesting.apis) {
          if (i.ID == this.addTesting.route) {
            api = i;
          }
        }
        ti = genTestFromApi(api);
      }
      this.currentTesting.Data.push(ti);
      this.activeTestingItem = ti.id;
    },
    removeTestingItem(tid) {
      console.log(tid);
      let index = null;
      for (let i in this.currentTesting.Data) {
        if (this.currentTesting.Data[i].id == tid) {
          index = i;
          break;
        }
      }
      this.currentTesting.Data.splice(index, 1);
      if (this.currentTesting.Data.length > 0) {
        this.activeTestingItem = this.currentTesting.Data[0].id;
      }
    },
    remove(id) {
      this.$confirm("确定删除 ?", "", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(() => {
          this.$tcli.delTesting(id).then((resp) => {
            console.log(resp);
            elmsg.successmsg("删除成功");
            this.resetCurrentTesting();
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
.choiced-testing {
  color: #66b1ff !important;
}

.testing-box-card {
  width: 70%;
}

.aside-testing {
  width: 100% !important;
}
.aside-testing-card {
  margin-left: 3%;
  width: 20% !important;
}
.aside-testing-filter-input {
  margin-top: 5px;
  height: 40px;
  width: 100%;
}
.aside-testing-list {
  margin-top: 10px;
  height: 450px;
}
.aside-testing-list-item {
  border: 1px solid #ebeef5;
  padding: 10px 10px;
  font-weight: 700;
  color: #909399;
  background: #fafafa;
  text-align: left;
  line-height: 1.5;
  cursor: pointer;
}
.aside-testing-create-button {
  height: 40px;
  margin-top: 5px;
  margin-left: 5px;
}
</style>