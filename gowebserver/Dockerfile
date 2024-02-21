##########################################################################################
# 镜像构建命令: docker build -t kedacom/mpuaps:v1.0 . || gf docker -t kedacom/mpuaps:v1.0
##########################################################################################

# 指定基础镜像
#FROM golang:latest
#FROM alpine:latest
FROM frolvlad/alpine-glibc

#时区
ENV TZ Asia/Shanghai
RUN apk add --update --no-cache \
	tzdata && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 维护者信息,镜像信息
LABEL maintainer="liyongjun@kedacom.com" version="1.0" description="MPU Application platform services"


#########################################################################################
#                                    INSTALLATION
#########################################################################################
# 设置固定的项目路径
ENV APPDIR /opt/kedacom/mpuaps
ENV APPVER v1.0.0

# 添加应用可执行文件，并设置执行权限
ADD ./bin/$APPVER/linux_amd64/mpuaps $APPDIR/mpuaps
ADD mpuaps.sh $APPDIR/mpuaps.sh
RUN chmod +x $APPDIR/mpuaps*

# 添加静态文件、配置文件 ...
ADD config   $APPDIR/config
ADD i18n     $APPDIR/i18n
ADD public   $APPDIR/public


########################################################################################
#                                       START
########################################################################################
WORKDIR $APPDIR
CMD $APPDIR/mpuaps.sh