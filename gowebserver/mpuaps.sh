#!/bin/bash

AppDir=/msp/mpuaps
if [ -f "/opt/kdm/mpuaps" ];then
    echo "海豚容器中启动程序"
    AppDir=/opt/kdm
fi
AppStartTime=`date +'%Y.%m.%d-%H:%M:%S'`
AppStartLog=$AppDir/log/exception.log
AppStartLogBack=$AppDir/log/exception_back_$AppStartTime.log
AppConfigFile=$AppDir/config/config-prod.toml

AppDBNameDefault=mpuaps
AppDBUserDefault=root
AppDBPassDefault=msp123
AppDBPortDefault=3316

AppDBBackFile=/msp/db/mpuaps_back_$AppStartTime.log

#校验Exception日志文件大小
CheckExceptionLog() {
    maxLogLength=10485760 #10*1024*1024
    if [ -e "${AppStartLog}" ];then
        logFileLength=$(ls -l ${AppStartLog} | awk '{print $5}')
        echo "\"日志文件大小 = ${logFileLength}\""
        if [ "${logFileLength}" -gt "${maxLogLength}" ]; then
            mv $AppStartLog $AppStartLogBack
        fi
    fi
}

#数据库权限配置 $1=user; $2=pass; $3=port
MysqlConfig() {
    #开放本地程序使用【root】用户访问数据库
    mysql -u$1 -p$2 -P$3 -e "grant all privileges on *.* to 'root'@'%' identified by '$2';flush privileges;"

    #代理数据库【mpuaps】权限赋予给【msp】账户
    mysql -u$1 -p$2 -P$3 -e "grant all privileges on *.* to 'msp'@'%' identified by '$2';flush privileges;"
    #mysql -u$1 -p$2 -P$3 -e "grant all privileges on $AppDBNameDefault.* to 'msp'@'%' identified by '$2';flush privileges;"
}

#修改程序配置文件
ModifyAppConfig() {
    #校验数据库默认账户密码
    MysqlConfig $AppDBUserDefault $AppDBPassDefault $AppDBPortDefault
    if [ $? = 0 ];then
        echo "默认账户密码连接数据库ok"
    else
        if [ ! -f "/msp/cfg/mpu.cfg" ];then
            echo "数据库无法连接"
            return "0"
        fi

        #读/mpu/cfg/mpu.cfg显控配置信息
        MpuCfgData=$(cat /msp/cfg/mpu.cfg | sed 's/[[:space:]]//g' | sed 's/"//g' | sed -r 's/,//g' | egrep -v '^[{}]' | sed 's/:/=/1')
        declare $MpuCfgData
        #校验显控配置文件中数据库配置
        MysqlConfig $AppDBUserDefault $mysql_pw $mysql_port
        if [ $? = 0 ];then
            echo "显控账户密码连接数据库ok"

            #修改配置文件数据库端口
            sed -i 's/'"${AppDBPortDefault}"'/'"${mysql_port}"'/' $AppConfigFile
            #修改配置文件数据库密码
            sed -i 's/'"${AppDBPassDefault}"'/'"${mysql_pw}"'/' $AppConfigFile

            AppDBPassDefault=$mysql_pw
            AppDBPortDefault=$mysql_port
        else
            echo "数据库无法连接"
        fi
    fi
}

#校验数据库是否需要重置
CheckNeedRebuildDatabase() {
    #B5_20221108至B5_20221123版本mpuaps表初始化sql脚本存在Bug.部分表初始化失败,9180级联功能无法正常使用
    #现在已有外部项目部署了B5-20221118版本
    #兼容对应版本升级至新版本,重建mpuaps数据库
    SQL_EXISTS=$(printf 'SHOW TABLES LIKE "%s"' "mpu_access_platform")
    if [[ $(mysql -u$AppDBUserDefault -p$AppDBPassDefault -P$AppDBPortDefault -e "$SQL_EXISTS" $AppDBNameDefault) ]]
    then
        echo "check mpuaps database ok"
    else
        echo "" >> $AppStartLog
        echo "" >> $AppStartLog
        echo "版本数据库异常,删除数据库重建" >> $AppStartLog
        echo "数据库备份路径: $AppDBBackFile" >> $AppStartLog
        mysqldump -u$AppDBUserDefault -p$AppDBPassDefault --databases $AppDBNameDefault > $AppDBBackFile
        mysql -u$AppDBUserDefault -p$AppDBPassDefault -P$AppDBPortDefault -e "drop database if exists mpuaps;"
    fi
}

#校验是否升级至支持权限管理版本
CheckSupportPermission() {
    if [ `grep -c "mpuaps_sys_menu.sql" $AppConfigFile` -eq '0' ];then
        echo "升级至权限管理版本" >> $AppStartLog
        echo "数据库备份路径: $AppDBBackFile" >> $AppStartLog
        mysqldump -u$AppDBUserDefault -p$AppDBPassDefault --databases $AppDBNameDefault > $AppDBBackFile

        mysql -u$AppDBUserDefault -p$AppDBPassDefault -P$AppDBPortDefault -e "use mpuaps;drop table if exists mpu_channel_group_to_channel;"
        mysql -u$AppDBUserDefault -p$AppDBPassDefault -P$AppDBPortDefault -e "use mpuaps;drop table if exists mpu_channel_group;"
        mysql -u$AppDBUserDefault -p$AppDBPassDefault -P$AppDBPortDefault -e "use mpuaps;drop table if exists mpu_channel;"
        mysql -u$AppDBUserDefault -p$AppDBPassDefault -P$AppDBPortDefault -e "use mpuaps;drop table sys_role;drop table sys_role_menu;drop table sys_user_role;"

        #修改配置文件
        OldSqlScriptFile=".\/config\/sql\/mysql\/mpuaps_system.sql;.\/config\/sql\/mysql\/mpuaps_mpusrv.sql"
        NewSqlScriptFile=".\/config\/sql\/mysql\/mpuaps_system.sql;.\/config\/sql\/mysql\/mpuaps_mpusrv.sql;.\/config\/sql\/mysql\/mpuaps_sys_menu.sql"
        sed -i 's/'${OldSqlScriptFile}'/'${NewSqlScriptFile}'/' $AppConfigFile
    fi
}


case $1 in
"start")
    cd $AppDir
    mkdir -p $AppDir/log

    ModifyAppConfig
    CheckExceptionLog
    CheckNeedRebuildDatabase
    CheckSupportPermission
    echo "" >> $AppStartLog
    echo "*********************************************" >> $AppStartLog
    echo "* $AppStartTime mpuaps start"                  >> $AppStartLog
    echo "*********************************************" >> $AppStartLog
    chmod +x $AppDir/mpuaps
    $AppDir/mpuaps >> $AppStartLog 2>&1 &

    if [ `ps -aux | grep mpuaps | wc -l` -gt 1 ];then
        echo "mpuaps start ok"
    else
        echo "mpuaps start failed"
    fi
;;

"stop")
    ps -ef | grep mpuaps | grep -v grep | awk '{print $2}' | xargs kill -9
;;
esac;

exit 0