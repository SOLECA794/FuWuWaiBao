<template>
  <div class="panel-box">
    <div class="upload-area">
      <el-upload
        class="upload-demo"
        drag
        action="#"
        :auto-upload="false"
        :on-change="handleChange"
        accept=".pdf,.pptx"
        :limit="1"
      >
        <i class="el-icon-upload"></i>
        <div class="el-upload__text">
          拖拽文件到此处，或<em>点击上传</em><br />
          <span style="font-size: 12px; color: #999">支持 PDF / PPTX 格式</span>
        </div>
      </el-upload>
      <el-button type="primary" @click="$emit('parse-knowledge')" :disabled="!uploadedFile" style="margin-top: 10px">
        开始拆解知识点
      </el-button>
    </div>

    <div class="knowledge-tree" v-if="knowledgeList.length > 0">
      <h4 style="margin: 10px 0">知识点结构（点击可定位）</h4>
      <el-tree
        :data="knowledgeList"
        :props="treeProps"
        node-key="id"
        @node-click="$emit('node-click', $event)"
        default-expand-all
      ></el-tree>
    </div>

    <el-alert
      v-if="isParsing"
      title="正在拆解知识点，请稍候..."
      type="info"
      show-icon
      style="margin-top: 20px"
    ></el-alert>
    <el-alert
      v-if="parseResult"
      :title="parseResult"
      type="success"
      show-icon
      closable
      style="margin-top: 20px"
    ></el-alert>
  </div>
</template>

<script setup>
/* eslint-disable no-undef */
const emit = defineEmits(['file-change', 'parse-knowledge', 'node-click'])

const handleChange = (file) => {
  emit('file-change', file)
}

defineProps({
  uploadedFile: {
    type: [Object, File],
    default: null
  },
  isParsing: {
    type: Boolean,
    default: false
  },
  parseResult: {
    type: String,
    default: ''
  },
  knowledgeList: {
    type: Array,
    default: () => []
  },
  treeProps: {
    type: Object,
    default: () => ({})
  }
})
</script>

<style scoped>
.panel-box {
  background: rgba(255, 255, 255, 0.96);
  border-radius: 14px;
  padding: 16px;
  border: 1px solid #e6ecf5;
  box-shadow: 0 8px 20px rgba(15, 23, 42, 0.06);
}
.upload-area {
  padding: 20px;
  border: 1px dashed #dcdfe6;
  border-radius: 8px;
  margin-bottom: 20px;
  text-align: center;
}
.knowledge-tree {
  margin-top: 20px;
  max-height: 400px;
  overflow-y: auto;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  padding: 10px;
}
:deep(.el-tree) {
  --el-tree-node-content-hover-bg-color: #e6f7ff;
}
</style>