<template>
  <el-container>
    <el-card class="aside-env-card" shadow="hover">
      <el-aside class="aside-env">
        <div style="display: flex">
          <el-input
            v-model="filterByName"
            @input="filter"
            placeholder="按名称过滤"
            class="aside-env-filter-input"
          ></el-input>
          <el-button
            type="primary"
            icon="el-icon-plus"
            class="aside-env-create-button"
            @click="add"
          ></el-button>
        </div>
        <div class="aside-env-list">
          <div
            v-for="env in envs"
            :key="env.ID"
            :class="{
              'aside-env-list-item': true,
              'choiced-env': choicedEnvs[env.ID],
            }"
            @click="choiceEnv(env.ID)"
          >
            <!-- <el-tooltip :content="item.route" placement="right" effect="light"> -->
            <span>{{ env.Name }}</span>
            <!-- </el-tooltip> -->
          </div>
        </div>
        <el-pagination
          small
          layout="prev, pager, next"
          :page-size="pageSize"
          :total="totalEnv"
          :current-page="currentPage"
          @current-change="pageChange"
        >
        </el-pagination>
      </el-aside>
    </el-card>

    <el-main v-if="currentEnv">
      <el-card class="env-box-card" shadow="hover">
        <el-form label-width="80px">
          <el-form-item label="环境名称">
            <div style="display: flex">
              <el-input
                v-model="currentEnv.Name"
                :disabled="currentEnv.ID != ''"
              ></el-input>
              <el-button
                v-if="currentEnv.ID != ''"
                type="danger"
                icon="el-icon-delete"
                @click="remove(currentEnv.ID)"
                style="margin-left: 20px"
              ></el-button>
            </div>
          </el-form-item>
          <el-form-item label="环境变量">
            <template v-for="(item, i) in currentEnv.Data" :key="i">
              <div style="display: flex; margin-top: 20px">
                <el-input v-model="item.key" placeholder="key"></el-input>
                <span style="margin: 0px 5px">:</span>
                <el-input v-model="item.value" placeholder="value"></el-input>
                <el-button style="margin-left: 10px" @click="removeVar(i)"
                  >删除</el-button
                >
              </div>
            </template>
          </el-form-item>
          <el-button @click="save" type="primary">保存</el-button>
          <el-button @click="appendVar" type="info">新建变量</el-button>
          <el-button @click="resetCurrentEnv">取消</el-button>
        </el-form>
      </el-card>
    </el-main>
  </el-container>
</template>

<script>
import { elmsg, debounce } from "../js/common";

export default {
  name: "Env",
  data() {
    return {
      pageSize: 9,
      totalEnv: 0,
      currentPage: 1,
      filterByName: "",
      currentEnv: null,
      envs: [],
      choicedEnvs: {},
    };
  },
  mounted() {
    this.filter();
  },
  methods: {
    filter: debounce(function (page = 1) {
      this.currentPage = page;
      this.queryEnv(this.filterByName, this.currentPage, this.pageSize);
    }, 500),
    pageChange(p) {
      this.filter(p);
    },
    choiceEnv(id) {
      for (let env of this.envs) {
        this.choicedEnvs[env.ID] = false;
        if (env.ID == id) {
          this.currentEnv = env;
          this.choicedEnvs[env.ID] = true;
        }
      }
    },
    add() {
      if (this.currentEnv != null && this.currentEnv.ID == "") {
        return;
      }
      this.currentEnv = {
        ID: "",
        Name: "",
        Data: [{ key: "BaseUrl", value: "" }],
      };
      this.choicedEnvs = {};
    },
    resetCurrentEnv() {
      this.currentEnv = null;
    },
    queryEnv(name, page, pagesize) {
      this.$tcli.env(name, page, pagesize).then((resp) => {
        console.log(resp);
        let envs = [];
        for (let env of resp.data.data.rows) {
          envs.push({
            ID: env.ID,
            Name: env.Name,
            Data: JSON.parse(env.Data),
          });
        }
        this.totalEnv = resp.data.data.total;
        this.envs = envs;
      });
    },
    appendVar() {
      this.currentEnv.Data.push({ key: "", value: "" });
    },
    removeVar(i) {
      this.currentEnv.Data.splice(i, 1);
    },
    save() {
      this.$tcli.saveEnv(this.currentEnv).then((resp) => {
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
          this.$tcli.delEnv(id).then((resp) => {
            console.log(resp);
            elmsg.successmsg("删除成功");
            this.resetCurrentEnv();
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
.choiced-env {
  color: #66b1ff !important;
}

.env-box-card {
  width: 70%;
}

.aside-env {
  width: 100% !important;
}
.aside-env-card {
  margin-left: 3%;
  width: 20% !important;
}
.aside-env-filter-input {
  margin-top: 5px;
  height: 40px;
  width: 100%;
}
.aside-env-list {
  margin-top: 10px;
  height: 450px;
}
.aside-env-list-item {
  border: 1px solid #ebeef5;
  padding: 10px 10px;
  font-weight: 700;
  color: #909399;
  background: #fafafa;
  text-align: left;
  line-height: 1.5;
  cursor: pointer;
}
.aside-env-create-button {
  height: 40px;
  margin-top: 5px;
  margin-left: 5px;
}
</style>