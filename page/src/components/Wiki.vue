<template>
  <article
    v-html="compiledMarkdown"
    style="
      text-align: left;
      margin-left: 16%;
      margin-right: 16%;
      margin-top: 20px;
    "
    class="markdown-body"
  ></article>
</template>

<script>
import * as marked from "marked";
import "github-markdown-css/github-markdown.css";

let readme = `
### 支持的 swagger 规则

- type: string, integer, number, boolean
- description
- maximum
- minimum
- default
- enum

`;

export default {
  name: "Wiki",
  data() {
    return {
      input: readme,
    };
  },
  computed: {
    compiledMarkdown() {
      return marked(this.input, { sanitize: true });
    },
  },
};
</script>

<style>
</style>