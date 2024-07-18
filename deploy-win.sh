#!/bin/bash
# 自动化部署脚本

SERVER="root@8.134.126.6"
SERVER_PWD="/app/douyu/go"
# 拉取代码
echo "---拉取最新代码---"
git pull
# 编译程序
echo "---开始编译---"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/run ./
echo "---开始压缩---"
time=$(date "+%Y%m%d_%H%M%S")
# 压缩
file=build-$time.zip
zip -q -r $file ./bin/run
rm -rf ./bin/run
echo "---开始上传---"
scp $file $SERVER:$SERVER_PWD/zips/
echo "---运行服务端脚本---"
ssh $SERVER "$SERVER_PWD/deploy.sh $file"
rm -rf $file

