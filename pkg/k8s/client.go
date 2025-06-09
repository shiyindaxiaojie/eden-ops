package k8s

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// Client Kubernetes客户端
type Client struct {
	clientset     *kubernetes.Clientset
	metricsClient *versioned.Clientset
}

// NewClient 创建Kubernetes客户端
func NewClient(configContent string) (*Client, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(configContent))
	if err != nil {
		return nil, fmt.Errorf("failed to create config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %v", err)
	}

	metricsClient, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	return &Client{
		clientset:     clientset,
		metricsClient: metricsClient,
	}, nil
}

// GetClusterInfo 获取集群信息
func (c *Client) GetClusterInfo(ctx context.Context) (map[string]interface{}, error) {
	version, err := c.clientset.Discovery().ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get server version: %v", err)
	}

	nodes, err := c.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %v", err)
	}

	pods, err := c.clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %v", err)
	}

	// 计算资源总量和使用量
	var cpuTotal, cpuUsed, memoryTotal, memoryUsed float64
	for _, node := range nodes.Items {
		cpu := node.Status.Capacity.Cpu().MilliValue()
		memory := node.Status.Capacity.Memory().Value()
		cpuTotal += float64(cpu) / 1000
		memoryTotal += float64(memory) / 1024 / 1024 / 1024 // 转换为GB
	}

	nodeMetrics, err := c.metricsClient.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
	if err == nil {
		for _, metric := range nodeMetrics.Items {
			cpu := metric.Usage.Cpu().MilliValue()
			memory := metric.Usage.Memory().Value()
			cpuUsed += float64(cpu) / 1000
			memoryUsed += float64(memory) / 1024 / 1024 / 1024 // 转换为GB
		}
	}

	return map[string]interface{}{
		"version":      version.String(),
		"node_count":   len(nodes.Items),
		"pod_count":    len(pods.Items),
		"cpu_total":    cpuTotal,
		"cpu_used":     cpuUsed,
		"memory_total": memoryTotal,
		"memory_used":  memoryUsed,
	}, nil
}

// GetWorkloads 获取工作负载信息
func (c *Client) GetWorkloads(ctx context.Context) ([]map[string]interface{}, error) {
	var workloads []map[string]interface{}

	// 获取Deployments
	deployments, err := c.clientset.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %v", err)
	}
	for _, d := range deployments.Items {
		workloads = append(workloads, map[string]interface{}{
			"name":           d.Name,
			"namespace":      d.Namespace,
			"kind":           "Deployment",
			"replicas":       *d.Spec.Replicas,
			"ready_replicas": d.Status.ReadyReplicas,
			"status":         getDeploymentStatus(d.Status),
		})
	}

	// 获取StatefulSets
	statefulsets, err := c.clientset.AppsV1().StatefulSets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list statefulsets: %v", err)
	}
	for _, s := range statefulsets.Items {
		workloads = append(workloads, map[string]interface{}{
			"name":           s.Name,
			"namespace":      s.Namespace,
			"kind":           "StatefulSet",
			"replicas":       *s.Spec.Replicas,
			"ready_replicas": s.Status.ReadyReplicas,
			"status":         getStatefulSetStatus(s.Status),
		})
	}

	// 获取DaemonSets
	daemonsets, err := c.clientset.AppsV1().DaemonSets("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list daemonsets: %v", err)
	}
	for _, d := range daemonsets.Items {
		workloads = append(workloads, map[string]interface{}{
			"name":           d.Name,
			"namespace":      d.Namespace,
			"kind":           "DaemonSet",
			"replicas":       d.Status.DesiredNumberScheduled,
			"ready_replicas": d.Status.NumberReady,
			"status":         getDaemonSetStatus(d.Status),
		})
	}

	// 获取Jobs
	jobs, err := c.clientset.BatchV1().Jobs("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %v", err)
	}
	for _, j := range jobs.Items {
		workloads = append(workloads, map[string]interface{}{
			"name":           j.Name,
			"namespace":      j.Namespace,
			"kind":           "Job",
			"replicas":       *j.Spec.Completions,
			"ready_replicas": j.Status.Succeeded,
			"status":         getJobStatus(j.Status),
		})
	}

	// 尝试使用 BatchV1 获取 CronJobs
	cronjobs, err := c.clientset.BatchV1().CronJobs("").List(ctx, metav1.ListOptions{})
	if err != nil {
		// 如果失败，尝试使用 BatchV1Beta1
		betaCronjobs, betaErr := c.clientset.BatchV1beta1().CronJobs("").List(ctx, metav1.ListOptions{})
		if betaErr != nil {
			// 如果两个版本都失败，返回原始错误
			return nil, fmt.Errorf("failed to list cronjobs: %v", err)
		}
		// 使用 beta 版本的 CronJobs
		for _, cj := range betaCronjobs.Items {
			workloads = append(workloads, map[string]interface{}{
				"name":      cj.Name,
				"namespace": cj.Namespace,
				"kind":      "CronJob",
				"status":    getCronJobStatusBeta(cj.Status),
			})
		}
	} else {
		// 使用 v1 版本的 CronJobs
		for _, cj := range cronjobs.Items {
			workloads = append(workloads, map[string]interface{}{
				"name":      cj.Name,
				"namespace": cj.Namespace,
				"kind":      "CronJob",
				"status":    getCronJobStatus(cj.Status),
			})
		}
	}

	return workloads, nil
}

// 获取各类资源的状态
func getDeploymentStatus(status appsv1.DeploymentStatus) string {
	if status.ReadyReplicas == status.Replicas {
		return "Running"
	}
	return "Progressing"
}

func getStatefulSetStatus(status appsv1.StatefulSetStatus) string {
	if status.ReadyReplicas == status.Replicas {
		return "Running"
	}
	return "Progressing"
}

func getDaemonSetStatus(status appsv1.DaemonSetStatus) string {
	if status.NumberReady == status.DesiredNumberScheduled {
		return "Running"
	}
	return "Progressing"
}

func getJobStatus(status batchv1.JobStatus) string {
	if status.Succeeded > 0 {
		return "Completed"
	}
	if status.Active > 0 {
		return "Running"
	}
	if status.Failed > 0 {
		return "Failed"
	}
	return "Pending"
}

func getCronJobStatus(status batchv1.CronJobStatus) string {
	if status.LastScheduleTime != nil {
		return "Scheduled"
	}
	return "Pending"
}

func getCronJobStatusBeta(status batchv1beta1.CronJobStatus) string {
	if status.LastScheduleTime != nil {
		return "Scheduled"
	}
	return "Pending"
}
