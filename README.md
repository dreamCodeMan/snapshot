# snapshot
支持对视频编码格式为H264的视频文件进行视频截帧，采用opencv实现快速截取；为https://help.aliyun.com/document_detail/64555.html的一种实现


http://localhost:8080/thumb?url=http%3a%2f%2fimg.ksbbs.com%2fasset%2fMon_1703%2f05cacb4e02f9d9e.mp4&w=848&h=480&t=1510


t	指定截图时间	[0,视频时长] 单位：ms

w	指定截图宽度，如果指定为0，则自动计算。	[0,视频宽度] 单位：像素（px）

h	指定截图高度，如果指定为0,则自动计算。如果w和h都为0，则输出为原视频宽高。	[0,视频高度] 单位：像素（px）

url 指定需要截取的视频的url或者服务所在机器的绝对路径的UrlEncode编码。
