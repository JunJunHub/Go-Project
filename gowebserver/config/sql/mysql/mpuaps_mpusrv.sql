set names utf8mb4;
set foreign_key_checks = 1;

-- drop database if exists mpuaps;
-- create database mpuaps;
-- use mpuaps;


-- ------------------------------------------------
-- 显控各型号槽位信息,物理槽位号与逻辑槽位号的映射关系
-- ------------------------------------------------
drop table if exists `mpu_slot_info`;
create table if not exists `mpu_slot_info`
(
    `MPU_box_type` varchar(64) not null comment '机箱型号',
    `physical_slot_no` int not null comment '物理槽位号',
    `logic_slot_no` varchar(64) not null comment '逻辑槽位号'
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='显控机箱槽位信息';
-- 初始化各型号物理槽位号与逻辑槽位映射关系
BEGIN;
-- N7
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 3,  'IPU-1');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 4,  'IPU-2');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 5,  'IPU-3');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 6,  'IPU-4');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 7,  'IPU-5');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 8,  'IPU-6');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 9,  'IPU-7');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 10, 'IPU-8');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 11, 'IPU-9');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 15, 'IPU-10');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 16, 'IPU-11');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 17, 'IPU-12');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 18, 'IPU-13');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 19, 'IPU-14');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 20, 'IPU-15');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 21, 'IPU-16');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 22, 'IPU-17');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 23, 'IPU-18');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 27, 'OPU-1');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 28, 'OPU-2');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 29, 'OPU-3');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 30, 'OPU-4');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 31, 'OPU-5');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 32, 'OPU-6');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 33, 'OPU-7');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 34, 'OPU-8');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 35, 'OPU-9');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 38, 'OPU-10');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 39, 'OPU-11');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 40, 'OPU-12');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 41, 'OPU-13');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 42, 'OPU-14');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 43, 'OPU-15');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 44, 'OPU-16');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 45, 'OPU-17');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS1', 46, 'OPU-18');

-- N5
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS2', 2,  'IPU-1');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS2', 3,  'IPU-2');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS2', 4,  'IPU-3');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS2', 5,  'IPU-4');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS2', 6,  'IPU-5');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS2', 9,  'IPU-6');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS2', 10, 'IPU-7');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS2', 11, 'IPU-8');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS2', 12, 'IPU-9');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS2', 13, 'IPU-10');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS2', 16, 'OPU-1');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS2', 17, 'OPU-2');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS2', 18, 'OPU-3');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS2', 19, 'OPU-4');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS2', 20, 'OPU-5');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS2', 22, 'OPU-6');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS2', 23, 'OPU-7');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS2', 24, 'OPU-8');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS2', 25, 'OPU-9');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS2', 26, 'OPU-10');

-- N3\N3-L(经济型)
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS3', 2,  'IPU-1');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS3', 3,  'IPU-2');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS3', 5,  'IPU-3');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS3', 6,  'IPU-4');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS3', 7,  'IPU-5');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS3', 9,  'OPU-1');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS3', 10, 'OPU-2');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS3', 11, 'OPU-3');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS3', 12, 'OPU-4');
INSERT INTO `mpu_slot_info` VALUES ('MSP1000-MS3', 13, 'OPU-5');

COMMIT;


-- -------------------------------------------------
-- 代理级联配置(各级代理配置信息)
-- -------------------------------------------------
create table if not exists `mpu_cascade_cfg`
(
    `cascade_id`                    varchar(64) unique not null comment '平台id',
    `parent_id`                     varchar(64) default '' comment '上级平台id',
    `cascade_alias`                 varchar(256) comment '级联描述',
    `cascade_type`                  smallint default '0' comment '级联方式：0=网络连接 1=串口连接',
    `cascade_plat_sn`               varchar(64) comment '平台唯一标识',
    `cascade_plat_regionalism_code` varchar(12) default '320501000000' comment '平台部署区划行政代码,参照《中华人名共和国行政区划代码》',
    `cascade_plat_business_code`    varchar(2)  default '98' comment '平台使用行业代码,默认: 98',
    `cascade_plat_dev_type`         varchar(3)  default '531' comment '平台设备类型,默认: 531',
    `cascade_plat_network`          varchar(1)  default '0' comment '平台部署网络环境: 0、1、2、3、4监控报警专网 5公安信息网 6政务网 7internet网 8社会资源接入网 9预留',
    `cascade_plat_alias`            varchar(256) comment '平台描述',
    `cascade_plat_ip`               varchar(50)  comment '平台ip地址',
    `cascade_plat_port`             int comment '平台端口',
    `cascade_plat_username`         varchar(30) comment '平台用户名',
    `cascade_plat_password`         varchar(50) comment '平台用户密码',
    `cascade_plat_url`              varchar(256) default '' comment '级联连接url',
    `cascade_plat_state`            smallint comment '平台状态: 0不在线, 1在线不可用, 2在线可用',
    `cascade_connect_state`         smallint default '0' comment '应用代理级联连接状态: 0连接失败 1连接成功 2连接中(资源同步中)',
    `enable`                        bool default true comment '级联配置是否启用',
    `is_current_platform_info`      bool default false comment '标识是否是代理自身信息',
    `create_by`                     varchar(64) default '' comment '创建者',
    `create_time`                   datetime default null comment '创建时间',
    `update_by`                     varchar(64) default '' comment '更新者',
    `update_time`                   datetime default null comment '更新时间',
    `deleted_at`                    datetime default null comment '删除时间',
    primary key (`cascade_id`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='级联配置信息';


-- -----------------------------------------------------
-- 级联通道配置(通过物理连线的方式,让上下级显控通道资源级联)
-- -----------------------------------------------------
create table if not exists `mpu_cascade_channel_cfg`
(
    `cfg_id` bigint not null comment '配置id',
    `cascade_id` varchar(64) not null comment '级联平台id,表示配置是哪一级显控上配置的',
    `superior_box_id` tinyint comment '上级平台机箱号(配置)',
    `superior_logic_slot_id` tinyint comment '上级单板逻辑槽位号(配置)',
    `superior_physical_slot_id` tinyint comment '上级单板对应物理槽位号',
    `superior_port` tinyint comment '上级单板端口号(配置)',
    `superior_channel_type` tinyint not null comment '上级通道类型: 0未知 1音频输入 2视频输入 3音频输出 4视频输出 5编码通道',
    `superior_channel_id` varchar(64) not null comment '上级通道id',
    `lower_box_id` tinyint comment '下级平台机箱号(配置)',
    `lower_logic_slot_id` tinyint comment '下级单板逻辑槽位号(配置)',
    `lower_physical_slot_id` tinyint comment '下级单板对应物理槽位号',
    `lower_port` tinyint comment '下级单板端口号(配置)',
    `lower_plat_ip` varchar(50) comment '下级平台ip地址',
    `lower_channel_id` varchar(64) not null comment '下级通道id',
    primary key (`cfg_id`),
    foreign key(cascade_id) references mpu_cascade_cfg(cascade_id) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='显控级联通道配置';


-- ---------------------------------------------------------------------------------------------------------
-- 接入平台信息: 显控、统一设备、流媒体、会议、监控;
-- 目前第三方平台都是接到显控平台上(网络信号).
-- ----------------------------------------------------------------------------------------------------------
-- drop table `mpu_access_platform`;
create table if not exists `mpu_access_platform`
(
    `access_plat_id` varchar(64) unique not null comment '接入平台id,本服务统一管理的接入平台id',
    `access_plat_id_real` bigint not null comment '从显控平台查到的id,显控平台使用管理. 0无效',
    `cascade_id` varchar(64) not null comment '级联id',
    `access_plat_regionalism_code` varchar(12) comment '平台所属行政区划代码,参照《中华人名共和国行政区划代码》',
    `access_plat_sn` varchar(64) comment '平台唯一标识',
    `access_plat_box_id` int comment '平台所属显控机箱ID',
    `access_plat_alias` varchar(256) comment '平台描述',
    `access_plat_state` smallint comment '平台状态: 0不在线, 1在线不可用, 2在线可用',
    `access_plat_ip` varchar(50) comment '平台ip地址',
    `access_plat_port` int comment '平台端口',
    `access_plat_username` varchar(30) comment '平台用户名',
    `access_plat_password` varchar(50) comment '平台用户密码',
    `access_plat_type` smallint comment '平台类型: 0显控 1老流媒体 2流媒体 3统一设备 4会议 5监控',
    `access_plat_keep_cnt` int comment '',
    `access_plat_handle` int comment '',
    `connect_state` smallint default '0' comment '连接状态: 0连接失败 1在线不可用 2在线可用',
    `enable` bool default true comment '接入配置是否生效',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`access_plat_id`),
    foreign key(cascade_id) references mpu_cascade_cfg(cascade_id) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='接入平台信息';


-- ----------------------------------------------------------------------
-- 使用的物理设备资源基本信息, 包含显控平台设备和通过第三方平台接入使用的设备信息(未使用)
-- ----------------------------------------------------------------------
create table if not exists `mpu_dev`
(
    `dev_id` varchar(64) unique not null comment '设备id,本服务统一管理的设备id',
    `dev_id_real` varchar(64) not null comment '在原平台上的id',
    `dev_sn` varchar(64) comment '设备唯一标识',
    `dev_gb_id` varchar(64) comment '设备国标id',
    `dev_alias` varchar(256) comment '设备别名',
    `dev_class` tinyint comment '物理设备类型: 0显控机箱 1输入板 2输出板 3接收器 4发送器 5解码板 6网管处理器 7视频处理器 8音频处理器 9dante板 10球型摄像机 11枪型摄像机 12半球型摄像机',
    `dev_model` varchar(64) comment '设备型号',
    `dev_source` tinyint not null comment '设备来源: 0虚拟 1显控平台 2监控平台 3会议平台 4统一设备',
    `dev_ip` varchar(50) comment '设备ip地址',
    `dev_state` tinyint comment '设备状态',
    primary key (`dev_id`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='物理设备资源信息';


-- ----------------------------------------------
-- 显控设备状态检测配置表 mpu_dev_keepalive_cfg
-- ----------------------------------------------
create table if not exists `mpu_dev_keepalive_cfg`
(
    `cfg_id` bigint unique not null auto_increment comment '配置id',
    `box_id` int not null comment '设备所在机箱号',
    `slot_id` int not null comment '设备槽位号',
    `ip` varchar(50) not null comment '设备ip地址',
    `reboot_time_consuming` int comment '设备重启耗费时间(s)',
    `check_time_interval` int comment '检测时间间隔(s)',
    `check_err_times` int comment '检测失败次数. 连续探测失败大于等于该值,认为状态异常',
    `enable` bool default true comment '配置是否生效',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`cfg_id`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='显控设备板状态检测配置表';


-- ---------------------------------------------------
-- 通道资源: 显控通道资源和通过第三方平台接入的通道资源
-- 设备<->通道: 一对多的关系
-- ---------------------------------------------------
create table if not exists `mpu_channel`
(
    `channel_id` varchar(64) unique not null comment '通道id,本服务统一生成管理的通道资源id',
    `cascade_id` varchar(64) not null comment '级联id,标识通道属于哪一级显控代理平台',
    `channel_id_real` varchar(64) comment '在显控平台上的通道id,显控分配的通道ID',
    `channel_alias` varchar(256)  comment '通道别名',
    `channel_source` tinyint not null default '0' comment '通道来源: 0未知 1显控平台资源(模拟信号) 2Umt平台(网络信号)',
    `channel_type` tinyint not null comment '通道类型: 0无效 1音频输入 2视频输入 4音频输出 8视频输出 16编码通道 32虚拟通道',
    `channel_physical_interfaces` tinyint comment '通道物理接口类型: 1=HDMI 2=VGA 3=YPbPr 4=SDI 5=C 6=S 7=DIV 8=DP 9=Camera 10=HDBaseT 11=SerDes 12=Fiber 13=Dante',
    `channel_state` bool comment '通道在线状态: 0离线(网络点位离线/模拟信号无信号) 1在线(网络点位在线/模拟信号有信号) 2已删除(板子拔走/网络点位删除) 3已收回(级联资源收回) 4不可用(级联未连接)',
    `channel_connect_cables` bool comment '通道是否连线(物理上)',
    `channel_inuse` bool comment  '通道是否被使用',
    `channel_pixels_width` smallint comment '视屏输出通道支持最大分辨率',
    `channel_pixels_height` smallint comment '视屏输出通道支持最大分辨率',
    `channel_fps` smallint comment '视频输出通道支持最高刷新率',
    `channel_dev_id` varchar(64) comment '通道设备id,本代理平台统一分配管理的设备ID,关联设备信息表,暂未使用',
    `channel_dev_type` varchar(256) comment '通道设备类型: IPC/MT/DECODER/VMS/MSP',
    `channel_dev_ip` varchar(50) comment '通道设备ip地址',
    `channel_dev_name` varchar(256) comment '通道设备名称',
    `channel_dev_manufacturer` varchar(256) comment '通道设备生产厂商',
    `channel_umt_id` int comment '网络点位对应UMT平台ID,显控平台接入Umt平台对应的UmtID',
    `channel_umt_dev_id` varchar(256) comment '网络点位对应设备ID,umt平台分配',
    `channel_umt_gb_id` varchar(64) comment '国标id,统一设备分配管理',
    `channel_umt_org_id` varchar(4096) comment '网络点位所属Umt组织节点ID(网络信号分组)',
    `channel_stream_id` varchar(256) comment '网络点位拉流ID, 监控点位为国标ID, 会议点位为会议C节点ID',
    `channel_umt_dev_net_id` int comment '网络点位穿网ID',
    primary key (`channel_id`),
    foreign key(cascade_id) references mpu_cascade_cfg(cascade_id) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='通道资源信息';
alter table `mpu_channel` modify column `channel_physical_interfaces` tinyint comment '通道物理接口类型: 1=HDMI 2=VGA 3=YPbPr 4=SDI 5=C 6=S 7=DIV 8=DP 9=Camera 10=HDBaseT 11=SerDes 12=Fiber 13=Dante';
alter table `mpu_channel` ADD INDEX channel_id_real_index (`channel_id_real`);
alter table `mpu_channel` ADD INDEX channel_source_index (`channel_source`);
alter table `mpu_channel` ADD INDEX channel_type_index (`channel_type`);
alter table `mpu_channel` ADD INDEX channel_umt_gb_id_index (`channel_umt_gb_id`);
-- alter table `mpu_channel` ADD FULLTEXT INDEX channel_umt_org_id_fulltext_index (`channel_umt_org_id`);
-- alter table `mpu_channel` ADD FULLTEXT INDEX channel_alias_fulltext_index (`channel_alias`);
alter table `mpu_channel` add column `deleted_at` datetime default null comment '删除时间';

-- 先删除未使用的索引
alter table `mpu_channel` drop index channel_umt_org_id_fulltext_index;
alter table `mpu_channel` drop index channel_alias_fulltext_index;


-- ------------------------------------------------------
-- 点位关联配置(记录显控物理点位录入统一设备平台时对应的网络信号ID)
-- ------------------------------------------------------
-- drop table `mpu_channel_association_map`;
create table if not exists `mpu_channel_association_map`
(
    `cfg_id` bigint not null auto_increment comment '配置ID',
    `analog_channel_id`  varchar(64) not null comment '模拟信号通道ID',
    `analog_channel_alias` varchar(256) comment '模拟信号通道名称',
    `analog_channel_state` int comment '模拟信号状态: 0离线(网络点位离线/模拟信号无信号) 1在线(网络点位在线/模拟信号有信号) 2已删除(板子拔走/网络点位删除) 3已收回(级联资源收回) 4不可用(级联未连接)',
    `net_channel_id` varchar(64) not null comment '网络信号通道ID',
    `net_channel_alias` varchar(256) comment '网络信号通道名称',
    `net_channel_state` int comment '网络信号状态: 0离线(网络点位离线/模拟信号无信号) 1在线(网络点位在线/模拟信号有信号) 2已删除(板子拔走/网络点位删除) 3已收回(级联资源收回) 4不可用(级联未连接)',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`cfg_id`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='物理通道与网络通道关联映射表';


-- ---------------------------------------------------------
-- 通道资源分组信息(大屏实时信号、模拟信号、网络信号),构建通道资源树
-- ---------------------------------------------------------
create table if not exists `mpu_channel_group`
(
    `group_id` varchar(64) unique not null comment '分组id,在本服务统一管理的分组id',
    `parent_id` varchar(64) default '' comment '父分组id',
    `cascade_id` varchar(64) not null comment '级联id,标识通道分组属于哪一级显控平台',
    `group_id_real` varchar(64) default '' comment '在原平台上的分组id,如同步Umt平台组织节点ID. 默认为0,参数无效',
    `parent_id_real` varchar(64) default '' comment '在原平台上的分组id,如同步Umt平台组织节点父节点ID. 默认为0,参数无效',
    `group_name` varchar(256) not null comment '分组名称',
    `is_default_group` int default '0' comment '是否是系统默认分组: 0用户自定义分组 1系统内置默认分组节点 2系统自生成角色绑定输出通道资源分组 3系统自生成大屏使用信号源资源分组 4系统同步一机一档组织节点 5系统自生成级联资源分组节点',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`group_id`),
    foreign key(cascade_id) references mpu_cascade_cfg(cascade_id) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='通道资源分组信息(通道资源树节点)';
alter table `mpu_channel_group` ADD INDEX parent_id_index (`parent_id`);
alter table `mpu_channel_group` ADD INDEX group_id_real_index (`group_id_real`);
alter table `mpu_channel_group` ADD INDEX is_default_group_index (`is_default_group`);
alter table `mpu_channel_group` ADD `group_member_number` bigint default null comment '分组内通道个数';




-- ----------------------------------------------------
-- 通道资源分组成员(实际为各大屏对应实时信号源)
-- ----------------------------------------------------
create table if not exists `mpu_channel_group_to_channel`
(
    `group_id` varchar(64) not null comment '分组id',
    `channel_id` varchar(64) not null comment '通道id',
    primary key (`group_id`,`channel_id`),
    foreign key(group_id) references mpu_channel_group(group_id) on update cascade on delete cascade,
    foreign key (channel_id) references mpu_channel(channel_id) on update cascade on delete cascade
)  engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='通道分组、通道关联表';
alter table mpu_channel_group_to_channel drop foreign key mpu_channel_group_to_channel_ibfk_2;
alter table `mpu_channel_group_to_channel` ADD INDEX group_id_index (`group_id`);

-- ------------------------------------------------------
-- 通道收藏夹管理
-- ------------------------------------------------------
create table if not exists `mpu_channel_favorites`
(
    `favorites_id` varchar(64) unique not null comment '收藏夹ID',
    `parent_id` varchar(64) default '' comment '父分组ID',
    `favorites_name` varchar(256) not null comment '收藏夹名称',
    `user_id` bigint not null comment '收藏夹所属用户id',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`favorites_id`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='通道收藏夹管理';


-- ----------------------------
-- 通道收藏夹次序管理
-- ----------------------------
CREATE table if not exists `mpu_channel_favorites_order`
(
    `favorites_id` varchar(64) NOT NULL COMMENT '收藏夹ID',
    `ind_in_par` int(11) unsigned DEFAULT NULL COMMENT '在父组里的次序',
    PRIMARY KEY (`favorites_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- -----------------------------------------------------
-- 通道收藏夹成员管理
-- -----------------------------------------------------
create table if not exists `mpu_channel_favorites_to_channel`
(
    `favorites_id` varchar(64) not null comment '收藏夹id',
    `channel_id` varchar(64) not null comment '通道id',
    `channel_alias` varchar(256) comment '通道名称',
    `channel_state` int comment '通道状态: 0离线(网络点位离线/模拟信号无信号) 1在线(网络点位在线/模拟信号有信号) 2已删除(板子拔走/网络点位删除) 3已收回(级联资源收回) 4不可用(级联未连接)',
    `channel_source` tinyint not null default '0' comment '通道来源: 0未知 1显控平台资源(模拟信号) 2Umt平台(网络信号)',
    `channel_type` tinyint not null comment '通道类型: 0无效 1音频输入 2视频输入 4音频输出 8视频输出 16编码通道 32虚拟通道',
    primary key (`favorites_id`,`channel_id`),
    foreign key(favorites_id) references mpu_channel_favorites(favorites_id) on update cascade on delete cascade
)  engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='通道收藏夹、通道关联表';


-- ------------------------------------------------------
-- 音频通道收藏夹管理
-- ------------------------------------------------------
create table if not exists `mpu_audio_channel_favorites`
(
    `favorites_id` varchar(64) unique not null comment '收藏夹ID',
    `parent_id` varchar(64) default '' comment '父分组ID',
    `favorites_name` varchar(256) not null comment '收藏夹名称',
    `user_id` bigint not null comment '收藏夹所属用户id',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`favorites_id`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='通道收藏夹管理';

-- ----------------------------
-- 音频通道收藏夹次序管理
-- ----------------------------
CREATE table if not exists `mpu_audio_channel_favorites_order`
(
    `favorites_id` varchar(64) NOT NULL COMMENT '音频收藏夹ID',
    `ind_in_par` int(11) unsigned DEFAULT NULL COMMENT '在父组里的次序',
    PRIMARY KEY (`favorites_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -----------------------------------------------------
-- 音频通道收藏夹成员管理
-- -----------------------------------------------------
create table if not exists `mpu_audio_channel_favorites_to_channel`
(
    `favorites_id` varchar(64) not null comment '收藏夹id',
    `channel_id` varchar(64) not null comment '通道id',
    `channel_alias` varchar(256) comment '通道名称',
    `channel_state` int comment '通道状态: 0离线(网络点位离线/模拟信号无信号) 1在线(网络点位在线/模拟信号有信号) 2已删除(板子拔走/网络点位删除) 3已收回(级联资源收回) 4不可用(级联未连接)',
    `channel_source` tinyint not null default '0' comment '通道来源: 0未知 1显控平台资源(模拟信号) 2Umt平台(网络信号)',
    `channel_type` tinyint not null comment '通道类型: 0无效 1音频输入 2视频输入 4音频输出 8视频输出 16编码通道 32虚拟通道',
    primary key (`favorites_id`,`channel_id`),
    foreign key(favorites_id) references mpu_audio_channel_favorites(favorites_id) on update cascade on delete cascade
)  engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='通道收藏夹、通道关联表';


-- ----------------------------------------------------------
-- 平台信息<->设备: 多对多的关系
-- ----------------------------------------------------------
create table if not exists `mpu_access_platform_to_dev`
(
    `access_plat_id` varchar(64) not null comment '接入平台id',
    `dev_id` varchar(64) not null comment '设备id',
    primary key (`access_plat_id`,`dev_id`),
    foreign key(access_plat_id) references mpu_access_platform(access_plat_id) on update cascade on delete cascade,
    foreign key(dev_id) references mpu_dev(dev_id) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='接入平台、设备信息映射表';


-- -------------------------------------------
-- 电视墙(大屏)基本信息
-- -------------------------------------------
create table if not exists `mpu_tv_wall`
(
    `tv_wall_id` varchar(64) unique not null comment '大屏id',
    `cascade_id` varchar(64) not null comment '级联id,标识大屏属于哪一级显控平台',
    `tv_wall_id_real` bigint not null comment '在原显控平台上的大屏id',
    `tv_wall_alias` varchar(256) comment '大屏别名',
    `tv_wall_mode` tinyint comment '大屏模式: 0简易模式 1自定义模式',
    `tv_wall_cell_num` tinyint comment '大屏单元格个数',
    `tv_wall_horizontal_cell` tinyint comment '水平单元格个数',
    `tv_wall_vertical_cell` tinyint comment '垂直单元格个数',
    `tv_wall_pixels_width` smallint comment 'w分辨率',
    `tv_wall_pixels_height` smallint comment 'h分辨率',
    `tv_wall_fps` smallint comment '屏幕刷新频率',
    `tv_wall_bufferW` int comment 'fpga缓冲大小',
    `tv_wall_bufferH` int comment 'fpga缓冲大小',
    `tv_wall_base_map_id` int comment '大屏底图id',
    `tv_wall_current_scene_id` varchar(64) comment '大屏当前加载的场景id',
    `tv_wall_current_scene_id_real` bigint comment '大屏当前加载的场景id(显控查到的场景id)',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`tv_wall_id`),
    foreign key(cascade_id) references mpu_cascade_cfg(cascade_id) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='电视墙(大屏)基本信息';



-- ------------------------------------------------
-- 大屏单元格信息.
-- 大屏<->大屏单元格 一对多
-- ------------------------------------------------
create table if not exists `mpu_tv_wall_cell`
(
    `tv_wall_id` varchar(64) not null comment '大屏id',
    `cell_id` tinyint not null comment '单元格id',
    `cell_channel_id` varchar(64) comment '单元格对应通道id',
    `cell_channel_id_real` varchar(64) comment '单元格对应通道id(显控平台上的通道id)',
    `layoutX` smallint comment '布局参数:起始横坐标',
    `layoutY` smallint comment '布局参数:起始纵坐标',
    `layoutW` smallint comment '布局参数:单元格宽',
    `layoutH` smallint comment '布局参数:单元格高',
    `cutX` smallint comment '裁剪参数:起始横坐标',
    `cutY` smallint comment '裁剪参数:起始纵坐标',
    `cutW` smallint comment '裁剪参数:宽',
    `cutH` smallint comment '裁剪参数:高',
    primary key (tv_wall_id,cell_id),
    foreign key(tv_wall_id) references mpu_tv_wall(tv_wall_id) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='电视墙(大屏)单元格信息';

-- ------------------------------------------------
-- 大屏多画面风格基本信息
-- ------------------------------------------------
create table if not exists `mpu_tv_wall_multi_win_style`
(
    `multi_win_style_id` bigint unique not null auto_increment comment '多画面风格配置id',
    `alias` varchar(256) comment '风格描述',
    `winNum` int comment '窗口个数',
    primary key (`multi_win_style_id`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci auto_increment=100 comment='大屏多画面风格基本信息';
-- 大屏多画面风格窗口布局信息
create table if not exists `mpu_tv_wall_multi_win_layout`
(
    `index` int not null comment '布局窗口index',
    `multi_win_style_id` bigint not null comment '多画面风格配置id',
    `layoutX_ratio` varchar(64) comment '横坐标占比：窗口起始横坐标占整个屏宽度的比值; 例: 0.3333',
    `layoutY_ratio` varchar(64) comment '纵坐标占比：窗口起始纵坐标占整个屏高度的比值; 例: 0.3333',
    `layoutW_ratio` varchar(64) comment '宽度占比：窗口宽度占整个屏宽度的比值; 例: 0.3333',
    `layoutH_ratio` varchar(64) comment '高度占比：窗口高度占整个屏高度的比值; 例: 0.3333',
    primary key (`index`,`multi_win_style_id`),
    foreign key(`multi_win_style_id`) references mpu_tv_wall_multi_win_style(`multi_win_style_id`) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='大屏多画面风格布局信息';

BEGIN;
-- 初始化默认风格信息(保留前100个ID为系统默认多画面风格布局配置)
DELETE FROM `mpu_tv_wall_multi_win_style` WHERE multi_win_style_id < 100;
DELETE FROM `mpu_tv_wall_multi_win_layout` WHERE multi_win_style_id < 100;
-- 1*1
INSERT INTO `mpu_tv_wall_multi_win_style` VALUES (1, '1*1', 1);
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (1, 1, '0.0000', '0.0000', '1.0000', '1.0000');
-- 2*1
INSERT INTO `mpu_tv_wall_multi_win_style` VALUES ('2', '2*1', 2);
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (1, 2, '0.0000', '0.0000', '0.5000', '1.0000');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (2, 2, '0.5000', '0.0000', '0.5000', '1.0000');
-- 3*1
INSERT INTO `mpu_tv_wall_multi_win_style` VALUES ('3', '3*1', 3);
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (1, 3, '0.0000', '0.0000', '0.3333', '1.0000');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (2, 3, '0.3333', '0.0000', '0.3333', '1.0000');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (3, 3, '0.6666', '0.0000', '0.3333', '1.0000');
-- 2*2
INSERT INTO `mpu_tv_wall_multi_win_style` VALUES ('4', '2*2', 4);
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (1, 4, '0.0000', '0.0000', '0.5000', '0.5000');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (2, 4, '0.5000', '0.0000', '0.5000', '0.5000');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (3, 4, '0.0000', '0.5000', '0.5000', '0.5000');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (4, 4, '0.5000', '0.5000', '0.5000', '0.5000');
-- 3*3
INSERT INTO `mpu_tv_wall_multi_win_style` VALUES ('5', '3*3', 9);
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (1, 5, '0.0000', '0.0000', '0.3333', '0.3333');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (2, 5, '0.3333', '0.0000', '0.3333', '0.3333');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (3, 5, '0.6666', '0.0000', '0.3333', '0.3333');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (4, 5, '0.0000', '0.3333', '0.3333', '0.3333');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (5, 5, '0.3333', '0.3333', '0.3333', '0.3333');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (6, 5, '0.6666', '0.3333', '0.3333', '0.3333');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (7, 5, '0.0000', '0.6666', '0.3333', '0.3333');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (8, 5, '0.3333', '0.6666', '0.3333', '0.3333');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (9, 5, '0.6666', '0.6666', '0.3333', '0.3333');
-- 标准五画面
INSERT INTO `mpu_tv_wall_multi_win_style` VALUES ('6', '五画面', 5);
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (1, 6, '0.0000', '0.0000', '0.5000', '0.6666');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (2, 6, '0.5000', '0.0000', '0.5000', '0.6666');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (3, 6, '0.0000', '0.6666', '0.3333', '0.3333');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (4, 6, '0.3333', '0.6666', '0.3333', '0.3333');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (5, 6, '0.6666', '0.6666', '0.3333', '0.3333');
-- 标准六画面
INSERT INTO `mpu_tv_wall_multi_win_style` VALUES ('7', '六画面', 6);
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (1, 7, '0.0000', '0.0000', '0.6666', '0.6666');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (2, 7, '0.6666', '0.0000', '0.3333', '0.3333');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (3, 7, '0.6666', '0.3333', '0.3333', '0.3333');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (4, 7, '0.0000', '0.6666', '0.3333', '0.3333');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (5, 7, '0.3333', '0.6666', '0.3333', '0.3333');
INSERT INTO `mpu_tv_wall_multi_win_layout` VALUES (6, 7, '0.6666', '0.6666', '0.3333', '0.3333');
COMMIT;

-- --------------------------------------------
-- 大屏使用场景信息(实时数据)
-- 大屏<->场景: 一对多的关系
-- --------------------------------------------
create table if not exists `mpu_tv_wall_scene`
(
    `scene_id` varchar(64) unique not null comment '场景id',
    `tv_wall_id` varchar(64) not null comment '大屏id',
    `cascade_id` varchar(64) not null comment '级联id,标识属于哪一级显控平台',
    `index` bigint comment '自定义大屏场景展示排列顺序',
    `multi_win_style_id` bigint default '0' comment '场景对应多画面风格配置id',

    `scene_id_real` bigint comment '在原平台上的场景id',
    `scene_alias` varchar(256) comment '使用场景描述',
    `scene_lamp_state` tinyint comment '场景跑马灯状态',
    `scene_sync_poll` tinyint comment '同步',
    `enable_use_back_pic` tinyint comment '是否启用底图',
    `scene_back_pic_id` int comment '底图id',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`scene_id`),
    foreign key(tv_wall_id) references mpu_tv_wall(tv_wall_id) on update cascade on delete cascade,
    foreign key(cascade_id) references mpu_cascade_cfg(cascade_id) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='电视墙(大屏)使用场景信息(实时数据)';

-- --------------------------------------------
-- 大屏开窗信息(实时数据)
-- 场景<->开窗: 一对多的关系
-- --------------------------------------------
create table if not exists `mpu_window`
(
    `window_id` varchar(64) unique not null comment '开窗id',
    `scene_id` varchar(64) not null comment '对应场景id',
    `cascade_id` varchar(64) not null comment '级联id,标识属于哪一级显控平台',
    `window_id_real` bigint not null comment '对应mpuserver侧窗口ID(非机箱冗余环境)',
    `window_sort` tinyint comment '窗口叠加时排列排序',
    `layoutX`  int comment '窗口布局: 起始横坐标',
    `layoutY`  int comment '窗口布局: 起始纵坐标',
    `layoutW`  int comment '窗口布局: 窗口宽',
    `layoutH`  int comment '窗口布局: 窗口高',
    `cutX`  int comment '裁剪参数: 起始横坐标',
    `cutY`  int comment '裁剪参数: 起始纵坐标',
    `cutW`  int comment '裁剪参数: 宽',
    `cutH`  int comment '裁剪参数: 高',
    `cut_enable` bool comment '裁剪使能',
    `audio_open` bool comment '音频开关',
    `multi_window_enable` bool comment '多画面使能,目前只有网络流的开窗才支持多画面功能',
    `style` tinyint comment '窗口风格: 1单画面 2两画面 3三画面 4四画面 ...',
    `result` int default null comment '整个窗口的开窗结果(机箱冗余环境任意机箱开窗成功即不报错)',
    `err_msg` varchar(256) default '' comment '开窗报错信息',
    primary key (`scene_id`,`window_id`),
    foreign key(`scene_id`) references mpu_tv_wall_scene(`scene_id`) on update cascade on delete cascade,
    foreign key(cascade_id) references mpu_cascade_cfg(cascade_id) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='场景开窗信息(实时数据)';
alter table `mpu_window` modify column `window_id_real` bigint comment '对应mpuserver侧窗口ID(非机箱冗余环境)';
alter table `mpu_window` add column `window_id_real_master` bigint comment '机箱冗余：主机箱mpuserver侧窗口ID';
alter table `mpu_window` add column `result_master` int default null comment '机箱冗余：主机箱mpuserver返回开窗错误码';
alter table `mpu_window` add column `window_id_real_Slave` bigint comment '机箱冗余：备机箱mpuserver侧窗口ID';
alter table `mpu_window` add column `result_Slave` int default null comment '机箱冗余：备机箱mpuserver返回开窗错误码';

-- -----------------------------------------------------
-- 窗口的子窗口信息(多画面窗口,每个子窗对应的信号源,开窗结果)(实时数据)
-- 窗口<->子窗口: 一对多
-- -----------------------------------------------------
create table if not exists `mpu_sub_window`
(
    `window_id` varchar(64) not null comment '对应主窗口id',
    `sub_window_id` bigint not null comment '子窗口id,从0开始',
    `input_channel_id` varchar(64) comment '输入通道id',
    `input_channel_id_real` varchar(64) comment '输入通道id(显控分配的通道id)',
    `input_channel_alias` varchar(256) comment '输入通道名称',
    `result` int default null comment '子窗口开窗结果, 对应显控定义错误码',
    `err_msg` varchar(256) default '' comment '每个子窗开窗失败的错误信息,默认为空',
    primary key (`window_id`,`sub_window_id`),
    foreign key(window_id) references mpu_window(window_id) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='窗口的子窗口信息(实时数据)';


-- --------------------------------------------
-- 大屏使用场景信息(备份数据)
-- 大屏<->场景: 一对多的关系
-- --------------------------------------------
create table if not exists `mpu_tv_wall_scene_backup`
(
    `scene_id` varchar(64) unique not null comment '场景id',
    `tv_wall_id` varchar(64) not null comment '大屏id',
    `cascade_id` varchar(64) not null comment '级联id,标识属于哪一级显控平台',
    `index` bigint comment '自定义大屏场景展示排列顺序',
    `multi_win_style_id` bigint default '0' comment '场景对应多画面风格配置id',

    `scene_id_real` bigint comment '在原平台上的场景id',
    `scene_alias` varchar(256) comment '使用场景描述',
    `scene_lamp_state` tinyint comment '场景跑马灯状态',
    `scene_sync_poll` tinyint comment '同步',
    `enable_use_back_pic` tinyint comment '是否启用底图',
    `scene_back_pic_id` int comment '底图id',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`scene_id`),
    foreign key(tv_wall_id) references mpu_tv_wall(tv_wall_id) on update cascade on delete cascade,
    foreign key(cascade_id) references mpu_cascade_cfg(cascade_id) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='电视墙(大屏)使用场景信息(备份数据)';
-- --------------------------------------------
-- 场景开窗信息(备份数据)
-- 场景<->开窗: 一对多的关系
-- --------------------------------------------
create table if not exists `mpu_window_backup`
(
    `window_id` varchar(64) unique not null comment '开窗id',
    `scene_id` varchar(64) not null comment '场景id',
    `cascade_id` varchar(64) not null comment '级联id,标识属于哪一级显控平台',
    `window_id_real` bigint not null comment '在原平台上的窗口id',
    `window_sort` tinyint comment '窗口叠加时排列排序',
    `layoutX`  int comment '窗口布局: 起始横坐标',
    `layoutY`  int comment '窗口布局: 起始纵坐标',
    `layoutW`  int comment '窗口布局: 窗口宽',
    `layoutH`  int comment '窗口布局: 窗口高',
    `cutX`  int comment '裁剪参数: 起始横坐标',
    `cutY`  int comment '裁剪参数: 起始纵坐标',
    `cutW`  int comment '裁剪参数: 宽',
    `cutH`  int comment '裁剪参数: 高',
    `cut_enable` bool comment '裁剪使能',
    `audio_open` bool comment '音频开关',
    `multi_window_enable` bool comment '多画面使能,目前只有网络流的开窗才支持多画面功能',
    `style` tinyint comment '窗口风格: 1单画面 2两画面 3三画面 4四画面 ...',
    `result` int default null comment '整个窗口的开窗结果',
    `err_msg` varchar(256) default '' comment '开窗报错信息',
    primary key (`scene_id`,`window_id`),
    foreign key(`scene_id`) references mpu_tv_wall_scene_backup(`scene_id`) on update cascade on delete cascade,
    foreign key(cascade_id) references mpu_cascade_cfg(cascade_id) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='场景开窗信息(备份数据)';
alter table `mpu_window_backup` modify column `window_id_real` bigint comment '对应mpuserver侧窗口ID(非机箱冗余环境)';
alter table `mpu_window_backup` add column `window_id_real_master` bigint comment '机箱冗余：主机箱mpuserver侧窗口ID';
alter table `mpu_window_backup` add column `result_master` int default null comment '机箱冗余：主机箱mpuserver返回开窗错误码';
alter table `mpu_window_backup` add column `window_id_real_Slave` bigint comment '机箱冗余：备机箱mpuserver侧窗口ID';
alter table `mpu_window_backup` add column `result_Slave` int default null comment '机箱冗余：备机箱mpuserver返回开窗错误码';
-- -----------------------------------------------------------------
-- 窗口的子窗口信息(多画面窗口,每个子窗对应的信号源,开窗结果)(备份数据)
-- 窗口<->子窗口: 一对多
-- -----------------------------------------------------------------
-- drop table `mpu_sub_window_backup`;
create table if not exists `mpu_sub_window_backup`
(
    `window_id` varchar(64) not null comment '对应主窗口id',
    `sub_window_id` bigint not null comment '子窗口id',
    `input_channel_id` varchar(64) comment '输入通道id(本代理分配管理的通道id)',
    `input_channel_id_real` varchar(64) comment '输入通道id(显控查到的通道id)',
    `input_channel_alias` varchar(256) comment '输入通道名称',
    `result` int default null comment '子窗口开窗结果, 对应显控定义错误码',
    `err_msg` varchar(256) default '' comment '每个子窗开窗失败的错误信息,默认为空',
    primary key (`window_id`,`sub_window_id`),
    foreign key(window_id) references mpu_window_backup(window_id) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='窗口的子窗口信息(备份数据)';



-- ---------------------------------------------------------------------------
-- 大屏轮循预案配置(大屏场景轮巡|窗口通道轮巡|批量信号源轮循)
-- ---------------------------------------------------------------------------
create table if not exists `mpu_tv_wall_loop_cfg`
(
    `loop_cfg_id` varchar(64) unique not null comment '轮循配置id',
    `loop_alias`  varchar(256) comment '轮循预案描述',
    `loop_type`   tinyint comment '轮循类型: 1大屏场景轮循 2窗口通道轮循 3多窗口批量轮循',
    `loop_dest_id` varchar(64) not null comment '轮循目标ID：大屏ID/窗口ID',
    `loop_count`  tinyint comment '轮循次数,为0时无限轮循需要手动停止,非0时达到轮循次数自动停止',
    `loop_duration` int comment '轮循时长,单位min,为0时无限轮循需要手动停止,非0时达到轮循时间自动停止',
    `loop_start_time_ms` bigint comment '轮循任务启动时刻,毫秒时间,用于重启时判断是否需要恢复轮循任务',
    `loop_interval_time_default` int comment '默认成员轮询时间间隔,单位s,每个轮询成员可配置自已的时间间隔,优先以轮巡成员配置的时间间隔为准',
    `loop_state` bool comment '轮循任务状态: true启用 false未启用. 对于同一个轮循目标,同时只有一个轮循配置是启用的',
    `loop_member_number` int comment '轮循成员个数',
    primary key (`loop_cfg_id`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='轮循配置';
alter table `mpu_tv_wall_loop_cfg` modify column loop_type tinyint comment '轮循类型: 1大屏场景轮循 2窗口通道轮循 3多窗口批量轮循';
alter table `mpu_tv_wall_loop_cfg` add column loop_mode tinyint comment '多窗口轮循模式: 1从头开始 2连续循环';

-- --------------------------------------------------------------------
-- 大屏轮循预案成员参数
-- --------------------------------------------------------------------
create table if not exists `mpu_tv_wall_loop_member`
(
    `loop_cfg_id` varchar(64) not null comment '轮循配置ID',
    `loop_member_index` bigint not null comment '轮循成员排序,查询轮循成员列表和轮循成员轮循顺序默认按index升序',
    `loop_member_id` varchar(64) not null comment '轮循成员ID：loop_type=1为场景ID; loop_type=2为通道ID; loop_type=3为窗口ID',
    `loop_member_alias` varchar(256) not null comment '轮循成员描述: 场景描述/通道描述',
    `loop_interval_time` int comment '成员轮循间隔,单位s',
    `loop_member_state` bool comment '轮循成员状态: false标识未在轮循,true标识轮循到当前成员. 对于同一个轮循预案,只有一个轮循成员在轮循中.',
    foreign key(loop_cfg_id) references mpu_tv_wall_loop_cfg(loop_cfg_id) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='大屏轮循预案成员参数';
alter table mpu_tv_wall_loop_member add primary key (loop_cfg_id, loop_member_index);

-- -------------------------------------------------------------------------
-- 大屏轮循预案配置(窗口通道轮巡配置数据备份表)
-- 场景切换、重置可恢复对应窗口的通道轮循配置
-- 注: 点击场景保存时,会将对应场景开窗窗口的轮循配置保存到此表
-- -------------------------------------------------------------------------
create table if not exists `mpu_tv_wall_loop_cfg_backup`
(
    `loop_cfg_id` varchar(64) unique not null comment '轮循配置id',
    `loop_alias`  varchar(256) comment '轮循预案描述',
    `loop_type`   tinyint comment '轮循类型: 1大屏场景轮循 2窗口通道轮循',
    `loop_dest_id` varchar(64) not null comment '轮循目标ID：大屏ID/窗口ID',
    `loop_count`  tinyint comment '轮循次数,为0时无限轮循需要手动停止,非0时达到轮循次数自动停止',
    `loop_duration` int comment '轮循时长,单位min,为0时无限轮循需要手动停止,非0时达到轮循时间自动停止',
    `loop_start_time_ms` bigint comment '轮循任务启动时刻,毫秒时间,用于重启时判断是否需要恢复轮循任务',
    `loop_interval_time_default` int comment '默认成员轮询时间间隔,单位s,每个轮询成员可配置自已的时间间隔,优先以轮巡成员配置的时间间隔为准',
    `loop_state` bool comment '轮循任务状态: true启用 false未启用. 对于同一个轮循目标,同时只有一个轮循配置是启用的',
    `loop_member_number` int comment '轮循成员个数',
    primary key (`loop_cfg_id`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='轮循配置';
-- --------------------------------------------------------------------
-- 大屏轮循预案成员参数(窗口通道轮巡配置数据备份表)
-- 场景切换、重置可恢复对应窗口的通道轮循配置
-- 注: 点击场景保存时,会将对应场景开窗窗口的轮循配置保存到此表
-- --------------------------------------------------------------------
create table if not exists `mpu_tv_wall_loop_member_backup`
(
    `loop_cfg_id` varchar(64) not null comment '轮循配置ID',
    `loop_member_index` bigint not null comment '轮循成员排序,查询轮循成员列表和轮循成员轮循顺序默认按index升序',
    `loop_member_id` varchar(64) not null comment '轮循成员ID：场景ID/通道ID',
    `loop_member_alias` varchar(256) not null comment '轮循成员描述: 场景描述/通道描述',
    `loop_interval_time` int comment '成员轮循间隔,单位s',
    `loop_member_state` bool comment '轮循成员状态: false标识未在轮循,true标识轮循到当前成员. 对于同一个轮循预案,只有一个轮循成员在轮循中.',
    foreign key(loop_cfg_id) references mpu_tv_wall_loop_cfg_backup(loop_cfg_id) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='大屏轮循预案成员参数';
alter table mpu_tv_wall_loop_member_backup add primary key (loop_cfg_id, loop_member_index);

-- --------------------------------------------------------
-- 用户批量轮循信号源配置
-- 用户级配置参数，配置与角色绑定，相同角色的账户共享配置.一个角色最多10个轮循配置
-- --------------------------------------------------------
create table if not exists `mpu_tv_wall_loop_chn_cfg`
(
    `cfg_id` bigint not null auto_increment comment '配置ID',
    `role_id` bigint not null comment '所属角色id',
    `alias`  varchar(256) comment '配置描述',
    `loop_chn_number` int comment '轮循信号源个数',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`cfg_id`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='批量轮循信号源配置';
create table if not exists `mpu_tv_wall_loop_chn_list`
(
    `cfg_id` bigint not null comment '配置id',
    `index` bigint not null comment '轮循信号源排序',
    `channel_id` varchar(64) not null comment '通道id',
    primary key (`cfg_id`,`channel_id`),
    foreign key(cfg_id) references mpu_tv_wall_loop_chn_cfg(cfg_id) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='批量轮循信号源列表';

-- ------------------------------------------------------
-- 级联点位推送配置,记录本级向直属下级推送通道的具体配置信息
-- ------------------------------------------------------
-- drop table `mpu_cascade_push_channel_cfg`;
create table if not exists `mpu_cascade_push_channel_cfg`
(
    `cfg_id` varchar(64) unique not null comment '推送配置ID,全局唯一',
    `cascade_id` varchar(64) not null comment '级联id,标识配置属于哪一级显控代理平台上的配置,方便级联时各级的配置信息汇聚',
    `lower_plat_cascade_id` varchar(64) not null comment '下级平台级联ID',
    `channel_id` varchar(64) not null comment '通道id,推送给下级的点位ID',
    `channel_alias` varchar(256) comment '通道名称',
    `channel_state` int comment '通道状态: 0离线(网络点位离线/模拟信号无信号) 1在线(网络点位在线/模拟信号有信号) 2已删除(板子拔走/网络点位删除) 3已收回(级联资源收回) 4不可用(级联未连接)',
    `push_channel_inuse_by_lower_plat` bool comment '下级平台是否在使用该通道开窗浏览 -- 未要求同步下级对该信号的使用情况,字段暂时保留',
    `push_channel_call_result_by_lower_plat` varchar(256) comment '下级平台使用该通道报错信息 -- 未要求同步下级对该信号的使用情况,字段暂时保留',
    `push_result` varchar(256) comment '推送结果',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '推送时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`cfg_id`),
    foreign key(cascade_id) references mpu_cascade_cfg(cascade_id) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='级联点位推送配置信息';


-- -----------------------------------------------------------------
-- 上级推送点位默认播放配置,记录某上级推送的点位默认在本级哪个大屏上播放浏览
-- -----------------------------------------------------------------
create table if not exists `mpu_cascade_push_channel_default_play_cfg`
(
    `cfg_id` varchar(64) unique not null comment '推送点位默认播放配置ID',
    `cascade_id` varchar(64) not null comment '本级级联id,标识配置属于哪一级显控代理平台上的配置,方便级联时各级的配置信息汇聚',
    `superior_plat_cascade_id` varchar(64) not null comment '上级级联id -- 未使用,暂不支持为多个上级单独配置多个默认播放配置',
    `tv_wall_id` varchar(64) not null comment '默认播放的大屏ID',
    `tv_wall_alias` varchar(256) not null comment '大屏别名',
    `layoutX` int comment '布局参数: 起始横坐标',
    `layoutY` int comment '布局参数: 起始纵坐标',
    `layoutW` int comment '布局参数: 宽',
    `layoutH` int comment '布局参数: 高',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`cfg_id`),
    foreign key(cascade_id) references mpu_cascade_cfg(cascade_id) on update cascade on delete cascade,
    foreign key(tv_wall_id) references mpu_tv_wall(tv_wall_id) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='推送点位默认播放配置信息';


-- -------------------------------------------------------
-- 显控矩阵调度预案分组(不需要支持多级分组)
-- -------------------------------------------------------
create table if not exists `mpu_matrix_dispatch_scheme_group`
(
    `group_id` bigint unique not null auto_increment comment '分组ID',
    `parent_id` bigint default '0' comment '父分组ID',
    `group_name` varchar(256) not null comment '分组名称',

    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`group_id`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='显控矩阵调度预案分组';

-- -------------------------------------------------------
-- 显控矩阵调度预案信息
-- -------------------------------------------------------
create table if not exists `mpu_matrix_dispatch_scheme`
(
    `scheme_id` bigint unique not null auto_increment comment '矩阵调度预案ID',
    `group_id` bigint default '0' comment '所属调度预案分组ID,不属于任何分组为0',
    `scheme_name` varchar(256) not null comment '矩阵调度预案名称',
    `load_state` bool comment '预案加载状态',
    `load_time_ms` bigint comment '矩阵预案启动时刻(毫秒时间戳)',

    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`scheme_id`)
) engine=innodb auto_increment=1025 default charset=utf8mb4 collate=utf8mb4_general_ci comment='显控矩阵调度预案信息';

-- -------------------------------------------------------
-- 显控矩阵调度分组(不需要支持多级分组)
-- -------------------------------------------------------
create table if not exists `mpu_matrix_dispatch_group`
(
    `group_id` bigint unique not null auto_increment comment '矩阵调度分组ID',
    `scheme_id` bigint not null comment '所属矩阵调度预案ID',
    `group_name` varchar(256) not null comment '分组名称',
    `index` bigint comment '同一个矩阵预案下的多个调度组的排列顺序(默认从1开始排序)',

    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`group_id`),
    foreign key(`scheme_id`) references mpu_matrix_dispatch_scheme(`scheme_id`) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='显控矩阵调度分组信息';



-- ----------------------------------------------------
-- 显控矩阵调度信息(配置数据)
-- ----------------------------------------------------
create table if not exists `mpu_matrix_dispatch`
(
    `dispatch_id` bigint unique not null auto_increment comment '矩阵调度ID',
    `group_id` bigint default '0' comment '所属矩阵调度分组ID',
    `scheme_id` bigint not null comment '所属矩阵调度预案ID',

    `dispatch_type` int not null comment '矩阵调度类型: 1视屏矩阵 2音频矩阵 3音视频矩阵 4混音矩阵',
    `dst_channel_id` varchar(64) not null comment '矩阵调度目的通道id,只能是本级通道',
    `dst_channel_id_real` varchar(64) comment '对应显控侧分配的通道ID',
    `dst_channel_alias` varchar(256) comment '矩阵调度目的通道名称',
    `dst_channel_state` int comment '目的通道状态: 0离线(网络点位离线/模拟信号无信号) 1在线(网络点位在线/模拟信号有信号) 2已删除(板子拔走/网络点位删除) 3已收回(级联资源收回) 4不可用(级联未连接)',
    `dst_channel_occupied` int comment '目的通道占用状态(是否冲突): 0不冲突 1抢占中',
    `dst_volume` int comment '输出音量(百分比1-100)',
    `dst_is_mute` bool comment '输出是否静音',

    `mixer_volume` int comment '混音器输出音量(百分比1-100)',
    `mixer_is_mute` bool comment '混音器输出是否静音',

    `enable` bool comment '调度启用状态',
    `err_code` int comment '加载状态码',
    `err_msg` varchar(256) comment '加载报错信息',
    `load_time_ms` bigint comment '矩阵调度加载时刻(毫秒时间戳)',

    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`dispatch_id`),
    foreign key(`scheme_id`) references mpu_matrix_dispatch_scheme(`scheme_id`) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='显控矩阵调度基本信息';
alter table `mpu_matrix_dispatch` add column associated_out_channel_id varchar(64) not null comment '相互关联的音视频输出通道ID. 同一个调度组内的【视频调度】和【音频调度】可以相互关联，此字段记录有相互关联的两个调度对应输出通道';

-- ---------------------------------------
-- 矩阵调度源信息(同一个矩阵调度可能有多个源)
-- ---------------------------------------
create table if not exists `mpu_matrix_dispatch_src`
(
    `dispatch_id` bigint unique not null comment '对应矩阵调度ID',
    `src_index` int comment '同一个调度的多个源index(默认从1开始排序)',
    `src_channel_id` varchar(64) not null comment '源通道id,可以是本级资源支持级联也可以是级联资源',
    `src_channel_id_real` varchar(64) comment '对应显控侧分配的通道ID',
    `src_channel_alias` varchar(256) comment '矩阵调度目的通道名称',
    `src_channel_state` int comment '源通道状态: 0离线(网络点位离线/模拟信号无信号) 1在线(网络点位在线/模拟信号有信号) 2已删除(板子拔走/网络点位删除) 3已收回(级联资源收回) 4不可用(级联未连接)',
    `src_volume` int comment '输入音量(百分比1-100)',
    `src_is_mute` bool comment '是否静音',

    `err_code` int comment '调度源调度错误状态码',
    `err_msg` varchar(256) comment '调度源调度报错信息',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`dispatch_id`,`src_channel_id`),
    foreign key(`dispatch_id`) references mpu_matrix_dispatch(`dispatch_id`) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='显控矩阵调度源';



-- ----------------------------------------------------
-- 显控矩阵调度信息(加载态数据)
-- ----------------------------------------------------
create table if not exists `mpu_matrix_dispatch_loading`
(
    `dispatch_id` bigint unique not null auto_increment comment '矩阵调度ID',
    `group_id` bigint default '0' comment '所属矩阵调度分组ID',
    `scheme_id` bigint not null comment '所属矩阵调度预案ID',

    `dispatch_type` int not null comment '矩阵调度类型: 1视屏矩阵 2音频矩阵 3音视频矩阵 4混音矩阵',
    `dst_channel_id` varchar(64) not null comment '矩阵调度目的通道id,只能是本级通道',
    `dst_channel_id_real` varchar(64) comment '对应显控侧分配的通道ID',
    `dst_channel_alias` varchar(256) comment '矩阵调度目的通道名称',
    `dst_channel_state` int comment '目的通道状态: 0离线(网络点位离线/模拟信号无信号) 1在线(网络点位在线/模拟信号有信号) 2已删除(板子拔走/网络点位删除) 3已收回(级联资源收回) 4不可用(级联未连接)',
    `dst_channel_occupied` int comment '目的通道占用状态(是否冲突): 0不冲突 1抢占中',
    `dst_volume` int comment '输出音量(百分比1-100)',
    `dst_is_mute` bool comment '输出是否静音',

    `mixer_volume` int comment '混音器输出音量(百分比1-100)',
    `mixer_is_mute` bool comment '混音器输出是否静音',

    `enable` bool comment '调度启用状态',
    `err_code` int comment '加载状态码',
    `err_msg` varchar(256) comment '加载报错信息',
    `load_time_ms` bigint comment '矩阵调度加载时刻(毫秒时间戳)',

    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`dispatch_id`),
    foreign key(`scheme_id`) references mpu_matrix_dispatch_scheme(`scheme_id`) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='显控矩阵调度基本信息';
alter table `mpu_matrix_dispatch_loading` add column associated_out_channel_id varchar(64) not null comment '相互关联的音视频输出通道ID. 同一个调度组内的【视频调度】和【音频调度】可以相互关联，此字段记录有相互关联的两个调度对应输出通道';

-- -------------------------------------------------------
-- 矩阵调度源信息(同一个矩阵调度可能有多个源)
-- --------------------------------------------------------
create table if not exists `mpu_matrix_dispatch_src_loading`
(
    `dispatch_id` bigint unique not null comment '对应矩阵调度ID',
    `src_index` int comment '同一个调度的多个源index(默认从1开始排序)',
    `src_channel_id` varchar(64) not null comment '源通道id,可以是本级资源支持级联也可以是级联资源',
    `src_channel_id_real` varchar(64) comment '对应显控侧分配的通道ID',
    `src_channel_alias` varchar(256) comment '矩阵调度目的通道名称',
    `src_channel_state` int comment '源通道状态: 0离线(网络点位离线/模拟信号无信号) 1在线(网络点位在线/模拟信号有信号) 2已删除(板子拔走/网络点位删除) 3已收回(级联资源收回) 4不可用(级联未连接)',
    `src_volume` int comment '输入音量(百分比1-100)',
    `src_is_mute` bool comment '是否静音',

    `err_code` int comment '调度源调度错误状态码',
    `err_msg` varchar(256) comment '调度源调度报错信息',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`dispatch_id`,`src_channel_id`),
    foreign key(`dispatch_id`) references mpu_matrix_dispatch_loading(`dispatch_id`) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='显控矩阵调度源';



-- ------------------------------------------------------
-- 显控矩阵调度默认源配置
-- 1、为目的通道配置默认输入源
-- 2、目的通道用于矩阵调度时,自动生成该配置,且默认源为空.
-- 3、多个矩阵调度为相同的目的通道,共用一个默认源配置
-- ------------------------------------------------------
create table if not exists `mpu_matrix_dispatch_default_src`
(
    `cfg_id` bigint unique not null auto_increment comment '配置id',
    `dispatch_type` int not null comment '矩阵调度类型: 1视屏矩阵 2音频矩阵 3音视频矩阵 4混音矩阵',
    `dst_channel_id` varchar(64) not null comment '目的通道ID',
    `dst_channel_id_real` varchar(64) comment '对应显控侧分配的通道ID',
    `dst_channel_alias` varchar(256) comment '矩阵调度目的通道名称',
    `dst_channel_state` int comment '目的通道状态: 0离线(网络点位离线/模拟信号无信号) 1在线(网络点位在线/模拟信号有信号) 2已删除(板子拔走/网络点位删除) 3已收回(级联资源收回) 4不可用(级联未连接)',

    `src_channel_id` varchar(64) not null comment '源通道id,可以是本级资源支持级联也可以是级联资源',
    `src_channel_id_real` varchar(64) comment '对应显控侧分配的通道ID',
    `src_channel_alias` varchar(256) comment '矩阵调度目的通道名称',
    `src_channel_state` int comment '源通道状态: 0离线(网络点位离线/模拟信号无信号) 1在线(网络点位在线/模拟信号有信号) 2已删除(板子拔走/网络点位删除) 3已收回(级联资源收回) 4不可用(级联未连接)',

    `effective` int default '0' comment '配置是否生效: 0未生效 1生效',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`cfg_id`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='显控矩阵调度源';

-- ------------------------------------------------------
-- 显控矩阵调度音视频源绑定配置
-- 【音视频矩阵调度同步切换】功能使用的【关联信号源配置】
-- ------------------------------------------------------
create table if not exists `mpu_matrix_AV_input_chn_bind_cfg`
(
    `cfg_id` bigint unique not null auto_increment comment '配置id',
    `alias` varchar(256) comment '配置描述',
    `index` bigint comment '排序',
    `video_channel_id` varchar(64) not null comment '视频通道ID,只能是本级[视频输入]',
    `audio_channel_id` varchar(64) not null comment '音频通道ID,只能是本级[音频输入][音频输出]',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`cfg_id`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='显控矩阵调度源';


-- ---------------------------------------------
-- 显控告警配置 mpu_alarm_config
-- ---------------------------------------------
create table if not exists `mpu_alarm_config`
(
    `alarm_event_type` int unique not null comment '告警行为(告警类型): 1信号传输异常 2板卡设备离线 3主控板状态异常 ',
    `whether_the_alarm` smallint comment '是否告警: 0不告警 1告警',
    `alarm_mode` smallint comment '告警方式: 1客户端提醒 2邮件提醒 4短信提醒 3客户端提醒+邮件提醒 5客户端提醒+短信提醒 6邮件提醒+短信提醒 7客户端提醒+邮件提醒+短信提醒',
    `alarm_level` smallint comment '告警等级: 1严重 2普通',
    `alarm_param` varchar(2048) default '' comment '告警参数(json格式告警参数)',
    primary key (`alarm_event_type`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='显控告警配置';

-- ---------------------------------------------
-- 显控告警记录 mpu_alarm_record
-- ---------------------------------------------
create table if not exists `mpu_alarm_record`
(
    `alarm_id` bigint unique not null auto_increment comment '告警ID',
    `alarm_event_type` int not null comment '告警事件(告警类型): 1信号传输异常 2板卡设备离线 3主控板状态异常 ',
    `alarm_mode` smallint comment '告警方式: 1客户端提醒 2邮件提醒 4短信提醒 3客户端提醒+邮件提醒 5客户端提醒+短信提醒 6邮件提醒+短信提醒 7客户端提醒+邮件提醒+短信提醒',
    `alarm_level` smallint comment '告警等级: 1严重 2普通',
    `alarm_param` varchar(2048) comment '告警参数(json格式告警参数)',
    `alarm_count` int not null comment '告警次数，针对本次警告，临时消警结束后再次触发告警时次数+1',
    `alarm_time` datetime default null comment '告警时间',
    `alarm_state` smallint default '0' comment '消警方式(告警状态)：0未消警 1忽略告警 2临时消警 3自动消警(告警解除)',
    `alarm_clean_time` datetime default null comment '消警时间',
    `alarm_detail_info` varchar(4096) comment '告警详情(json格式告警详情)',
    primary key (`alarm_id`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='显控告警记录';
alter table `mpu_alarm_record` add column alarm_extend_msg varchar (2048) comment '扩展字段; alarm_event_type=1时该值为通道ID,用于联表查询有权限的音视频信号告警';
-- 声明修改告警信息表PROCEDURE
# drop procedure if exists `alterAlarmRecordTable`;
# DELIMITER $$
# create procedure `alterAlarmRecordTable`()
# top:begin
#         if not exists (select 1 from information_schema.columns where table_name='mpu_alarm_record' and column_name='alarm_extend_msg') then
#             alter table `mpu_alarm_record` add column alarm_extend_msg varchar (2048) comment '扩展字段; alarm_event_type=1时该值为通道ID,用于联表查询有权限的音视频信号告警';
#         end if;
#     end $$
# DELIMITER
#
# DELIMITER
# call alterAlarmRecordTable()
# DELIMITER

-- ---------------------------------------------
-- 临时消警记录表 mpu_alarm_temporary_clean_record
-- ---------------------------------------------
create table if not exists `mpu_alarm_temporary_clean_record`
(
    `alarm_temporary_clean_id` bigint not null comment '临时消警处理ID',
    `alarm_id` bigint not null comment '告警ID',
    `handle_time_ms` bigint comment '处理时间,毫秒时间戳',
    `duration` int not null comment '临时消警时长',
    `optUser` varchar(128) not null comment '操作者',
    primary key (`alarm_id`,`alarm_temporary_clean_id`),
    foreign key(`alarm_id`) references mpu_alarm_record(`alarm_id`) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='临时消警记录表';

-- ---------------------------------------------
-- 显控机箱内设备上下线记录 mpu_device_online_record
-- ---------------------------------------------
create table if not exists `mpu_device_online_record`
(
    `id` bigint unique not null auto_increment comment '记录ID',
    `device_model` varchar(256) not null comment '设备型号: MSP1000-IF08-SM、MSP1000-OF08-LMP、MSP1000-DE08-LG ...',
    `device_alias` varchar(256) comment '设备名称: 输入板、输出板、解码板 ...',
    `device_serial_number` varchar(256) comment '设备序列号',
    `device_box_number` smallint comment '机箱编号',
    `device_slot_number` smallint comment '设备所在槽位号',
    `device_state` smallint comment '设备状态: 0未知 1上线 2离线',
    `time` datetime default null comment '时间',
    primary key (`id`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='显控设备上下线记录';

-- -------------------------------------
-- 显控环境设备表 mpu_appliance_appdev
-- -------------------------------------
CREATE table if not exists `mpu_appliance_appdev` (
                                                      `dev_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '环境设备ID',
                                                      `group_id` bigint(20) DEFAULT NULL COMMENT '所属分组ID',
                                                      `dev_name` varchar(256) CHARACTER SET utf8 DEFAULT NULL COMMENT '环境设备名称',
                                                      `dev_type` tinyint(4) DEFAULT NULL COMMENT '设备类型',
                                                      `dev_module` varchar(256) CHARACTER SET utf8 DEFAULT NULL COMMENT '设备型号',
                                                      `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
                                                      `create_time` datetime DEFAULT NULL COMMENT '更新时间',
                                                      `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
                                                      `update_time` datetime DEFAULT NULL COMMENT '更新时间',
                                                      `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                                      PRIMARY KEY (`dev_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='环境设备信息表';

-- -------------------------------------------
-- 显控环境设备继电器表 mpu_appliance_relay
-- -------------------------------------------
CREATE table if not exists `mpu_appliance_relay` (
                                                     `relay_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '继电器ID',
                                                     `relay_name` varchar(256) CHARACTER SET utf8 DEFAULT NULL COMMENT '继电器名称',
                                                     `relay_type` tinyint(4) DEFAULT NULL COMMENT '继电器型号',
                                                     `ctrl_protocol` tinyint(4) DEFAULT NULL COMMENT '通讯协议',
                                                     `switch_out_num` tinyint(4) DEFAULT NULL COMMENT '开关数目',
                                                     `trx_ipaddr` varchar(50) CHARACTER SET utf8 DEFAULT NULL COMMENT '盒子IP地址',
                                                     `channel_id` varchar(64) CHARACTER SET utf8 DEFAULT NULL COMMENT '盒子对应通道ID',
                                                     `channel_id_real` varchar(64) CHARACTER SET utf8 DEFAULT NULL COMMENT '显控平台上的通道ID',
                                                     `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
                                                     `create_time` datetime DEFAULT NULL COMMENT '创建时间',
                                                     `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
                                                     `update_time` datetime DEFAULT NULL COMMENT '更新时间',
                                                     `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                                     PRIMARY KEY (`relay_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='继电器信息表';

-- ------------------------------------------------------
-- 显控环境设备控制模块表 mpu_appliance_appdev_ctrl_module
-- ------------------------------------------------------
CREATE table if not exists `mpu_appliance_appdev_ctrl_module` (
                                                                  `dev_id` bigint(20) NOT NULL COMMENT '环境设备ID',
                                                                  `ctrl_module_id` bigint(20) NOT NULL COMMENT '控制模块ID',
                                                                  `ctrl_module_type` tinyint(4) DEFAULT NULL COMMENT '控制模块类型',
                                                                  `ctrl_module_param` int(11) DEFAULT NULL COMMENT '控制模块参数',
                                                                  PRIMARY KEY (`dev_id`,`ctrl_module_id`),
                                                                  CONSTRAINT `mpu_appliance_ctrl_module_to_dev` FOREIGN KEY (`dev_id`) REFERENCES `mpu_appliance_appdev` (`dev_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='环境设备所用控制模块表';

-- -----------------------------------------------
-- 显控环境设备分组表 mpu_appliance_group
-- -----------------------------------------------
CREATE table if not exists `mpu_appliance_group` (
                                                     `group_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '分组ID',
                                                     `parent_id` bigint(20) DEFAULT '0' COMMENT '父组ID',
                                                     `group_name` varchar(256) CHARACTER SET utf8 DEFAULT NULL COMMENT '分组名称',
                                                     `create_by` varchar(64) DEFAULT '' COMMENT '创建者',
                                                     `create_time` datetime DEFAULT NULL COMMENT '创建时间',
                                                     `update_by` varchar(64) DEFAULT '' COMMENT '更新者',
                                                     `update_time` datetime DEFAULT NULL COMMENT '更新时间',
                                                     `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
                                                     PRIMARY KEY (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='环境设备分组表';

-- ----------------------------------------------------
-- MCU平台信息
-- ----------------------------------------------------
create table if not exists `mcu_platform_cfg`
(
    `mcu_id` bigint unsigned not null auto_increment comment 'MCU配置ID',
    `mcu_alias` varchar(256) default '' comment 'MCU描述',
    `mcu_ver` int comment 'MCU版本枚举: 01=kd_v4.7; 10=kd_v5.0; 20=kd_v6.0; 21=kd_v6.1; 30=kd_v7.0; 40=kd_v8.0',
    `mcu_ip` varchar(50) comment 'MCU对接ip地址',
    `mcu_port` int comment 'MCU对接端口',
    `mcu_oauth_consumer_key` varchar(256) default '' comment 'MCU软件认证KEY',
    `mcu_oauth_consumer_secret` varchar(256) default '' comment 'MCU软件认证密钥',
    `mcu_SSL_cert_file` varchar(1024) default '' comment 'SSL证书文件(全路径)',
    `mcu_SSL_key_file` varchar(1024) default '' comment 'SSL密钥文件(全路径)',
    `last_err_msg` varchar(1024) default '' comment '对接MCU最后的报错信息,无报错为空',
    primary key (`mcu_id`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='MCU平台信息';

-- -----------------------------------------------------
-- MCU账户信息
-- -----------------------------------------------------
create table if not exists `mcu_platform_user`
(
    `mcu_id` bigint unsigned not null comment '所属MCU',
    `mcu_username` varchar(256) not null comment 'MCU用户名',
    `mcu_user_password` varchar(256) not null comment 'MCU用户密码',
    `mcu_user_domain_moid` varchar(256) not null comment 'MCU用户域moid',
    `last_err_msg` varchar(1024) default '' comment '对接MCU最后的报错信息,无报错为空',
    primary key (`mcu_id`,`mcu_username`),
    foreign key(mcu_id) references mcu_platform_cfg(mcu_id) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='MCU平台账户';

-- ------------------------------------------------------
-- 会议模板信息
-- ------------------------------------------------------
create table if not exists `mpu_meet_template`
(
    `meet_template_id` bigint unsigned not null auto_increment comment '显控分配的会议模板id',
    `mcu_id` bigint unsigned not null comment '所属MCU',
    `mcu_username` varchar(256) not null comment '所属MCU账户',
    `mcu_template_id` varchar(256) comment '模板id',
    `mcu_template_name` varchar(256) not null comment '模板名称',
    `mcu_conf_type` tinyint comment '会议类型：0=传统会议，1=端口会议，2=SFU纯转发会议',
    `mcu_encrypted_type` tinyint comment '传输加密了类型:0=不加密，2=AES加密，3=商密SM4，4=商密SM1',
    `max_join_mt` int comment '最大与会终端数：8=小型8方会议，32=32方会议，64=64方会议，192=大型192方会议',
    `number_of_invite_member` int comment '参会人数',
    `e164` varchar(64) comment '预分配会议号(MCU侧会议号)',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`meet_template_id`),
    foreign key(mcu_id) references mcu_platform_cfg(mcu_id) on update cascade on delete cascade
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='MCU会议模板信息';

-- ---------------------------------------------------------
-- 会议终端信息 -- 废弃
-- ---------------------------------------------------------
# create table if not exists `mpu_meet_terminal`
# (
#     `meet_terminal_id` bigint unsigned not null auto_increment comment '显控分配的会议终端id',
#     `mcu_id` bigint not null comment '所属MCU',
#     `meet_terminal_moid` varchar(64) default '' comment '会议终端moid',
#     `meet_terminal_name` varchar(64) default '' comment '会议终端名称',
#     `meet_terminal_E164` varchar(64) default '' comment '设备E164号',
#     `meet_terminal_online` varchar(64) default '' comment '设备是否在线',
#     `meet_terminal_ip` varchar(64) default '' comment '设备IP',
#     `domain_moid` varchar(64) default '' comment '所属用户域',
#     `deleted_at` datetime default null comment '删除时间',
#     primary key (`meet_terminal_id`),
#     foreign key(mcu_id) references mcu_platform_cfg(mcu_id) on update cascade on delete cascade
# ) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='MCU会议终端信息';


-- -----------------------------------------
-- 会议终端与显控通道关联配置
-- -----------------------------------------
create table if not exists `mcu_meet_terminal_to_channel`
(
    `cfg_id` bigint not null auto_increment comment '配置ID',
    `channel_id` varchar(64) not null comment '通道ID',
    `channel_alias` varchar(256) comment '模拟信号通道名称',
    `channel_state` int comment '模拟信号状态: 0离线(网络点位离线/模拟信号无信号) 1在线(网络点位在线/模拟信号有信号) 2已删除(板子拔走/网络点位删除) 3已收回(级联资源收回) 4不可用(级联未连接)',
    `meet_terminal_E164` varchar(64) default '' comment '设备E164号',
    `meet_terminal_name` varchar(256) default '' comment '会议终端名称',
    primary key (`cfg_id`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='会议终端与显控输出输出通道绑定关系';


-- --------------------------------------
-- 显控会议调度分组
-- --------------------------------------
create table if not exists `mpu_meet_dispatch_group`
(
    `group_id` bigint unique not null auto_increment comment '分组ID',
    `parent_id` bigint default '0' comment '父分组ID',
    `group_name` varchar(256) not null comment '分组名称',

    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`group_id`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_general_ci comment='显控会议调度分组';


-- --------------------------------------
-- 会议调度信息
-- --------------------------------------
create table if not exists `mpu_meet_dispatch`
(
    `meet_dispatch_id` bigint unique not null auto_increment comment '会议调度ID',
    `group_id` bigint default '0' comment '所属显控会议调度分组,不属于任何分组为0',
    `meet_dispatch_alias` varchar(256) not null comment '会议调度描述',

    `meet_template_id` bigint unsigned not null comment '关联会议模板ID',
    `matrix_scheme_id` bigint not null comment '关联矩阵调度预案ID',

    `meeting_id` varchar(64) comment '关联会议(MCU侧会议号)',
    `dispatch_state` bool comment '会议调度加载状态',
    `dispatch_load_time_ms` bigint comment '会议调度启动时刻(毫秒时间戳)',

    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    `deleted_at` datetime default null comment '删除时间',
    primary key (`meet_dispatch_id`)
) engine=innodb auto_increment=1025 default charset=utf8mb4 collate=utf8mb4_general_ci comment='显控会议调度信息';