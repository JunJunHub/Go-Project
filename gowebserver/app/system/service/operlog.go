// =================================================================================
// 系统操作日志相关业务逻辑(日志查询、导出、删除...)
// =================================================================================

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"gowebserver/app/common/global"
	"gowebserver/app/common/utils/convert"
	"gowebserver/app/common/utils/excel"
	"gowebserver/app/common/utils/ip"
	"gowebserver/app/common/utils/page"
	"gowebserver/app/common/utils/response"
	"gowebserver/app/system/dao"
	"gowebserver/app/system/define"
	"gowebserver/app/system/model"
	"strings"
)

var OperLog = &operLogService{}

type operLogService struct{}

// GetList 分页列表
func (s *operLogService) GetList(ctx context.Context, param *define.OperlogApiSelectPageReq) *define.OperlogServiceList {
	m := dao.SysOperLog.Ctx(ctx).As("t")
	if param != nil {
		if param.Title != "" {
			m = m.Where("t.title like ?", "%"+param.Title+"%")
		}
		if param.OperName != "" {
			m = m.Where("t.oper_name like ?", "%"+param.OperName+"%")
		}
		if param.Status != "" {
			m = m.Where("t.status = ?", param.Status)
		}
		if param.BusinessType != "" {
			m = m.Where("t.business_type = ?", param.BusinessType)
		}

		if param.BeginTime != "" {
			m = m.Where("date_format(t.oper_time,'%y%m%d') >= date_format(?,'%y%m%d')", param.BeginTime)
		}

		if param.EndTime != "" {
			m = m.Where("date_format(t.oper_time,'%y%m%d') <= date_format(?,'%y%m%d')", param.EndTime)
		}
	}
	total, err := m.Count()
	m = m.Fields("oper_id, title, business_type, method, request_method, operator_type, oper_name, dept_name, oper_url, oper_ip, oper_location, oper_param, json_result, status, error_msg, oper_time")
	if err != nil {
		return nil
	}
	pageInfo := page.CreatePaging(param.PageNum, param.PageSize, total)
	m = m.Order("t.oper_time desc")
	m = m.Limit(pageInfo.StartNum, pageInfo.PageSize)
	result := &define.OperlogServiceList{
		Page:  pageInfo.PageNum,
		Size:  pageInfo.PageSize,
		Total: pageInfo.Total,
	}
	if err = m.Scan(&result.List); err != nil {
		return nil
	}
	return result
}

// Create 新增记录
func (s *operLogService) Create(r *ghttp.Request, title, inContent string, outContent *response.Response, businessType global.BusinessType) {
	var operLog model.SysOperLog

	user := Context.GetUser(r.Context())
	if user == nil {
		operLog.OperName = "未登录"
	} else {
		operLog.OperName = user.LoginName
	}

	//操作执行结果
	outContent.Data = ""
	outJson, _ := json.Marshal(outContent)
	outJsonResult := string(outJson)

	operLog.Title = title
	operLog.OperParam = inContent
	operLog.JsonResult = outJsonResult
	operLog.BusinessType = gconv.Int(businessType)
	if operLog.BusinessType == 0 {
		operLog.BusinessType = s.GetBusinessTypeByMethod(r.Method)
	}

	//操作类别
	if strings.Contains(r.RequestURI, "cascadeApi") {
		operLog.OperatorType = 2
	} else {
		operLog.OperatorType = 1
	}

	//操作状态（0正常 -1异常）
	if outContent.Code == 0 || outContent.Code == 200 {
		operLog.Status = 0
	} else {
		operLog.Status = -1
	}

	operLog.RequestMethod = r.Method
	operLog.OperUrl = r.RequestURI
	operLog.Method = r.RequestURI
	operLog.OperIp = r.GetClientIp()
	operLog.OperLocation = ip.GetCityByIp(operLog.OperIp)
	operLog.OperTime = gtime.Now()

	if _, err := dao.SysOperLog.Ctx(context.Background()).Data(operLog).Insert(); err != nil {
		g.Log().Line().Error(err)
	}
}

// Delete 批量删除记录
func (s *operLogService) Delete(ctx context.Context, ids string) int64 {
	idArr := convert.ToInt64Array(ids, ",")
	result, err := dao.SysOperLog.Ctx(ctx).Delete(fmt.Sprintf("%s in(?)", dao.SysOperLog.Columns.OperId), idArr)
	if err != nil {
		return 0
	}
	nums, _ := result.RowsAffected()
	return nums
}

// Clean 清空记录
func (s *operLogService) Clean(ctx context.Context) int64 {
	result, err := dao.SysOperLog.Ctx(ctx).Delete(fmt.Sprintf("%s > ?", dao.SysOperLog.Columns.OperId), "0")
	if err != nil {
		return 0
	}
	nums, _ := result.RowsAffected()
	return nums
}

// Export 导出excel
func (s *operLogService) Export(ctx context.Context, param *define.OperlogApiSelectPageReq) (string, error) {
	//"日志主键", "模块标题", "业务类型", "方法名称", "请求方式", "操作类别", "操作人员", "部门名称", "请求URL", "主机地址", "操作地点", "请求参数", "返回参数", "操作结果", "错误信息", "操作时间"
	m := dao.SysOperLog.Ctx(ctx).Fields("oper_id, title, business_type, method, request_method, operator_type, oper_name, dept_name, oper_url, oper_ip, oper_location, oper_param, json_result, status, error_msg, oper_time")
	if param != nil {
		if param.Title != "" {
			m = m.Where("title like ?", "%"+param.Title+"%")
		}

		if param.OperName != "" {
			m = m.Where("oper_name like ?", "%"+param.OperName+"%")
		}

		if param.Status != "" {
			m = m.Where("status = ?", param.Status)
		}

		if param.BusinessType != "" {
			m = m.Where("business_type = ?", param.BusinessType)
		}

		if param.BeginTime != "" {
			m = m.Where("date_format(oper_time,'%y%m%d') >= date_format(?,'%y%m%d')", param.BeginTime)
		}

		if param.EndTime != "" {
			m = m.Where("date_format(oper_time,'%y%m%d') <= date_format(?,'%y%m%d')", param.EndTime)
		}
	}
	result, err := m.All()
	if err != nil {
		return "", err
	}
	head := []string{"日志主键", "模块标题", "业务类型", "方法名称", "请求方式", "操作类别", "操作人员", "部门名称", "请求URL", "主机地址", "操作地点", "请求参数", "返回参数", "操作结果", "错误信息", "操作时间"}
	key := []string{"oper_id", "title", "business_type", "method", "request_method", "operator_type", "oper_name", "dept_name", "oper_url", "oper_ip", "oper_location", "oper_param", "json_result", "status", "error_msg", "oper_time"}
	url, err := excel.DownloadExcel(head, key, result)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *operLogService) GetBusinessTypeByMethod(method string) int {
	switch method {
	case "GET":
		return gconv.Int(global.BusinessGet)
	case "POST":
		return gconv.Int(global.BusinessAdd)
	case "PUT":
		return gconv.Int(global.BusinessEdit)
	case "DELETE":
		return gconv.Int(global.BusinessDel)
	default:
		return gconv.Int(global.BusinessOther)
	}
}

// CleanOperlog
// @summary 清理操作日志
//          根据配置的日志存储策略,清理日志
func (s *operLogService) CleanOperlog() error {
	//查询操作日志存储配置
	logStorageCfg, err := InterfaceSysParamMgr().GetLogStorageStrategy(context.TODO(), define.ELogTypeOperlog)
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
		_, err = dao.SysOperLog.Ctx(context.TODO()).
			WhereLTE(dao.SysOperLog.Columns.OperTime, gtime.NewFromTimeStamp(gtime.Now().Timestamp()-60*24*60*60)).
			Unscoped().Delete()
	} else {
		_, err = dao.SysOperLog.Ctx(context.TODO()).
			WhereLTE(dao.SysOperLog.Columns.OperTime, gtime.NewFromTimeStamp(gtime.Now().Timestamp()-int64(logStorageCfg.LogMaxDay*24*60*60))).
			Unscoped().Delete()
	}
	if err != nil {
		g.Log().Error(err)
	}
	return err
}
