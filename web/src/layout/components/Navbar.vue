<template>
  <div class="navbar">
    <hamburger
      id="hamburger-container"
      :is-active="sidebar.opened"
      class="hamburger-container"
      @toggleClick="toggleSideBar"
    />

    <breadcrumb id="breadcrumb-container" class="breadcrumb-container" />

    <div class="right-menu">
      <template v-if="device !== 'mobile'">
        <el-tooltip content="全屏" effect="dark" placement="bottom">
          <screenfull id="screenfull" class="right-menu-item hover-effect" />
        </el-tooltip>
      </template>

      <el-dropdown class="avatar-container right-menu-item hover-effect" trigger="click">
        <div class="avatar-wrapper">
          <img :src="avatar" class="user-avatar">
          <span class="user-name">{{ name }}</span>
          <el-icon><caret-bottom /></el-icon>
        </div>

        <template #dropdown>
          <el-dropdown-menu>
            <router-link to="/profile">
              <el-dropdown-item>个人中心</el-dropdown-item>
            </router-link>
            <el-dropdown-item divided @click="logout">
              <span>退出登录</span>
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import { CaretBottom } from '@element-plus/icons-vue'
import Breadcrumb from '@/components/Breadcrumb/index.vue'
import Hamburger from '@/components/Hamburger/index.vue'
import Screenfull from '@/components/Screenfull/index.vue'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'

const userStore = useUserStore()
const appStore = useAppStore()
const router = useRouter()

const sidebar = computed(() => appStore.sidebar)
const device = computed(() => appStore.device)
const avatar = computed(() => userStore.userInfo?.avatar || require('@/assets/avatar.gif'))
const name = computed(() => userStore.userInfo?.nickname || userStore.userInfo?.username || '用户')

const toggleSideBar = () => {
  appStore.toggleSideBar()
}

const logout = async () => {
  await ElMessageBox.confirm('确认退出系统吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
  
  try {
    // 使用Pinia中的logout方法
    await userStore.logout()
  } catch (error) {
    console.error('退出登录失败:', error)
    // 错误处理
    ElMessageBox.alert('退出登录失败，将刷新页面', '错误', {
      confirmButtonText: '确定',
      callback: () => {
        window.location.reload()
      }
    })
  }
}
</script>

<style lang="scss" scoped>
.navbar {
  height: 50px;
  overflow: hidden;
  position: relative;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0,21,41,.08);

  .hamburger-container {
    line-height: 46px;
    height: 100%;
    float: left;
    cursor: pointer;
    transition: background .3s;
    -webkit-tap-highlight-color:transparent;

    &:hover {
      background: rgba(0, 0, 0, .025)
    }
  }

  .breadcrumb-container {
    float: left;
  }

  .right-menu {
    float: right;
    height: 100%;
    line-height: 50px;
    display: flex;
    align-items: center;

    &:focus {
      outline: none;
    }

    .right-menu-item {
      display: inline-block;
      padding: 0 8px;
      height: 100%;
      font-size: 18px;
      color: #5a5e66;
      vertical-align: text-bottom;

      &.hover-effect {
        cursor: pointer;
        transition: background .3s;

        &:hover {
          background: rgba(0, 0, 0, .025)
        }
      }
    }

    .avatar-container {
      margin-right: 30px;

      .avatar-wrapper {
        position: relative;
        display: flex;
        align-items: center;
        
        .user-avatar {
          cursor: pointer;
          width: 30px;
          height: 30px;
          border-radius: 50%;
          margin-right: 8px;
        }

        .user-name {
          color: #333;
          font-size: 14px;
        }

        .el-icon {
          font-size: 12px;
          margin-left: 4px;
        }
      }
    }
  }
}
</style> 