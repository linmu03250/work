### linux 开启 tcp bbr ...

开启 TCP BBR
开机后 uname -r 看看是不是内核 ≥ 4.9

执行

echo "net.core.default_qdisc=fq" >> /etc/sysctl.conf
echo "net.ipv4.tcp_congestion_control=bbr" >> /etc/sysctl.conf
保存生效

sysctl -p
执行

sysctl net.ipv4.tcp_available_congestion_control
sysctl net.ipv4.tcp_congestion_control
如果结果都有bbr, 则证明你的内核已开启 TCP BBR！

执行

lsmod | grep bbr

看到有 tcp_bbr 模块即说明bbr已启动

修改 net.ipv4.tcp_congestion_control=cubic， sysctl -p 重新设置成cubic
