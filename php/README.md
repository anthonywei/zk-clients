说明：
===========================
这是一个php调用配置中心的例子

需要：
===========================
在使用之前需要安装php的zookeeper扩展，请确认php -m查看是否已经安装可zk的模块，如果没有，可以找运维同学安装

介绍：
===========================
    mtconfigserver.php 库文件
    test.php 调用的例子

代码本身不需要做过多的修改了，做了进程内的缓存以及会实时更新zk server端的变化。

其他：
===========================
使用过程中有其他问题，可以和weishouyang@meituan.com联系