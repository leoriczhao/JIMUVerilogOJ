<template>
  <footer class="footer">
    <div class="footer-container">
      <div class="footer-content">
        <!-- 友链部分 -->
        <div class="footer-section">
          <h3>友情链接</h3>
          <div class="links">
            <a href="https://www.verilog.com" target="_blank" rel="noopener">Verilog官网</a>
            <a href="https://github.com" target="_blank" rel="noopener">GitHub</a>
            <a href="https://stackoverflow.com" target="_blank" rel="noopener">Stack Overflow</a>
            <a href="https://www.edaplayground.com" target="_blank" rel="noopener">EDA Playground</a>
          </div>
        </div>

        <!-- 关于我们 -->
        <div class="footer-section">
          <h3>关于我们</h3>
          <div class="links">
            <a href="/about">关于项目</a>
            <a href="/contact">联系我们</a>
            <a href="/help">帮助文档</a>
            <a href="/feedback">意见反馈</a>
          </div>
        </div>

        <!-- 技术支持 -->
        <div class="footer-section">
          <h3>技术支持</h3>
          <div class="links">
            <a href="/api-docs">API文档</a>
            <a href="/tutorial">使用教程</a>
            <a href="/faq">常见问题</a>
            <a href="/changelog">更新日志</a>
          </div>
        </div>

        <!-- 联系方式 -->
        <div class="footer-section">
          <h3>联系方式</h3>
          <div class="contact-info">
            <p><el-icon><Message /></el-icon> support@verilog-oj.com</p>
            <p><el-icon><Location /></el-icon> 中国</p>
            <p><el-icon><Clock /></el-icon> 24/7 在线支持</p>
          </div>
        </div>
      </div>

      <!-- 版权声明 -->
      <div class="footer-bottom">
        <div class="copyright">
          <p>&copy; 2024 Verilog OJ. All rights reserved.</p>
          <p>
            <a href="/privacy">隐私政策</a> | 
            <a href="/terms">服务条款</a> | 
            <a href="/license">开源协议</a>
          </p>
        </div>
        <div class="powered-by">
          <p>Powered by Vue.js + Element Plus</p>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup lang="ts">
// 无需额外的逻辑
</script>

<style scoped>
.footer {
  background: #2c3e50 !important;
  color: #ecf0f1 !important;
  padding: 40px 0 20px !important;
  margin-top: 60px;
  border-top: none !important;
  width: 100% !important;
}

/* 确保Footer占满屏幕宽度 */
.footer {
  display: block;
  position: relative;
}

/* 确保所有Footer内容都有正确的背景色 */
.footer,
.footer * {
  background: #2c3e50 !important;
}

.footer-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 40px;
  background: #2c3e50;
}

.footer-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 40px;
  margin-bottom: 30px;
}

.footer-section h3 {
  color: #3498db;
  margin-bottom: 20px;
  font-size: 20px;
  font-weight: 600;
}

.links {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.links a {
  color: #bdc3c7;
  text-decoration: none;
  transition: color 0.3s;
  font-size: 14px;
}

.links a:hover {
  color: #3498db;
}

.contact-info p {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  color: #bdc3c7;
  font-size: 14px;
}

.contact-info .el-icon {
  color: #3498db;
}

.footer-bottom {
  border-top: 1px solid #34495e;
  padding-top: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 20px;
}

.copyright p {
  margin: 0;
  color: #95a5a6;
  font-size: 14px;
}

.copyright a {
  color: #bdc3c7;
  text-decoration: none;
  transition: color 0.3s;
}

.copyright a:hover {
  color: #3498db;
}

.powered-by p {
  margin: 0;
  color: #95a5a6;
  font-size: 14px;
}

@media (max-width: 768px) {
  .footer-content {
    grid-template-columns: 1fr;
    gap: 30px;
  }
  
  .footer-bottom {
    flex-direction: column;
    text-align: center;
  }
}
</style>