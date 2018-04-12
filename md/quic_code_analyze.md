[1.chrominum源码的顶级目录src下的子目录说明](#1)  
[2.quic的主要数据机构说明](#jump). 

[XXXX](#jump)
<h3 id="1">1.chrominum源码的顶级目录src下的子目录说明</h3> 

+ src  
	+ apps : chromium中用来支持网页应用chrome package apps的代码  
	+ base : 被所有项目共享的代码。例如字符串操作等，这个文件夹下只能添加那些必须被多个顶层项目共享的代码。  
	+ breakpad : 谷歌开源的崩溃上报项目。它直接从谷歌的 svn 仓库中拉取。  
	+ build : 编译相关的配置、安装编译依赖的脚本等，由所有项目共享。  
	+ cc : chromium compositor(合成器) 实现。  
	+ chrome : Chromium browser。  
	+ chrome/test/data : 测试用数据。  
	+ content : 包含建立 多进程浏览器 所需要的核心代码。  
	+ device : 对底层硬件接口进行抽象，使其可以跨平台调用。  
	+ net : 网络库。  
	+ sandbox : 沙箱项目。其目的是防止 renderer 进程被 hack 之后攻击系统。  
	+ skia : 谷歌的 skia 图形库。这块项目是给 Android 使用的。ui/gfx 中对它进行了封装。  
	+ sql : 对 sqlite 的封装。  
	+ testing : 谷歌开源的测试工具 GTest。用于进行单元测试。  
	+ third_party : 一系列第三方库，例如 图片解码，解压缩库。chrome/third_party 里包含了一些专门给 chrome 用的第三方库。  
	+ ui/gix : 共享绘图类，基于 Chromium 的 UI 绘图库。   
	+ ui/views : 进行 ui 开发的简单框架，提供了渲染、布局、事件处理机制。大部分的浏览器 ui 都基于这个框架来实现。这个目录下包含了基本对象，另外一些跟 浏览器 相关的对象包含在 chrome/  
	+ browser/ui/views 目录中。   
	+ url : 谷歌开源的 url 解析和标准化库。   
	+ v8 : V8 JavaScript 库。它直接从谷歌的 svn 仓库中拉取   
	+ buildtools: 不同平台的编译工具   
	+ crypto：各种加解密、哈希等实现   
	+ quic依赖及源码结构说明   
+ 其中编译quic需要涉及到的模块有： 
	+ base build   
	+ crypto   
	+ net   
	+ testing  
	+ url  
	+ third_party 中的部分模块:
		+ Brotli:   google 开源的压缩算法，是一种全新的数据格式，提供比zopfli高20%的压缩比      
		+ Icu:      International Component for Unicode, 详见http://site.icu-project.org/ 。这里依赖icui18n，是一个国际标准化库  
		+ modp_b64:   chrominum中实现的base64编解码   
		+ boringssl: OpenSSL的高危漏洞Heartbleed曝光之后,Google创建的OpenSSL分支boringssl 。 网上有一些 基于chrominum的项目，有将boringssl使用openssl替换的例子  
		+ Zlib：压缩算法 http://zlib.net/  
		+ Ced : Compact Encoding Detection，是一个google开源的c++库。 It scans given raw bytes and detect the most likely text encoding.  
		+ Protobuf: google高效的结构化数据存储格式

+ 另外src/base/third_party 还有几个需要的第三方库： 
 	+ symbolize: google-glog's symbolization library  
 	+ dynamic_annotations: http://code.google.com/p/data-race-test/wiki/DynamicAnnotations  
 	+ xdg_mime: 解析MIME spec的模块  
 	+ libevent  
 	+ xdg_user_dirs: http://www.freedesktop.org/wiki/Software/xdg-user-dir
        

<h3 id="abc">2.quic的主要数据机构说明</h3>

<span id="jump">Hello World</span>

+ quic的核心代码位于src/net/quic/core 目录下，包含如下几个子目录:  
	+ crypto： 一些加解密的实现、加密握手的实现等  
	+ congesion_control: 包含拥塞控制相关的模块   
	+ frames： quic所有涉及到的frame类型、结构的定义，在这个目录下  
	+ proto：protobuf定义

+ src/net/tools/quic 下包含了1个实现上层quic client、server的例子程序所需要的实现文件

+ QuicServerId: 用于标识一个session , 包括ip, port , mode

~~~

net::QuicServerId server_id(server_host, server_port(), net::PRIVACY_MODE_DISABLED);
QuicClient::QuicClient(IPEndPoint server_address,
const QuicServerId& server_id,
const QuicVersionVector& supported_versions,
                       EpollServer* epoll_server)
    : server_address_(server_address),
      server_id_(server_id),
      epoll_server_(epoll_server),
      ...
      supported_versions_(supported_versions) {
} 
~~~

+ QuicConnection  
 	+ quic的connection由一个64位的connection id 标识，为了保证在服务端唯一，有2种方案：
 	 	+ 1是client 使用uuid 作为connect id； 
 	 	+ 1是client使用服务端指定的connect id
ProcessUdpPacket 用于处理接收到的raw packet，并将packet传给frame层  
+ QuicSession  
	+ client、server之间要通信，需要首先创建一个session，产生会话秘钥进行通信；session和connection是一一对应的，session id也就是connect id。某种意义上，connection是物理上的意义，session是逻辑上的表达。quic支持多路复用，每个session管理者多个不同状态的stream。
quicSession和QuicConnection 一一对应，创建session，需要先创建一个connection，作为session构造函数的参数。

~~~
QuicServerSessionBase* QuicSimpleDispatcher::CreateQuicSession(
            QuicConnectionId connection_id,
            const QuicSocketAddress& client_address,
            QuicStringPiece /*alpn*/) {
            // The QuicServerSessionBase takes ownership of |connection| below.
           QuicConnection* connection = new QuicConnection(
           connection_id, client_address, helper(), alarm_factory(),
           CreatePerConnectionWriter(),
            true, Perspective::IS_SERVER, GetSupportedVersions());
          QuicServerSessionBase* session = new QuicSimpleServerSession(
               config(), connection, this, session_helper(), crypto_config(),
               compressed_certs_cache(), response_cache_);
               session->Initialize();
       return session;
}
~~~  
	+  session包含多个流，根据流处于的不同状态，有如下分类：
static_stream_map_： id 为1的crypto stream、id为3 的header stream
dynamic_stream_map_：包括进入对端创建和本端创建2种，incoming_streams_和outgoing_streams_
draining_streams_：已经收到或发出fin，但是recv还有未被上层消费的，归为这一类
zombie_streams_：stream 已经被关闭了，但是发送的数据还有未收到ack的，归为这一类
flow_controller_: 用来进行session级别的流控 
QuicStream
quic支持多路复用，就是通过每个session包含多个stream实现的。在发送和接收数据之前，stream需要被创建，区分为进入和出去的流，需要继承CreateIncomingDynamicStream（）和CreateOutgoingDynamicStream（）函数来实现流的创建

~~~
QuicServerSessionWraper::CreateIncomingDynamicStream(QuicStreamId id) 
{
       QuartcStream* stream =new QuartcStream(id, this);
       stream->SetDelegate(visitor_);
       ActivateStream(QuicWrapUnique(stream));
       return stream;
}
QuicStream* QuicServerSessionWraper::CreateOutgoingDynamicStream() 
{
      QuartcStream* stream = new QuartcStream(
     GetNextOutgoingStreamId(), this);
     stream->SetDelegate(visitor_);
     ActivateStream(QuicWrapUnique(stream));
     return stream;
}
QuicPacketWriter
QuicPacketWriter是个纯虚类， 上层要发送数据需要继承QuicPacketWriter，实现它的纯虚函数WritePacket（）来发送数据。client和server发送数据都需要创建QuicPacketWriter 类。
WriteResult QuicDefaultPacketWriter::WritePacket(
    const char* buffer,
    size_t buf_len,
    const IPAddressNumber& self_address,
    const IPEndPoint& peer_address) {
  WriteResult result = QuicSocketUtils::WritePacket(
      fd_, buffer, buf_len, self_address, peer_address);
  if (result.status == WRITE_STATUS_BLOCKED) {
    write_blocked_ = true;
  }
  return result;
}
~~~

+ QuicStreamSequencerBuffer  
   + 这是接收缓冲区类，每个stream都对应一个。它是一个vector，item是内存指针。支持随机写和顺序读。新的数据到来是，数据和其在stream中的偏移会保存下来，连续的内存会给上层消费。  
   
~~~
QuicStreamSequencerBuffer有2个重要的函数，供QuicStreamSequencer调用：
  QuicErrorCode Readv(const struct iovec* dest_iov,
                      size_t dest_count,
                      size_t* bytes_read,
                      std::string* error_details);
  int GetReadableRegions(struct iovec* iov, int iov_len) const;
~~~
+ QuicStreamSendBuffer 实现发送缓存的类
+ QuicStreamSequencer
             这是接收方保证通信数据顺序的被上层接收的一个类，包含一个  
QuicStreamSequencerBuffer 成员buffered_frames_。有新接收的数据会触发OnStreamFrame。  
+ QuicConnectionStats
连接统计类，主要是连接的一些统计数据，接收、发送字节等  
+ QuicPacketCreator
主要用于创建packet，将1个或多个frame序列化成packet。
它的一个重要的数据成员是 QuicFrames queued_frames_; 上层需要发送到对方的消息，会先以frame的形式，append到frame 队列queued_frames_。
当上层调用到它的Flush()时，其SerializePacket（）成员函数负责将queued_frames_里面的frames序列化成1个个的packet，在发送出去  
+ QuicPacketGenerator  
+ QuicSentPacketManager
发送packet的管理类，里面会记录一下发包的统计信息  
+ QuicReceivedPacketManager
接收packet的管理类，里面会记录一些接收包的统计信息  
+ QuicCryptoClientHandshaker
负责加密握手的类
+ QuicFramer
负责帧处理的类。在quic中，帧是最小的数据单元，多个帧组合成一个packet，发送到接收方，接收方先将packet解成1个个的frame，再根据不同的frame类型，进行不同的处理
+ QuicSpdySession  
+ QuicSpdyStream
