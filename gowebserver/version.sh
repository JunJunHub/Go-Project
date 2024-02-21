#!/bin/bash  

########################################
#版本号: v1.0.0
#1、主版本号
#   大版本迭代时修改     
#2、次版本号
#   功能增删改时更新
#3、编译版本
#   根据提交记录生成的编译版本ID
########################################

#声明主次版本号



#根据git提交记录生成编译版本ID
git rev-list HEAD | sort > config.git-hash  
LocalVersion=`wc -l config.git-hash | awk '{print $1}'`
if [ $LocalVersion \> 1 ] ; then
    VER=`git rev-list origin/master | sort | join config.git-hash - | wc -l | awk '{print $1}'`  
    if [ $VER != $LocalVersion ] ; then
        VER="$VER+$(($LocalVersion-$VER))"
    fi  
    if git status | grep -q "modified:" ; then  
        VER="${VER}M"  
    fi  
    VER="$VER $(git rev-list HEAD -n 1 | cut -c 1-7)"  
    GIT_VERSION=r$VER  
else  
    GIT_VERSION=  
    VER="x"  
fi  
rm -f config.git-hash  

echo "$GIT_VERSION $VER"

cat version.go.template | sed "s/\$FULL_VERSION/$GIT_VERSION/g" > app/common/global/version.go  
   
echo "update version.go" 
