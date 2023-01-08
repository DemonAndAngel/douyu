#!/bin/bash
ZIP_PATH=/app/douyu/go/zips
RUN_PATH=/app/douyu/go
if [ "$1" = "" ]; then
  echo "请输入正确的压缩文件名"
  exit
fi
file=$ZIP_PATH/$1
if [ ! -f $file ];then
  # 不存在项目 拉取 创建
  echo "找不到上传的文件,请检查"
  exit
fi
# 开始解压
echo "解压缩开始"
sudo unzip -o $file -d $RUN_PATH/
# 删除文件
echo "删除压缩文件"
rm -rf $file
echo "权限修改"
chown -R app:app $RUN_PATH/bin
chmod +x $RUN_PATH/bin/run
# 重新运行程序
supervisorctl restart qymall