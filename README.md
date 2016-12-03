DOORMAN
--------

Doorman 是一个 HTTP 的单点登录系统。


Install
=======

```
	go build
```


Config
=======

	{
		// doorman 监听的地址与端口
		// :6252 监听所有网络接口的 6252 端口
		// 127.0.0.1:6252 监听 127.0.0.1 的 6252 端口
		"Listen": ":6252",
		
		// doorman 的接收用户名密码的接口
		// 当你的用户系统登录时，向此接口发送用户名与密码，返回一个 token
		// 通过多次调用此接口，相同的用户名与密码可以生产不同的 token
		// 例如：
		// >>> POST /token HTTP/1.1
		// >>> Host: doorman.example.com
		// >>> Content-Type: application/x-www-form-urlencoded
		// >>> Content-Length: 34
		// >>>
		// >>> username=myname&password=mypassword
		// <<< HTTP/1.1 201 Created
		// <<< Content-Type: text/plain
		// <<< Content-Length: 40
		// <<<
		// <<< 22e86969ee47cc4d4a0e6a5a12149867d6820100
		// 上例的 22e86969ee47cc4d4a0e6a5a12149867d6820100 就是返回的 token
		"TokenEndpoint": "/token",
		
		// 使用 token 进行登录，除 token 外，
		// 它还可以接收一个 url 可选参数，用于登录成功的重定向跳转。
		// 例如：
		// >>> GET /doorman-login?token=22e86969ee47cc4d4a0e6a5a12149867d6820100 HTTP/1.1
		// >>> Host: oa.example.com
		// <<< HTTP/1.1 302 Found
		// <<< Content-Length: 0
		// <<< Location: http://oa.example.com/default.jsp
		"LoginEndpoint": "/doorman-login",
		
		// 配置需要使用 token 登录的系统的登录请求
		// 登录请求一般是在登录页点击【登录】按钮后发起的 POST 请求
		// 这里配置的目的，就是模拟这个请求
		// 每个登录请求要以该系统的 HTTP Host 头作为 Key，
		// 其值为包含 Method, Url, Header, Body 的对象。
		// Method 是需要模拟的登录请求的方法
		// Url 是需要模拟的登录请求的 Url,
		// Header 是需要模拟的登录请求的头字段，一般都要有一个 Content-Type
		// Body 是需要模拟的登录请求的 Body，一般是用户名与密码。
		// 其中 URL 与 Body 在这个配置中是一个模板，可以使用
		// username, password, redirectUrl 这三个模板变量
		"LoginRequestes": {
			"oa.example.com": {
				"Method": "POST",
				"Url": "http://oa.example.com/login?url={redirectUrl}",
				"Header": {
					"Content-Type": ["application/x-www-form-urlencoded"]
				},
				"Body": "username={username}&password={password}"
			}
		},
		
		// 日志配置，参见 https://github.com/heiing/logs
		"Logs": {
			"types": ["debug", "info", "warn", "error"],
			"files": {
				"{AppPath}/info.log": ["debug", "info", "warn"],
				"{AppPath}/error.log": ["error"],
				"STDOUT": ["debug", "info", "warn"],
				"STDERR": ["error"]
			}
		}
	}


Usage
=====

doorman 是一个无入侵的轻量、简单的单点登录系统。在使用 doorman 之前，
应当先了解你的应用场景是否适用。

doorman 适用于使用用户名密码登录的 HTTP 系统，特别适用于登录第三方系统，
而无需修改第三方系统的源码或者要求第三方系统支持。所有系统必需能够使用相同
的用户名密码登录。

假定你拥有的系统如下：

- 协同办公系统，地址为 http://oa.example.com，登录请求地址为：
  http://oa.example.com/login.jsp，服务监听地址为 192.168.1.10:8010
- 客户关系管理系统，地址为 http://crm.example.com，登录请求地址为：
  http://crm.example.com/login.php，服务监听地址为 192.168.1.20:8020
- 员工中心，地址为 http://my.example.com，登录请求地址为：
  http://my.example.com/login

现在希望当员工登录员工中心后，能在员工中心单点登录到 oa 与 crm 系统。

为实现无入侵的单点登录，我们还需要一个反向代理服务器，用于代理 oa 与 crm
的请求，比如 nginx。下面的例子使用 nginx。

首先配置好 doorman，LoginRequestes 的配置要与你实际的登录请求想匹配即可。

	{
	    "Listen": "192.168.1.30:6252",
	    "TokenEndpoint": "/token",
	    "LoginEndpoint": "/doorman-login",

	    "LoginRequestes": {
	        "oa.example.com": {
	            "Method": "POST",
	            "Url": "http://oa.example.com/login.jsp",
	            "Header": {
	                "Content-Type": ["application/x-www-form-urlencoded"]
	            },
	            "Body": "username={username}&password={password}"
	        },
			"crm.example.com": {
	            "Method": "POST",
	            "Url": "http://crm.example.com/login.php",
	            "Header": {
	                "Content-Type": ["application/x-www-form-urlencoded"]
	            },
	            "Body": "username={username}&password={password}"
	        }
	    },
	    "Logs": {
	        "types": ["info", "warn", "error"],
	        "files": {
	            "{AppPath}/info.log": ["info", "warn"],
	            "{AppPath}/error.log": ["error"],
	            "STDOUT": ["info", "warn"],
	            "STDERR": ["error"] 
	        }
	    }
	}

第二，配置 oa.example.com 与 crm.example.com 的反向代理

	server {
	    server_name oa.example.com;
		# 此处的 location 为 LoginEndpoint 配置的值
        # 将这个请求反向代理到 doorman
	    location /doorman-login {
	        proxy_set_header Host $http_host;
	        proxy_pass http://192.168.1.30:6252;
	    }
	    location / {
	        proxy_set_header Host $http_host;
	        proxy_pass http://192.168.1.10:8010;
	    }
	}
	server {
	    server_name crm.example.com;
		# 此处的 location 为 LoginEndpoint 配置的值
        # 将这个请求反向代理到 doorman
	    location /doorman-login {
	        proxy_set_header Host $http_host;
	        proxy_pass http://192.168.1.30:6252;
	    }
	    location / {
	        proxy_set_header Host $http_host;
	        proxy_pass http://192.168.1.20:8020;
	    }
	}

最后，当你的员工中心在用户登录成功之后，将用户名密码发送到 doorman，
让 doorman 记住账号密码，并生成一个用于登录的 token。你可以将 token 
保存在员工中心的 session 或浏览器的 cookies 中，并显示 oa 与 crm 
的登录链接。

在员工中心登录成功后，调用 doorman 的 /token (配置的 TokenEndpoint)
生成 token。

例如发起请求：

	POST /token HTTP/1.1
	Host: 192.168.1.30:6252
	Content-Type: application/x-www-form-urlencoded
	Content-Length: 34

	username=myname&password=mypassword

得到以下响应

	HTTP/1.1 201 Created
	Content-Type: text/plain
	Content-Length: 40

	22e86969ee47cc4d4a0e6a5a12149867d6820100

你将获得的 token，生成 oa.example.com 与 crm.example.com 的登录链接：

- http://oa.example.com/doorman-login?token=22e86969ee47cc4d4a0e6a5a12149867d6820100
- http://crm.example.com/doorman-login?token=22e86969ee47cc4d4a0e6a5a12149867d6820100

当用户在员工中心点击这两个链接之后，即可登录 oa 与 crm。