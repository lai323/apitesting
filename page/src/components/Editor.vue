<template>
  <pre ref="editor"></pre>
</template>

<script>
import ace from "ace-builds";
import "ace-builds/webpack-resolver";
import "ace-builds/src-noconflict/ext-language_tools";
import "ace-builds/src-noconflict/mode-graphqlschema";
export default {
  name: "Editor",
  props: {
    modelValue: {
      type: String,
      default: "",
    },
    mode: {
      type: String,
      default: "ace/mode/javascript",
    },
    maxLines: {
      type: Number,
      default: 20,
    },
    minLines: {
      type: Number,
      default: 15,
    },
    fontSize: {
      type: String,
      default: "18px",
    },
    readOnly: {
      type: Boolean,
      default: false,
    },
    tabSize: {
      type: Number,
      default: 4,
    },
  },
  watch: {
    modelValue(value) {
      if (value !== this.cacheValue) {
        this.editor.session.setValue(value);
        this.cacheValue = value;
      }
    },
  },
  data() {
    return {
      editor: null,
      cacheValue: "",
    };
  },
  mounted() {
    const editor = ace.edit(this.$refs.editor, {
      mode: this.mode,
      enableBasicAutocompletion: true,
      enableLiveAutocompletion: true,
      maxLines: this.maxLines,
      minLines: this.minLines,
      fontSize: this.fontSize,
      autoScrollEditorIntoView: true,
      showPrintMargin: false,
      useWorker: false,
      readOnly: this.readOnly,
    });

    editor.setTheme("ace/theme/clouds");
    editor.session.setTabSize(this.tabSize);

    if (this.modelValue) editor.setValue(this.modelValue, 1);

    this.editor = editor;
    this.cacheValue = this.modelValue;

    if (!this.readOnly) {
      this.editor.on("change", () => {
        const content = editor.getValue();
        this.$emit("update:modelValue", content);
        this.cacheValue = content;
      });
    }
  },
};
</script>