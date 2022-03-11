这是一个基于Gin的编辑流程图的网站，将流程图保存到阿里云OSS存储
## 交互流程图

![draw-io-with-current-site](https://image.ahsup.top//draw/draw-io-with-current-site.png)

## 特性
网站有如下特点：
* 支持jwt鉴权：账号密码存储在配置文件中；
* 无任何其他依赖：二进程文件，只依赖阿里云oss；
* 对服务器要求低：云服务器1C0.5G即可满足；
* 数据完全可控：数据存储在服务器本地和OSS中；