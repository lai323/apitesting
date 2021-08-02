import { ElMessage } from 'element-plus'


export const elmsg = {
  errmsg(msg) {
    ElMessage({
      showClose: true,
      message: msg,
      type: 'error',
      offset: 10,
    });
  },
  warnmsg(msg) {
    ElMessage({
      showClose: true,
      message: msg,
      type: 'warning',
      offset: 10,
    });
  },
  successmsg(msg) {
    ElMessage({
      showClose: true,
      message: msg,
      type: 'success',
      offset: 10,
    });
  },

  findGetParameter(parameterName) {
    var result = '', tmp = [];
    location.search
      .substr(1)
      .split("&")
      .forEach(function (item) {
        tmp = item.split("=");
        if (tmp[0] === parameterName) {
          result = decodeURIComponent(tmp[1]);
        }
      });
    return result;
  }
}


export const debounce = (func, delay) => {
  let inDebounce
  return function () {
    const context = this
    const args = arguments
    clearTimeout(inDebounce)
    inDebounce = setTimeout(() => func.apply(context, args), delay)
  }
}

export function randomInt(min, max) {
  min = Math.ceil(min);
  max = Math.floor(max);
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

const chars = '1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ'
export function randomStr(min, max) {
  let length = randomInt(min, max)
  let ret = ''
  for (let i = 0; i < length; i++) {
    ret += chars[randomInt(0, (chars.length - 1))]
  }
  return ret
}

const phoneHeaderNums = new Array("139", "138", "137", "136", "135", "134", "159", "158", "157", "150", "151", "152", "188", "187", "182", "183", "184", "178", "130", "131", "132", "156", "155", "186", "185", "176", "133", "153", "189", "180", "181", "177");
export function randomPhone() {
  const headerNum = phoneHeaderNums[randomInt(0, (phoneHeaderNums.length - 1))]
  const bodyNum = Math.random().toString().replace('0.', '').slice(0, 8)
  return headerNum + bodyNum
}

export function sleep(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}
