Description: 一个内网高性能针对性的爬虫(?)实现

flag:
	erp       erp账户
	pwd       erp密码
	file      读入文件
	n         请求goroutine数量,读写chan缓存size,字节缓存池大小2n
	proxy     是否开启代理
	proxyurl  代理地址

util.go: 一个利用chan实现的高效文本文件读写方法
	@param file 读入file文件(每行一个caseID)
	@param chr  通过for-range 从其中取出一个caseID,出错重试传入caseID
	@param chw  将处理结果传入chw,最后写入file.out.csv
	@return 返回已读入缓存数量与已写入缓存数量
	func FileDuplex(file string, chr, chw chan string) (chan uint64, chan uint64)
	
1.自动处理cookie,实现SSO登录,支持代理
2.利用chan实现非阻塞的bytes.Buffer池(思路来自effective go chan章节),高效利用字节缓存
3.QPS统计


Notice:公司内网-相关地址隐藏