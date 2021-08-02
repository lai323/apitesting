import axios from 'axios';
import { runScript, parseTemplateString } from '../js/script';

export default {
  async test(t, variables) {
    runScript(t.preScript, variables)
    let resp = await this.request(t, variables)
    let report = runScript(t.testScript, variables, resp, t.respSchem)
    return {
      resp: JSON.stringify(resp.data, null, 2),
      report: report
    }
  },
  async request(t, variables) {
    let url = parseTemplateString(t.route, variables)
    let headers = {}
    for (let h of t.headers) {
      if (h.active) {
        headers[h.key.toLowerCase()] = parseTemplateString(h.value, variables)
      }
    }
    let params = {}
    for (let p of t.params) {
      if (p.active) {
        params[p.key] = parseTemplateString(p.value, variables)
      }
    }
    let data = {}
    for (let d of t.bodyParams) {
      if (d.active) {
        data[d.key] = parseTemplateString(d.value, variables)
      }
    }

    return await axios({
      method: t.method,
      url: url,
      headers: headers,
      params: params,
      data: data,
      transformRequest: [function (data, headers) {
        if (headers["content-type"] == undefined) {
          return
        }
        if (headers["content-type"].toLowerCase() == 'application/x-www-form-urlencoded') {
          const params = new URLSearchParams();
          for (let item of Object.entries(data)) {
            params.append(item[0], item[1]);
          }
          return params.toString()
        }
        return data;
      }]
    })
  }

}
