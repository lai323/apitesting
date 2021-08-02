<template>
  <el-container>
    <el-card class="aside-api-card" shadow="hover">
      <el-aside class="aside-api">
        <div style="display: flex">
          <el-input
            v-model="filterByRoute"
            @input="filter()"
            placeholder="按路由过滤"
            class="aside-api-filter-input"
          ></el-input>
          <el-button
            type="primary"
            icon="el-icon-plus"
            class="aside-api-create-button"
            @click="add"
          ></el-button>
        </div>
        <el-select
          clearable
          filterable
          placeholder="按项目过滤"
          v-model="filterByProject"
          @change="filter()"
          class="aside-api-filter-input"
        >
          <el-option v-for="p in projects" :key="p" :label="p" :value="p">
          </el-option>
        </el-select>

        <div class="aside-api-list">
          <div
            v-for="api in apis"
            :key="api.ID"
            :class="{
              'aside-api-list-item': true,
              'choiced-api': choicedApis[api.ID],
            }"
            @click="choiceApi(api.ID)"
          >
            <el-tooltip :content="api.Name" placement="right" effect="light">
              <span v-if="api.Route.length < 23">{{ api.Route }}</span>
              <span v-else>{{ api.Route.slice(0, 23) + "..." }}</span>
            </el-tooltip>
          </div>
        </div>
        <el-pagination
          small
          :pager-count="pageCount"
          layout="prev, pager, next"
          :page-size="pageSize"
          :total="totalApi"
          :current-page="currentPage"
          @current-change="pageChange"
        >
        </el-pagination>
      </el-aside>
    </el-card>

    <el-main v-if="currentApi && currentApi.ID == null">
      <el-card shadow="hover">
        <el-form label-width="80px">
          <el-form-item label="路由">
            <div style="display: flex">
              <el-select v-model="currentEnvId" filterable placeholder="环境">
                <el-option
                  v-for="(env, id) in envs"
                  :key="id"
                  :label="env.Name"
                  :value="id"
                >
                </el-option>
              </el-select>
              <el-input v-model="currentApi.Route" style="margin-left: 10px">
                <template #prepend v-if="currentEnvId">{{
                  envs[currentEnvId].BaseUrl
                }}</template>
              </el-input>
              <!-- <el-button
                type="primary"
                icon="el-icon-attract"
                @click="autoImport"
                style="margin-left: 10px"
                >自动导入</el-button
              > -->
            </div>
          </el-form-item>

          <el-form-item label="项目">
            <el-select
              clearable
              v-model="currentApi.Project"
              filterable
              allow-create
              style="width: 100%"
            >
              <el-option v-for="p in projects" :key="p" :label="p" :value="p">
              </el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="名称">
            <el-input v-model="currentApi.Name"></el-input>
          </el-form-item>
          <el-form-item label="定义">
            <Editor
              v-model="currentApi.Define"
              :mode="'ace/mode/yaml'"
              :maxLines="15"
              :minLines="15"
              :tabSize="2"
            />
          </el-form-item>
          <el-button @click="createapi" type="primary">保存</el-button>
          <el-button @click="resetCurrentApi">取消</el-button>
        </el-form>
      </el-card>
    </el-main>

    <el-main v-if="currentApi && currentApi.ID">
      <el-card shadow="hover">
        <div style="display: flex">
          <el-descriptions :column="3" border style="width: 100%">
            <el-descriptions-item
              ><template #label>接口名称</template>
              {{ currentApi.Name }}
            </el-descriptions-item>
            <el-descriptions-item
              ><template #label>所属项目</template>
              {{ currentApi.Project }}
            </el-descriptions-item>
          </el-descriptions>
          <el-select
            v-model="selectEnvId"
            filterable
            placeholder="环境"
            style="margin-left: 15%; width: 40%"
            medium="medium"
          >
            <el-option
              v-for="(env, id) in envs"
              :key="id"
              :label="env.Name"
              :value="id"
            >
            </el-option>
          </el-select>
          <el-button
            style="margin-left: 20px; height: 40px"
            type="primary"
            @click="save"
            :loading="saveing"
            >保存</el-button
          >
        </div>

        <template v-if="currentTesting.Data">
          <TestingItem
            :env="currentEnv"
            :testing="currentTesting.Data"
            @remove="remove"
          />
        </template>
      </el-card>
    </el-main>
  </el-container>
</template>

<script>
import Editor from "./Editor.vue";
import TestingItem from "./TestingItem.vue";
import axios from "axios";
import { genFormPostOpenapi, genTestFromApi, validate } from "../js/openapi";
import { elmsg, debounce } from "../js/common";
const jsyaml = require("js-yaml");

const httpMethods = [
  "GET",
  "POST",
  "HEAD",
  "PUT",
  "DELETE",
  "CONNECT",
  "OPTIONS",
  "TRACE",
  "PATCH",
  "CUSTOM",
];

export default {
  name: "Api",
  components: {
    Editor,
    TestingItem,
  },
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
  },
  data() {
    return {
      pageSize: 9,
      // 页码有几个
      pageCount: 5,
      totalApi: 0,
      currentPage: 1,
      filterByRoute: "",
      filterByProject: null,
      currentApi: null,
      currentTesting: null,
      apis: [],
      envs: {},
      projects: [],
      currentEnvId: null,
      choicedApis: {},
      saveing: false,
    };
  },
  watch: {
    currentTesting: {
      handler: debounce(function () {
        this.save(null, true);
      }, 650),
      deep: true,
    },
  },
  mounted() {
    this.filter();
    this.initEnv();
    this.initProjects();
    this.resetCurrentTesting();
  },
  methods: {
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
    initProjects() {
      this.$tcli.apiProjects().then((resp) => {
        this.projects = resp.data.data;
        console.log("projects", this.projects);
      });
    },
    initTesting() {
      this.resetCurrentTesting();
      this.$tcli.getTestingByApiId(this.currentApi.ID).then((resp) => {
        if (resp.data.data.rows < 1) {
          elmsg.errmsg("未找到默认测试");
        }
        let t = resp.data.data.rows[0];
        this.currentTesting = {
          ID: t.ID,
          EnvId: t.EnvId,
          ApiId: t.ApiId,
          UserId: t.UserId,
          Name: t.Name,
          Project: t.Project,
          Data: JSON.parse(t.Data),
        };
      });
    },
    filter: debounce(function (page = 1) {
      this.currentPage = page;
      this.queryApi(
        this.filterByRoute,
        this.filterByProject,
        this.currentPage,
        this.pageSize
      );
    }, 500),
    pageChange(p) {
      this.filter(p);
    },
    choiceApi(id) {
      for (let api of this.apis) {
        this.choicedApis[api.ID] = false;
        if (api.ID == id) {
          this.currentApi = api;
          this.choicedApis[api.ID] = true;
        }
      }
      this.initTesting();
    },
    add() {
      if (this.currentApi != null && this.currentApi.ID == null) {
        return;
      }
      this.currentApi = {
        ID: null,
        Project: "",
        Name: "",
        Route: "",
        Define: "",
      };
      this.choicedApis = {};
      this.currentEnvId = null;
    },
    resetCurrentTesting() {
      this.currentTesting = {
        ID: null,
        EnvId: "",
        Data: null,
      };
    },
    resetCurrentApi() {
      this.currentApi = null;
      this.currentEnvId = null;
    },
    queryApi(route, project, page, pagesize) {
      this.$tcli.api(route, project, page, pagesize).then((resp) => {
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
        this.totalApi = resp.data.data.total;
        this.apis = apis;
        console.log(this.apis);
      });
    },
    async createapi() {
      try {
        let pass = await this.validateDefine();
        if (!pass) {
          return;
        }
      } catch (error) {
        console.error(error);
        elmsg.errmsg("接口定义不合法");
        return;
      }
      let testingTitle = this.currentApi.Name.trim();
      let testingData = genTestFromApi({
        Define: this.currentApi.Define,
        Name: testingTitle,
      });
      let testing = {
        EnvId: this.currentEnvId,
        ApiId: null,
        UserId: null,
        Name: testingTitle,
        Project: this.currentApi.Project,
        Data: JSON.stringify(testingData),
      };

      this.$tcli.saveApi(this.currentApi, testing).then((resp) => {
        console.log(resp);
        elmsg.successmsg("保存成功");
        this.resetCurrentApi();
        this.filterByRoute = "";
        this.filterByProject = null;
        this.filter();
        this.initEnv();
      });
    },
    save(_, auto) {
      console.log("currentTesting", this.currentTesting);
      if (this.currentTesting == null || this.currentTesting.ID == null) {
        return;
      }
      if (
        (this.currentTesting.EnvId == "" ||
          this.currentTesting.EnvId == null) &&
        !auto
      ) {
        elmsg.errmsg("环境为空");
      }
      this.saveing = true;
      this.$tcli
        .saveTesting(this.currentTesting)
        .then((resp) => {
          console.log(resp);
          if (!auto) {
            elmsg.successmsg("保存成功");
          }
          this.saveing = false;
        })
        .catch((error) => {
          console.log(error);
          elmsg.errmsg("保存失败");
          this.saveing = false;
        });
    },
    remove() {
      this.$confirm(`确定删除 ${this.currentApi.Name} ?`, "", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
      })
        .then(() => {
          this.$tcli.delApi(this.currentApi.ID).then((resp) => {
            console.log(resp);
            elmsg.successmsg("删除成功");
            this.resetCurrentApi();
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
    autoImport() {
      if (this.currentEnvId == null) {
        elmsg.errmsg("环境为空");
        return;
      }
      if (this.currentApi.Route == "") {
        elmsg.errmsg("路由为空");
        return;
      }
      axios
        .get(
          `${
            this.envs[this.currentEnvId].BaseUrl
          }/doc/apiinfo.php?uri=${this.currentApi.Route.trim()}`
        )
        .then((response) => {
          console.log(response);
          if (response.data.length == 0) {
            elmsg.warnmsg("未找到接口");
            return;
          }
          let api = response.data[0];
          this.currentApi.Name = api.desc;
          this.currentApi.Define = this.convRuleOpenApi(api);
        })
        .catch((error) => {
          console.log(error);
          elmsg.errmsg("自动导入失败");
        });
    },
    convRuleOpenApi(rule) {
      if (rule.method == "POST") {
        return genFormPostOpenapi(rule);
      }
      // TODO 处理 get parameters
      elmsg.errmsg("无法自动导入 GET 请求");
      return "";
    },
    async validateDefine() {
      try {
        var data = jsyaml.load(this.currentApi.Define);
      } catch (error) {
        console.error(error);
        elmsg.errmsg("接口定义不是合法的 yaml");
        return false;
      }
      if (data["openapi"] != "3.0.0") {
        elmsg.errmsg("openapi 版本必须是 3.0.0");
        return false;
      }
      console.log(data["paths"]);
      if (
        data["paths"] === undefined ||
        data["paths"] === null ||
        Object.entries(data["paths"]).length == 0
      ) {
        elmsg.errmsg("paths 未定义");
        return false;
      }
      let paths = Object.entries(Object.entries(data["paths"])[0][1]);
      if (paths.length > 1) {
        elmsg.errmsg("只能定义一个方法");
        return false;
      }
      let method = paths[0][0];
      let detail = paths[0][1];
      console.log(method, detail);
      if (!httpMethods.includes(method.toUpperCase())) {
        elmsg.errmsg("http method 不合法");
        return false;
      }
      await validate(data);
      return true;
    },
  },
};
</script>

<style>
.choiced-api {
  color: #66b1ff !important;
}

.aside-api {
  width: 100% !important;
}
.aside-api-card {
  margin-left: 3%;
  width: 20% !important;
}
.aside-api-filter-input {
  margin-top: 5px;
  height: 40px;
  width: 100%;
}
.aside-api-list {
  margin-top: 10px;
  height: 450px;
}
.aside-api-list-item {
  border: 1px solid #ebeef5;
  padding: 10px 10px;
  font-weight: 700;
  color: #909399;
  background: #fafafa;
  text-align: left;
  line-height: 1.5;
  cursor: pointer;
}
.aside-api-create-button {
  height: 40px;
  margin-top: 5px;
  margin-left: 5px;
}
</style>