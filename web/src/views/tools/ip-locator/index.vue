<template>
  <div class="app-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>IP定位器</span>
        </div>
      </template>

      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="IP地址" prop="ip">
          <el-input
            v-model="form.ip"
            placeholder="请输入IP地址"
            clearable
            style="width: 300px"
          >
            <template #append>
              <el-button @click="handleLocate" :loading="loading">查询</el-button>
            </template>
          </el-input>
        </el-form-item>
      </el-form>

      <el-divider v-if="locationInfo" />

      <div v-if="locationInfo" class="location-info">
        <el-descriptions title="位置信息" :column="2" border>
          <el-descriptions-item label="IP地址">{{ locationInfo.ip }}</el-descriptions-item>
          <el-descriptions-item label="国家/地区">{{ locationInfo.country }}</el-descriptions-item>
          <el-descriptions-item label="省份">{{ locationInfo.province }}</el-descriptions-item>
          <el-descriptions-item label="城市">{{ locationInfo.city }}</el-descriptions-item>
          <el-descriptions-item label="运营商">{{ locationInfo.isp }}</el-descriptions-item>
          <el-descriptions-item label="查询时间">{{ locationInfo.queryTime }}</el-descriptions-item>
        </el-descriptions>
      </div>

      <div v-if="locationInfo" class="query-history">
        <el-divider>查询历史</el-divider>
        <el-table :data="queryHistory" style="width: 100%" border stripe>
          <el-table-column prop="ip" label="IP地址" width="180" />
          <el-table-column prop="location" label="位置" />
          <el-table-column prop="isp" label="运营商" width="180" />
          <el-table-column prop="queryTime" label="查询时间" width="180" />
        </el-table>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'

interface LocationInfo {
  ip: string
  country: string
  province: string
  city: string
  isp: string
  queryTime: string
}

interface QueryRecord {
  ip: string
  location: string
  isp: string
  queryTime: string
}

const formRef = ref<FormInstance>()
const loading = ref(false)
const locationInfo = ref<LocationInfo | null>(null)
const queryHistory = ref<QueryRecord[]>([])

const form = reactive({
  ip: ''
})

const rules: FormRules = {
  ip: [
    { required: true, message: '请输入IP地址', trigger: 'blur' },
    { 
      pattern: /^(\d{1,3}\.){3}\d{1,3}$/, 
      message: '请输入正确的IP地址格式', 
      trigger: 'blur' 
    }
  ]
}

const handleLocate = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        // TODO: Call backend API to get location info
        // Mock response for now
        await new Promise(resolve => setTimeout(resolve, 1000))
        const mockResponse = {
          ip: form.ip,
          country: '中国',
          province: '广东省',
          city: '深圳市',
          isp: '中国联通',
          queryTime: new Date().toLocaleString()
        }
        
        locationInfo.value = mockResponse
        queryHistory.value.unshift({
          ip: mockResponse.ip,
          location: `${mockResponse.country} ${mockResponse.province} ${mockResponse.city}`,
          isp: mockResponse.isp,
          queryTime: mockResponse.queryTime
        })
        
        // Keep only last 10 records
        if (queryHistory.value.length > 10) {
          queryHistory.value.pop()
        }
      } catch (error) {
        ElMessage.error('查询失败，请稍后重试')
      } finally {
        loading.value = false
      }
    }
  })
}
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

.location-info {
  margin-top: 20px;
}

.query-history {
  margin-top: 30px;
}
</style> 