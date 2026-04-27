<template>
  <n-layout style="height: 100vh;">
    <Navbar />
    <n-layout has-sider style="height: calc(100vh - 48px);">
      <Sidebar />
      <n-layout-content style="padding: 20px; overflow: auto;">
        <router-view />
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<script setup lang="ts">
import { NLayout, NLayoutContent } from 'naive-ui'
import Navbar from './Navbar.vue'
import Sidebar from './Sidebar.vue'
</script>
