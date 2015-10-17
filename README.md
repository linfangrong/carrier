# carrier

通过carrier + redis cluster 替换 twemproxy(nutcracker) + redis standalone方式。
[redis cluster](http://redis.io/topics/cluster-tutorial)集群的代理。


## 和twemproxy差别
+ twemproxy不支持动态扩容和缩容。redis cluster支持了动态扩容和缩容，通过carrier屏蔽。
+ twemproxy是单进程单线程。carrier则通过go的协程多并发处理，效率较高。
+ 支持读取命令从slave节点读取。

## 安装
sh install.sh

## 感谢
+ https://github.com/3xian
+ https://github.com/samuelduann
