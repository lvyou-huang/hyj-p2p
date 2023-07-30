#P2P
## 用户鉴权
- 使用公钥加密钥进行加密每个节点都有一个公钥和一个私钥。当节点A向节点B发送消息时，节点A使用节点B的公钥对消息进行加密，然后将加密后的消息发送给节点B。节点B使用自己的私钥对消息进行解密，以验证消息的真实性
## 用户属性
- 用户第一次链接的时候应该
## 客户端
- 正常情况下客户端需要一直运行listen.go监听端口
- 运行p2p.go向中心服务器询问在线客户进行udp链接
## 服务端
- 服务端维护mysql存储账密和hub存储在线客户，并向来询问的客户端给出所有在线用户
## 哈希
- 获得文件的哈希值，进行校验，保证文件传输过程的完整性
## processfile
- 处理文件，对文件进行分片和合并操作
## 不足之处
- 没有处理并发。