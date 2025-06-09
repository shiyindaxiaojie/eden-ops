<template>
  <el-dropdown trigger="click" @command="handleThemeChange">
    <div class="theme-switch">
      <i class="el-icon-magic-stick"></i>
    </div>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item command="default">
          <i class="el-icon-sunny"></i>
          默认蓝
        </el-dropdown-item>
        <el-dropdown-item command="light">
          <i class="el-icon-light-rain"></i>
          清新蓝
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script>
export default {
  name: 'ThemeSwitch',
  methods: {
    handleThemeChange(theme) {
      document.documentElement.setAttribute('data-theme', theme)
      localStorage.setItem('theme', theme)
      this.$message.success('主题切换成功')
    }
  },
  mounted() {
    const theme = localStorage.getItem('theme') || 'default'
    this.handleThemeChange(theme)
  }
}
</script>

<style lang="scss" scoped>
.theme-switch {
  padding: 0 15px;
  height: 100%;
  font-size: 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  
  &:hover {
    background: var(--primary-light-color);
  }
}
</style> 