package service

import (
	"context"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/gogf/gf/crypto/gmd5"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	defCom "gowebserver/app/common/define"
	"gowebserver/app/common/errcode"
	"gowebserver/app/common/global"
	comSrv "gowebserver/app/common/service"
	"gowebserver/app/common/utils/convert"
	"gowebserver/app/common/utils/excel"
	"gowebserver/app/common/utils/page"
	"gowebserver/app/common/utils/random"
	"gowebserver/app/common/utils/valueCheck"
	"gowebserver/app/system/dao"
	"gowebserver/app/system/define"
	"gowebserver/app/system/model"
	"gowebserver/app/system/service/publish"
	"os"
	"strings"
	"sync"
	"time"
)

// ISysUser 用户管理业务接口
type ISysUser interface {
}

// User 用户管理服务对象
var User = &userServiceImpl{
	AvatarUploadPath:      `/public/upload/avatar/`,
	AvatarUploadUrlPrefix: `/public/upload/avatar/`,
}

type userServiceImpl struct {
	AvatarUploadPath      string     // 头像上传路径
	AvatarUploadUrlPrefix string     // 头像上传对应的URL前缀
	mutexUserUpdate       sync.Mutex //用户信息更新锁
}

// SysUserInit
// @summary 初始化默认账户信息
func (s *userServiceImpl) SysUserInit() {
	s.createDefaultUser()
}

// createDefaultUser
// @summary 创建默认账户
func (s *userServiceImpl) createDefaultUser() {
	record, err := dao.SysUser.Ctx(context.Background()).Where(g.Map{
		dao.SysUser.Columns.UserId: 1,
	}).FindOne()
	if err != nil {
		g.Log().Error(err)
	}
	if record == nil {
		//读取默认配置账户信息
		defUserLoginName := gconv.String(g.Cfg().Get("defaultUser.UserLoginName"))
		defUserName := gconv.String(g.Cfg().Get("defaultUser.UserNiceName"))
		defUserEmail := gconv.String(g.Cfg().Get("defaultUser.UserEmail"))
		defUserPhoneNumber := gconv.String(g.Cfg().Get("defaultUser.UserPhone"))
		defPasswd := gconv.String(g.Cfg().Get("defaultUser.UserPasswd"))
		if !s.CheckUserLoginNameComplexity(defUserLoginName) {
			defUserLoginName = "SysAdmin"
		}
		if !s.CheckPasswordComplexity(defPasswd) {
			defPasswd = "kedacom123!"
		}
		if defUserName == "" {
			defUserName = "系统管理员"
		}
		if defUserEmail == "" {
			defUserEmail = "https://www.kedacom.com/cn"
		}
		if defUserPhoneNumber == "" {
			defUserPhoneNumber = "4008282866"
		}
		//默认系统用户信息
		defSysUser := model.SysUser{
			UserId:      1,
			LoginName:   defUserLoginName,
			UserName:    defUserName,
			Email:       defUserEmail,
			PhoneNumber: defUserPhoneNumber,
			UserType:    "0",
			Sex:         "0",
			Avatar:      "",
			Password:    "",
			Salt:        "",
			Status:      "0",
			LoginIp:     "",
			LoginDate:   nil,
			CreateBy:    defUserLoginName,
			CreateTime:  gtime.Now(),
			UpdateBy:    "",
			UpdateTime:  nil,
			Remark:      "系统管理员",
		}
		defSysUser.Password, defSysUser.Salt = s.EncryptPassword(defSysUser.LoginName, defPasswd, "")
		_, err = dao.SysUser.Ctx(context.Background()).Save(defSysUser)
		if err != nil {
			g.Log().Error(err)
			return
		}
	}

	//默认账户绑定角色
	var roleIds []int64
	roleIds = append(roleIds, 1)
	if err = s.editUserRole(context.TODO(), roleIds, 1); err != nil {
		g.Log().Critical("默认系统管理员绑定角色失败", err)
	}
}

// IsDefSysUser
// @summary 是否为系统默认账户
// @return1
func (s *userServiceImpl) IsDefSysUser(uid int64) bool {
	return uid == 1
}

// IsSysAdmin
// @summary 判断用户是否为系统管理员
// @param1  uid int64 "用户ID"
// @return1 bool ""
func (s *userServiceImpl) IsSysAdmin(uid int64) bool {
	roleIds, err := s.GetUserRoleIds(context.TODO(), uid)
	if err != nil {
		g.Log().Error(err)
		return false
	}
	for _, roleId := range roleIds {
		if roleId == 1 {
			return true
		}
	}
	return false
}

// GetUserByLoginName
// @summary 根据用户登录名查询用户信息
func (s *userServiceImpl) GetUserByLoginName(ctx context.Context, loginName string) (user *model.SysUser, err error) {
	record, err := dao.SysUser.Ctx(ctx).Where(g.Map{
		dao.SysUser.Columns.LoginName: loginName,
	}).One()
	if err != nil {
		return nil, err
	}
	err = record.Struct(&user)
	return
}

// GetUser
// @summary 获取登录用户基本信息
func (s *userServiceImpl) GetUser(ctx context.Context) (*model.SysUser, error) {
	customCtx := Context.GetCtx(ctx)
	if customCtx != nil && customCtx.Uid != 0 {
		return s.GetInfo(ctx, customCtx.Uid)
	}
	return nil, gerror.NewCode(errcode.ErrSysUserNotLogin)
}

// GetUsers
// @summary 批量查询用户信息
func (s *userServiceImpl) GetUsers(ctx context.Context, userIds []int64) (users []model.SysUser, err error) {
	err = g.Try(func() {
		err = dao.SysUser.Ctx(ctx).WhereIn(dao.SysUser.Columns.UserId, userIds).FindScan(&users)
		valueCheck.ErrIsNil(ctx, err)
	})
	return
}

// GetProfile
// @summary 获取登录用户详细信息
//          1、用户信息
//          2、角色信息
func (s *userServiceImpl) GetProfile(ctx context.Context) (userInfo *model.SysUserExtend, err error) {
	err = g.Try(func() {
		//获取用户基本信息
		user, e := s.GetUser(ctx)
		valueCheck.ErrIsNil(ctx, e)
		valueCheck.ValueIsNil(user, "GetProfile failed")
		_ = gconv.Struct(user, &userInfo)

		//获取用户角色信息
		userInfo.Roles, err = s.GetUserRoles(ctx, user.UserId)
		valueCheck.ErrIsNil(ctx, err)
	})
	return
}

// GetList
// @summary 获取用户列表
func (s *userServiceImpl) GetList(ctx context.Context, param *define.UserApiSelectPageReq) *define.UserServiceList {
	if param == nil {
		return nil
	}
	m := dao.SysUser.Ctx(ctx).As("u")
	if param.LoginName != "" {
		//m = m.Where("u.login_name like ?", "%"+param.LoginName+"%")
		param.LoginName = strings.Replace(param.LoginName, "'", "''", -1)
		m = m.Where(fmt.Sprintf("locate('%s', %s) > 0", param.LoginName, dao.SysUser.Columns.LoginName))
	}
	if param.Phonenumber != "" {
		//m = m.Where("u.phone_number like ?", "%"+param.Phonenumber+"%")
		m = m.Where(fmt.Sprintf("locate('%s', %s) > 0", param.Phonenumber, dao.SysUser.Columns.PhoneNumber))
	}
	if param.Status != "" {
		m = m.Where("u.status = ?", param.Status)
	}
	if param.BeginTime != "" {
		m = m.Where("date_format(u.create_time,'%y%m%d') >= date_format(?,'%y%m%d')", param.BeginTime)
	}
	if param.EndTime != "" {
		m = m.Where("date_format(u.create_time,'%y%m%d') <= date_format(?,'%y%m%d')", param.EndTime)
	}
	total, err := m.Count()
	if err != nil {
		return nil
	}
	if param.OrderByColumn != "" && param.IsAsc != "" {
		param.OrderByColumn = gstr.CaseSnake(param.OrderByColumn)
		m = m.Order(param.OrderByColumn + " " + param.IsAsc)
	}
	pageInfo := page.CreatePaging(param.PageNum, param.PageSize, total)
	m = m.Limit(pageInfo.StartNum, pageInfo.PageSize)
	result := &define.UserServiceList{
		Page:  pageInfo.PageNum,
		Size:  pageInfo.PageSize,
		Total: pageInfo.Total,
	}
	if err = m.Scan(&result.List); err != nil {
		return nil
	}

	//获取用户角色ID
	for idx := range result.List {
		var roleIds []int64
		roleIds, err = s.GetUserRoleIds(ctx, result.List[idx].UserId)
		if err != nil || len(roleIds) == 0 {
			//用户未绑定角色
			g.Log().Warningf("用户【%d-%s】未绑定角色", result.List[idx].UserId, result.List[idx].LoginName)
			continue
		}
		result.List[idx].RoleId = roleIds[0]
	}

	return result
}

// SelectExportList
// @summary 导出用户信息
func (s *userServiceImpl) SelectExportList(param *define.UserApiSelectPageReq) (gdb.Result, error) {
	m := dao.SysUser.Ctx(context.Background()).As("u").
		Fields("u.login_name, u.user_name, u.email, u.phone_number, u.sex, u.status, u.create_by, u.create_time, u.remark")
	if param != nil {
		if param.LoginName != "" {
			m = m.Where("u.login_name like ?", "%"+param.LoginName+"%")
		}
		if param.Phonenumber != "" {
			m = m.Where("u.phone_number like ?", "%"+param.Phonenumber+"%")
		}
		if param.Status != "" {
			m = m.Where("u.status = ?", param.Status)
		}
		if param.BeginTime != "" {
			m = m.Where("date_format(u.create_time,'%y%m%d') >= date_format(?,'%y%m%d')", param.BeginTime)
		}
		if param.EndTime != "" {
			m = m.Where("date_format(u.create_time,'%y%m%d') <= date_format(?,'%y%m%d')", param.EndTime)
		}
	}
	return m.All()
}

// Create
// @summary 创建用户
func (s *userServiceImpl) Create(ctx context.Context, req *define.UserApiCreateReq) (int64, error) {
	if req.LoginName == "" {
		return 0, gerror.NewCode(errcode.ErrSysUserLoginNameInvalid)
	}

	//校验用户登录名复杂度
	if !s.CheckUserLoginNameComplexity(req.LoginName) {
		return 0, gerror.NewCode(errcode.ErrSysUserLoginNameInvalid)
	}
	//校验用户登录名是否重复
	if s.CheckLoginNameDuplicate(0, req.LoginName) {
		return 0, gerror.NewCode(errcode.ErrSysUserLoginNameDuplicate)
	}
	//校验用户登录面复杂度
	if !s.CheckPasswordComplexity(req.Password) {
		return 0, gerror.NewCode(errcode.ErrSysUserPassInvalid)
	}

	//if s.CheckPhoneUniqueAll(req.Phonenumber) {
	//	return 0, gerror.New("手机号码已经存在")
	//}
	//if s.CheckEmailUniqueAll(req.Email) {
	//	return 0, gerror.New("邮箱已经存在")
	//}

	var user model.SysUser
	user.Password, user.Salt = s.EncryptPassword(req.LoginName, req.Password, "")
	user.LoginName = req.LoginName
	user.CreateTime = gtime.Now()
	user.CreateBy = Context.GetUserName(ctx)

	var editReq *define.UserApiUpdateReq
	if err := gconv.Struct(req, &editReq); err != nil {
		return 0, err
	}
	return s.saveUpdate(ctx, &user, editReq)
}

// Update
// @summary 更新用户信息
func (s *userServiceImpl) Update(ctx context.Context, req *define.UserApiUpdateReq) (int64, error) {
	s.mutexUserUpdate.Lock()
	defer s.mutexUserUpdate.Unlock()

	//检查用户是否存在
	if !s.CheckUserExist(req.UserId) {
		return 0, gerror.NewCode(errcode.ErrSysUserInvalidUid)
	}

	//一个手机号/邮箱只能绑定一个账户
	//if s.CheckPhoneUnique(req.UserId, req.Phonenumber) {
	//	return 0, gerror.New("手机号码已绑定其他用户")
	//}
	//if s.CheckEmailUnique(req.UserId, req.Email) {
	//	return 0, gerror.New("邮箱已绑定其他用户")
	//}

	//用户信息
	user, err := s.GetInfo(ctx, req.UserId)
	if err != nil {
		return 0, err
	}

	//判断用户角色是否变更
	bUserRoleChange := false
	newRoleIds := strings.Split(req.RoleIds, ",")
	oldRoleIds, _ := s.GetUserRoleIds(ctx, req.UserId)
	if len(newRoleIds) == 0 || len(oldRoleIds) == 0 {
		return 0, gerror.NewCode(errcode.ErrCommonInvalidParameter)
	}
	if gconv.Int64(newRoleIds[0]) != oldRoleIds[0] {
		//目前用户只有一个角色
		bUserRoleChange = true
	}

	//修改用户登录名时需重新生成密码
	if user.LoginName != req.LoginName {
		//校验用户登录名是否重复
		if s.CheckLoginNameDuplicate(req.UserId, req.LoginName) {
			return 0, gerror.NewCode(errcode.ErrSysUserLoginNameDuplicate)
		}
		//校验用户登录名复杂度
		if !s.CheckUserLoginNameComplexity(req.LoginName) {
			return 0, gerror.NewCode(errcode.ErrSysUserLoginNameInvalid)
		}
		//校验密码
		if !s.CheckPasswordRight(user.LoginName, user.Password, user.Salt, req.Password) {
			return 0, gerror.NewCode(errcode.ErrSysUserPasswdErr)
		}
		//重新生成密码
		user.Password, user.Salt = s.EncryptPassword(req.LoginName, req.Password, user.Salt)
	}

	//禁止修改默认账户角色类型
	if s.IsDefSysUser(user.UserId) && user.UserType != req.UserType {
		return 0, gerror.NewCode(errcode.ErrSysUserModifyDefSysUserRoleForbid)
	}

	//保存用户信息
	if _, err = s.saveUpdate(ctx, user, req); err != nil {
		return 0, err
	}

	//用户角色变更通知在线账户权限变更
	if bUserRoleChange {
		s.NotifyOnlineUserPermissionUpdate(req.UserId)
	}
	return user.UserId, err
}

//NotifyOnlineUserPermissionUpdate
//@summary 通知在线用户权限变更(角色变更)
func (s *userServiceImpl) NotifyOnlineUserPermissionUpdate(userId int64) {
	onlineUserList := OnlineUser.GetAllOnlineUser()
	for _, onlineUser := range onlineUserList {
		if onlineUser.UserId == userId {
			publish.WsPublish.PublishMsgToSpecifiedClient(onlineUser.Token, defCom.ENotifyForceLogout.String(), define.OnlineUserForceLogoutNotify{
				Token:  onlineUser.Token,
				Reason: define.EForceLogOut_PermissionChange,
			})
		}
	}
}

// saveUpdate
// @summary 保存更新用户信息
// @param1  ctx context.Context "上下文信息"
// @param2  user *model.SysUser "用户基本信息"
// @param3  req *define.UserApiUpdateReq "要更新的用户数据"
// @return1 uid int64 "用户ID"
// @return2 err error "返回报错信息"
func (s *userServiceImpl) saveUpdate(ctx context.Context, user *model.SysUser, req *define.UserApiUpdateReq) (int64, error) {
	user.LoginName = req.LoginName
	user.UserName = req.UserName
	user.Email = req.Email
	user.PhoneNumber = req.Phonenumber
	user.Status = req.Status
	user.Sex = req.Sex
	user.Remark = req.Remark
	user.UserType = req.UserType
	if user.UserType == "" {
		user.UserType = "2"
	}

	//用户角色ID
	var roleIds []int64
	params := strings.Split(req.RoleIds, ",")
	for _, roleId := range params {
		//校验角色ID是否有效
		roleInfo, err := InterfaceSysRole().Get(ctx, gconv.Int64(roleId))
		if err != nil || roleInfo == nil {
			return user.UserId, gerror.NewCode(errcode.ErrSysUserBindRoleUnknow)
		}
		roleIds = append(roleIds, gconv.Int64(roleId))
	}

	//更新保存用户数据
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		return g.Try(func() {
			if user.UserId == 0 {
				//新增用户
				result, err := dao.SysUser.Ctx(ctx).TX(tx).Data(user).Save()
				valueCheck.ErrIsNil(ctx, err)

				user.UserId, err = result.LastInsertId()
				valueCheck.ErrIsNil(ctx, err)

				//绑定角色
				err = s.addUserRole(ctx, roleIds, user.UserId)
				valueCheck.ErrIsNil(ctx, err)

			} else {
				//更新用户信息
				_, err := dao.SysUser.Ctx(ctx).TX(tx).WherePri(user.UserId).Update(g.Map{
					dao.SysUser.Columns.LoginName:   user.LoginName,
					dao.SysUser.Columns.UserName:    user.UserName,
					dao.SysUser.Columns.Email:       user.Email,
					dao.SysUser.Columns.PhoneNumber: user.PhoneNumber,
					dao.SysUser.Columns.Status:      user.Status,
					dao.SysUser.Columns.Sex:         user.Sex,
					dao.SysUser.Columns.Remark:      user.Remark,
					dao.SysUser.Columns.UserType:    user.UserType,
					dao.SysUser.Columns.Password:    user.Password,
					dao.SysUser.Columns.Salt:        user.Salt,
					dao.SysUser.Columns.UpdateTime:  gtime.Now(),
					dao.SysUser.Columns.UpdateBy:    Context.GetUserName(ctx),
				})
				valueCheck.ErrIsNil(ctx, err)

				//绑定角色
				err = s.editUserRole(ctx, roleIds, user.UserId)
				valueCheck.ErrIsNil(ctx, err)
			}
		})
	})
	return user.UserId, err
}

// SaveMpuUser
// @summary 保存显控账户信息
//          1、重新分配账户ID
func (s *userServiceImpl) SaveMpuUser(ctx context.Context, user *model.SysUser) error {
	//校验用户登录名是否重复
	if s.CheckLoginNameDuplicate(0, user.LoginName) {
		return gerror.NewCode(errcode.ErrSysUserLoginNameDuplicate)
	}

	user.Password, user.Salt = s.EncryptPassword(user.LoginName, user.Password, "")
	user.CreateTime = gtime.Now()
	user.CreateBy = Context.GetUserName(ctx)
	if user.UserType == "" {
		user.UserType = "2"
	}
	result, err := dao.SysUser.Ctx(ctx).Data(user).Save()
	if err != nil {
		return err
	}
	user.UserId, err = result.LastInsertId()
	return err
}

// Delete
// @summary 删除用户信息
func (s *userServiceImpl) Delete(ctx context.Context, ids string) (int64, error) {
	idArr := convert.ToInt64Array(ids, ",")
	for _, uid := range idArr {
		//禁止删除默认账户
		if s.IsDefSysUser(uid) {
			return 0, gerror.NewCode(errcode.ErrSysUserDelDefSysUserForbid)
		}

		//用户登录使用中,禁止删除
		userInfo, _ := s.GetInfo(ctx, uid)
		if userInfo != nil && nil != OnlineUser.GetOnlineUserInfoByLoginName(userInfo.LoginName) {
			return 0, gerror.NewCode(errcode.ErrSysUserDelOnlineUserForbid)
		}
	}

	result, err := dao.SysUser.Ctx(context.Background()).
		Where(fmt.Sprintf("%s in(?)", dao.SysUser.Columns.UserId), idArr).
		Unscoped().Delete()
	if err != nil {
		return 0, err
	}

	//删除用户绑定角色
	for _, userId := range idArr {
		err = s.editUserRole(ctx, nil, userId)
		if err != nil {
			g.Log().Error(err)
		}
	}

	delUserNums, _ := result.RowsAffected()
	return delUserNums, nil
}

// GetInfo
// @summary 根据用户ID获取用户信息
func (s *userServiceImpl) GetInfo(ctx context.Context, uid int64) (user *model.SysUser, err error) {
	err = g.Try(func() {
		err = dao.SysUser.Ctx(context.Background()).Where(dao.SysUser.Columns.UserId, uid).FindScan(&user)
		valueCheck.ErrIsNil(ctx, err)
	})
	return
}

// ResetPassword
// @summary 用户重置密码
func (s *userServiceImpl) ResetPassword(ctx context.Context, req *define.UserApiResetPwdReq) (bool, error) {
	s.mutexUserUpdate.Lock()
	defer s.mutexUserUpdate.Unlock()

	//获取用户信息
	user, _ := s.GetInfo(ctx, req.UserId)
	if user == nil {
		return false, gerror.NewCode(errcode.ErrSysUserInvalidUid)
	}

	//校验旧密码
	if !s.CheckPasswordRight(user.LoginName, user.Password, user.Salt, req.OldPassword) {
		return false, gerror.NewCode(errcode.ErrSysUserUpdatePassInputOldPassErr)
	}

	//校验新密码复杂度
	if !s.CheckPasswordComplexity(req.NewPassword) {
		return false, gerror.NewCode(errcode.ErrSysUserPassInvalid)
	}

	//新密码加密
	user.Password, user.Salt = s.EncryptPassword(user.LoginName, req.NewPassword, user.Salt)

	//更新用户密码
	if _, err := dao.SysUser.Ctx(ctx).WherePri(user.UserId).Update(g.Map{
		dao.SysUser.Columns.Password: user.Password,
		dao.SysUser.Columns.Salt:     user.Salt,
	}); err != nil {
		return false, err
	}
	return true, nil
}

// UpdateProfile
// @summary 登录用户修改个人信息
func (s *userServiceImpl) UpdateProfile(ctx context.Context, profile *define.UserApiProfileReq) error {
	s.mutexUserUpdate.Lock()
	defer s.mutexUserUpdate.Unlock()

	//获取登录用户信息
	user, err := s.GetUser(ctx)
	if err != nil {
		return gerror.NewCode(errcode.ErrSysUserNotLogin)
	}

	//一个手机号/邮箱只能绑定一个账户
	//if s.CheckPhoneUnique(user.UserId, profile.PhoneNumber) {
	//	return gerror.NewCode(errcode.ErrSysUserPhonenumberDuplicate)
	//}
	//if s.CheckEmailUnique(user.UserId, profile.Email) {
	//	return gerror.NewCode(errcode.ErrSysUserEmailDuplicate)
	//}

	if profile.UserName != "" {
		user.UserName = profile.UserName
	}
	if profile.Email != "" {
		user.Email = profile.Email
	}
	if profile.PhoneNumber != "" {
		user.PhoneNumber = profile.PhoneNumber
	}
	if profile.Sex != "" {
		user.Sex = profile.Sex
	}

	_, err = dao.SysUser.Ctx(ctx).WherePri(user.UserId).Update(g.Map{
		dao.SysUser.Columns.UserName:    user.UserName,
		dao.SysUser.Columns.Email:       user.Email,
		dao.SysUser.Columns.PhoneNumber: user.PhoneNumber,
		dao.SysUser.Columns.Sex:         user.Sex,
		dao.SysUser.Columns.UpdateTime:  gtime.Now(),
		dao.SysUser.Columns.UpdateBy:    Context.GetUserName(ctx),
	})
	if err != nil {
		return err
	}
	return nil
}

// UpdateAvatar
// @summary 登录用户修改个人头像
func (s *userServiceImpl) UpdateAvatar(ctx context.Context, r *define.UserApiAvatarUploadReq) (string, error) {
	user, err := s.GetUser(ctx)
	if err != nil {
		return "", gerror.NewCode(errcode.ErrSysUserNotLogin)
	}
	curDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	saveDir := curDir + s.AvatarUploadPath + gconv.String(user.UserId) + "/"
	filename, err := r.AvatarFile.Save(saveDir, true)
	if err != nil {
		return "", err
	}
	avatar := s.AvatarUploadPath + gconv.String(user.UserId) + "/" + filename
	if avatar != "" {
		user.Avatar = avatar
	}
	if _, err = dao.SysUser.Ctx(ctx).Data(user).Save(); err != nil {
		return "", err
	}
	return avatar, nil
}

// UpdatePassword
// @summary 登录用户修改个人密码
func (s *userServiceImpl) UpdatePassword(ctx context.Context, req *define.UserApiReSetPasswordReq) error {
	s.mutexUserUpdate.Lock()
	defer s.mutexUserUpdate.Unlock()

	user, err := s.GetUser(ctx)
	if err != nil {
		return gerror.NewCode(errcode.ErrSysUserNotLogin)
	}
	if strings.EqualFold(req.OldPassword, "") {
		return gerror.NewCode(errcode.ErrCommonInvalidParameter)
	}
	if strings.EqualFold(req.NewPassword, "") {
		return gerror.NewCode(errcode.ErrCommonInvalidParameter)
	}
	if strings.EqualFold(req.NewPassword, req.OldPassword) {
		return gerror.NewCode(errcode.ErrSysUserUpdateNewPassSameAsOldPass)
	}
	if !strings.EqualFold(req.NewPassword, req.Confirm) {
		return gerror.NewCode(errcode.ErrSysUserUpdatePassNotSameAsConfirm)
	}

	//校验新密码复杂度
	if !s.CheckPasswordComplexity(req.NewPassword) {
		return gerror.NewCode(errcode.ErrSysUserPassInvalid)
	}

	//校验旧密码
	if !s.CheckPasswordRight(user.LoginName, user.Password, user.Salt, req.OldPassword) {
		return gerror.NewCode(errcode.ErrSysUserUpdatePassInputOldPassErr)
	}

	//新密码加密
	user.Password, user.Salt = s.EncryptPassword(user.LoginName, req.NewPassword, user.Salt)

	//更新用户密码
	if _, err = dao.SysUser.Ctx(ctx).WherePri(user.UserId).Update(g.Map{
		dao.SysUser.Columns.Password: user.Password,
		dao.SysUser.Columns.Salt:     user.Salt,
	}); err != nil {
		return err
	}
	return nil
}

// UpdateLastLoginInfo
// @summary 更新用户最后登录点
func (s *userServiceImpl) UpdateLastLoginInfo(userId int64, lastLoginIp string) {
	_, _ = dao.SysUser.Ctx(context.Background()).Data(g.Map{
		dao.SysUser.Columns.LoginIp:   lastLoginIp,
		dao.SysUser.Columns.LoginDate: gtime.Now(),
	}).Where(dao.SysUser.Columns.UserId, userId).Update()
}

// ChangeStatus
// @summary 修改用户状态
func (s *userServiceImpl) ChangeStatus(userId int64, status string) error {
	if s.IsDefSysUser(userId) {
		return gerror.NewCode(errcode.ErrSysUserLockDefSysUserForbid)
	}
	if _, err := dao.SysUser.Ctx(context.Background()).Where(dao.SysUser.Columns.UserId, userId).Data(g.Map{
		dao.SysUser.Columns.Status: status,
	}).Update(); err != nil {
		return err
	}
	return nil
}

// Export
// @summary 导出excel
func (s *userServiceImpl) Export(param *define.UserApiSelectPageReq) (string, error) {
	userList, err := s.SelectExportList(param)
	if err != nil {
		return "", err
	}
	head := []string{"用户名", "昵称", "Email", "电话号码", "性别", "状态", "创建人", "创建时间", "备注"}
	key := []string{"login_name", "user_name", "email", "phone_number", "sex", "status", "create_by", "create_time", "remark"}
	url, err := excel.DownloadExcel(head, key, userList)
	if err != nil {
		return "", err
	}
	return url, nil
}

// CheckUserExist
// @summary 检查用户是否存在, 存在返回true, 不存在返回false
func (s *userServiceImpl) CheckUserExist(userId int64) bool {
	if i, err := dao.SysUser.Ctx(context.Background()).FindCount(g.Map{
		dao.SysUser.Columns.UserId: userId,
	}); err != nil {
		return false
	} else {
		return i > 0
	}
}

// CheckLoginNameDuplicate
// @summary 检查登陆名是否存在, 存在返回true, 不存在返回false
func (s *userServiceImpl) CheckLoginNameDuplicate(userId int64, loginName string) bool {
	if i, err := dao.SysUser.Ctx(context.Background()).
		Where(dao.SysUser.Columns.LoginName, loginName).
		WhereNot(dao.SysUser.Columns.UserId, userId).
		FindCount(); err != nil {
		return false
	} else {
		return i > 0
	}
}

// CheckEmailUnique
// @summary 检查邮箱是否已使用
func (s *userServiceImpl) CheckEmailUnique(userId int64, email string) bool {
	if i, err := dao.SysUser.Ctx(context.Background()).FindCount(g.Map{
		dao.SysUser.Columns.Email:                        email,
		fmt.Sprintf("%s <>", dao.SysUser.Columns.UserId): userId,
	}); err != nil {
		return false
	} else {
		return i > 0
	}
}

// CheckEmailUniqueAll
// @summary 检查邮箱是否存在,存在返回true,否则false
func (s *userServiceImpl) CheckEmailUniqueAll(email string) bool {
	if i, err := dao.SysUser.Ctx(context.Background()).FindCount(g.Map{
		dao.SysUser.Columns.Email: email,
	}); err != nil {
		return false
	} else {
		return i > 0
	}
}

// CheckPhoneUnique
// @summary 检查手机号是否已使用,存在返回true,否则false
func (s *userServiceImpl) CheckPhoneUnique(userId int64, phone string) bool {
	if i, err := dao.SysUser.Ctx(context.Background()).FindCount(g.Map{
		dao.SysUser.Columns.PhoneNumber:                  phone,
		fmt.Sprintf("%s <>", dao.SysUser.Columns.UserId): userId,
	}); err != nil {
		return false
	} else {
		return i > 0
	}
}

// CheckPhoneUniqueAll
// @summary 检查手机号是否已使用 ,存在返回true,否则false
func (s *userServiceImpl) CheckPhoneUniqueAll(phone string) bool {
	if i, err := dao.SysUser.Ctx(context.Background()).FindCount(g.Map{
		dao.SysUser.Columns.PhoneNumber: phone,
	}); err != nil {
		return false
	} else {
		return i > 0
	}
}

// CheckSysAdminPasswd
// @summary 校验系统管理员密码
func (s *userServiceImpl) CheckSysAdminPasswd(inputPasswd string) bool {
	sysAdminUser, _ := s.GetInfo(context.TODO(), 1)
	return s.CheckPasswordRight(sysAdminUser.LoginName, sysAdminUser.Password, sysAdminUser.Salt, inputPasswd)
}

// CheckPasswordRight
// @summary 校验密码是否正确
// @param1  loginName   string 登录名
// @param2  passwd      string 正确密码
// @param3  salt        string 加密盐
// @param4  inputPasswd string 输入密码
// @return1 true = 正确; false = 错误;
func (s *userServiceImpl) CheckPasswordRight(loginName, passwd, salt string, inputPasswd string) bool {
	pwdEncryptKey := loginName + inputPasswd + salt
	encryptPasswd := gmd5.MustEncryptString(pwdEncryptKey)
	return strings.EqualFold(passwd, encryptPasswd)
}

// EncryptPassword
// @summary 密码加密
// @param1  loginName   string 登录名
// @param2  inputPasswd string 输入密码
// @return1 encryptPasswd string 加密后的密码
// @return2 salt 加密盐
func (s *userServiceImpl) EncryptPassword(loginName, inputPasswd, inputSalt string) (encryptPasswd string, salt string) {
	if inputSalt == "" {
		inputSalt = random.GenerateSubId(6)
	}
	salt = inputSalt
	pwdEncryptKey := loginName + inputPasswd + salt
	encryptPasswd = gmd5.MustEncryptString(pwdEncryptKey)
	return
}

// CheckPasswordComplexity
// @summary 校验密码复杂度
//          用户密码复杂度要求(参考《科达产品开发安全规范》)
//          1、组成字符: 数字、字母、特殊字符、大小写敏感
//          2、长度要求: 8个字符及以上
//          3、密码强度:
//                     弱-仅包含数字、字母中的一种或两种  -- 禁止使用弱口令
//                     中-包含数字、字母和任一特殊字符
//                     强-包含数字、字母和多个特殊字符
//          4、不能包含的特殊字符: 单引号(')、双引号(")、分号(;)、逗号(,)、星号(*)、问号(?)、双横杠（--）、等号（=）
//          正则表达式说明:
//          (?=.*?[A-Za-z])      表示至少有一个字符是字母
//          (?=.*?[0-9])         表示至少有一个字符是数字
//			(?=.*[._~!@#$^&+])   表示至少有一个字符是特殊字符(中括号中列出的字符)
//          [A-Za-z0-9._~!@#$^&+]表示密码只能是中括号中列出来的字符
//          {8,}                 表示必须在8位以上
func (s *userServiceImpl) CheckPasswordComplexity(passwd string) bool {
	regexPattern := "^(?=.*?[A-Za-z])(?=.*?[0-9])(?=.*[._~!@#$^&+])[A-Za-z0-9._~!@#$^&+]{8,}$"

	//由于golang不支持预查之类的正则匹配规则,regex包不支持?=之类的格式,复杂的正则匹配难以实现
	//使用开源包regex2
	/*
		if err := gregex.Validate(regexPattern); err != nil {
			g.Log().Errorf("invalidate regexPattern: %s\n%v", regexPattern, err)
			return false
		}
		return gregex.IsMatchString(regexPattern, passwd)
	*/

	reg, _ := regexp2.Compile(regexPattern, 0)
	m, _ := reg.FindStringMatch(passwd)
	return m != nil
}

// CheckUserLoginNameComplexity
// @summary 校验用户登录名复杂度
//          用户密码复杂度要求(参考《科达产品开发安全规范》)
//          1、组成字符: 数字、字母、大小写敏感
//          2、长度要求: 5-30个字符
func (s *userServiceImpl) CheckUserLoginNameComplexity(loginName string) bool {
	regexPattern := "^[a-zA-Z0-9_-]{5,30}$"
	reg, _ := regexp2.Compile(regexPattern, 0)
	m, _ := reg.FindStringMatch(loginName)
	return m != nil
}

// CheckCanSyncMpuUser
// @summary 校验是否可以直接同步原显控平台账户数据
//          1、只有代理侧不存在账户数据时返回true
func (s *userServiceImpl) CheckCanSyncMpuUser() bool {
	userCount, err := dao.SysUser.Ctx(context.TODO()).Count()
	if err != nil || userCount > 2 {
		return false
	}
	return true
}

// SetPasswordCounts
// @summary 记录密码尝试次数
func (s *userServiceImpl) SetPasswordCounts(loginName string) int {
	curTimes := 0

	// 读取已输出错误次数
	curTimeObj, _ := gcache.Get(global.CacheUserNoPassTimePrefix + loginName)
	if curTimeObj != nil {
		curTimes = gconv.Int(curTimeObj)
	}

	// 输出错误次数+1 (缓存时间1min, 1min后自动重置尝试次数)
	curTimes = curTimes + 1
	err := gcache.Set(global.CacheUserNoPassTimePrefix+loginName, curTimes, 1*time.Minute)
	if err != nil {
		g.Log().Error(err)
	}

	// 超过最大次数, 锁定用户
	if curTimes >= global.UserPasswdErrCountMax {
		s.LoginLock(loginName)
	}
	return curTimes
}

// RemovePasswordCounts
// @summary 移除密码错误次数
func (s *userServiceImpl) RemovePasswordCounts(loginName string) {
	_, err := gcache.Remove(global.CacheUserNoPassTimePrefix + loginName)
	if err != nil {
		g.Log().Error(err)
		return
	}
}

// LoginLock
// @summary 锁定账号30分钟
func (s *userServiceImpl) LoginLock(loginName string) {
	err := gcache.Set(global.CacheUserLockPrefix+loginName, true, global.UserLockTime*time.Minute)
	if err != nil {
		g.Log().Error(err)
		return
	}
}

// LoginUnLock
// @summary 移除登录锁定状态
func (s *userServiceImpl) LoginUnLock(loginName string) {
	_, err := gcache.Remove(global.CacheUserLockPrefix + loginName)
	if err != nil {
		g.Log().Error(err)
	}
}

// CheckLoginLockState
// @summary 校验是否锁定
func (s *userServiceImpl) CheckLoginLockState(loginName string) bool {
	result := false
	rs, _ := gcache.Get(global.CacheUserLockPrefix + loginName)
	if rs != nil {
		result = true
	}
	return result
}

// CheckAccessConnectionsNumExceedMax
// @summary 判断访问人数是否达到限制上限
// @return1 bool "true=达到上限 false=未达到上限"
func (s *userServiceImpl) CheckAccessConnectionsNumExceedMax() bool {
	onlineUserList := OnlineUser.GetAllOnlineUser()

	accessConnLimitParam, _ := InterfaceSysParamMgr().GetSysParamCfg(context.TODO(), global.SysParamKey_UserLoginAccessConnectionsNumMax)
	if accessConnLimitParam != nil {
		return len(onlineUserList) >= gconv.Int(accessConnLimitParam.ParamValue)
	}
	return true
}

// addUserRole
// @summary 绑定用户角色
func (s *userServiceImpl) addUserRole(ctx context.Context, roleIds []int64, userId int64) (err error) {
	err = g.Try(func() {
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)
		for _, v := range roleIds {
			_, e = enforcer.AddGroupingPolicy(fmt.Sprintf("u_%d", userId), gconv.String(v))
			valueCheck.ErrIsNil(ctx, e)
		}
	})
	return
}

// editUserRole
// @summary 修改用户绑定角色
func (s *userServiceImpl) editUserRole(ctx context.Context, roleIds []int64, userId int64) (err error) {
	err = g.Try(func() {
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)

		//删除用户旧角色信息
		_, e = enforcer.RemoveFilteredGroupingPolicy(0, fmt.Sprintf("u_%d", userId))
		valueCheck.ErrIsNil(ctx, e)

		//绑定用户角色信息
		for _, v := range roleIds {
			_, err = enforcer.AddGroupingPolicy(fmt.Sprintf("u_%d", userId), gconv.String(v))
			valueCheck.ErrIsNil(ctx, err)
		}
	})
	return
}

// GetUserRoleIds
// @summary 获取用户关联角色ID
func (s *userServiceImpl) GetUserRoleIds(ctx context.Context, userId int64) (roleIds []int64, err error) {
	enforcer, e := comSrv.CasbinEnforcer(ctx)
	if e != nil {
		err = e
		return
	}
	//查询关联角色规则
	groupPolicy := enforcer.GetFilteredGroupingPolicy(0, fmt.Sprintf("u_%d", userId))
	if len(groupPolicy) > 0 {
		roleIds = make([]int64, len(groupPolicy))
		//得到角色id的切片
		for k, v := range groupPolicy {
			roleIds[k] = gconv.Int64(v[1])
		}
	}
	return
}

// GetUserRoles
// @summary 获取用户关联角色信息
func (s *userServiceImpl) GetUserRoles(ctx context.Context, userId int64) (roles []model.SysRole, err error) {
	err = g.Try(func() {
		var roleIds []int64
		roleIds, err = s.GetUserRoleIds(ctx, userId)
		valueCheck.ErrIsNil(ctx, err)

		roles, err = InterfaceSysRole().GetRoles(ctx, roleIds)
		valueCheck.ErrIsNil(ctx, err)
	})
	return
}

// CheckUserRole
// @summary 校验用户是否关联到角色
func (s *userServiceImpl) CheckUserRole(ctx context.Context, userId, roleId int64) bool {
	enforcer, err := comSrv.CasbinEnforcer(ctx)
	if err != nil {
		g.Log().Error(err)
		return false
	}
	return enforcer.HasGroupingPolicy(fmt.Sprintf("u_%d", userId), gconv.String(roleId))
}

// GetUserIdsByRoleId
// @summary 获取角色关联用户ID
func (s *userServiceImpl) GetUserIdsByRoleId(ctx context.Context, roleId int64) (userIds []int64, err error) {
	enforcer, e := comSrv.CasbinEnforcer(ctx)
	if e != nil {
		err = e
		return
	}
	//查询关联角色规则
	groupPolicy := enforcer.GetFilteredGroupingPolicy(1, fmt.Sprintf("%d", roleId))
	if len(groupPolicy) > 0 {
		userIds = make([]int64, len(groupPolicy))
		//得到用户id的切片
		for k, v := range groupPolicy {
			params := strings.Split(v[0], "_")
			if len(params) == 2 {
				userIds[k] = gconv.Int64(params[1])
			} else {
				g.Log().Warning("用户角色权限规则错误:", v)
			}
		}
	}
	return
}

// GetUsersByRoleId
// @summary 获取角色关联账户信息
func (s *userServiceImpl) GetUsersByRoleId(ctx context.Context, roleId int64) (users []model.SysUser, err error) {
	err = g.Try(func() {
		var userIds []int64
		userIds, err = s.GetUserIdsByRoleId(ctx, roleId)
		valueCheck.ErrIsNil(ctx, err)

		users, err = s.GetUsers(ctx, userIds)
		valueCheck.ErrIsNil(ctx, err)
	})
	return
}
