
import { v4 as uuidv4 } from "uuid";
import SwaggerParser from "@apidevtools/swagger-parser";
import { randomInt, randomStr } from "./common";
import { defaultPreScript, defaultTestScript } from "./script";

const jsyaml = require("js-yaml");

export async function validate(schema) {
  console.log(schema)
  try {
    let api = await SwaggerParser.validate(schema);
    console.log("API name: %s, Version: %s", api.info.title, api.info.version);
  } catch (error) {
    console.log(error.message)
    throw error
  }
}

export function genJsonPostOpenapi(rule) {
  return genPostOpenapi(rule, "application/json")
}

export function genFormPostOpenapi(rule) {
  return genPostOpenapi(rule, "application/x-www-form-urlencoded")
}

function genObjectSchemaFromRule(rule) {
  if (Array.isArray(rule)) {
    if (rule.length == 0) {
      return
    }
  }

  let schema = {
    "type": "object",
    "required": [],
    "properties": {}
  }

  for (const name in rule) {
    let r = rule[name]
    let desc = r.desc
    if (desc == undefined) {
      desc = ""
    }
    let p = {
      "type": "string",
      "description": desc
    }
    if (r.numeric) {
      p.type = "number"
    }
    if (r.integer) {
      p.type = "integer"
    }
    if (r.boolean) {
      p.type = "boolean"
    }
    if (r.required) {
      schema["required"].push(name)
    }
    if (r.in) {
      p.enum = r.in.split(",")
    }
    if (r.default) {
      p.default = r.default
      let i = 0
      if (r.integer) {
        i = parseInt(r.default)
        if (!isNaN(i)) {
          p.default = i
        }
      }
      if (r.numeric) {
        i = parseFloat(r.default)
        if (!isNaN(i)) {
          p.default = i
        }
      }
    }
    if (r.min) {
      let i = 0
      if (r.integer) {
        i = parseInt(r.min)
      }
      if (r.numeric) {
        i = parseFloat(r.min)
      }
      if (isNaN(i)) {
        i = 0
      }
      p.minimum = i
    }
    if (r.max) {
      let i = 0
      if (r.integer) {
        i = parseInt(r.max)
      }
      if (r.numeric) {
        i = parseFloat(r.max)
      }
      if (isNaN(i)) {
        i = 0
      }
      p.maximum = i
    }
    schema["properties"][name] = p
  }
  if (schema['required'].length == 0) {
    delete schema['required']
  }
  if (Object.entries(schema["properties"]).length == 0) {
    delete schema['properties']
  }
  return schema
}


// php 的返回和 go 的返回是一致的
// {
//   "ret": 200,
//   "data": {
//       "result": 4,
//       "userId": "213"
//   },
//   "code": "0000",
//   "message": "请求成功"
// }
function genPostOpenapi(rule, contenttype) {
  let operationId = [rule.method, rule.uri].join("").replaceAll("/", "_")
  let paths = {}
  let pathitem = {
    "post": {
      "operationId": operationId,
      "description": rule.desc,
      "requestBody": {
        "content": {
          [contenttype]: {
            // "schema": {
            //   "type": "object",
            //   "required": [],
            //   "properties": {}
            // }
          }
        }
      },
      "responses": {
        "200": {
          "description": "success",
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "required": [
                  "ret",
                  "data",
                  "code",
                  "message",
                ],
                "properties": {
                  "ret": {
                    "type": "integer"
                  },
                  "code": {
                    "type": "string"
                  },
                  "message": {
                    "type": "string"
                  },
                }
              }
            }
          }
        }
      },
    }
  }
  let reqSchem = genObjectSchemaFromRule(rule.req_rule)
  if (reqSchem == null) {
    pathitem["post"]["requestBody"]["content"][contenttype]["schema"] = {
      "type": "object",
      "required": [],
      "properties": {}
    }
  } else {
    pathitem["post"]["requestBody"]["content"][contenttype]["schema"] = reqSchem
  }

  let dataSchema = genObjectSchemaFromRule(rule.resp_rule)
  if (dataSchema != null) {
    pathitem["post"]["responses"]["200"]["content"]["application/json"]["schema"]["properties"]["data"] = dataSchema
  }

  paths[rule.uri] = pathitem
  let openapi = {
    "openapi": "3.0.0",
    "info": {
      "title": `auto import of ${rule.uri}`,
      "version": "1.0.0",
    },
    "paths": paths,
  }
  return jsyaml.dump(openapi, { indent: 2 });
}

function genParameter(spec) {
  if (spec.schema.default != undefined) {
    return spec.schema.default
  }
  if (spec.schema.enum != undefined && spec.schema.enum.length != 0) {
    return spec.schema.enum[randomInt(0, spec.schema.enum.length - 1)]
  }
  if (spec.schema.type == "integer" || spec.schema.type == "number") {
    let maximum = 0
    let minimum = 10000
    if (spec.schema.maximum != undefined) {
      minimum = spec.schema.maximum
    }
    if (spec.schema.maximum != undefined) {
      minimum = spec.schema.maximum
    }
    return randomInt(maximum, minimum)
  }
  if (spec.schema.type == "boolean") {
    // return [true, false][randomInt(0, 1)]
    // php 中只识别 0 1
    return [1, 0][randomInt(0, 1)]
  }
  return randomStr(randomInt(5, 10))
}

function genBodyParam(spec) {
  if (spec.default != undefined) {
    return spec.default
  }
  if (spec.enum != undefined && spec.enum.length != 0) {
    return spec.enum[randomInt(0, spec.enum.length - 1)]
  }
  if (spec.type == "integer" || spec.type == "number") {
    let maximum = 0
    let minimum = 10000
    if (spec.maximum != undefined) {
      minimum = spec.maximum
    }
    if (spec.maximum != undefined) {
      minimum = spec.maximum
    }
    return randomInt(maximum, minimum)
  }
  if (spec.type == "boolean") {
    // return [true, false][randomInt(0, 1)]
    // php 中只识别 0 1
    return [1, 0][randomInt(0, 1)]
  }
  return randomStr(5, 10)
}

export function genTestFromApi(api) {
  let rule = jsyaml.load(api.Define)
  let pathObj = Object.entries(rule["paths"])[0]
  let route = pathObj[0]
  let methodObj = Object.entries(pathObj[1])[0]
  let method = methodObj[0]
  let apidetail = methodObj[1]["requestBody"]
  let headers = apidetail["headers"]
  let parameters = apidetail["parameters"]
  let content = apidetail["content"]

  console.log(headers)
  console.log(parameters)
  console.log(content)
  // TODO 默认生成 value
  let apiHeaders = []
  let apiParams = []
  let apiBodyParams = []
  if (headers !== undefined) {
    for (let k in headers) {
      apiHeaders.push({
        key: k,
        active: true,
        value: "",
      })
    }
  }
  if (parameters !== undefined) {
    for (let pname in parameters) {
      let p = parameters[pname]
      apiParams.push({
        key: p.name,
        active: true,
        value: genParameter(p),
      })
    }
  }

  // 这里没有处理 json 嵌套
  if (content !== undefined) {
    let contentObj = Object.entries(content)[0]
    let contentType = contentObj[0]
    console.log(contentObj)
    let contentProperties = contentObj[1]["schema"]["properties"]
    apiHeaders.unshift({
      key: "Content-Type",
      active: true,
      value: contentType,
    })
    for (let pname in contentProperties) {
      let p = contentProperties[pname]
      apiBodyParams.push({
        key: pname,
        active: true,
        value: genBodyParam(p),
      })
    }
  }

  let ret = {
    id: uuidv4(),
    name: api.Name,
    route: `{{BaseUrl}}${route}`,
    method: method.toUpperCase(),
    headers: apiHeaders,
    params: apiParams,
    bodyParams: apiBodyParams,
    // rawParams: "",
    // rawInput: true,
    preScript: defaultPreScript,
    testScript: defaultTestScript,
    resp: "",
  }
  let respSchem = methodObj[1]["responses"]
  if (respSchem != undefined) {
    ret['respSchem'] = { "responses": respSchem }
  }
  return ret
}

export function genEmptyTest(name = "空测试") {
  return {
    id: uuidv4(),
    name: name,
    route: "",
    method: "",
    headers: [],
    params: [],
    bodyParams: [],
    // rawParams: "",
    // rawInput: true,
    preScript: defaultPreScript,
    testScript: defaultTestScript,
    resp: "",
  }
}