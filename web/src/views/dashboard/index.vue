<template>
  <div class="app-container">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>用户总数</span>
              <el-icon><User /></el-icon>
            </div>
          </template>
          <div class="card-value">{{ statistics.userCount }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>角色总数</span>
              <el-icon><UserFilled /></el-icon>
            </div>
          </template>
          <div class="card-value">{{ statistics.roleCount }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>菜单总数</span>
              <el-icon><Menu /></el-icon>
            </div>
          </template>
          <div class="card-value">{{ statistics.menuCount }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>今日IP查询</span>
              <el-icon><Location /></el-icon>
            </div>
          </template>
          <div class="card-value">{{ statistics.todayQueries }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>系统信息</span>
            </div>
          </template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="系统名称">eden*</el-descriptions-item>
            <el-descriptions-item label="系统版本">1.0.0</el-descriptions-item>
            <el-descriptions-item label="服务器地址">{{ systemInfo.serverUrl }}</el-descriptions-item>
            <el-descriptions-item label="操作系统">{{ systemInfo.os }}</el-descriptions-item>
            <el-descriptions-item label="Go版本">{{ systemInfo.goVersion }}</el-descriptions-item>
            <el-descriptions-item label="启动时间">{{ systemInfo.startTime }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>最近登录记录</span>
            </div>
          </template>
          <el-table :data="loginRecords" style="width: 100%" :max-height="300">
            <el-table-column prop="username" label="用户名" width="120" />
            <el-table-column prop="ip" label="IP地址" width="120" />
            <el-table-column prop="location" label="登录地点" />
            <el-table-column prop="loginTime" label="登录时间" width="180" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import { Location, User, UserFilled, Menu } from '@element-plus/icons-vue'

interface Statistics {
  userCount: number
  roleCount: number
  menuCount: number
  todayQueries: number
}

interface SystemInfo {
  serverUrl: string
  os: string
  goVersion: string
  startTime: string
}

interface LoginRecord {
  username: string
  ip: string
  location: string
  loginTime: string
}

const statistics = ref<Statistics>({
  userCount: 0,
  roleCount: 0,
  menuCount: 0,
  todayQueries: 0
})

const systemInfo = ref<SystemInfo>({
  serverUrl: 'http://localhost:8080',
  os: 'Linux',
  goVersion: 'go1.20',
  startTime: '2024-01-01 00:00:00'
})

const loginRecords = ref<LoginRecord[]>([
  {
    username: 'admin',
    ip: '192.168.1.1',
    location: '广东省深圳市',
    loginTime: '2024-01-01 12:00:00'
  },
  {
    username: 'user1',
    ip: '192.168.1.2',
    location: '北京市',
    loginTime: '2024-01-01 11:30:00'
  }
])

onMounted(async () => {
  // TODO: Load dashboard data from API
  // Mock data for now
  statistics.value = {
    userCount: 100,
    roleCount: 10,
    menuCount: 50,
    todayQueries: 1000
  }
})
</script>

<style scoped>
.app-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-value {
  font-size: 24px;
  font-weight: bold;
  text-align: center;
  color: #409EFF;
}

.el-card {
  margin-bottom: 20px;
}
</style>