package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// K8sPod K8s Pod模型
type K8sPod struct {
	ID            int64      `json:"id" gorm:"primaryKey;autoIncrement;comment:Pod ID"`
	ConfigID      int64      `json:"config_id" gorm:"not null;index;comment:K8s配置ID"`
	WorkloadID    *int64     `json:"workload_id" gorm:"index;comment:工作负载ID"`
	Name          string     `json:"name" gorm:"size:255;not null;comment:Pod名称"`
	Namespace     string     `json:"namespace" gorm:"size:100;not null;index;comment:命名空间"`
	WorkloadName  string     `json:"workload_name" gorm:"size:255;comment:工作负载名称"`
	WorkloadKind  string     `json:"workload_kind" gorm:"size:50;comment:工作负载类型"`
	Status        string     `json:"status" gorm:"size:50;not null;index;comment:Pod状态"`
	Phase         string     `json:"phase" gorm:"size:50;comment:Pod阶段"`
	NodeName      string     `json:"node_name" gorm:"size:255;comment:节点名称"`
	PodIP         string     `json:"pod_ip" gorm:"size:45;comment:Pod IP"`
	HostIP        string     `json:"host_ip" gorm:"size:45;comment:主机IP"`
	InstanceIP    string     `json:"instance_ip" gorm:"size:45;comment:实例IP"`
	CPURequest    *string    `json:"cpu_request" gorm:"size:20;comment:CPU请求"`
	CPULimit      *string    `json:"cpu_limit" gorm:"size:20;comment:CPU限制"`
	MemoryRequest *string    `json:"memory_request" gorm:"size:20;comment:内存请求"`
	MemoryLimit   *string    `json:"memory_limit" gorm:"size:20;comment:内存限制"`
	RestartCount  int        `json:"restart_count" gorm:"default:0;comment:重启次数"`
	StartTime     *time.Time `json:"start_time" gorm:"comment:启动时间"`
	CreatedAt     time.Time  `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"index;comment:删除时间"`
	Config        *K8sConfig `json:"config,omitempty" gorm:"foreignKey:ConfigID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;-"`
}

// TableName 指定表名
func (K8sPod) TableName() string {
	return "infra_k8s_pod"
}

// K8sPodResponse Pod响应结构
type K8sPodResponse struct {
	ID                  int64     `json:"id"`
	ConfigID            int64     `json:"config_id"`
	WorkloadID          *int64    `json:"workload_id"`
	Name                string    `json:"name"`
	Namespace           string    `json:"namespace"`
	WorkloadName        string    `json:"workload_name"`
	WorkloadKind        string    `json:"workload_kind"`
	Status              string    `json:"status"`
	Phase               string    `json:"phase"`
	NodeName            string    `json:"node_name"`
	PodIP               string    `json:"pod_ip"`
	HostIP              string    `json:"host_ip"`
	InstanceIP          string    `json:"instance_ip"`
	CPURequestLimits    string    `json:"cpu_request_limits"`    // CPU Request/Limits
	MemoryRequestLimits string    `json:"memory_request_limits"` // 内存 Request/Limits
	RestartCount        int       `json:"restart_count"`
	RunningTime         string    `json:"running_time"` // 运行时间
	StartTime           *time.Time `json:"start_time"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// ToResponse 转换为响应结构
func (p *K8sPod) ToResponse() *K8sPodResponse {
	return &K8sPodResponse{
		ID:                  p.ID,
		ConfigID:            p.ConfigID,
		WorkloadID:          p.WorkloadID,
		Name:                p.Name,
		Namespace:           p.Namespace,
		WorkloadName:        p.WorkloadName,
		WorkloadKind:        p.WorkloadKind,
		Status:              p.Status,
		Phase:               p.Phase,
		NodeName:            p.NodeName,
		PodIP:               p.PodIP,
		HostIP:              p.HostIP,
		InstanceIP:          p.InstanceIP,
		CPURequestLimits:    p.GetCPUResource(),
		MemoryRequestLimits: p.GetMemoryResource(),
		RestartCount:        p.RestartCount,
		RunningTime:         p.GetRunningTime(),
		StartTime:           p.StartTime,
		CreatedAt:           p.CreatedAt,
		UpdatedAt:           p.UpdatedAt,
	}
}

// GetCPUResource 获取CPU资源显示
func (p *K8sPod) GetCPUResource() string {
	var request, limit string
	
	if p.CPURequest != nil {
		request = formatCPUResource(*p.CPURequest)
	} else {
		request = "-"
	}
	
	if p.CPULimit != nil {
		limit = formatCPUResource(*p.CPULimit)
	} else {
		limit = "-"
	}
	
	return fmt.Sprintf("%s/%s", request, limit)
}

// GetMemoryResource 获取内存资源显示
func (p *K8sPod) GetMemoryResource() string {
	var request, limit string
	
	if p.MemoryRequest != nil {
		request = *p.MemoryRequest
	} else {
		request = "-"
	}
	
	if p.MemoryLimit != nil {
		limit = *p.MemoryLimit
	} else {
		limit = "-"
	}
	
	return fmt.Sprintf("%s/%s", request, limit)
}

// GetRunningTime 获取运行时间
func (p *K8sPod) GetRunningTime() string {
	if p.StartTime == nil {
		return "-"
	}
	
	duration := time.Since(*p.StartTime)
	
	days := int(duration.Hours()) / 24
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	
	if days > 0 {
		return fmt.Sprintf("%d天%d小时", days, hours)
	} else if hours > 0 {
		return fmt.Sprintf("%d小时%d分钟", hours, minutes)
	} else {
		return fmt.Sprintf("%d分钟", minutes)
	}
}

// formatCPUResource 格式化CPU资源显示
func formatCPUResource(cpu string) string {
	if cpu == "" {
		return "-"
	}
	
	// 处理毫核 (m)
	if strings.HasSuffix(cpu, "m") {
		if milliStr := strings.TrimSuffix(cpu, "m"); milliStr != "" {
			if milli, err := strconv.ParseFloat(milliStr, 64); err == nil {
				cores := milli / 1000
				return fmt.Sprintf("%.2f核", cores)
			}
		}
	}
	
	// 处理千核 (k)
	if strings.HasSuffix(cpu, "k") {
		if kiloStr := strings.TrimSuffix(cpu, "k"); kiloStr != "" {
			if kilo, err := strconv.ParseFloat(kiloStr, 64); err == nil {
				cores := kilo * 1000
				if cores >= 10000 {
					return fmt.Sprintf("%.2f万核", cores/10000)
				}
				return fmt.Sprintf("%.2f核", cores)
			}
		}
	}
	
	// 处理普通核数
	if cores, err := strconv.ParseFloat(cpu, 64); err == nil {
		if cores >= 10000 {
			return fmt.Sprintf("%.2f万核", cores/10000)
		}
		return fmt.Sprintf("%.2f核", cores)
	}
	
	return cpu
}
