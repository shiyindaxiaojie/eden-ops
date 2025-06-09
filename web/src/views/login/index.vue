<template>
  <div class="login-container">
    <el-form
      ref="loginFormRef"
      :model="loginForm"
      :rules="loginRules"
      class="login-form"
      autocomplete="on"
      label-position="left"
    >
      <div class="title-container">
        <h3 class="title">eden*</h3>
      </div>

      <el-form-item prop="username">
        <el-input
          ref="usernameRef"
          v-model="loginForm.username"
          placeholder="用户名"
          name="username"
          type="text"
          tabindex="1"
          autocomplete="on"
          prefix-icon="User"
        />
      </el-form-item>

      <el-form-item prop="password">
        <el-input
          ref="passwordRef"
          v-model="loginForm.password"
          :type="passwordType"
          placeholder="密码"
          name="password"
          tabindex="2"
          autocomplete="on"
          prefix-icon="Lock"
          @keyup.enter="handleLogin"
        >
          <template #suffix>
            <el-icon 
              class="show-pwd" 
              @click="showPwd"
            >
              <View v-if="passwordType === 'password'" />
              <Hide v-else />
            </el-icon>
          </template>
        </el-input>
      </el-form-item>

      <el-button 
        :loading="loading" 
        type="primary" 
        @click.prevent="handleLogin"
      >
        登录
      </el-button>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { validateUsername } from '@/utils/validate'
import { ElMessage } from 'element-plus'
import { User, Lock, View, Hide } from '@element-plus/icons-vue'
import { SHA256 } from 'crypto-js'

const userStore = useUserStore()
const router = useRouter()
const route = useRoute()

const redirect = computed(() => {
  return route.query.redirect || '/'
})

const loginForm = ref({
  username: '',
  password: ''
})

const loginRules = {
  username: [{ required: true, trigger: 'blur', validator: validateUsername }],
  password: [{ required: true, trigger: 'blur', message: '请输入密码' }]
}

const loading = ref(false)
const passwordType = ref('password')
const loginFormRef = ref(null)

function showPwd() {
  passwordType.value = passwordType.value === 'password' ? '' : 'password'
}

async function handleLogin() {
  if (!loginFormRef.value) return
  
  try {
    await loginFormRef.value.validate()
    loading.value = true
    
    // 对密码进行SHA256加密
    const loginData = {
      username: loginForm.value.username,
      password: SHA256(loginForm.value.password).toString()
    }
    
    // 登录并获取用户信息
    await userStore.loginAction(loginData.username, loginData.password)
    
    // 登录成功
    ElMessage.success('登录成功')
    
    // 跳转到目标页面
    const targetPath = redirect.value
    try {
      await router.replace(targetPath)
    } catch (err) {
      console.error('路由跳转失败，尝试跳转到首页', err)
      await router.replace('/')
    }
  } catch (error: any) {
    console.error('登录失败:', error)
    ElMessage.error(error.message || '登录失败，请重试')
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss" scoped>
$bg: #283443;
$light_gray: #fff;
$cursor: #fff;

@supports (-webkit-mask: none) and (not (cater-color: $cursor)) {
  .login-container .el-input input {
    color: $cursor;
  }
}

.login-container {
  min-height: 100vh;
  width: 100%;
  background: linear-gradient(135deg, #1f4160 0%, #2d3a4b 100%);
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;

  .login-form {
    position: relative;
    width: 420px;
    max-width: 100%;
    padding: 40px 35px;
    margin: 0 auto;
    overflow: hidden;
    background: rgba(255, 255, 255, 0.9);
    border-radius: 8px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);

    .title-container {
      text-align: center;
      margin-bottom: 40px;

      .title {
        font-size: 26px;
        color: #2d3a4b;
        margin: 0;
        font-weight: bold;
      }
    }

    .el-form-item {
      border: 1px solid rgba(0, 0, 0, 0.1);
      border-radius: 4px;
      margin-bottom: 24px;
      
      &:hover {
        border-color: #409eff;
      }

      .el-input {
        height: 48px;
        
        input {
          height: 48px;
          padding: 12px 5px 12px 15px;
          background: transparent;
          border: 0;
          border-radius: 4px;
          
          &:-webkit-autofill {
            box-shadow: 0 0 0 1000px #fff inset !important;
          }
        }
      }
    }

    .el-button {
      width: 100%;
      height: 48px;
      margin-top: 10px;
      font-size: 16px;
    }

    .show-pwd {
      position: absolute;
      right: 10px;
      top: 7px;
      font-size: 16px;
      color: #889aa4;
      cursor: pointer;
      user-select: none;
    }
  }
}
</style>