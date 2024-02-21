set names utf8mb4;
set foreign_key_checks = 1;

# drop database if exists mpuaps;
# create database mpuaps;
# use mpuaps;

-- -------------------------------
-- casbin规则表 casbin_rule
-- -------------------------------
create table if not exists `casbin_rule`
(
    `ptype` varchar(10) default null,
    `v0` varchar(256) default null,
    `v1` varchar(256) default null,
    `v2` varchar(256) default null,
    `v3` varchar(256) default null,
    `v4` varchar(256) default null,
    `v5` varchar(256) default null
) engine=innodb default charset=utf8mb4 collate=utf8mb4_bin comment='casbin规则表';
alter table `casbin_rule` ADD INDEX ptype_index (`ptype`);
alter table `casbin_rule` ADD INDEX v0_index (`v0`);
alter table `casbin_rule` ADD INDEX v1_index (`v1`);


-- ----------------------------
-- 角色信息表 sys_role
-- ----------------------------
create table if not exists `sys_role`
(
    `role_id` bigint not null auto_increment comment '角色id',
    `role_name` varchar(128) not null comment '角色名称',
    `role_key` varchar(128) not null comment '角色标签',
    `role_sort` int not null comment '显示顺序',
    `role_data_scope` tinyint default 2 comment '数据范围标识(1：全部数据权限 2：自定数据权限)',
    `role_status` tinyint not null default 1 comment '角色状态(0停用 1正常)',
    `remark` varchar(512) default null comment '备注信息',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    primary key (`role_id`)
) engine=innodb auto_increment=11 default charset=utf8mb4 collate=utf8mb4_bin comment='角色信息表';


-- ----------------------------
-- 用户信息表 sys_user
-- ----------------------------
create table if not exists `sys_user`
(
    `user_id` bigint not null unique auto_increment comment '用户id',
    `login_name` varchar(128) not null unique comment '登录账号',
    `user_name` varchar(128) not null comment '用户昵称',
    `user_type` varchar(2) default '2' comment '用户类型(0系统管理员 1管理员 2操作员)|20221020实现角色管理功能废弃',
    `email` varchar(64) default '' comment '用户邮箱',
    `phone_number` varchar(11) default '' comment '手机号码',
    `sex` char(1) default '0' comment '用户性别(0保密 1男 2女)',
    `avatar` varchar(128) default '' comment '头像路径',
    `password` varchar(64) default '' comment '密码',
    `salt` varchar(20) default '' comment '加密盐',
    `status` char(1) default '0' comment '帐号状态（0正常 1停用 2未验证）',
    `login_ip` varchar(64) default '' comment '最后登陆ip',
    `login_date` datetime default null comment '最后登陆时间',
    `remark` varchar(512) default null comment '备注信息',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    primary key (`user_id`)
) engine=innodb auto_increment=10 default charset=utf8mb4 collate=utf8mb4_bin comment='用户信息表';


-- ----------------------------
-- 在线用户信息表 sys_user_online
-- ----------------------------
create table if not exists `sys_user_online`
(
    `token` varchar(255) not null default '' comment '用户会话token',
    `user_id` bigint not null comment '用户id',
    `login_name` varchar(128) default '' comment '登录账号',
    `dept_name` varchar(128) default '' comment '部门名称',
    `ipaddr` varchar(64) default '' comment '登录ip地址',
    `login_location` varchar(255) default '' comment '登录地点',
    `browser` varchar(50) default '' comment '浏览器类型',
    `os` varchar(50) default '' comment '操作系统',
    `status` varchar(10) default '' comment '在线状态on_line在线off_line离线',
    `start_timestamp` datetime default null comment '创建时间',
    `last_access_time` datetime default null comment '最后访问时间',
    `expire_time` int default '0' comment '超时时间，单位为分钟',
    primary key (`token`) using btree
) engine=innodb default charset=utf8mb4 collate=utf8mb4_bin comment='在线用户记录';


-- -------------------------------------
-- 系统参数信息表 sys_param
-- -------------------------------------
create table if not exists `sys_param`
(
    `param_id` bigint not null unique auto_increment comment '参数配置id',
    `param_name` varchar(128) default '' comment '参数名称',
    `param_key` varchar(128) default '' comment '参数键名',
    `param_value` varchar(512) default '' comment '参数值',
    `param_type` tinyint default '1' comment '参数类型: 1系统内置 2非系统内置参数',
    `remark` varchar(512) default null comment '备注信息',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    primary key (`param_id`)
) engine=innodb auto_increment=100 default charset=utf8mb4 collate=utf8mb4_bin comment='系统参数信息';


-- ----------------------------
-- 定时任务表 sys_job
-- ----------------------------
create table if not exists `sys_job`
(
    `job_id` bigint not null auto_increment comment '任务id',
    `job_name` varchar(64) not null default '' comment '任务名称',
    `job_params` varchar(255) default null comment '任务参数',
    `job_group` varchar(64) not null default 'default' comment '任务组名',
    `invoke_target` varchar(500) not null comment '调用目标字符串',
    `cron_expression` varchar(255) default '' comment 'cron执行表达式',
    `misfire_policy` varchar(20) default '1' comment '计划执行策略（1多次执行 2执行一次）',
    `concurrent` char(1) default '1' comment '是否并发执行（0允许 1禁止）',
    `status` char(1) default '0' comment '状态（0正常 1暂停）',
    `remark` varchar(500) default '' comment '备注信息',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    primary key (`job_id`,`job_name`,`job_group`)
) engine=innodb auto_increment=100 default charset=utf8mb4 collate=utf8mb4_bin comment='定时任务调度表';


-- ----------------------------
-- 定时任务日志表 sys_job_log
-- ----------------------------
create table if not exists `sys_job_log`
(
    `job_log_id` bigint not null auto_increment comment '任务日志id',
    `job_name` varchar(64) not null comment '任务名称',
    `job_group` varchar(64) not null comment '任务组名',
    `invoke_target` varchar(512) not null comment '调用目标字符串',
    `job_message` varchar(512) default null comment '日志信息',
    `status` char(1) default '0' comment '执行状态（0正常 1失败）',
    `exception_info` varchar(2000) default '' comment '异常信息',
    `create_time` datetime default null comment '创建时间',
    primary key (`job_log_id`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_bin comment='定时任务调度日志表';


-- -------------------------------
-- 系统访问记录表 sys_login_history
-- -------------------------------
create table if not exists `sys_login_history`
(
    `info_id` bigint not null auto_increment comment '访问id',
    `login_name` varchar(128) default '' comment '登录账号',
    `ipaddr` varchar(64) default '' comment '登录ip地址',
    `login_location` varchar(255) default '' comment '登录地点',
    `browser` varchar(50) default '' comment '浏览器类型',
    `os` varchar(64) default '' comment '操作系统',
    `status` char(24) default '0' comment '登录状态（0成功 其他错误代码）',
    `msg` varchar(255) default '' comment '提示消息',
    `login_time` datetime default null comment '访问时间',
    primary key (`info_id`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_bin comment='系统访问记录';


-- ----------------------------
-- 系统操作日志表 sys_oper_log
-- ----------------------------
create table if not exists `sys_oper_log`
(
    `oper_id` bigint not null auto_increment comment '日志主键',
    `title` varchar(64) default '' comment '模块标题',
    `business_type` int default '0' comment '业务类型（0其它 1新增 2修改 3删除 4查询 5授权 6导出 7导入 8强退 9清空）',
    `method` varchar(255) default '' comment '方法名称',
    `request_method` varchar(10) default '' comment '请求方式',
    `operator_type` int default '0' comment '操作类别（0其它 1前台用户 2级联操作）',
    `oper_name` varchar(128) default '' comment '操作人员',
    `dept_name` varchar(128) default '' comment '部门名称',
    `oper_url` varchar(255) default '' comment '请求url',
    `oper_ip` varchar(64) default '' comment '主机地址',
    `oper_location` varchar(255) default '' comment '操作地点',
    `oper_param` text comment '请求参数',
    `json_result` text comment '返回参数',
    `status` int default '0' comment '操作结果（0正常 -1异常）',
    `error_msg` varchar(2000) default '' comment '错误消息',
    `oper_time` datetime default null comment '操作时间',
    primary key (`oper_id`)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_bin comment='操作日志记录';


-- ------------------------------------------------
-- 菜单信息表 sys_menu
-- 当前版本不实现前端编辑修改菜单数据,菜单数据由后端生成
-- ------------------------------------------------
drop table if exists sys_menu;
create table if not exists `sys_menu`
(
    `id` bigint not null auto_increment comment '菜单id',
    `pid` bigint default null comment '父菜单id',
    `title` varchar(128) not null comment '菜单标题',
    `title_tag` varchar(64) not null comment '菜单标签',
    `type` tinyint unsigned default '1' comment '菜单类型(0目录 1菜单 2按钮)',
    `component` varchar(255) default null comment '前端组件路径',
    `icon` varchar(128) default '' comment '菜单图标',
    `router` varchar(128) default '' comment '后端路由规则',
    `rest_interface` varchar(2048) default '' comment '请求接口,一个菜单功能项可对应多个接口,多个接口 | 分隔',
    `weigh` int default '0' comment '菜单权重(菜单显示顺序)',
    `is_hide` tinyint unsigned default '0' comment '是否隐藏: 0不隐藏 1隐藏',
    `remark` varchar(512) default '' comment '备注',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    primary key (`id`) using btree
) engine=innodb auto_increment=1000 default charset=utf8mb4 collate=utf8mb4_bin comment='菜单信息表';



-- ----------------------------
-- 字典数据表 sys_dict_data
-- ----------------------------
drop table if exists sys_dict_data;
# create table if not exists `sys_dict_data`
# (
#     `dict_code` bigint not null auto_increment comment '字典编码',
#     `dict_sort` int default '0' comment '字典排序',
#     `dict_label` varchar(100) default '' comment '字典标签',
#     `dict_value` varchar(100) default '' comment '字典键值',
#     `dict_type` varchar(100) default '' comment '字典类型',
#     `status` char(1) default '0' comment '状态（0正常 1停用）',
#     `create_by` varchar(64) default '' comment '创建者',
#     `create_time` datetime default null comment '创建时间',
#     `update_by` varchar(64) default '' comment '更新者',
#     `update_time` datetime default null comment '更新时间',
#     `remark` varchar(500) default null comment '备注',
#     primary key (`dict_code`)
# ) engine=innodb auto_increment=100 default charset=utf8mb4 collate=utf8mb4_bin comment='字典数据表';


-- ----------------------------
-- 字典类型表 sys_dict_type
-- ----------------------------
drop table if exists sys_dict_type;
# create table if not exists `sys_dict_type`
# (
#     `dict_id` bigint not null auto_increment comment '字典主键',
#     `dict_name` varchar(100) default '' comment '字典名称',
#     `dict_type` varchar(100) default '' comment '字典类型',
#     `status` char(1) default '0' comment '状态（0正常 1停用）',
#     `create_by` varchar(64) default '' comment '创建者',
#     `create_time` datetime default null comment '创建时间',
#     `update_by` varchar(64) default '' comment '更新者',
#     `update_time` datetime default null comment '更新时间',
#     `remark` varchar(500) default null comment '备注',
#     primary key (`dict_id`),
#     unique key `dict_type` (`dict_type`)
# ) engine=innodb auto_increment=100 default charset=utf8mb4 collate=utf8mb4_bin comment='字典类型表';


-- ------------------------------------
-- 系统访问黑白名单信息 sys_black_list
-- ------------------------------------
create table if not exists `sys_black_list`
(
    `id` int not null auto_increment comment '配置ID',
    `type` tinyint unsigned default '1' comment '1黑名单 2白名单',
    `ip_addr` varchar(128) default '' comment '访问IP地址',
    `create_by` varchar(64) default '' comment '创建者',
    `create_time` datetime default null comment '创建时间',
    `update_by` varchar(64) default '' comment '更新者',
    `update_time` datetime default null comment '更新时间',
    primary key (id)
) engine=innodb default charset=utf8mb4 collate=utf8mb4_bin comment='系统访问黑白名单';