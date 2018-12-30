#!/bin/bash
# 上传m3u8到云存储，覆盖旧的m3u8, 并刷新cdn
if [ $# -lt 1 ]
then
  echo "USAGE: $0 stream_id"
  echo "stream_id:  如果 直播id 是 mtz1.mt-meipai.220627c192ca4415383578325106，\
    对应stream_id 是220627c192ca4415383578325106"
  exit 1
fi
stream_id=$1
local_file="mtz1${stream_id}.m3u8"
remote_file="mtz1.mt-meipai.${stream_id}/mtz1${stream_id}.m3u8"
./mtshell_linux delete fragments ${remote_file}  http://up.bx.m.com
./mtshell_linux rput fragments ${remote_file} ${local_file} true http://up.bx.m.com
curl "http://cdn.ops.m.com/?target=http://fragments.zone1.meitudata.com/${remote_file}?start=-1&end=-1"
curl "http://cdn.ops.m.com/?target=http://fragments.zone1.meitudata.com/${remote_file}"
