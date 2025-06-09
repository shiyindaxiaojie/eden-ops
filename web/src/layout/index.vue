<template>
  <el-container class="layout-container">
    <el-aside width="200px">
      <div class="logo">
        <h1>云原生应用平台</h1>
      </div>
      <el-menu
        :default-active="route.path"
        class="el-menu-vertical"
        :router="true"
        :collapse="isCollapse"
      >
        <el-menu-item index="/dashboard">
          <el-icon><Odometer /></el-icon>
          <template #title>仪表盘</template>
        </el-menu-item>
        <el-sub-menu index="/tools">
          <template #title>
            <el-icon><Tools /></el-icon>
            <span>常用工具</span>
          </template>
          <el-menu-item index="/tools/ip-locator">
            <el-icon><Location /></el-icon>
            <template #title>IP定位器</template>
          </el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="/infrastructure">
          <template #title>
            <el-icon><Monitor /></el-icon>
            <span>基础设施</span>
          </template>
          <el-menu-item index="/infrastructure/kubernetes">
            <el-icon><Connection /></el-icon>
            <template #title>Kubernetes</template>
          </el-menu-item>
          <el-menu-item index="/infrastructure/cloud-account">
            <el-icon><Cloudy /></el-icon>
            <template #title>云账号管理</template>
          </el-menu-item>
          <el-menu-item index="/infrastructure/database">
            <el-icon><DataLine /></el-icon>
            <template #title>数据库配置</template>
          </el-menu-item>
          <el-menu-item index="/infrastructure/server">
            <el-icon><Monitor /></el-icon>
            <template #title>服务器配置</template>
          </el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="/system">
          <template #title>
            <el-icon><Setting /></el-icon>
            <span>系统管理</span>
          </template>
          <el-menu-item index="/system/user">
            <el-icon><User /></el-icon>
            <template #title>用户管理</template>
          </el-menu-item>
          <el-menu-item index="/system/role">
            <el-icon><UserFilled /></el-icon>
            <template #title>角色管理</template>
          </el-menu-item>
          <el-menu-item index="/system/menu">
            <el-icon><Menu /></el-icon>
            <template #title>菜单管理</template>
          </el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header>
        <div class="header-left">
          <el-icon
            class="fold-button"
            @click="toggleSidebar"
          >
            <Fold v-if="!isCollapse" />
            <Expand v-else />
          </el-icon>
        </div>
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              {{ userInfo.username }}
              <el-icon><CaretBottom /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Odometer,
  Location,
  Fold,
  Expand,
  CaretBottom,
  Setting,
  User,
  UserFilled,
  Menu,
  Monitor,
  Cloudy,
  Connection,
  DataLine,
  Tools
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const isCollapse = ref(false)
const userInfo = ref({
  username: '',
})

onMounted(() => {
  // 从userStore获取用户信息
  if (userStore.userInfo && userStore.userInfo.username) {
    userInfo.value = userStore.userInfo
  } else {
    // 兼容从localStorage获取
    const userStr = localStorage.getItem('userInfo')
    if (userStr) {
      userInfo.value = JSON.parse(userStr)
    }
  }
})

const toggleSidebar = () => {
  isCollapse.value = !isCollapse.value
}

const handleCommand = async (command: string) => {
  if (command === 'logout') {
    try {
      await userStore.logout()
      ElMessage.success('已退出登录')
    } catch (error) {
      console.error('退出登录失败:', error)
      ElMessage.error('退出登录失败，请重试')
    }
  }
}
</script>

<style lang="scss" scoped>
.layout-container {
  height: 100vh;
  display: flex;
}

.el-aside {
  background-color: rgb(56, 56, 56) !important;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgb(56, 56, 56);
  color: #fff;
}

.logo h1 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: rgb(255, 255, 255);
}

:deep(.el-menu) {
  border-right: none;
  background-color: rgb(56, 56, 56) !important;
}

:deep(.el-menu-item),
:deep(.el-sub-menu__title) {
  color: rgb(255, 255, 255) !important;
  background-color: rgb(56, 56, 56) !important;
}

:deep(.el-menu-item:hover),
:deep(.el-sub-menu__title:hover) {
  background-color: rgba(56, 56, 56, 0.8) !important;
}

:deep(.el-menu-item.is-active) {
  background-color: rgb(33, 116, 255) !important;
  color: rgb(255, 255, 255) !important;
}

:deep(.el-menu-item .el-icon),
:deep(.el-sub-menu__title .el-icon) {
  color: rgb(255, 255, 255) !important;
}

.el-main {
  padding: 0;
  background-color: #f0f2f5;
}

.el-header {
  color: rgb(255, 255, 255) !important;
  background-color: rgb(56, 56, 56) !important;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;

  .fold-button {
    font-size: 20px;
    cursor: pointer;
    &:hover {
      color: #409eff;
    }
  }

  .user-info {
    display: flex;
    align-items: center;
    cursor: pointer;
    color:rgb(255, 255, 255);

    .el-icon {
      margin-left: 5px;
    }
  }
}
</style>