
file=$1
dir=$2
cat $file | while read line
logf=${dir}/log.txt
mkdir -p $dir
do
    url="http://fragments.zone1.meitudata.com/mtz1.mt-meipai.${line}/mtz1${line}.m3u8"
    wget $url -O ${dir}/$line.m3u8 -o 1.log
    cat 1.log >> ${logf}
done
echo `grep "200 OK" $logf  |wc -l`
echo `grep "404 Not Found" $logf  |wc -l`
