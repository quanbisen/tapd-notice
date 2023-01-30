### TAPD消息通知到钉钉应用号
#### 1.前景

​        TAPD是公司内部项目工程管理的主要工具，其需求、任务、缺陷等的变更，需要有对应的消息通知，提高工作时效性。若公司内部使用的是企业微信协同办公软件的话，可以直接在TAPD上面关联公司即可实现消息通知功能，但因我们公司内部是使用的钉钉协同办公软件，需要实现相同的功能，需要开发一个通知服务。

#### 2.实现原理

​         TAPD拥有自动化任务执行，在每个项目里面的`设置`->`自动化助手`里面添加针对需求、任务、缺陷的自动化任务，任务执行的动作选择webhook，每当变更时TAPD会把变更的消息体发送到webhook地址上，由此webhook地址的服务负责消息的通知逻辑，本程序就是做的这个通知逻辑。

#### 3.支持的自动化任务类型

* 需求创建、状态变更、评论变更
* 任务创建、状态变更、评论变更
* 缺陷创建、状态变更、评论变更

​        其中评论支持@人单独通知

#### 4.配置说明

​        `config/settings.yml`是项目的配置文件

```yaml
settings:
  application:
    host: 0.0.0.0
    port: 8080
  database:
    source: root:10010@tcp(172.10.23.92:3306)/tapd-notice?charset=utf8&parseTime=True&loc=Local&timeout=1000ms
    maxidleconns: 5
    maxopenconns: 20
  dingding:
    appkey:                               # 钉钉应用的key
    appsecret:                            # 钉钉应用的密钥
    agentid:                              # 钉钉应用通知的agentid
  tapd:
    companyid:                            # TAPD的公司id
    apiuser:                              # TAPD的apiuser
    apipassword:                          # TAPD的apipassword
```

* 数据库

​        数据库只需创建数据库即可，程序运行时，使用的`gorm`框架会自动创建数据库表。

```sql
CREATE DATABASE `tapd-notice` CHARACTER SET 'utf8mb4';
```

* 钉钉

​        钉钉需要创建一个应用号，点击应用号可以获取到对应的`appkey`、`appsecret`及其`agentid`，除此，还需要开通对应的API权限，如获取部门、用户信息，发送消息等权限，请打开。因我本人的开发者权限已经被收回，无法提供截图了，请参考钉钉开发文档。

* TAPD

​        TAPD需要有API账号管理权限，点击`设置`->`API账号管理`即可查看，如下图所示：

![](https://images.cnblogs.com/cnblogs_com/quanbisen/2262984/o_230110094214_TAPD-API.png)

* TAPD自动化任务配置

​        以需求创建为例，创建自动化任务，如下图所示：

![](https://images.cnblogs.com/cnblogs_com/quanbisen/2262984/o_230110095517_TAPD-AUTO.png)

#### 5.实现效果

​        最终实现的效果如下图所示：

![](https://images.cnblogs.com/cnblogs_com/quanbisen/2262984/o_230110100004_Dingding-Comment1.png)

![](https://images.cnblogs.com/cnblogs_com/quanbisen/2262984/o_230110100015_Dingding-Comment2.png)

#### 6.源码说明

* Github地址：[tapd-notice](https://github.com/quanbisen/tapd-notice.git)

* 定时任务说明：项目里通知的用户数据，使用定时任务抓取钉钉API接口保存的，若第一下运行，请修改源码里面的定时任务的执行时间，先初始化用户的信息数据。

​        定时任务配置路径`cmd/cobra.go:74-76`

```go
goCron.AddJob("0 0 * * *", deptCronJob)
goCron.AddJob("30 0 * * *", userCronJob)
goCron.AddJob("0 0 * * *", projectCronJob)
```