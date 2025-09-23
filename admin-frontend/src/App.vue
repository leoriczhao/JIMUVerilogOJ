<template>
  <div class="app">
    <el-container style="height:100vh">
      <el-aside width="200px">
        <el-menu router>
          <el-menu-item index="/">
            仪表盘
          </el-menu-item>
          <el-menu-item index="/users">用户管理</el-menu-item>
          <el-menu-item index="/problems">题目管理</el-menu-item>
          <el-menu-item index="/submissions">评测管理</el-menu-item>
        </el-menu>
      </el-aside>
      <el-container>
        <el-header>Verilog OJ 管理后台</el-header>
        <el-main>
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
</script>

<style>
body, html, #app { height: 100%; margin: 0; }
</style>

