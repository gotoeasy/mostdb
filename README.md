<div align=center>
<img src="https://gotoeasy.github.io/screenshots/mostdb/mostdb.png"/>
</div>


## `mostdb`

`mostdb`是个`golang`编写的分布式KV数据库，基于`goleveldb`保存数据，提供`get`、`set`、`del`等http接口。<br>
分布式算法非常简单，使用过半策略，只要大部分节点成功就认为是成功，顾名思义`mostdb`
<br>

[![Docker Pulls](https://img.shields.io/docker/pulls/gotoeasy/mostdb)](https://hub.docker.com/r/gotoeasy/mostdb)
[![GitHub release](https://img.shields.io/github/release/gotoeasy/mostdb.svg)](https://github.com/gotoeasy/mostdb/releases/latest)
<br>


## 特点
- [x] 使用`goleveldb`保存数据，适用于节省内存资源保存大量数据的场景
- [x] 基于过半策略的分布式集群支持，提供数据冗余性、高可靠性，配置部署极其简单
- [x] 提供`get`、`set`、`del`等http接口，开箱即用


## `docker`集群部署模式简易示例
```shell
# 以下3台以集群方式启动，配置本节点地址及关联节点地址即可，无中心节点，部署极其简单
# （也可以不配置GLC_SERVER_URL，仅配置MOSTDB_CLUSTER_URLS为不含本节点的其他节点地址即可）

# 节点1
docker run -d -p 5379:5379 MOSTDB_SERVER_URL=http://172.27.59.51:5379 \
       -e MOSTDB_CLUSTER_URLS=http://172.27.59.51:5379;http://172.27.59.52:5379;http://172.27.59.53:5379 \
       gotoeasy/mostdb

# 节点2
docker run -d -p 5379:5379 MOSTDB_SERVER_URL=http://172.27.59.52:5379 \
       -e MOSTDB_CLUSTER_URLS=http://172.27.59.51:5379;http://172.27.59.52:5379;http://172.27.59.53:5379 \
       gotoeasy/mostdb

# 节点3
docker run -d -p 5379:5379 MOSTDB_SERVER_URL=http://172.27.59.53:5379 \
       -e MOSTDB_CLUSTER_URLS=http://172.27.59.51:5379;http://172.27.59.52:5379;http://172.27.59.53:5379 \
       gotoeasy/mostdb
```


## 环境变量
- [x] `MOSTDB_STORE_ROOT`数据存储目录，默认`/opt/mostdb`
- [x] `MOSTDB_SERVER_PORT`服务端口，默认`5379`
- [x] `MOSTDB_SERVER_URL`本机节点服务地址，默认“”
- [x] `MOSTDB_CLUSTER_URLS`集群节点服务地址，默认“”，默认时等同于单机方式


## 更新履历

### 开发版`latest`

- [ ] 节点心跳监控
- [ ] 配置参数更新
- [ ] 审计日志
- [ ] 看板管理页面

### 初版`0.1.0`

- [x] 使用`goleveldb`保存数据，适用于节省内存资源保存大量数据的场景
- [x] 基于过半策略的分布式集群支持，提供数据冗余性、高可靠性，配置部署极其简单
- [x] 提供`get`、`set`、`del`等http接口，开箱即用
- [x] 提供`http`协议的`rest`服务接口`/mostdb/set`，存操作
- [x] 提供`http`协议的`rest`服务接口`/mostdb/get`，取操作
- [x] 提供`http`协议的`rest`服务接口`/mostdb/del`，删操作
