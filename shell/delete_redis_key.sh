#!/bin/bash
# 删除 key "pontus_cache_segment_mtz1.mt-meipai.26631323d8dfc91533097087166643" 和
# "pontus_cache_stream_mtz1.mt-meipai.26631323d8dfc91533097087166643"
file=$1
cat $file | while read line
do
  echo “del $file”
  redis-cli -h m7702.bx.redis.m.com -p 7702 -a rAtPMZ*Ia374TF28 del pontus_cache_segment_${stream_id}
  redis-cli -h m7702.bx.redis.m.com -p 7702 -a rAtPMZ*Ia374TF28 del pontus_cache_stream_${stream_id}
  sleep 1
done
