<template>
  <div style="margin-top: 20px">
    <div style="display: flex">
      <el-input
        v-model="t.name"
        style="margin-left: 5px; width: 40%"
        placeholder="名称"
      >
      </el-input>
      <el-select
        v-model="t.method"
        filterable
        placeholder="方法"
        style="margin-left: 5px; width: 40%"
      >
        <el-option v-for="m in httpMethods" :key="m" :label="m" :value="m">
        </el-option>
      </el-select>
      <el-input v-model="t.route" style="margin-left: 5px" placeholder="url">
        <!-- <template #prepend v-if="env">{{ env.BaseUrl }}</template> -->
      </el-input>

      <el-button
        v-if="showSend"
        type="primary"
        icon="el-icon-s-promotion"
        @click="selfsend"
        style="margin-left: 10px"
        >发送</el-button
      >
      <el-button
        type="danger"
        icon="el-icon-delete"
        @click="$emit('remove', t.id)"
        style="margin-left: 10px"
        >删除</el-button
      >
    </div>
  </div>

  <el-tabs v-model="activeReqTab" class="tabs-box" style="margin-top: 25px">
    <el-tab-pane label="请求头" name="headers">
      <template v-for="(h, i) in t.headers" :key="i">
        <div style="display: flex; margin-bottom: 5px">
          <el-input
            style="width: 25%"
            v-model="h.key"
            placeholder="key"
          ></el-input>
          <el-input
            style="width: 30%"
            v-model="h.value"
            placeholder="value"
          ></el-input>
          <el-button style="width: 6%" @click="activeHeader(h.key, i)">
            <i v-if="h.active" class="el-icon-check"></i>
            <i v-else class="el-icon-close" style="color: red"></i>
          </el-button>
          <el-button
            icon="el-icon-delete"
            style="width: 6%; margin: 0px"
            @click="delHeader(h.key, i)"
          ></el-button>
        </div>
      </template>
      <div style="display: flex; margin-bottom: 10px">
        <el-button
          icon="el-icon-plus"
          style="width: 67%"
          @click="addHeader"
        ></el-button>
      </div>
    </el-tab-pane>
    <el-tab-pane label="参数" name="params">
      <template v-for="(p, i) in t.params" :key="i">
        <div style="display: flex; margin-bottom: 5px">
          <el-input
            style="width: 25%"
            v-model="p.key"
            placeholder="key"
          ></el-input>
          <el-input
            style="width: 30%"
            v-model="p.value"
            placeholder="value"
          ></el-input>
          <el-button style="width: 6%" @click="activeParam(p.key, i)">
            <i v-if="p.active" class="el-icon-check"></i>
            <i v-else class="el-icon-close" style="color: red"></i>
          </el-button>
          <el-button
            icon="el-icon-delete"
            style="width: 6%; margin: 0px"
            @click="delParam(p.key, i)"
          ></el-button>
        </div>
      </template>
      <div style="display: flex; margin-bottom: 10px">
        <el-button
          icon="el-icon-plus"
          style="width: 67%"
          @click="addParam"
        ></el-button>
      </div>
    </el-tab-pane>
    <el-tab-pane label="请求体" name="body">
      <template v-for="(bp, i) in t.bodyParams" :key="i">
        <div style="display: flex; margin-bottom: 5px">
          <el-input
            style="width: 25%"
            v-model="bp.key"
            placeholder="key"
          ></el-input>
          <el-input
            style="width: 30%"
            v-model="bp.value"
            placeholder="value"
          ></el-input>
          <el-button style="width: 6%" @click="activeBodyParams(bp.key, i)">
            <i v-if="bp.active" class="el-icon-check"></i>
            <i v-else class="el-icon-close" style="color: red"></i>
          </el-button>
          <el-button
            icon="el-icon-delete"
            style="width: 6%; margin: 0px"
            @click="delBodyParams(bp.key, i)"
          ></el-button>
        </div>
      </template>
      <div style="display: flex; margin-bottom: 10px">
        <el-button
          icon="el-icon-plus"
          style="width: 67%"
          @click="addBodyParams"
        ></el-button>
      </div>
    </el-tab-pane>
    <el-tab-pane label="预请求脚本" name="preScript">
      <Editor v-model="t.preScript" :mode="'ace/mode/javascript'" />
    </el-tab-pane>
    <el-tab-pane label="测试脚本" name="testScript">
      <Editor v-model="t.testScript" :mode="'ace/mode/javascript'" />
    </el-tab-pane>
  </el-tabs>

  <el-tabs v-model="activeRespTab" class="tabs-box" style="margin-top: 25px">
    <el-tab-pane label="响应" name="resp">
      <Editor
        v-if="resp"
        v-model="resp"
        :mode="'ace/mode/javascript'"
        :readOnly="true"
      />
    </el-tab-pane>
    <el-tab-pane label="测试报告" name="report">
      <div style="text-align: left">
        <el-divider></el-divider>
        <template v-for="(r, i) in report" :key="i">
          <el-row :class="r.styles.class">
            <el-col :span="1">
              <i :class="r.styles.icon"></i>
            </el-col>
            <el-col :span="2">
              {{ r.result }}
            </el-col>
            <el-col :span="5">
              {{ r.name }}
            </el-col>
            <el-col :span="16">
              <ul v-if="r.message instanceof Array">
                <li v-for="(m, i) in r.message" :key="i">{{ m }}</li>
              </ul>
              <span v-else> {{ r.message }} </span>
            </el-col>
          </el-row>
          <el-divider></el-divider>
        </template>
      </div>
    </el-tab-pane>
  </el-tabs>
</template>

<script>
import Editor from "./Editor.vue";
import test from "../js/test";
import { elmsg } from "../js/common";
import { PASS, FAIL, ERROR } from "../js/script";

export default {
  components: {
    Editor,
  },
  props: {
    showSend: {
      type: Boolean,
      default: true,
    },
    env: {
      type: Object,
      default: () => {
        return null;
      },
    },
    testing: {
      type: Object,
      required: true,
    },
  },
  emits: ["remove"],
  data() {
    return {
      t: this.testing,
      report: [],
      resp: "",
      activeRespTab: "resp",
      activeReqTab: "body",
      httpMethods: [
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
      ],
    };
  },
  mounted() {},
  methods: {
    async selfsend() {
      this.cleanRet();
      let vars = {};
      for (let v of this.env.Data) {
        vars[v.key] = v.value;
      }
      console.log(vars);
      let allpass = await this.send(vars);
      if (allpass) {
        elmsg.successmsg(`<${this.t.name}> 🙏🙏🙏 测试都成功了`);
      } else {
        elmsg.warnmsg(`<${this.t.name}> 👨‍🍳👨‍🍳👨‍🍳 YOU SUCK ! 测试未通过 !`);
      }
    },
    async send(vars) {
      let allpass = true;
      try {
        let ret = await test.test(this.t, vars);
        this.report = ret.report;
        this.resp = ret.resp;
      } catch (error) {
        console.log(error);
        elmsg.errmsg("测试运行异常");
        return;
      }

      for (let r of this.report) {
        if (r.result == PASS) {
          continue;
        }
        if (r.result == FAIL) {
          allpass = false;
          break;
        }
        if (r.result == ERROR) {
          allpass = false;
          break;
        }
      }
      return allpass;
    },
    cleanRet() {
      console.log("cleanRet", this.t.name);
      this.report = [];
      this.resp = "";
    },
    addHeader() {
      this.t.headers.push({
        key: "",
        value: "",
        active: true,
      });
    },
    delHeader(k, index) {
      let newHeader = [];
      for (let i = 0; i < this.t.headers.length; i++) {
        if (i == index) {
          continue;
        }
        newHeader.push(this.t.headers[i]);
      }
      this.t.headers = newHeader;
    },
    activeHeader(k, index) {
      for (let i = 0; i < this.t.headers.length; i++) {
        if (i == index) {
          this.t.headers[i].active = !this.t.headers[i].active;
        }
      }
    },
    addParam() {
      this.t.params.push({
        key: "",
        value: "",
        active: true,
      });
    },
    delParam(k, index) {
      let newHeader = [];
      for (let i = 0; i < this.t.params.length; i++) {
        if (i == index) {
          continue;
        }
        newHeader.push(this.t.params[i]);
      }
      this.t.params = newHeader;
    },
    activeParam(k, index) {
      for (let i = 0; i < this.t.params.length; i++) {
        if (i == index) {
          this.t.params[i].active = !this.t.params[i].active;
        }
      }
    },
    addBodyParams() {
      this.t.bodyParams.push({
        key: "",
        value: "",
        active: true,
      });
    },
    delBodyParams(k, index) {
      let newHeader = [];
      for (let i = 0; i < this.t.bodyParams.length; i++) {
        if (i == index) {
          continue;
        }
        newHeader.push(this.t.bodyParams[i]);
      }
      this.t.bodyParams = newHeader;
    },
    activeBodyParams(k, index) {
      for (let i = 0; i < this.t.bodyParams.length; i++) {
        if (i == index) {
          this.t.bodyParams[i].active = !this.t.bodyParams[i].active;
        }
      }
    },
  },
};
</script>

<style>
.success-report {
  color: #4ade80;
}

.fail-report {
  color: #f87171;
}

.error-report {
  color: #f87171;
}
</style>