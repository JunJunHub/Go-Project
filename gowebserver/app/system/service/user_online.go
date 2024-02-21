package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"gowebserver/app/common/errcode"
	"gowebserver/app/common/global"
	"gowebserver/app/common/utils/page"
	"gowebserver/app/system/dao"
	"gowebserver/app/system/define"
	"gowebserver/app/system/model"
	"strings"
)

var OnlineUser = &onlineUserService{}

type onlineUserService struct{}

// GetList 根据条件分页查询数据
func (s *onlineUserService) GetList(param *define.OnlineUserApiSelectPageReq) *define.OnlineUserServiceList {
	m := dao.SysUserOnline.Ctx(context.Background()).As("t")
	if param != nil {
		if param.Token != "" {
			m = m.Where("t.token = ?", param.Token)
		}

		if param.LoginName != "" {
			m = m.Where("t.login_name like ?", "%"+param.LoginName+"%")
		}

		if param.DeptName != "" {
			m = m.Where("t.dept_name like ?", "%"+param.DeptName+"%")
		}

		if param.Ipaddr != "" {
			m = m.Where("t.ipaddr = ?", param.Ipaddr)
		}

		if param.LoginLocation != "" {
			m = m.Where("t.login_location = ?", param.LoginLocation)
		}

		if param.Browser != "" {
			m = m.Where("t.browser = ?", param.Browser)
		}

		if param.Os != "" {
			m = m.Where("t.os = ?", param.Os)
		}

		if param.Status != "" {
			m = m.Where("t.status = ?", param.Status)
		}

		if param.BeginTime != "" {
			m = m.Where("date_format(t.create_time,'%y%m%d') >= date_format(?,'%y%m%d') ", param.BeginTime)
		}

		if param.EndTime != "" {
			m = m.Where("date_format(t.create_time,'%y%m%d') <= date_format(?,'%y%m%d') ", param.EndTime)
		}
	}
	total, err := m.Count()
	if err != nil {
		return nil
	}
	pageInfo := page.CreatePaging(param.PageNum, param.PageSize, total)
	m = m.Limit(pageInfo.StartNum, pageInfo.PageSize)
	if param.OrderByColumn != "" {
		m = m.Order(param.OrderByColumn + " " + param.IsAsc)
	}
	result := &define.OnlineUserServiceList{
		Page:  pageInfo.PageNum,
		Size:  pageInfo.PageSize,
		Total: pageInfo.Total,
	}
	if err = m.Scan(&result.List); err != nil {
		return nil
	}
	return result
}

//UpdateOnlineUserState 更新在线用户状态
func (s *onlineUserService) UpdateOnlineUserState(token string) error {
	_, err := dao.SysUserOnline.Ctx(context.Background()).Data(g.Map{
		dao.SysUserOnline.Columns.Status: global.UserStatusOffline,
	}).Where(dao.SysUserOnline.Columns.Token, token).Update()
	if err != nil {
		return err
	}
	return nil
}

//GetOnlineUserInfoByLoginName 获取登录用户登录信息
func (s *onlineUserService) GetOnlineUserInfoByLoginName(loginName string) []model.SysUserOnline {
	record, _ := dao.SysUserOnline.Ctx(context.Background()).
		Where(dao.SysUserOnline.Columns.LoginName, loginName).
		Order("start_timestamp asc").FindAll()

	var onlineUsers []model.SysUserOnline
	err := record.Structs(&onlineUsers)
	if err != nil {
		return nil
	}
	return onlineUsers
}

//GetOnlineUserInfoByLoginToken 获取登录用户登录信息
func (s *onlineUserService) GetOnlineUserInfoByLoginToken(loginToken string) *model.SysUserOnline {
	record, _ := dao.SysUserOnline.Ctx(context.Background()).
		Where(dao.SysUserOnline.Columns.Token, loginToken).
		Order("start_timestamp asc").FindOne()

	var onlineUser *model.SysUserOnline
	err := record.Struct(&onlineUser)
	if err != nil {
		return nil
	}
	return onlineUser
}

func (s *onlineUserService) GetAllOfflineUser() []model.SysUserOnline {
	record, _ := dao.SysUserOnline.Ctx(context.Background()).
		Where(dao.SysUserOnline.Columns.Status, global.UserStatusOffline).
		Order("start_timestamp asc").FindAll()

	var users []model.SysUserOnline
	err := record.Structs(&users)
	if err != nil {
		return nil
	}
	return users
}

func (s *onlineUserService) GetAllOnlineUser() []model.SysUserOnline {
	record, _ := dao.SysUserOnline.Ctx(context.Background()).
		Where(dao.SysUserOnline.Columns.Status, global.UserStatusOnline).
		Order("start_timestamp asc").FindAll()

	var users []model.SysUserOnline
	err := record.Structs(&users)
	if err != nil {
		return nil
	}
	return users
}

//GetOnlineUserAcceptLanguage 获取登录客户端接收语言类型
func (s *onlineUserService) GetOnlineUserAcceptLanguage(loginToken string) string {
	onlineUser := s.GetOnlineUserInfoByLoginToken(loginToken)
	if onlineUser == nil {
		return "zh_CN"
	}
	params := strings.Split(onlineUser.Os, ":")
	if len(params) > 1 {
		return params[1]
	}
	return "zh_CN"
}

//AddOnlineUser
//@summary 用户上线,创建用户在线状态保活协程
//         1、用户上线必须同时创建ws订阅，如果ws订阅链路断开，强制用户登出
func (s *onlineUserService) AddOnlineUser(onlineUser model.SysUserOnline) {
	_, _ = dao.SysUserOnline.Ctx(context.Background()).Save(onlineUser)
}

//DeleteOnlineByToken 删除用户在线状态操作
func (s *onlineUserService) DeleteOnlineByToken(token string) error {
	result, err := dao.SysUserOnline.Ctx(context.Background()).Delete(dao.SysUserOnline.Columns.Token, token)
	if err != nil {
		g.Log().Error(err)
		return err
	}
	delNum, _ := result.RowsAffected()
	if delNum == 0 {
		g.Log().Errorf("删除在线用户信息失败,token不存在[%s]", token)
		return gerror.NewCode(errcode.ErrCommonOperationFailed, errcode.ErrCommonOperationFailed.Message()+"token unknown")
	} else {
		g.Log().Noticef("已删除在线用户信息,token[%s]\n %v", token, g.Log().GetStack(1))
	}
	return nil
}

//Clean 清空在线用户记录,系统重启时调用
func (s *onlineUserService) Clean() int64 {
	result, err := dao.SysUserOnline.Ctx(context.Background()).Delete(fmt.Sprintf("%s != ?", dao.SysUserOnline.Columns.LoginName), "")
	if err != nil {
		return 0
	}
	nums, _ := result.RowsAffected()

	return nums
}
