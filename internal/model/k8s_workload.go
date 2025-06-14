package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// K8sWorkload Kubernetes工作负载
type K8sWorkload struct {
	ID            int64      `json:"id" gorm:"primaryKey"`
	ConfigID      int64      `json:"config_id" gorm:"not null"`
	Name          string     `json:"name" gorm:"size:100;not null"`
	Namespace     string     `json:"namespace" gorm:"size:63;not null"`
	Kind          string     `json:"kind" gorm:"size:20;not null"`
	Replicas      int        `json:"replicas" gorm:"default:0"`
	ReadyReplicas int        `json:"ready_replicas" gorm:"default:0"`
	Status        string     `json:"status" gorm:"size:20"`
	Labels        *string    `json:"labels" gorm:"type:text"`
	Selector      *string    `json:"selector" gorm:"type:text"`
	Images        *string    `json:"images" gorm:"type:text"`
	CPURequest    *string    `json:"cpu_request" gorm:"size:20"`
	CPULimit      *string    `json:"cpu_limit" gorm:"size:20"`
	MemoryRequest *string    `json:"memory_request" gorm:"size:20"`
	MemoryLimit   *string    `json:"memory_limit" gorm:"size:20"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"index"`
	Config        *K8sConfig `json:"config,omitempty" gorm:"foreignKey:ConfigID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;-"`
}

// GetCPUResource 获取CPU资源配置（请求/限制）
func (w *K8sWorkload) GetCPUResource() string {
	if w.CPURequest != nil && w.CPULimit != nil {
		// 获取两个值的核数
		requestCores := w.parseCPUToCores(*w.CPURequest)
		limitCores := w.parseCPUToCores(*w.CPULimit)

		// 选择合适的单位（如果任一值>=1万核，都用万核）
		if requestCores >= 10000 || limitCores >= 10000 {
			requestWan := requestCores / 10000
			limitWan := limitCores / 10000
			return fmt.Sprintf("%.2f万核/%.2f万核", requestWan, limitWan)
		} else {
			return fmt.Sprintf("%.2f核/%.2f核", requestCores, limitCores)
		}
	}
	if w.CPURequest != nil {
		return w.formatCPU(*w.CPURequest)
	}
	if w.CPULimit != nil {
		return w.formatCPU(*w.CPULimit)
	}
	return "-"
}

// GetMemoryResource 获取内存资源配置（请求/限制）
func (w *K8sWorkload) GetMemoryResource() string {
	if w.MemoryRequest != nil && w.MemoryLimit != nil {
		return fmt.Sprintf("%s/%s", *w.MemoryRequest, *w.MemoryLimit)
	}
	if w.MemoryRequest != nil {
		return *w.MemoryRequest
	}
	if w.MemoryLimit != nil {
		return *w.MemoryLimit
	}
	return "-"
}

// GetPodStatus 获取Pod状态（运行/期望）
func (w *K8sWorkload) GetPodStatus() string {
	return fmt.Sprintf("%d/%d", w.ReadyReplicas, w.Replicas)
}

// GetLabelsMap 获取标签映射
func (w *K8sWorkload) GetLabelsMap() map[string]string {
	if w.Labels == nil {
		return nil
	}
	var labels map[string]string
	if err := json.Unmarshal([]byte(*w.Labels), &labels); err != nil {
		return nil
	}
	return labels
}

// GetSelectorMap 获取选择器映射
func (w *K8sWorkload) GetSelectorMap() map[string]string {
	if w.Selector == nil {
		return nil
	}
	var selector map[string]string
	if err := json.Unmarshal([]byte(*w.Selector), &selector); err != nil {
		return nil
	}
	return selector
}

// GetImagesList 获取镜像列表
func (w *K8sWorkload) GetImagesList() []string {
	if w.Images == nil {
		return nil
	}
	var images []string
	if err := json.Unmarshal([]byte(*w.Images), &images); err != nil {
		return nil
	}
	return images
}

// K8sWorkloadResponse 工作负载响应结构
type K8sWorkloadResponse struct {
	ID                  int64             `json:"id"`
	ConfigID            int64             `json:"config_id"`
	Name                string            `json:"name"`
	Namespace           string            `json:"namespace"`
	Kind                string            `json:"kind"`
	Replicas            int               `json:"replicas"`
	ReadyReplicas       int               `json:"ready_replicas"`
	PodStatus           string            `json:"pod_status"` // 运行/期望Pod数量
	Status              string            `json:"status"`
	Labels              map[string]string `json:"labels,omitempty"`
	Selector            map[string]string `json:"selector,omitempty"`
	Images              []string          `json:"images,omitempty"`
	CPURequestLimits    string            `json:"cpu_request_limits"`    // CPU Request/Limits
	MemoryRequestLimits string            `json:"memory_request_limits"` // 内存 Request/Limits
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`
}

// ToResponse 转换为响应结构
func (w *K8sWorkload) ToResponse() *K8sWorkloadResponse {
	return &K8sWorkloadResponse{
		ID:                  w.ID,
		ConfigID:            w.ConfigID,
		Name:                w.Name,
		Namespace:           w.Namespace,
		Kind:                w.Kind,
		Replicas:            w.Replicas,
		ReadyReplicas:       w.ReadyReplicas,
		PodStatus:           w.GetPodStatus(),
		Status:              w.Status,
		Labels:              w.GetLabelsMap(),
		Selector:            w.GetSelectorMap(),
		Images:              w.GetImagesList(),
		CPURequestLimits:    w.GetCPUResource(),
		MemoryRequestLimits: w.GetMemoryResource(),
		CreatedAt:           w.CreatedAt,
		UpdatedAt:           w.UpdatedAt,
	}
}

// parseCPUToCores 解析CPU字符串为核数（不带单位）
func (w *K8sWorkload) parseCPUToCores(cpuStr string) float64 {
	if cpuStr == "" {
		return 0
	}

	var cores float64
	var err error

	// 处理毫核格式 (如 27500m)
	if strings.HasSuffix(cpuStr, "m") {
		milliCores := strings.TrimSuffix(cpuStr, "m")
		if value, parseErr := strconv.ParseFloat(milliCores, 64); parseErr == nil {
			cores = value / 1000
		} else {
			return 0
		}
	} else if strings.HasSuffix(cpuStr, "k") {
		// 处理千核格式 (如 5760k)
		kiloCores := strings.TrimSuffix(cpuStr, "k")
		if value, parseErr := strconv.ParseFloat(kiloCores, 64); parseErr == nil {
			cores = value * 1000
		} else {
			return 0
		}
	} else {
		// 处理已经是核数的格式 (如 27.5)
		if cores, err = strconv.ParseFloat(cpuStr, 64); err != nil {
			return 0
		}
	}

	return cores
}

// formatCPU 格式化CPU值，将毫核转换为核数
func (w *K8sWorkload) formatCPU(cpuStr string) string {
	if cpuStr == "" {
		return ""
	}

	var cores float64
	var err error

	// 处理毫核格式 (如 27500m)
	if strings.HasSuffix(cpuStr, "m") {
		milliCores := strings.TrimSuffix(cpuStr, "m")
		if value, parseErr := strconv.ParseFloat(milliCores, 64); parseErr == nil {
			cores = value / 1000
		} else {
			return cpuStr
		}
	} else if strings.HasSuffix(cpuStr, "k") {
		// 处理千核格式 (如 5760k)
		kiloCores := strings.TrimSuffix(cpuStr, "k")
		if value, parseErr := strconv.ParseFloat(kiloCores, 64); parseErr == nil {
			cores = value * 1000
		} else {
			return cpuStr
		}
	} else {
		// 处理已经是核数的格式 (如 27.5)
		if cores, err = strconv.ParseFloat(cpuStr, 64); err != nil {
			return cpuStr
		}
	}

	// 根据核数大小选择合适的单位
	if cores >= 10000 {
		// 大于等于1万核，显示为万核
		wanCores := cores / 10000
		return fmt.Sprintf("%.2f万核", wanCores)
	} else {
		// 小于1万核，显示为核
		return fmt.Sprintf("%.2f核", cores)
	}
}

// TableName 表名
func (K8sWorkload) TableName() string {
	return "infra_k8s_workload"
}
