说明：
===========================
go语言使用zookeeper的例子，为什么要再做一次封装呢。主要有以下几个原因

1.修改从http://config.sankuai.com/api/zkserverlist获取zk的host
2.增加一个本地的进程内的代码缓存，不要每次都去取zk的配置，提高性能，当server端更新后，再去server拉取一次，覆盖进程内的cache
3.对zk输出的结果进行额外处理，使得业务可以直接获取到某个配置的value

mtconfig文件夹
==========================
这个文件夹是实现的MtConfigServer.go的方法

zk文件夹
==========================
这个是zk原生的golang的library，但是由于原生的一些限制导致，每次获取的时候都是走的网络获取，性能不佳

test.go
=========================
这个是测试的代码，业务使用的时候，可以参考这个文件