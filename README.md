<div align=center>
<img src="https://gotoeasy.github.io/screenshots/mostdb/mostdb.png"/>
</div>


# MostDB

`MostDB`是个分布式KV数据库，基于`goleveldb`保存数据，提供`get`、`set`、`del`http接口。<br>
分布式算法非常简单，只要大部分节点成功就认为是成功，顾名思义`MostDB`
<br>

[![Docker Pulls](https://img.shields.io/docker/pulls/gotoeasy/glc)](https://hub.docker.com/r/gotoeasy/glc)
[![GitHub release](https://img.shields.io/github/release/gotoeasy/glogcenter.svg)](https://github.com/gotoeasy/glogcenter/releases/latest)
<br>


## 特点
- [x] 使用`goleveldb`保存数据，适用于节省内存资源保存大量数据的场景
- [x] 分布式支持，配置部署极其简单，提供数据冗余性、高可靠性
- [x] 提供`get`、`set`、`del`http接口，开箱即用


## 部署示例
```shell
# 以下3台以集群方式启动，配置关联节点地址即可

# 节点1： http://172.27.59.51:5379
export GLC_CLUSTER_URLS='http://172.27.59.52:5379;http://172.27.59.53:5379'
mostdb

# 节点2： http://172.27.59.52:5379
export GLC_CLUSTER_URLS='http://172.27.59.51:5379;http://172.27.59.53:5379'
mostdb

# 节点3： http://172.27.59.53:5379
export GLC_CLUSTER_URLS='http://172.27.59.51:5379;http://172.27.59.52:5379'
mostdb
```

## 环境变量
- [x] `MOSTDB_STORE_ROOT`数据存储目录，默认`/opt/mostdb`
- [x] `MOSTDB_SERVER_PORT`服务端口，默认`5379`
- [x] `MOSTDB_SERVER_URL`本机服务地址，默认“”
- [x] `MOSTDB_CLUSTER_URLS`集群节点服务地址，默认“”
