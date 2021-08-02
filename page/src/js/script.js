import { randomInt, randomStr, randomPhone } from "./common";
import OpenAPIResponseValidator from "openapi-response-validator";

// const OpenAPIResponseValidator = require('openapi-response-validator');



export const defaultPreScript = `
// 设置环境变量
// t.env.set('strVariable', t.randomStr(5, 10))
// t.env.set('intVariable', t.randomInt(0, 100))
// t.env.set('phoneVariable', t.randomPhone())

`
export const defaultTestScript = `
// 值检查
// t.expect('variable').toBe('value')
// t.expect(true).not.toBe(false);
// 使用环境变量进行值检测
// t.expect('variable').toBe(t.env.get('variable')) 

t.expect(t.resp.status, 'http code 检查').toBe(200)
t.expect(t.resp.data.ret, 'resp ret 检查').toBe(200)
t.expect(t.resp.data.code, 'resp code 检查').toBe('0000')
// 通过接口的 response 的结构定义验证返回是否合法
t.validate()

`

export const PASS = "PASS"
export const FAIL = "FAIL"
export const ERROR = "ERROR"

const styles = {
  [PASS]: { icon: "el-icon-check", class: "success-report" },
  [FAIL]: { icon: "el-icon-close", class: "fail-report" },
  [ERROR]: { icon: "el-icon-close", class: "error-report" },
}

export function runScript(script, variables, resp, respSchem) {
  let pw = {
    _testReports: [],
    resp: resp,
    expect(value, name) {
      try {
        return expect(value, this._testReports, name)
      } catch (e) {
        pw._testReports.push({ name: name, result: ERROR, message: e, styles: styles[ERROR] })
      }
    },
    env: {
      set: (key, value) => (variables[key] = value),
      get: (key) => { return variables[key] },
    },
    randomInt: randomInt,
    randomStr: randomStr,
    randomPhone: randomPhone,
    validate() {
      if (respSchem != undefined && respSchem != null) {
        let v = new OpenAPIResponseValidator(respSchem)
        var verror = v.validateResponse(200, this.resp.data);
        if (verror == undefined) {
          pw._testReports.push({ name: "validate", result: PASS, styles: styles[PASS] })
          return
        }
        console.log(verror)
        let msg = []
        for (let e of verror.errors) {
          msg.push(` ${e.path} : ${e.message} `)
        }
        pw._testReports.push({ name: "validate", result: FAIL, styles: styles[FAIL], message: msg })
      }
    }
  }

  new Function("t", script)(pw)
  return pw._testReports
}

function expect(expectValue, _testReports, name) {
  return new Expectation(expectValue, null, _testReports, name)
}

class Expectation {
  constructor(expectValue, _not, _testReports, name) {
    this.name = name
    this.expectValue = expectValue
    this.not = _not || new Expectation(this.expectValue, true, _testReports)
    this._testReports = _testReports // this values is used within Test.it, which wraps Expectation and passes _testReports value.
    this._satisfies = function (expectValue, targetValue) {
      // Used for testing if two values match the expectation, which could be === OR !==, depending on if not
      // was used. Expectation#_satisfies prevents the need to have an if(this.not) branch in every test method.
      // Signature is _satisfies([expectValue,] targetValue): if only one argument is given, it is assumed the targetValue, and expectValue is set to this.expectValue
      if (!targetValue) {
        targetValue = expectValue
        expectValue = this.expectValue
      }
      if (this.not === true) {
        // test the inverse. this.not is always truthly, but an Expectation that is inverted will always be strictly `true`
        return expectValue !== targetValue
      } else {
        return expectValue === targetValue
      }
    }
  }
  _fmtNot(message) {
    // given a string with "(not)" in it, replaces with "not" or "", depending if the expectation is expecting the positive or inverse (this._not)
    if (this.not === true) {
      return message.replace("(not)", "not ")
    } else {
      return message.replace("(not)", "")
    }
  }
  _fail(message) {
    return this._testReports.push({ name: this.name, result: FAIL, message, styles: styles[FAIL] })
  }
  _pass() {
    return this._testReports.push({ name: this.name, result: PASS, styles: styles[PASS] })
  }
  // TEST METHODS DEFINED BELOW
  // these are the usual methods that would follow expect(...)
  toBe(value) {
    return this._satisfies(value)
      ? this._pass()
      : this._fail(this._fmtNot(`Expected ${this.expectValue} (not)to be ${value}`))
  }
  toHaveProperty(value) {
    return this._satisfies(Object.prototype.hasOwnProperty.call(this.expectValue, value), true)
      ? this._pass()
      : this._fail(
        this._fmtNot(`Expected object ${this.expectValue} to (not)have property ${value}`)
      )
  }
  toBeLevel2xx() {
    const code = parseInt(this.expectValue, 10)
    if (Number.isNaN(code)) {
      return this._fail(`Expected 200-level status but could not parse value ${this.expectValue}`)
    }
    return this._satisfies(code >= 200 && code < 300, true)
      ? this._pass()
      : this._fail(this._fmtNot(`Expected ${this.expectValue} to (not)be 200-level status`))
  }
  toBeLevel3xx() {
    const code = parseInt(this.expectValue, 10)
    if (Number.isNaN(code)) {
      return this._fail(`Expected 300-level status but could not parse value ${this.expectValue}`)
    }
    return this._satisfies(code >= 300 && code < 400, true)
      ? this._pass()
      : this._fail(this._fmtNot(`Expected ${this.expectValue} to (not)be 300-level status`))
  }
  toBeLevel4xx() {
    const code = parseInt(this.expectValue, 10)
    if (Number.isNaN(code)) {
      return this._fail(`Expected 400-level status but could not parse value ${this.expectValue}`)
    }
    return this._satisfies(code >= 400 && code < 500, true)
      ? this._pass()
      : this._fail(this._fmtNot(`Expected ${this.expectValue} to (not)be 400-level status`))
  }
  toBeLevel5xx() {
    const code = parseInt(this.expectValue, 10)
    if (Number.isNaN(code)) {
      return this._fail(`Expected 500-level status but could not parse value ${this.expectValue}`)
    }
    return this._satisfies(code >= 500 && code < 600, true)
      ? this._pass()
      : this._fail(this._fmtNot(`Expected ${this.expectValue} to (not)be 500-level status`))
  }
  toHaveLength(expectedLength) {
    const actualLength = this.expectValue.length
    return this._satisfies(actualLength, expectedLength)
      ? this._pass()
      : this._fail(
        this._fmtNot(
          `Expected length to be ${expectedLength} but actual length was ${actualLength}`
        )
      )
  }
  toBeType(expectedType) {
    const actualType = typeof this.expectValue
    if (
      ![
        "string",
        "boolean",
        "number",
        "object",
        "undefined",
        "bigint",
        "symbol",
        "function",
      ].includes(expectedType)
    ) {
      return this._fail(
        this._fmtNot(
          `Argument for toBeType should be "string", "boolean", "number", "object", "undefined", "bigint", "symbol" or "function"`
        )
      )
    }
    return this._satisfies(actualType, expectedType)
      ? this._pass()
      : this._fail(
        this._fmtNot(`Expected type to be "${expectedType}" but actual type was "${actualType}"`)
      )
  }
}

export function parseTemplateString(string, variables) {
  if (!variables || !string) {
    return string
  }
  const searchTerm = /{{([^}]*)}}/g // {{myVariable}}
  return decodeURI(encodeURI(string)).replace(searchTerm, (match, p1) => variables[p1] || "")
}
