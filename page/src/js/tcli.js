
import { elmsg } from './common';
import axios from 'axios';


const codeTokenExpire = '6001'
const codeSuccess = '0000'

export function tcliPlugin(app, token, baseurl) {
  app.config.globalProperties.$tcli = createCli(token, baseurl, function () {
    app.$store.commit('logout')
  })
}

export function createCli(token, baseurl, onExpire) {
  if (token == null) {
    token = ""
  }
  return {
    token: token,
    async commonProcess(req, processErr = true) {
      try {
        var resp = await req
      } catch (e) {
        if (processErr) {
          elmsg.errmsg(e.toString())
          throw e
        }
      }

      if (resp.data.code == codeTokenExpire) {
        if (onExpire != undefined) {
          onExpire()
        }
      } else {
        if (resp.data.code != codeSuccess && processErr) {
          let respmsg = `${resp.data.message}`
          elmsg.errmsg(respmsg)
          throw respmsg
        }
      }
      return resp
    },

    async login(name, password, processErr = true) {
      let resp = await this.commonProcess(
        axios({
          method: 'Post',
          url: `${baseurl}/testing/admin/login`,
          headers: { 'content-type': 'application/json' },
          data: {
            Name: name,
            Password: password
          },
        }), processErr
      )
      this.token = resp.data.data.Token
      return resp
    },
    async env(name, page, pagesize, processErr = true) {
      let where = {}
      let op = {}
      if (name != "") {
        where['Name'] = name
        op['Name'] = 'LIKE'
      }
      let limit = pagesize
      let offset = (page - 1) * pagesize
      return await this.filter(`${baseurl}/testing/admin/env`, where, op, offset, limit, processErr)
    },
    async allenv(processErr = true) {
      return await this.filter(`${baseurl}/testing/admin/env`, null, null, 0, 0, processErr)
    },
    async user(name, page, pagesize, processErr = true) {
      let where = {}
      let op = {}
      if (name != "") {
        where['Name'] = name
        op['Name'] = 'LIKE'
      }
      let limit = pagesize
      let offset = (page - 1) * pagesize
      return await this.filter(`${baseurl}/testing/admin/user`, where, op, offset, limit, processErr)
    },
    async api(route, project, page, pagesize, processErr = true) {
      let where = {}
      let op = {}
      if (route != "") {
        where['Route'] = route
        op['Route'] = 'LIKE'
      }
      if (project != null && project != "") {
        where['Project'] = project
      }
      let limit = pagesize
      let offset = (page - 1) * pagesize
      return await this.filter(`${baseurl}/testing/admin/api`, where, op, offset, limit, processErr)
    },
    async testing(userid, name, project, page, pagesize, processErr = true) {
      let where = { UserId: userid }
      let op = {}
      if (name != "") {
        where['Name'] = name
        op['Name'] = 'LIKE'
      }
      if (project != null && project != "") {
        where['Project'] = project
      }
      let limit = pagesize
      let offset = (page - 1) * pagesize
      return await this.filter(`${baseurl}/testing/admin/testing`, where, op, offset, limit, processErr)
    },
    async getTestingByApiId(id, processErr = true) {
      return await this.filter(`${baseurl}/testing/admin/testing`, { ApiId: id }, null, 0, 0, processErr)
    },
    async filter(url, where, op, offset, limit, processErr = true) {
      let f = {
        // Sort: 'UpdatedAt',
        Sort: 'CreatedAt',
        Order: "desc"
      }
      if (where != null) {
        f['Where'] = where
      }
      if (op != null) {
        f['Op'] = op
      }
      if (limit > 0) {
        f['Limit'] = limit
      }
      if (offset > 0) {
        f['Offset'] = offset
      }

      return await this.commonProcess(
        axios({
          method: 'Post',
          url: url,
          headers: { 
            'content-type': 'application/json',
            'token': this.token,
          },
          data: {
            Filter: f,
          },
        }), processErr
      )
    },
    async delEnv(id, processErr = true) {
      return await this.delById(`${baseurl}/testing/admin/env/del`, id, processErr)
    },
    async delUser(id, processErr = true) {
      return await this.delById(`${baseurl}/testing/admin/user/del`, id, processErr)
    },
    async delApi(id, processErr = true) {
      return await this.delById(`${baseurl}/testing/admin/api/del`, id, processErr)
    },
    async delTesting(id, processErr = true) {
      return await this.delById(`${baseurl}/testing/admin/testing/del`, id, processErr)
    },
    async delById(url, id, processErr = true) {
      return await this.commonProcess(
        axios({
          method: 'Post',
          url: url,
          headers: { 
            'content-type': 'application/json',
            'token': this.token,
          },
          data: {
            ID: id,
          }
        }), processErr
      )
    },
    async saveEnv(rawenv, processErr = true) {
      var errmsg = ""
      let env = Object.assign({}, rawenv)
      if (env.Name == "" || env.Name == undefined) {
        errmsg = '名称为空'
      }
      if (env.Data == "" || env.Data == undefined) {
        errmsg = '环境变量为空'
      }
      let baseUrl = ""
      for (let item of env.Data) {
        if (item.key == "BaseUrl") {
          baseUrl = item.value
        }
      }
      if (baseUrl == "") {
        errmsg = '缺少 BaseUrl'
      }
      if (errmsg != "") {
        elmsg.errmsg(errmsg)
        throw errmsg
      }
      env.Data = JSON.stringify(env.Data)
      return await this.save(`${baseurl}/testing/admin/env/save`, env, processErr)
    },
    async saveUser(user, processErr = true) {
      var errmsg = ""
      if (user.Name == "" || user.Name == undefined) {
        errmsg = '名称为空'
      }
      if (user.Password == "" || user.Password == undefined) {
        errmsg = '密码为空'
      }
      if (errmsg != "") {
        elmsg.errmsg(errmsg)
        throw errmsg
      }
      return await this.save(`${baseurl}/testing/admin/user/save`, user, processErr)
    },
    async saveApi(api, testing, processErr = true) {
      var errmsg = ""
      if (api.Name == "" || api.Name == undefined) {
        errmsg = '名称为空'
      }
      if (api.Project == "" || api.Project == undefined) {
        errmsg = '项目为空'
      }
      if (api.Define == "" || api.Define == undefined) {
        errmsg = '接口定义为空'
      }
      if (api.Route == "" || api.Route == undefined) {
        errmsg = '路由为空'
      }
      if (errmsg != "") {
        elmsg.errmsg(errmsg)
        throw errmsg
      }
      api.Project = api.Project.trim()
      api.Route = api.Route.trim()
      return await this.commonProcess(
        axios({
          method: 'Post',
          url: `${baseurl}/testing/admin/api/save`,
          headers: { 
            'content-type': 'application/json',
            'token': this.token,
          },
          data: {
            Api: api,
            Testing: testing,
          }
        }), processErr
      )
    },
    async saveTesting(rawt, processErr = true) {
      let t = Object.assign({}, rawt)
      var errmsg = ""
      if (t.EnvId == "" || t.EnvId == undefined) {
        errmsg = '环境为空'
      }
      if (t.Project == "" || t.Project == undefined) {
        errmsg = '项目为空'
      }
      if (t.Name == "" || t.Name == undefined) {
        errmsg = '名称为空'
      }
      if (errmsg != "") {
        elmsg.errmsg(errmsg)
        throw errmsg
      }
      t.Data = JSON.stringify(t.Data)
      return await this.save(`${baseurl}/testing/admin/testing/save`, t, processErr)
    },
    async save(url, obj, processErr = true) {
      return await this.commonProcess(
        axios({
          method: 'Post',
          url: url,
          headers: { 
            'content-type': 'application/json',
            'token': this.token,
          },
          data: {
            Entity: obj,
          }
        }), processErr
      )
    },
    async apiProjects(processErr = true) {
      return await this.commonProcess(
        axios({
          method: 'Post',
          url: `${baseurl}/testing/admin/api/projects`,
          headers: { 
            'content-type': 'application/json',
            'token': this.token,
          },
        }), processErr
      )
    },
    async testingProjects(processErr = true) {
      return await this.commonProcess(
        axios({
          method: 'Post',
          url: `${baseurl}/testing/admin/testing/projects`,
          headers: { 
            'content-type': 'application/json',
            'token': this.token,
          },
        }), processErr
      )
    },
    async syncapi(project, entities, processErr = true) {
      return await this.commonProcess(
        axios({
          method: 'Post',
          url: `${baseurl}/testing/admin/syncapi`,
          headers: { 
            'content-type': 'application/json',
            'token': this.token,
          },
          data: {
            Project: project,
            Entities: entities
          }
        }), processErr
      )
    },
    async adminid(processErr = true) {
      return await this.commonProcess(
        axios({
          method: 'Post',
          url: `${baseurl}/testing/admin/adminid`,
          headers: { 
            'content-type': 'application/json',
            'token': this.token,
          },
        }), processErr
      )
    }
  }
}