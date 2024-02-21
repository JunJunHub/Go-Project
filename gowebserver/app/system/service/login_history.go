package service

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/mssola/user_agent"
	"gowebserver/app/common/utils/convert"
	"gowebserver/app/common/utils/excel"
	"gowebserver/app/common/utils/ip"
	"gowebserver/app/common/utils/page"
	"gowebserver/app/system/dao"
	"gowebserver/app/system/define"
	"gowebserver/app/system/model"

	"context"
	"fmt"
)

var LoginHistory = &loginHistoryService{}

type loginHistoryService struct{}

// GetList 根据条件分页查询列表
func (s *loginHistoryService) GetList(ctx context.Context, param *define.LoginHistoryApiSelectPageReq) *define.LoginHistoryServiceList {

	m := dao.SysLoginHistory.Ctx(ctx).As("t")

	if param != nil {
		if param.LoginName != "" {
			m = m.Where("login_name like ?", "%"+param.LoginName+"%")
		}

		if param.Ipaddr != "" {
			m = m.Where("ipaddr like ?", "%"+param.Ipaddr+"%")
		}

		if param.Status != "" {
			m = m.Where("status = ?", param.Status)
		}

		if param.BeginTime != "" {
			m = m.Where("date_format(login_time,'%y%m%d') >= date_format(?,'%y%m%d')", param.BeginTime)
		}
		if param.EndTime != "" {
			m = m.Where("date_format(login_time,'%y%m%d') <= date_format(?,'%y%m%d')", param.EndTime)
		}
	}
	total, err := m.Count()
	if err != nil {
		return nil
	}
	page := page.CreatePaging(param.PageNum, param.PageSize, total)
	m = m.Fields("info_id,login_name,ipaddr,login_location,browser,os,status,msg,login_time")
	m = m.Order("login_time desc")
	m = m.Limit(page.StartNum, page.PageSize)
	result := &define.LoginHistoryServiceList{
		Page:  page.PageNum,
		Size:  page.PageSize,
		Total: page.Total,
	}
	if err = m.Scan(&result.List); err != nil {
		return nil
	}
	return result
}

// Create 新增登录历史记录
func (s *loginHistoryService) Create(ctx context.Context, status, username, ipaddr, userAgent, msg string) {
	var loginHistory model.SysLoginHistory
	loginHistory.Status = status
	loginHistory.LoginName = username
	loginHistory.Ipaddr = ipaddr
	ua := user_agent.New(userAgent)
	os := ua.OS()
	browser, _ := ua.Browser()
	loginLocation := ip.GetCityByIp(ipaddr)
	loginHistory.Os = os
	loginHistory.Browser = browser
	loginHistory.LoginTime = gtime.Now()
	loginHistory.LoginLocation = loginLocation
	loginHistory.Msg = msg
	if _, err := dao.SysLoginHistory.Ctx(ctx).Insert(loginHistory); err != nil {
		g.Log().Line().Warning(err)
	}
}

// Delete 批量删除记录
func (s *loginHistoryService) Delete(ctx context.Context, ids string) int64 {
	idArr := convert.ToInt64Array(ids, ",")
	result, err := dao.SysLoginHistory.Ctx(ctx).Delete(fmt.Sprintf("%s in(?)", dao.SysLoginHistory.Columns.InfoId), idArr)
	if err != nil {
		return 0
	}
	nums, _ := result.RowsAffected()
	return nums
}

// Clean 清空记录
func (s *loginHistoryService) Clean(ctx context.Context) int64 {
	result, err := dao.SysLoginHistory.Ctx(ctx).Delete(fmt.Sprintf("%s > ?", dao.SysLoginHistory.Columns.InfoId), "0")
	if err != nil {
		return 0
	}
	nums, _ := result.RowsAffected()

	return nums
}

// Export 导出excel
func (s *loginHistoryService) Export(ctx context.Context, param *define.LoginHistoryApiSelectPageReq) (string, error) {
	//"访问编号", "用户名称", "登录地址", "登录地点", "浏览器", "操作系统", "登录状态", "操作信息", "登录日期"
	m := dao.SysLoginHistory.Ctx(ctx).Fields("info_id,login_name,ipaddr,login_location,browser,os,status,msg,login_time")
	if param != nil {
		if param.LoginName != "" {
			m = m.Where("login_name like ?", "%"+param.LoginName+"%")
		}

		if param.Ipaddr != "" {
			m = m.Where("ipaddr like ?", "%"+param.Ipaddr+"%")
		}

		if param.Status != "" {
			m = m.Where("status = ?", param.Status)
		}

		if param.BeginTime != "" {
			m = m.Where("date_format(login_time,'%y%m%d') >= date_format(?,'%y%m%d')", param.BeginTime)
		}
		if param.EndTime != "" {
			m = m.Where("date_format(login_time,'%y%m%d') <= date_format(?,'%y%m%d')", param.EndTime)
		}
	}
	result, err := m.All()
	if err != nil {
		return "", err
	}
	head := []string{"访问编号", "用户名称", "登录地址", "登录地点", "浏览器", "操作系统", "登录状态", "操作信息", "登录日期"}
	key := []string{"info_id", "login_name", "ipaddr", "login_location", "browser", "os", "status", "msg", "login_time"}
	url, err := excel.DownloadExcel(head, key, result)
	if err != nil {
		return "", err
	}
	return url, nil
}

// CleanLoginLog
// @summary 清理登录日志
//          根据配置的日志存储策略,清理日志
func (s *loginHistoryService) CleanLoginLog() error {
	//查询登录日志存储配置
	logStorageCfg, err := InterfaceSysParamMgr().GetLogStorageStrategy(context.TODO(), define.ELogTypeLoginHistory)
	if err != nil {
		g.Log().Error(err)
		return err
	}
	if logStorageCfg.LogMaxDay == 0 && logStorageCfg.LogMaxNum != 0 {
		//不限制存储日期,只限制日志存储条数
		return nil
	}

	if logStorageCfg.LogMaxDay == 0 {
		//配置错误,默认保存 60 天日志
		_, err = dao.SysLoginHistory.Ctx(context.TODO()).
			WhereLTE(dao.SysLoginHistory.Columns.LoginTime, gtime.NewFromTimeStamp(gtime.Now().Timestamp()-60*24*60*60)).
			Unscoped().Delete()
	} else {
		_, err = dao.SysLoginHistory.Ctx(context.TODO()).
			WhereLTE(dao.SysLoginHistory.Columns.LoginTime, gtime.NewFromTimeStamp(gtime.Now().Timestamp()-int64(logStorageCfg.LogMaxDay*24*60*60))).
			Unscoped().Delete()
	}
	if err != nil {
		g.Log().Error(err)
	}
	return err
}
