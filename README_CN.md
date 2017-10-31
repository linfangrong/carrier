# carrier

通过carrier + redis cluster 替换 twemproxy(nutcracker) + redis standalone方式。
[redis cluster](http://redis.io/topics/cluster-tutorial)集群的代理。


## 和twemproxy差别
+ twemproxy不支持动态扩容和缩容。redis cluster支持了动态扩容和缩容，通过carrier屏蔽。
+ twemproxy是单进程单线程。carrier则通过go的协程多并发处理，效率较高。
+ 支持读取命令从slave节点读取。
+ 支持的命令略多。

## 安装
+ 预先安装包管理glide
```bash
go get github.com/Masterminds/glide
go install github.com/Masterminds/glide
```
+ 安装carrier
```bash
make glide
make
```

## 支持命令如下(大部分是多key，还有管理命令)
+ String: 不支持BITOP、MSETNX。
+ Lists: 不支持BLPOP、BRPOP、BRPOPLPUSH、RPOPLPUSH。
+ Connection: 不支持AUTH、QUIT。(SELECT对于redis cluster无意义，故永远返回OK)
+ Server: 全部不支持。(包括BGREWRITEAOF、BGSAVE、CLIENT、COMMAND、CONFIG、DBSIZE、DEBUG、FLUSHALL、FLUSHDB、INFO、LASTSAVE、MONITOR、ROLE、SAVE、SHUTDOWN、SLAVEOF、SYNC、TIME)。
+ Cluster: 全部不支持。(没有必要)。
+ Keys: 不支持KEYS、MIGRATE、MOVE、OBJECT、RANDOMKEY、RENAME、RENAMENX、SCAN、WAIT。(EXISTS暂时只支持单key)。
+ Transactions: 全部不支持。(包括DISCARD、EXEC、MULTI、UNWATCH、WATCH)。
+ Scripting: 全部不支持。(包括EVAL、EVALSHA、SCRIPT)。
+ Geo: 支持。
+ Hashes: 支持。
+ HyperLogLog: 不支持PFMERGE。(PFCOUNT只支持单key返回)。
+ Pub/Sub: 由于集群模式pub/sub实现相当于广播，暂不考虑支持。(可用第三方消息队列，如beanstalk、rabbitmq)。
+ Sets: 不支持SDIFF、SDIFFSTORE、SINTER、SINTERSTORE、SMOVE、SUNION、SUNIONSTORE。
+ Sorted Sets: 不支持ZINTERSTORE、ZUNIONSTORE。


## 感谢
+ https://github.com/3xian
+ https://github.com/samuelduann
