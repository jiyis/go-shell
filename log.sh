#!/bin/bash

sourcepath="/opt/logs/all-logs"
exportpath="/jiyi/www/php/sys-common/storage/logexport"
downloadpath="/jiyi/www/php/sys-common/storage/logarchive"

# 删除十五天前的文件
find "$downloadpath" -type f -mtime +15 -exec rm {} \;

# 找出当天
find $sourcepath -mtime $1 -type f ! -name "audit.*" ! -path "*/*brm*/*" ! -path "*/*feesync*/*"  -name "messages" -o -name "app.log" -o -name "nginx_error.log" -o -name "miscroservice.log" | while read path
do
   # extract the file  name from the path
   fn="${path##*/}"
   # echo "Operation: $path -> $exportpath/$fn"
   # if the destination file already exists
   cp --parents "$path" "$exportpath"
done

current=`date "+%Y-%m-%d"`

tar -czPf "$downloadpath"/all-logs-"$current".tar.gz "$exportpath""$sourcepath"

rm -rf "$exportpath"/*

echo "all-logs-"$current".tar.gz"