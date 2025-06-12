const https = require('https');

// 登录获取token
const loginData = JSON.stringify({
  username: 'admin',
  password: '21232f297a57a5a743894a0e4a801fc3'
});

const loginOptions = {
  hostname: 'localhost',
  port: 8080,
  path: '/api/v1/login',
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Content-Length': Buffer.byteLength(loginData)
  }
};

const loginReq = https.request(loginOptions, (res) => {
  let data = '';
  res.on('data', (chunk) => {
    data += chunk;
  });
  res.on('end', () => {
    try {
      const response = JSON.parse(data);
      console.log('登录响应:', JSON.stringify(response, null, 2));
      
      if (response.code === 200 && response.data && response.data.token) {
        const token = response.data.token;
        console.log('获取到token:', token);
        
        // 测试命名空间API
        testNamespaceAPI(token);
      } else {
        console.error('登录失败:', response);
      }
    } catch (error) {
      console.error('解析登录响应失败:', error);
      console.log('原始响应:', data);
    }
  });
});

loginReq.on('error', (error) => {
  console.error('登录请求失败:', error);
});

loginReq.write(loginData);
loginReq.end();

function testNamespaceAPI(token) {
  const options = {
    hostname: 'localhost',
    port: 8080,
    path: '/api/v1/k8s-namespaces?configId=1',
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  };

  const req = https.request(options, (res) => {
    let data = '';
    res.on('data', (chunk) => {
      data += chunk;
    });
    res.on('end', () => {
      try {
        const response = JSON.parse(data);
        console.log('\n命名空间API响应:', JSON.stringify(response, null, 2));
        
        if (response.code === 200 && response.data) {
          console.log('\n命名空间列表:');
          response.data.forEach((namespace, index) => {
            console.log(`${index + 1}. ${namespace}`);
          });
          console.log(`\n总计: ${response.data.length} 个命名空间`);
        }
      } catch (error) {
        console.error('解析命名空间响应失败:', error);
        console.log('原始响应:', data);
      }
    });
  });

  req.on('error', (error) => {
    console.error('命名空间API请求失败:', error);
  });

  req.end();
}
