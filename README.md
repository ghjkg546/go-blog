# 一站式夸克网盘管理后台

## 简介
近来，自媒体平台上各种给夸克引流的博主，本项目是一个夸克网盘资源管理后台，方便博主们建立自己的资源站

包含了，资源抓取、资源转存、批量创建资源分享并生成web页面、telegram搜索等功能，因为是业余搞的，目前有诸多不足之处，会慢慢完善

## 使用方法

1.修改config_example.yaml为 config.yaml,并填入你的自己的数据库配置
2.编译go二进制文件,并上传到你的服务器 

```
go build -o main main.go

chmod +x main
./main
```
访问localhost:8080,即可看到前台页面

3.编译后台文件，进入 https://gitee.com/ghjkg546/ziyuan 里的element-admin目录
```
pnpm install 
pnpmbuild
```

（可选）4.本项目包含一个简单的原生apk编译apk，进入flutter_application 
替换flutter_launcher_icons.yaml里的图标路径为你自己的图标路径
修改utils/request.dart里请求域名为你的域名
```
dart run flutter_launcher_icons
flutter build apk
```

得到dist目录上传服务器

4.如果要支持搜索引擎，需要配置zinsearch,用这个 docker-compose up -d
```
version: '3.8'

services:
zincsearch:
image: public.ecr.aws/zinclabs/zincsearch:latest
container_name: zincsearch
ports:
- "4080:4080"
environment:
- ZINC_DATA_PATH=/data
- ZINC_FIRST_ADMIN_USER=admin222
- ZINC_FIRST_ADMIN_PASSWORD=bbb222
volumes:
- ./data:/data
```
访问ip:4080,打开zinsearch后台
创建索引教程
https://blog.csdn.net/lzcs1/article/details/143569012?spm=1001.2014.3001.5502

zinsearch相关配置需要填入config.yaml
search:
    url: 'ip:4080'
    user_name: '101111' #后台用户的id
    password: 'cccdd111'

## 📑 详细功能目录
- [抓取资源](./docs/1_爬虫抓取文件.md)
- [转存资源](./docs/2_转存资源到你的网盘.md)
- [分享资源](./docs/3_分享资源.md)
- [资源列表](./docs/4_资源列表.md)





