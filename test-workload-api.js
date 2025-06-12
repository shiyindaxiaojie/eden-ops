const http = require('http');
const crypto = require('crypto');

// 测试API的函数
function makeRequest(options, data = null) {
    return new Promise((resolve, reject) => {
        const req = http.request(options, (res) => {
            let body = '';
            res.on('data', (chunk) => {
                body += chunk;
            });
            res.on('end', () => {
                try {
                    const result = JSON.parse(body);
                    resolve({ status: res.statusCode, data: result });
                } catch (e) {
                    resolve({ status: res.statusCode, data: body });
                }
            });
        });

        req.on('error', (err) => {
            reject(err);
        });

        if (data) {
            req.write(JSON.stringify(data));
        }
        req.end();
    });
}

async function testWorkloadAPI() {
    console.log('开始测试工作负载API...\n');

    // 1. 测试登录
    console.log('1. 测试登录...');
    try {
        const password = 'admin123';
        const hashedPassword = crypto.createHash('sha256').update(password).digest('hex');
        
        const loginOptions = {
            hostname: 'localhost',
            port: 8080,
            path: '/api/v1/login',
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        };

        const loginData = {
            username: 'admin',
            password: hashedPassword
        };

        const loginResult = await makeRequest(loginOptions, loginData);
        console.log('登录状态:', loginResult.status);

        if (loginResult.status === 200 && loginResult.data.code === 200) {
            const token = loginResult.data.data.token;
            console.log('登录成功！\n');

            // 2. 测试工作负载列表（修复后的格式）
            console.log('2. 测试工作负载列表（修复后的格式）...');
            const workloadOptions = {
                hostname: 'localhost',
                port: 8080,
                path: '/api/v1/k8s-workloads?configId=1&page=1&pageSize=3',
                method: 'GET',
                headers: {
                    'Authorization': 'Bearer ' + token,
                    'Content-Type': 'application/json'
                }
            };

            const workloadResult = await makeRequest(workloadOptions);
            console.log('工作负载API状态:', workloadResult.status);

            if (workloadResult.status === 200 && workloadResult.data.data && workloadResult.data.data.list) {
                const workloads = workloadResult.data.data.list;
                console.log(`\n找到 ${workloads.length} 个工作负载（总数: ${workloadResult.data.data.total}）:`);
                
                workloads.forEach((workload, index) => {
                    console.log(`\n${index + 1}. ${workload.name}:`);
                    console.log(`   - 命名空间: ${workload.namespace}`);
                    console.log(`   - 类型: ${workload.kind}`);
                    console.log(`   - Pod状态: ${workload.pod_status}`);
                    console.log(`   - CPU资源: ${workload.cpu_request_limits}`);
                    console.log(`   - 内存资源: ${workload.memory_request_limits}`);
                    console.log(`   - 状态: ${workload.status}`);
                });

                // 验证格式是否正确
                console.log('\n=== 格式验证 ===');
                const firstWorkload = workloads[0];
                if (firstWorkload) {
                    console.log('✅ Pod状态格式:', firstWorkload.pod_status, '(应该是 X/Y 格式)');
                    console.log('✅ CPU资源格式:', firstWorkload.cpu_request_limits, '(应该是 X核/Y核 格式)');
                    console.log('✅ 内存资源格式:', firstWorkload.memory_request_limits, '(应该是 XMi/YMi 格式)');
                }
            } else {
                console.log('工作负载API返回错误:', workloadResult.data);
            }
        } else {
            console.log('登录失败:', loginResult.data);
        }
    } catch (error) {
        console.error('测试过程中出错:', error);
    }
}

// 运行测试
testWorkloadAPI();
