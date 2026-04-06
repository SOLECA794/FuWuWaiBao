<template>
  <div class="review-panel">
    <h3>📚 专项复习包</h3>
    <el-button type="primary" @click="handleGenerate" :loading="loading">
      智能生成复习包
    </el-button>

    <div v-if="currentPackage" class="package-info">
      <p>包名: {{ currentPackage.name }}</p>
      <p>生成时间: {{ currentPackage.generated_at }}</p>
      <el-button size="small" @click="handleExport">导出为 PDF</el-button>
    </div>
  </div>
</template>

<script>
import { generateReviewPackage, exportReviewPackage } from '@/services/v1/reviewApi';

export default {
  data() {
    return {
      loading: false,
      currentPackage: null
    };
  },
  methods: {
    async handleGenerate() {
      this.loading = true;
      try {
        // 假设你从 Vuex 或父组件拿到了 studentId 和 courseId
        const res = await generateReviewPackage({
          studentId: this.$store.state.userInfo?.id || 'test_student', 
          courseId: this.$route.params.courseId || 'course_001'
        });
        if (res.code === 200) {
          this.currentPackage = res.data;
          this.$message.success('复习包生成成功！');
        }
      } catch (err) {
        this.$message.error(err.message || '生成失败');
      } finally {
        this.loading = false;
      }
    },
    async handleExport() {
      if (!this.currentPackage) return;
      const res = await exportReviewPackage(this.currentPackage.id);
      if (res.code === 200) {
        window.open(res.data.export_url, '_blank');
      }
    }
  }
};
</script>

<style scoped>
.review-panel { padding: 20px; border: 1px solid #eee; border-radius: 8px; }
.package-info { margin-top: 15px; padding: 10px; background: #f9f9f9; }
</style>