# 阿里云ddns
动态更新本地外网ip至域名解析服务中，你懂的...
 

### 一、首先你得有一个外网ip

具体申请方式可以看以下教程
https://www.bilibili.com/video/av64112633



### 二、修改相关核心配置
1. recordDomain: 指的是具体的根域名，例如xxx.com
2. recordRr: 指的是解析记录，例如要解析@.exmaple.com，主机记录要填写”@”，全部填*。
3. recordId: 记录id可以从控制台获取[点我跳转](https://api.aliyun.com/?spm=a2c1g.8271268.10000.127.412cdf25k57N6a#/?product=Alidns&version=2015-01-09&api=DescribeDomainRecords&params={}&tab=DEMO&lang=JAVA)
4. cronExp: 扫描策略，cron表达式，例如0/5 * * * * *



### 三、系统打包
1. 无压缩打包：运行go build即可
2. go自带打包：go build -ldflags '-w -s'
3. go自带打包 + upx压缩，可以看以下文章[点我跳转](https://www.jianshu.com/p/cd3c766b893c)



### 四、系统运行
1. 开发环境，go run main.go
2. 生产环境只需拷贝aliyunddns， config.yaml文件至同一目录，运行nohup ./aliyunddns > start.log 2> &1 &