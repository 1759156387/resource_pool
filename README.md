# resource_pool
一个golang的资源的管理池，可以作为数据库连接的管理，减少连接套接字的次数，增加系统的响应能力！
平时写程序也有这么一个困惑的地方，要是每次打开连接速度肯定是很慢的，要是共享一个连接，并发能力不是太好
而且要是这个连接关闭了，什么时候去重新打开这个时间点也不好定义！
在测试中发现和只有一个数据库连接性能差不多，但是好歹可以自动管理了，后期需要更多的测试！