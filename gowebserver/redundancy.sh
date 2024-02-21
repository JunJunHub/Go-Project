#!/bin/bash

#机箱冗余Mysql主从配置脚本

#参数
#$1 调用方式[initMaster | initSlave | stop | status | reset]
#$2 主机箱IP
#$3 备机箱IP
if [ ! -n "$1" ] || [ ! -n "$2" ]; then
    echo "call [./redundancy $1 $2 $3 $4] input param err"
    exit -1
fi


#时间
Time=`date +'%Y.%m.%d-%H:%M:%S'`

#脚本调用方式
CallMode=$1

#主机箱
MasterIP=$2
MasterDBUser=root
MasterDBPass=msp@123.com
MasterDBPort=3316

#备机箱
SlaveIP=$3
SlaveDBUser=root
SlaveDBPass=msp@123.com
SlaveDBPort=3316

#数据同步账户信息
DBSyncDataUser=repl
DBSyncDataPass=mpuredu123!

#脚本执行LOG存储路径
ReduLog=/var/log/mpu.printf
ReduLogBack=/var/log/mpu_$Time.printf


#校验ReduLog日志文件大小
CheckReduLogFileSize() {
    maxLogLength=10485760 #10*1024*1024
    if [ -e "${ReduLog}" ];then
        logFileLength=$(ls -l ${ReduLog} | awk '{print $5}')
        if [ "${logFileLength}" -gt "${maxLogLength}" ]; then
            mv $ReduLog $ReduLogBack
        fi
    fi
}

#调试日志
PrintLog() {
   CheckReduLogFileSize

   content="$(date '+%Y-%m-%d %H:%M:%S') $@"
   echo $content >> $ReduLog && echo -e "\033[32m"  ${content}  "\033[0m"
}

#校验主备数据库密码
CheckDBPass() {
    #校验主库密码
    if [ -n "$MasterIP" ]; then
        mysql -h$MasterIP -uroot -pmsp@123.com -P3316 -e "use mysql; select user,host,authentication_string from user;"
        if [ $? -ne 0 ]; then
        	mysql -h $MasterIP -uroot -pmsp123 -P3316 -e "use mysql; select user,host,authentication_string from user;"
        	if [ $? -ne 0 ]; then
                PrintLog "Unable to connect to master database[$MasterIP $MasterDBUser $MasterDBPass $MasterDBPort]"
                exit -1
        	fi
        	MasterDBPass=msp123
            echo "$MasterIP $MasterDBUser $MasterDBPort $MasterDBPass"
        fi
    fi


    #校验备库密码
    if [ -n "$SlaveIP" ]; then
      mysql -h$SlaveIP -uroot -pmsp@123.com -P3316 -e "use mysql; select user,host,authentication_string from user;"
        if [ $? -ne 0 ]; then
        	mysql -h$SlaveIP -uroot -pmsp123 -P3316 -e "use mysql; select user,host,authentication_string from user;"
        	if [ $? -ne 0 ]; then
                PrintLog "Unable to connect to slave database[$SlaveIP $SlaveDBUser $SlaveDBPass $SlaveDBPort]"
                exit -1
        	fi
        	SlaveDBPass=msp123
        	echo "$SlaveIP $SlaveDBUser $SlaveDBPass $SlaveDBPort"
        fi
    fi
}

#创建数据库同步账户
#$1=IP; $2=user; $3=pass; $4=port
MysqlCreateSyncDataUser() {
    #创建Mysql数据同步账户
    PrintLog "mysql -h$1 -u$2 -p$3 -P$4 -e \"create user '$DBSyncDataUser'@'%';\""
    mysql -h$1 -u$2 -p$3 -P$4 -e "create user '$DBSyncDataUser'@'%';"

    #更新Mysql数据同步账户权限
    PrintLog "mysql -h$1 -u$2 -p$3 -P$4 -e \"grant all privileges on *.* to '$DBSyncDataUser'@'%' identified by '$DBSyncDataPass';flush privileges;\""
    mysql -h$1 -u$2 -p$3 -P$4 -e "grant all privileges on *.* to '$DBSyncDataUser'@'%' identified by '$DBSyncDataPass';flush privileges;"
}

#停止主从同步
#$1=IP; $2=user; $3=pass; $4=port
MysqlStopSlave(){
   PrintLog "mysql -h$1 -u$2 -p$3 -P$4 -e \"stop slave;reset slave;reset master;flush logs;\""
   mysql -h$1 -u$2 -p$3 -P$4 -e "stop slave;reset slave;reset master;flush logs;"
}

#重新配置Mysql数据同步
#参数:$1=IP $2=DBUser; $3=DBPass; $4=DBPort $5=peerDBIp $6=peerDBPort
MysqlResetSlave(){
    #配置主从复制
    PrintLog "mysql -h$1 -u$2 -p$3 -P$4 -e \"change master to master_host='$5',master_port=$6,master_user='$DBSyncDataUser',master_password='$DBSyncDataPass',master_auto_position=1;start slave;\""
    mysql -h$1 -u$2 -p$3 -P$4 -e "change master to master_host='$5',master_port=$6,master_user='$DBSyncDataUser',master_password='$DBSyncDataPass',master_auto_position=1;start slave;"
}

#导入对端数据
#参数: $1=masterIP $2=masterDBPort $3=slaveIp $4=slaveDBPort
MysqlImportPeerData(){
    dataSyncTime=$Time
    dbBackDir=/msp/mpuaps/public/download/dbback
    mkdir -p $dbBackDir

    #判断数据备份目录已使用空间
    dbBackMaxSize=5396480 #5*1024*1024
    dbBackDirUsedSize=`ls -lt $dbBackDir | head -n 1 | awk '{print $2}'`
    while [ ${dbBackDirUsedSize} -gt ${dbBackMaxSize} ]; do
        #空间超过5G删除之前的备份数据
        rm $dbBackDir/`ls -lt $dbBackDir | tail -n 1 | awk '{print $9}'`
        dbBackDirUsedSize=`ls -lt $dbBackDir | head -n 1 | awk '{print $2}'`
    done

    #导出备机箱数据备份
    slaveDataBackFile=$dbBackDir/mpuaps_backup_$3_$dataSyncTime.sql
    PrintLog "mysqldump -h$3 -u$DBSyncDataUser -p$DBSyncDataPass -P$4 --databases mpuaps --set-gtid-purged=off > $slaveDataBackFile"
    mysqldump -h$3 -u$DBSyncDataUser -p$DBSyncDataPass -P$4 --databases mpuaps --set-gtid-purged=off > $slaveDataBackFile

    #导出主机箱mpuaps库数据
    masterDataFile=$dbBackDir/mpuaps_master_$1_$dataSyncTime.sql
    PrintLog "mysqldump -h$1 -u$DBSyncDataUser -p$DBSyncDataPass -P$2 --databases mpuaps --ignore-table=mpuaps.sys_user_online --ignore-table=mpuaps.sys_login_history --ignore-table=mpuaps.sys_oper_log --ignore-table=mpuaps.mpu_device_online_record --ignore-table=mpuaps.mpu_alarm_record --ignore-table=mpuaps.mpu_alarm_temporary_clean_record --ignore-table=mpuaps.mpuaps_access_platform --ignore-table=mpuaps.mpu_channel --set-gtid-purged=off > $masterDataFile"
    mysqldump -h$1 -u$DBSyncDataUser -p$DBSyncDataPass -P$2 --databases mpuaps --ignore-table=mpuaps.sys_user_online --ignore-table=mpuaps.sys_login_history --ignore-table=mpuaps.sys_oper_log --ignore-table=mpuaps.mpu_device_online_record --ignore-table=mpuaps.mpu_alarm_record --ignore-table=mpuaps.mpu_alarm_temporary_clean_record --ignore-table=mpuaps.mpuaps_access_platform --ignore-table=mpuaps.mpu_channel --set-gtid-purged=off > $masterDataFile

    #备机箱导入主机箱数据
    PrintLog "mysql -h$3 -u$DBSyncDataUser -p$DBSyncDataPass -P$4 --database mpuaps < $masterDataFile"
    mysql -h$3 -u$DBSyncDataUser -p$DBSyncDataPass -P$4 --database mpuaps < $masterDataFile
}

#校验Mysql主从同步状态
#$1=IP; $2=user; $3=pass; $4=port
CheckMysqlSalveStatus(){
    mysql -h $1 -u$2 -p$3 -P$4 -e "show slave status\G"
	if [ $? -ne 0 ]; then
	    PrintLog "cmd [mysql -h $1 -u$2 -p$3 -P$4 -e \"show slave status\G\"] exec failed!"
	fi

    mysql -h $1 -u$2 -p$3 -P$4 -e "show master status\G"
    if [ $? -ne 0 ]; then
        PrintLog "cmd [mysql -h $1 -u$2 -p$3 -P$4 -e \"show master status\G\"] exec failed!"
    fi
}


case $CallMode in
"initMaster" | "initSlave")
    ##########################
    #1、修改Mysql配置文件
    #2、重启Mysql
    #3、创建数据同步账户
    #########################
    PrintLog "call [./redundancy $1 $2 $3 $4]"

    #校验数据库密码
    CheckDBPass
    if [ $? -ne 0 ]; then
        exit -1
    fi

    #Mysql创建数据同步账户
    MysqlCreateSyncDataUser $MasterIP $MasterDBUser $MasterDBPass $MasterDBPort
    if [ $? -ne 0 ]; then
        PrintLog "master[$MasterIP] mysql create sync data user[repl] err"
        #exit -1    
    fi

    #Mysql数据库配置
    #备份默认数据库配置文件
    cp /etc/mysql/mysql.conf.d/mysqld.cnf /etc/mysql/mysql.conf.d/mysqld.cnf_bak
    #替换数据库配置文件
    if [ $CallMode == "initMaster" ]; then
        sed -i 's/.*server-id.*/server-id=1/' /msp/mpuaps/config/mysqld_slave_template.cnf
        sed -i 's/.*auto-increment-offset.*/auto-increment-offset=1/' /msp/mpuaps/config/mysqld_slave_template.cnf
    else
        sed -i 's/.*server-id.*/server-id=2/' /msp/mpuaps/config/mysqld_slave_template.cnf
        sed -i 's/.*auto-increment-offset.*/auto-increment-offset=2/' /msp/mpuaps/config/mysqld_slave_template.cnf
    fi
    cp /msp/mpuaps/config/mysqld_slave_template.cnf /etc/mysql/mysql.conf.d/mysqld.cnf

    #重启Mysql服务
    systemctl restart mysql.service
    if [ $? -ne 0 ]; then
        #服务启动失败,恢复Mysql配置文件
        PrintLog "master box mysql config err! please check '/msp/mpuaps/config/mysqld_slave_template.cnf'"
        cp /etc/mysql/mysql.conf.d/mysqld.cnf_bak /etc/mysql/mysql.conf.d/mysqld.cnf
        systemctl restart mysql.service
        exit -1
    fi

#    MysqlCreateSyncDataUser $SlaveIP $SlaveDBUser $SlaveDBPass $SlaveDBPort
#    if [ $? -ne 0 ]; then
#        PrintLog "slave[$SlaveIP] mysql create sync data user[repl] err"
#        #exit -1
#    fi

    PrintLog "master box mysql config ok!"
    exit 0
;;

"reset")
    #############################
    #1、停止主备Mysql数据同步
    #2、备机重新导入主机数据
    #3、重置 Mysql Gtid
    #4、开启主备Mysql数据同步
    #
    #注：谁的数据最新，谁作为主机箱
    #############################
    PrintLog "call [./redundancy $1 $2 $3 $4]"

    #校验数据库密码
    CheckDBPass
    if [ $? -ne 0 ]; then
        exit -1
    fi

    #停止Mysql数据同步
    MysqlStopSlave $MasterIP $MasterDBUser $MasterDBPass $MasterDBPort
    MysqlStopSlave $SlaveIP $SlaveDBUser $SlaveDBPass $SlaveDBPort

    #备机箱导入主机箱数据
    PrintLog "start slave[$SlaveIP] import master[$MasterIP] data"
    MysqlImportPeerData $MasterIP $MasterDBPort $SlaveIP $SlaveDBPort
    if [ $? -ne 0 ]; then
        PrintLog "slave[$SlaveIP] import master[$MasterIP] data err"
        #exit -1
    fi
    PrintLog "end slave[$SlaveIP] import master[$MasterIP] data"

    #停止Mysql数据同步(重置Mysql Gtid)
    MysqlStopSlave $MasterIP $MasterDBUser $MasterDBPass $MasterDBPort
    MysqlStopSlave $SlaveIP $SlaveDBUser $SlaveDBPass $SlaveDBPort

    #重新配置数据库主从复制
    MysqlResetSlave $SlaveIP $SlaveDBUser $SlaveDBPass $SlaveDBPort $MasterIP $MasterDBPort
    if [ $? -ne 0 ]; then
        PrintLog "mysql reset slave mysql status err"
        #exit -1
    fi
    MysqlResetSlave $MasterIP $MasterDBUser $MasterDBPass $MasterDBPort $SlaveIP $SlaveDBPort
    if [ $? -ne 0 ]; then
        PrintLog "mysql reset master mysql status err"
        #exit -1
    fi
    exit 0
;;

"status")
    #PrintLog "call [./redundancy $1 $2 $3 $4]"

    #校验数据库密码
    CheckDBPass
    if [ $? -ne 0 ]; then
        exit -1
    fi

    CheckMysqlSalveStatus $MasterIP $MasterDBUser $MasterDBPass $MasterDBPort
    exit 0
;;

"stop")
    #校验数据库密码
    CheckDBPass
    if [ $? -ne 0 ]; then
        exit -1
    fi

    MysqlStopSlave $MasterIP $MasterDBUser $MasterDBPass $MasterDBPort
    MysqlStopSlave $SlaveIP $SlaveDBUser $SlaveDBPass $SlaveDBPort
    exit 0
;;

esac;

PrintLog "call [./redundancy $1 $2 $3 $4] input param err"
exit -1
