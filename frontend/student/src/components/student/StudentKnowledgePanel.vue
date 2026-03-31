<template>
  <div class="panel-box">
    <div class="panel-head">
      <div class="eyebrow">知识点拆解</div>
      <h3>上传讲义并自动生成知识树</h3>
    </div>

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
          <span class="hint">支持 PDF / PPTX 格式</span>
        </div>
      </el-upload>
      <el-button type="primary" @click="$emit('parse-knowledge')" :disabled="!uploadedFile" style="margin-top: 10px">
        开始拆解知识点
      </el-button>
    </div>

    <div class="knowledge-tree" v-if="knowledgeList.length > 0">
      <h4>知识点结构（点击可定位）</h4>
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
  background: linear-gradient(180deg, #ffffff 0%, #f7faf8 100%);
  border-radius: 20px;
  padding: 16px;
  border: 1px solid #d9e7df;
  box-shadow: 0 16px 30px rgba(33, 61, 54, 0.08);
}

.panel-head {
  margin-bottom: 12px;
}

.eyebrow {
  font-size: 12px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #6f867d;
  font-weight: 700;
}

h3 {
  margin-top: 4px;
  color: #23463f;
  font-size: 18px;
}

.upload-area {
  padding: 20px;
  border: 1px dashed #b8cdc2;
  border-radius: 12px;
  margin-bottom: 20px;
  text-align: center;
  background: #f4f8f6;
}

.hint {
  font-size: 12px;
  color: #7a9187;
}

.knowledge-tree {
  margin-top: 20px;
  max-height: 400px;
  overflow-y: auto;
  border: 1px solid #d8e7df;
  border-radius: 10px;
  padding: 10px;
  background: #fff;
}

.knowledge-tree h4 {
  margin: 10px 0;
  color: #2b4d45;
}

:deep(.el-tree) {
  --el-tree-node-content-hover-bg-color: #e8f0ec;
}
</style>