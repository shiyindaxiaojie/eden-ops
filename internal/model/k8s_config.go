package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// K8sConfig Kubernetes配置模型
type K8sConfig struct {
	ID            int64      `gorm:"primaryKey" json:"id"`
	Name          string     `gorm:"type:varchar(100);not null" json:"name"`
	Description   string     `gorm:"type:text" json:"description"`
	Kubeconfig    string     `gorm:"type:text;not null" json:"kubeconfig"`
	ProviderId    *int64     `gorm:"column:provider_id" json:"providerId"`
	ProviderName  string     `gorm:"-" json:"providerName"`
	Status        int        `gorm:"type:tinyint;default:1" json:"status"`
	SyncInterval  int        `gorm:"type:int;default:30;comment:同步间隔(秒)" json:"syncInterval"`
	Version       string     `gorm:"type:varchar(20)" json:"version"`
	Context       string     `gorm:"type:varchar(100)" json:"context"`
	NodeCount     int        `gorm:"type:int;default:0" json:"nodeCount"`
	PodCount      int        `gorm:"type:int;default:0" json:"podCount"`
	CPUTotal      string     `gorm:"type:varchar(20)" json:"cpuTotal"`
	CPUUsed       string     `gorm:"type:varchar(20)" json:"cpuUsed"`
	MemoryTotal   string     `gorm:"type:varchar(20)" json:"memoryTotal"`
	MemoryUsed    string     `gorm:"type:varchar(20)" json:"memoryUsed"`
	WorkloadCount int        `gorm:"type:int;default:0;comment:工作负载数量" json:"workloadCount"`
	LastSyncTime  *time.Time `json:"lastSyncTime"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `gorm:"index" json:"-"`
}

// K8sWorkloadInfo Kubernetes工作负载信息
type K8sWorkloadInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Kind      string `json:"kind"`
	Replicas  int    `json:"replicas"`
	Status    string `json:"status"`
}

// TableName specifies the table name for K8sConfig
func (K8sConfig) TableName() string {
	return "infra_k8s_config"
}

// K8sConfigResponse K8s配置响应结构（包含工作负载统计）
type K8sConfigResponse struct {
	ID            int64      `json:"id"`
	Name          string     `json:"name"`
	Description   string     `json:"description"`
	Kubeconfig    string     `json:"kubeconfig"`
	ProviderId    *int64     `json:"providerId"`
	ProviderName  string     `json:"providerName"`
	Status        int        `json:"status"`
	SyncInterval  int        `json:"syncInterval"`
	Version       string     `json:"version"`
	Context       string     `json:"context"`
	NodeCount     int        `json:"nodeCount"`
	PodCount      int        `json:"podCount"`
	CPUTotal      string     `json:"cpuTotal"`
	CPUUsed       string     `json:"cpuUsed"`
	MemoryTotal   string     `json:"memoryTotal"`
	MemoryUsed    string     `json:"memoryUsed"`
	WorkloadCount int64      `json:"workloadCount"`
	LastSyncTime  *time.Time `json:"lastSyncTime"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

// ToResponse 转换为响应结构
func (c *K8sConfig) ToResponse() *K8sConfigResponse {
	return &K8sConfigResponse{
		ID:            c.ID,
		Name:          c.Name,
		Description:   c.Description,
		Kubeconfig:    c.Kubeconfig,
		ProviderId:    c.ProviderId,
		ProviderName:  c.ProviderName,
		Status:        c.Status,
		SyncInterval:  c.SyncInterval,
		Version:       c.Version,
		Context:       c.Context,
		NodeCount:     c.NodeCount,
		PodCount:      c.PodCount,
		CPUTotal:      c.FormatCPU(c.CPUTotal),
		CPUUsed:       c.FormatCPU(c.CPUUsed),
		MemoryTotal:   c.MemoryTotal,
		MemoryUsed:    c.MemoryUsed,
		WorkloadCount: int64(c.WorkloadCount),
		LastSyncTime:  c.LastSyncTime,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
	}
}

// FormatCPU 格式化CPU值，将毫核转换为核数
func (c *K8sConfig) FormatCPU(cpuStr string) string {
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
