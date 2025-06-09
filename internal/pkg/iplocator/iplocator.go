package iplocator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"eden-ops/pkg/config"

	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"

	cfs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	es "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
)

// ResourceType 定义资源类型
type ResourceType string

const (
	ResourceTypeCLB     ResourceType = "CLB"
	ResourceTypeCVM     ResourceType = "CVM"
	ResourceTypeCFS     ResourceType = "CFS"
	ResourceTypeMariaDB ResourceType = "MariaDB"
	ResourceTypeRedis   ResourceType = "Redis"
	ResourceTypeES      ResourceType = "ES"
	ResourceTypeCKafka  ResourceType = "CKafka"
)

// ResourceInfo 定义资源信息结构
type ResourceInfo struct {
	Type        ResourceType `json:"type"`
	InstanceID  string       `json:"instance_id"`
	Name        string       `json:"name"`
	PrivateIP   string       `json:"private_ip,omitempty"`
	PublicIP    string       `json:"public_ip,omitempty"`
	VIP         string       `json:"vip,omitempty"`
	Region      string       `json:"region"`
	Status      string       `json:"status"`
	CreateTime  *time.Time   `json:"create_time,omitempty"`
	Description string       `json:"description,omitempty"`
}

// TencentIPLocator 腾讯云IP定位器
type TencentIPLocator struct {
	SecretID  string
	SecretKey string
	Region    string
	Logger    *logrus.Logger

	clbClient     *clb.Client
	cvmClient     *cvm.Client
	cfsClient     *cfs.Client
	mariadbClient *mariadb.Client
	redisClient   *redis.Client
	esClient      *es.Client
	ckafkaClient  *ckafka.Client
	credential    *common.Credential
	cpf           *profile.ClientProfile
}

// NewTencentIPLocator 创建新的腾讯云IP定位器实例
func NewTencentIPLocator(secretID, secretKey, region string, logger *logrus.Logger) *TencentIPLocator {
	cfg := &config.TencentConfig{
		SecretID:  secretID,
		SecretKey: secretKey,
		Region:    region,
	}
	cfg.LoadFromEnv()
	if logger == nil {
		logger = logrus.New()
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:          true,
			DisableLevelTruncation: true,
			DisableTimestamp:       false,
			DisableSorting:         false,
			ForceColors:            true,
			DisableQuote:           false,
			PadLevelText:           true,
		})
	}

	cred := common.NewCredential(cfg.SecretID, cfg.SecretKey)
	cpf := profile.NewClientProfile()

	locator := &TencentIPLocator{
		SecretID:   cfg.SecretID,
		SecretKey:  cfg.SecretKey,
		Region:     cfg.Region,
		Logger:     logger,
		credential: cred,
		cpf:        cpf,
	}

	// 初始化所有客户端
	if err := locator.InitClients(); err != nil {
		logger.WithError(err).Error("初始化腾讯云客户端失败")
	}

	return locator
}

// LocateIP 定位IP地址
func (t *TencentIPLocator) LocateIP(ctx context.Context, ip string) ([]ResourceInfo, error) {
	// 创建一个带取消的上下文
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 创建通道
	resultCh := make(chan []ResourceInfo, len(t.queryFuncs()))
	errCh := make(chan error, len(t.queryFuncs()))

	// 使用WaitGroup来跟踪所有goroutine
	var wg sync.WaitGroup

	// 获取所有查询函数
	queryFuncs := t.queryFuncs()

	// 启动所有查询
	for _, queryFunc := range queryFuncs {
		wg.Add(1)
		go func(query func(context.Context, string, chan<- []ResourceInfo, chan<- error)) {
			defer wg.Done()
			// 为每个查询创建单独的通道
			ch := make(chan []ResourceInfo, 1)
			qErrCh := make(chan error, 1)

			// 执行查询
			query(ctx, ip, ch, qErrCh)

			// 处理查询结果
			select {
			case result := <-ch:
				resultCh <- result
			case err := <-qErrCh:
				errCh <- err
			case <-ctx.Done():
				errCh <- ctx.Err()
			}
		}(queryFunc)
	}

	// 在后台等待所有goroutine完成
	go func() {
		wg.Wait()
		close(resultCh)
		close(errCh)
	}()

	// 收集结果
	var allResources []ResourceInfo
	var errs []error

	// 从通道读取结果直到通道关闭
	for {
		select {
		case resources, ok := <-resultCh:
			if !ok {
				// 所有查询已完成
				if len(errs) > 0 {
					return allResources, fmt.Errorf("查询过程中发生错误: %v", errs[0])
				}
				return allResources, nil
			}
			allResources = append(allResources, resources...)
		case err, ok := <-errCh:
			if !ok {
				continue
			}
			errs = append(errs, err)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// queryFuncs 返回所有查询函数
func (t *TencentIPLocator) queryFuncs() []func(context.Context, string, chan<- []ResourceInfo, chan<- error) {
	return []func(context.Context, string, chan<- []ResourceInfo, chan<- error){
		t.queryCLB,
		t.queryCVM,
		t.queryCFS,
		t.queryMariaDB,
		t.queryRedis,
		t.queryES,
		t.queryCKafka,
	}
}

// 初始化所有客户端
func (t *TencentIPLocator) InitClients() error {
	var err error

	t.clbClient, err = clb.NewClient(t.credential, t.Region, t.cpf)
	if err != nil {
		return fmt.Errorf("创建CLB客户端失败: %v", err)
	}

	t.cvmClient, err = cvm.NewClient(t.credential, t.Region, t.cpf)
	if err != nil {
		return fmt.Errorf("创建CVM客户端失败: %v", err)
	}

	t.cfsClient, err = cfs.NewClient(t.credential, t.Region, t.cpf)
	if err != nil {
		return fmt.Errorf("创建CFS客户端失败: %v", err)
	}

	t.mariadbClient, err = mariadb.NewClient(t.credential, t.Region, t.cpf)
	if err != nil {
		return fmt.Errorf("创建MariaDB客户端失败: %v", err)
	}

	t.redisClient, err = redis.NewClient(t.credential, t.Region, t.cpf)
	if err != nil {
		return fmt.Errorf("创建Redis客户端失败: %v", err)
	}

	t.esClient, err = es.NewClient(t.credential, t.Region, t.cpf)
	if err != nil {
		return fmt.Errorf("创建ES客户端失败: %v", err)
	}

	t.ckafkaClient, err = ckafka.NewClient(t.credential, t.Region, t.cpf)
	if err != nil {
		return fmt.Errorf("创建CKafka客户端失败: %v", err)
	}

	return nil
}

// 以下是各种资源查询方法的实现
func (t *TencentIPLocator) queryCLB(ctx context.Context, ip string, ch chan<- []ResourceInfo, errCh chan<- error) {
	// 使用已初始化的CLB客户端
	// 创建请求对象
	request := clb.NewDescribeLoadBalancersRequest()

	// 先查询公网CLB
	request.LoadBalancerVips = common.StringPtrs([]string{ip})
	request.LoadBalancerType = common.StringPtr("OPEN")

	// 发送请求
	response, err := t.clbClient.DescribeLoadBalancers(request)
	if err != nil {
		errCh <- fmt.Errorf("查询公网CLB失败: %v", err)
		return
	}

	// 处理响应
	var resources []ResourceInfo
	for _, lb := range response.Response.LoadBalancerSet {
		resources = append(resources, ResourceInfo{
			Type:       ResourceTypeCLB,
			InstanceID: *lb.LoadBalancerId,
			Name:       *lb.LoadBalancerName,
			VIP:        *lb.LoadBalancerVips[0],
			Region:     t.Region,
			Status:     fmt.Sprintf("%v", *lb.Status),
		})
	}

	// 如果公网CLB没有找到，继续查询内网CLB
	if len(resources) == 0 {
		request.LoadBalancerType = common.StringPtr("INTERNAL")
		response, err = t.clbClient.DescribeLoadBalancers(request)
		if err != nil {
			errCh <- fmt.Errorf("查询内网CLB失败: %v", err)
			return
		}

		for _, lb := range response.Response.LoadBalancerSet {
			resources = append(resources, ResourceInfo{
				Type:       ResourceTypeCLB,
				InstanceID: *lb.LoadBalancerId,
				Name:       *lb.LoadBalancerName,
				VIP:        *lb.LoadBalancerVips[0],
				Region:     t.Region,
				Status:     fmt.Sprintf("%v", *lb.Status),
			})
		}
	}

	// 发送结果
	select {
	case <-ctx.Done():
		errCh <- ctx.Err()
	case ch <- resources:
	}
}

func (t *TencentIPLocator) queryCVM(ctx context.Context, ip string, ch chan<- []ResourceInfo, errCh chan<- error) {
	// 使用已初始化的CVM客户端
	// 创建请求对象
	request := cvm.NewDescribeInstancesRequest()

	// 先查询私有IP
	request.Filters = []*cvm.Filter{
		{
			Name:   common.StringPtr("private-ip-address"),
			Values: common.StringPtrs([]string{ip}),
		},
	}

	// 发送请求
	response, err := t.cvmClient.DescribeInstances(request)
	if err != nil {
		errCh <- fmt.Errorf("查询CVM实例失败: %v", err)
		return
	}

	// 处理响应
	var resources []ResourceInfo
	for _, instance := range response.Response.InstanceSet {
		createTime, _ := time.Parse("2006-01-02T15:04:05Z", *instance.CreatedTime)
		resources = append(resources, ResourceInfo{
			Type:       ResourceTypeCVM,
			InstanceID: *instance.InstanceId,
			Name:       *instance.InstanceName,
			PrivateIP:  *instance.PrivateIpAddresses[0],
			PublicIP:   getPtrValue(instance.PublicIpAddresses),
			Region:     t.Region,
			Status:     *instance.InstanceState,
			CreateTime: &createTime,
		})
	}

	// 如果私有IP没有找到，继续查询公网IP
	if len(resources) == 0 {
		request.Filters = []*cvm.Filter{
			{
				Name:   common.StringPtr("public-ip-address"),
				Values: common.StringPtrs([]string{ip}),
			},
		}

		response, err = t.cvmClient.DescribeInstances(request)
		if err != nil {
			errCh <- fmt.Errorf("查询CVM实例失败: %v", err)
			return
		}

		for _, instance := range response.Response.InstanceSet {
			createTime, _ := time.Parse("2006-01-02T15:04:05Z", *instance.CreatedTime)
			resources = append(resources, ResourceInfo{
				Type:       ResourceTypeCVM,
				InstanceID: *instance.InstanceId,
				Name:       *instance.InstanceName,
				PrivateIP:  *instance.PrivateIpAddresses[0],
				PublicIP:   getPtrValue(instance.PublicIpAddresses),
				Region:     t.Region,
				Status:     *instance.InstanceState,
				CreateTime: &createTime,
			})
		}
	}

	// 发送结果
	select {
	case <-ctx.Done():
		errCh <- ctx.Err()
	case ch <- resources:
	}
}

// getPtrValue 安全地获取字符串指针的值
func getPtrValue(ptr []*string) string {
	if len(ptr) > 0 && ptr[0] != nil {
		return *ptr[0]
	}
	return ""
}

func (t *TencentIPLocator) queryCFS(ctx context.Context, ip string, ch chan<- []ResourceInfo, errCh chan<- error) {
	// 使用已初始化的CFS客户端
	request := cfs.NewDescribeCfsFileSystemsRequest()
	response, err := t.cfsClient.DescribeCfsFileSystems(request)
	if err != nil {
		errCh <- fmt.Errorf("查询CFS文件系统失败: %v", err)
		return
	}

	// 处理响应
	var resources []ResourceInfo
	for _, fs := range response.Response.FileSystems {
		// 查询挂载点信息
		mountRequest := cfs.NewDescribeMountTargetsRequest()
		mountRequest.FileSystemId = fs.FileSystemId
		mountResponse, err := t.cfsClient.DescribeMountTargets(mountRequest)
		if err != nil {
			errCh <- fmt.Errorf("查询CFS挂载点失败: %v", err)
			continue
		}

		// 检查挂载点IP是否匹配
		for _, mp := range mountResponse.Response.MountTargets {
			if *mp.IpAddress == ip {
				createTime, _ := time.Parse("2006-01-02 15:04:05", *fs.CreationTime)
				resources = append(resources, ResourceInfo{
					Type:        ResourceTypeCFS,
					InstanceID:  *fs.FileSystemId,
					Name:        *fs.FsName,
					PrivateIP:   *mp.IpAddress,
					Region:      t.Region,
					Status:      *fs.LifeCycleState,
					CreateTime:  &createTime,
					Description: "",
				})
				break
			}
		}
	}

	// 发送结果
	select {
	case <-ctx.Done():
		errCh <- ctx.Err()
	case ch <- resources:
	}
}

func (t *TencentIPLocator) queryMariaDB(ctx context.Context, ip string, ch chan<- []ResourceInfo, errCh chan<- error) {
	try := func() error {
		// 使用已初始化的MariaDB客户端
		request := mariadb.NewDescribeDBInstancesRequest()
		response, err := t.mariadbClient.DescribeDBInstances(request)
		if err != nil {
			return fmt.Errorf("查询MariaDB实例失败: %v", err)
		}

		// 处理响应
		var resources []ResourceInfo
		for _, instance := range response.Response.Instances {
			// 检查实例IP是否匹配
			if *instance.Vip == ip {
				createTime, _ := time.Parse("2006-01-02T15:04:05Z", *instance.CreateTime)
				resources = append(resources, ResourceInfo{
					Type:        ResourceTypeMariaDB,
					InstanceID:  *instance.InstanceId,
					Name:        *instance.InstanceName,
					PrivateIP:   *instance.Vip,
					Region:      t.Region,
					Status:      *instance.StatusDesc,
					CreateTime:  &createTime,
					Description: getPtrValue([]*string{instance.StatusDesc}),
				})
			}
		}

		ch <- resources
		return nil
	}

	if err := try(); err != nil {
		t.Logger.WithError(err).Error("查询MariaDB资源失败")
		errCh <- err
	}
	close(ch)
	close(errCh)
}

func (t *TencentIPLocator) queryRedis(ctx context.Context, ip string, ch chan<- []ResourceInfo, errCh chan<- error) {
	// 使用已初始化的Redis客户端
	request := redis.NewDescribeInstancesRequest()
	response, err := t.redisClient.DescribeInstances(request)
	if err != nil {
		errCh <- fmt.Errorf("查询Redis实例失败: %v", err)
		return
	}

	// 处理响应
	var resources []ResourceInfo
	for _, instance := range response.Response.InstanceSet {
		// 检查实例IP是否匹配
		if *instance.WanIp == ip || *instance.IPv6 == ip {
			createTime, _ := time.Parse("2006-01-02T15:04:05Z", *instance.Createtime)
			resources = append(resources, ResourceInfo{
				Type:        ResourceTypeRedis,
				InstanceID:  *instance.InstanceId,
				Name:        *instance.InstanceName,
				PrivateIP:   *instance.WanIp,
				PublicIP:    *instance.WanIp,
				Region:      t.Region,
				Status:      fmt.Sprintf("%v", *instance.Status),
				CreateTime:  &createTime,
				Description: getPtrValue([]*string{instance.InstanceTitle}),
			})
		}
	}

	// 发送结果
	select {
	case <-ctx.Done():
		errCh <- ctx.Err()
	case ch <- resources:
	}
}

func (t *TencentIPLocator) queryES(ctx context.Context, ip string, ch chan<- []ResourceInfo, errCh chan<- error) {
	// 使用已初始化的ES客户端
	request := es.NewDescribeInstancesRequest()
	response, err := t.esClient.DescribeInstances(request)
	if err != nil {
		errCh <- fmt.Errorf("查询ES实例失败: %v", err)
		return
	}

	// 处理响应
	var resources []ResourceInfo
	for _, instance := range response.Response.InstanceList {
		// 检查实例IP是否匹配
		if instance.PublicAccess != nil && *instance.PublicAccess == ip {
			createTime, _ := time.Parse("2006-01-02T15:04:05Z", *instance.CreateTime)
			resources = append(resources, ResourceInfo{
				Type:        ResourceTypeES,
				InstanceID:  *instance.InstanceId,
				Name:        *instance.InstanceName,
				PublicIP:    *instance.PublicAccess,
				PrivateIP:   *instance.VpcUid,
				Region:      t.Region,
				Status:      fmt.Sprintf("%v", *instance.Status),
				CreateTime:  &createTime,
				Description: "",
			})
		}
	}

	// 发送结果
	select {
	case <-ctx.Done():
		errCh <- ctx.Err()
	case ch <- resources:
	}
}

func (t *TencentIPLocator) queryCKafka(ctx context.Context, ip string, ch chan<- []ResourceInfo, errCh chan<- error) {
	// 获取CKafka实例列表
	listRequest := ckafka.NewDescribeInstancesRequest()
	listResponse, err := t.ckafkaClient.DescribeInstances(listRequest)
	if err != nil {
		errCh <- fmt.Errorf("查询CKafka实例列表失败: %v", err)
		return
	}

	// 处理响应
	var resources []ResourceInfo
	for _, instance := range listResponse.Response.Result.InstanceList {
		// 获取实例属性
		attrRequest := ckafka.NewDescribeInstanceAttributesRequest()
		attrRequest.InstanceId = instance.InstanceId
		attrResponse, err := t.ckafkaClient.DescribeInstanceAttributes(attrRequest)
		if err != nil {
			errCh <- fmt.Errorf("查询CKafka实例属性失败: %v", err)
			continue
		}

		// 检查VIP列表中的IP是否匹配
		for _, vip := range attrResponse.Response.Result.VipList {
			if *vip.Vip == ip {
				createTime := time.Unix(int64(*attrResponse.Response.Result.CreateTime), 0)
				resources = append(resources, ResourceInfo{
					Type:        ResourceTypeCKafka,
					InstanceID:  *instance.InstanceId,
					Name:        *instance.InstanceName,
					PrivateIP:   *vip.Vip,
					Region:      t.Region,
					Status:      fmt.Sprintf("%v", *instance.Status),
					CreateTime:  &createTime,
					Description: fmt.Sprintf("Port: %s", *vip.Vport),
				})
				break
			}
		}
	}

	// 发送结果
	select {
	case <-ctx.Done():
		errCh <- ctx.Err()
	case ch <- resources:
	}
}

// 调整日志格式化配置
func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   false,
	})
}
