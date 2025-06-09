<template>
  <div :class="{'hidden': hidden}" class="pagination-container">
    <el-pagination
      :background="background"
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :layout="layout"
      :page-sizes="pageSizes"
      :total="total"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </div>
</template>

<script lang="ts">
import { defineComponent, computed } from 'vue'
import { scrollTo } from '@/utils/scroll-to'

export default defineComponent({
  name: 'Pagination',
  props: {
    total: {
      required: true,
      type: Number
    },
    page: {
      type: Number,
      default: 1
    },
    limit: {
      type: Number,
      default: 20
    },
    pageSizes: {
      type: Array,
      default: () => [10, 20, 30, 50]
    },
    layout: {
      type: String,
      default: 'total, sizes, prev, pager, next, jumper'
    },
    background: {
      type: Boolean,
      default: true
    },
    autoScroll: {
      type: Boolean,
      default: true
    },
    hidden: {
      type: Boolean,
      default: false
    }
  },
  emits: ['update:page', 'update:limit', 'pagination'],
  setup(props, { emit }) {
    const currentPage = computed({
      get: () => props.page,
      set: (val: number) => {
        emit('update:page', val)
      }
    })

    const pageSize = computed({
      get: () => props.limit,
      set: (val: number) => {
        emit('update:limit', val)
      }
    })

    const handleSizeChange = (val: number) => {
      emit('pagination', { page: currentPage.value, limit: val })
      if (props.autoScroll) {
        scrollTo(0)
      }
    }

    const handleCurrentChange = (val: number) => {
      emit('pagination', { page: val, limit: pageSize.value })
      if (props.autoScroll) {
        scrollTo(0)
      }
    }

    const scrollTo = (top: number) => {
      window.scrollTo({
        top,
        behavior: 'smooth'
      })
    }

    return {
      currentPage,
      pageSize,
      handleSizeChange,
      handleCurrentChange
    }
  }
})
</script>

<style lang="scss" scoped>
.pagination-container {
  background: var(--background-color);
  padding: 32px 16px;
  margin-top: 16px;
  border-radius: 4px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);

  &.hidden {
    display: none;
  }
}
</style> 