package utils

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

// KubeconfigContext kubeconfig上下文结构
type KubeconfigContext struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

// KubeconfigContextItem kubeconfig上下文项
type KubeconfigContextItem struct {
	Context KubeconfigContext `yaml:"context"`
	Name    string            `yaml:"name"`
}

// KubeconfigCluster kubeconfig集群结构
type KubeconfigCluster struct {
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
	Server                   string `yaml:"server"`
}

// KubeconfigClusterItem kubeconfig集群项
type KubeconfigClusterItem struct {
	Cluster KubeconfigCluster `yaml:"cluster"`
	Name    string            `yaml:"name"`
}

// Kubeconfig kubeconfig结构
type Kubeconfig struct {
	APIVersion     string                   `yaml:"apiVersion"`
	Kind           string                   `yaml:"kind"`
	CurrentContext string                   `yaml:"current-context"`
	Contexts       []KubeconfigContextItem  `yaml:"contexts"`
	Clusters       []KubeconfigClusterItem  `yaml:"clusters"`
}

// ParseKubeconfig 解析kubeconfig内容
func ParseKubeconfig(kubeconfigContent string) (*Kubeconfig, error) {
	var config Kubeconfig
	err := yaml.Unmarshal([]byte(kubeconfigContent), &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse kubeconfig: %v", err)
	}
	return &config, nil
}

// GetClusterIDFromKubeconfig 从kubeconfig中获取集群ID
func GetClusterIDFromKubeconfig(kubeconfigContent string) (string, error) {
	config, err := ParseKubeconfig(kubeconfigContent)
	if err != nil {
		return "", err
	}

	// 如果没有当前上下文，返回错误
	if config.CurrentContext == "" {
		return "", fmt.Errorf("no current-context found in kubeconfig")
	}

	// 查找当前上下文对应的集群
	var clusterName string
	for _, ctx := range config.Contexts {
		if ctx.Name == config.CurrentContext {
			clusterName = ctx.Context.Cluster
			break
		}
	}

	if clusterName == "" {
		return "", fmt.Errorf("cluster not found for current context: %s", config.CurrentContext)
	}

	return clusterName, nil
}

// GetContextFromKubeconfig 从kubeconfig中获取上下文信息
func GetContextFromKubeconfig(kubeconfigContent string) (string, error) {
	config, err := ParseKubeconfig(kubeconfigContent)
	if err != nil {
		return "", err
	}

	return config.CurrentContext, nil
}

// GetClusterServerFromKubeconfig 从kubeconfig中获取集群服务器地址
func GetClusterServerFromKubeconfig(kubeconfigContent string) (string, error) {
	config, err := ParseKubeconfig(kubeconfigContent)
	if err != nil {
		return "", err
	}

	// 获取集群ID
	clusterID, err := GetClusterIDFromKubeconfig(kubeconfigContent)
	if err != nil {
		return "", err
	}

	// 查找集群服务器地址
	for _, cluster := range config.Clusters {
		if cluster.Name == clusterID {
			return cluster.Cluster.Server, nil
		}
	}

	return "", fmt.Errorf("cluster server not found for cluster: %s", clusterID)
}
