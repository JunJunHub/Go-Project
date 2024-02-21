#!/bin/bash

#set -x

echo =====================================
echo =   mpuaps amd64 version build
echo =====================================

#声明GO环境变量
export GOPROXY=https://goproxy.io,direct
export GOSUMDB="sum.golang.google.cn"

#打印编译环境信息
echo "print compiling environment ..."
pwd && uname -a && go version && go env && gcc -v && g++ -v && export

SourceDir=`pwd`
WebResources=$SourceDir/public/web
PackageDir=/tmp/pkgdir/mpuaps
PackageName=$SourceDir/target/mpuaps.tar.gz
BinOutputDir=$SourceDir/bin/v1.0.0/linux_amd64/
BuildTime=`date +'%Y.%m.%d-%H:%M:%S'`
BuildGoVersion=`go version`
LDFlags=" \
	-X 'gowebserver/app/common/global.BinBuildTime=${BuildTime}' \
	-X 'gowebserver/app/common/global.BinBuildGoVersion=${BuildGoVersion}' \
"

#编译
echo ""
echo "compile start ..."
rm mpuaps
rm ${BinOutputDir}/mpuaps
#go mod tidy
go build -ldflags "$LDFlags"
cp mpuaps ${BinOutputDir}

if [ ! -f $SourceDir/mpuaps ];then
    echo "compile failed! mpuaps not exist!"
    exit 0
fi
echo "compile ok ..."
echo ""
exit 0

#海豚打包(依据此压缩包构造Docker镜像)
echo "build dolphin packet start ..."
rm -rf $PackageDir
mkdir -p $PackageDir
mkdir -p $PackageDir/log
mkdir -p $PackageDir/config
mkdir -p $PackageDir/config/sql/mysql
mkdir -p $PackageDir/public/web
mkdir -p $PackageDir/public/web/manage_client
mkdir -p $PackageDir/public/web/config_client

cp $SourceDir/mpuaps $PackageDir
cp $SourceDir/mpuaps.sh $PackageDir
chmod +x $PackageDir/mpuaps $PackageDir/mpuaps.sh

cp $SourceDir/config/config.toml $PackageDir/config
cp $SourceDir/config/config-prod.toml $PackageDir/config
cp $SourceDir/document/sql/mysql/* $PackageDir/config/sql/mysql
cp $SourceDir/i18n $PackageDir/ -r

tar -xvf $WebResources/config_client/config.tar -C $PackageDir/public/web/config_client
mv $PackageDir/public/web/config_client/config/static $PackageDir/public/web/config_client
rm -r $PackageDir/public/web/config_client/config
tar -xvf $WebResources/manage_client/client.tar -C $PackageDir/public/web/manage_client
mv $PackageDir/public/web/manage_client/client/static $PackageDir/public/web/manage_client
rm -r $PackageDir/public/web/manage_client/client

cd $PackageDir && tar -cvzf $PackageName * && cd -

exit 0
