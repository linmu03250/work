#!/bin/bash
# 从云存储下载 streaam_id 对应的回放m3u8文件； 检查其中是否存在大于20 的 ts
# 直播 id 是

if [ $# -lt 1 ]
then
  echo "USAGE: $0 stream_id"
  echo "stream_id:  如果 直播id 是 mtz1.mt-meipai.220627c192ca4415383578325106，\
    对应stream_id 是220627c192ca4415383578325106"
  exit 1
fi
stream_id=$1
file=mtz1${stream_id}.m3u8
if [ -e $file ]
# 检查文件是否存在
then
  echo "$file file exist"
else
  url=http://fragments.zone1.meitudata.com/mtz1.mt-meipai.${stream_id}/$file
  echo "download $url"
  wget $url -O $file
fi

if [ ! -s $file ]
# 检查文件是否有内容
then
  echo "file not exist and download fail"
  rm $file
  exit 1
fi

ret=`awk -F: '$2 > 20 && /INF/ {print $2}' mtz1${stream_id}.m3u8`
echo $ret
if [ ! -z "$ret" ]
then
  # string $ret 非空
  echo "found ts EXTINF > 20"
else
  echo "EXTINF >20 not found"
fi
