package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	defCom "gowebserver/app/common/define"
	"gowebserver/app/common/errcode"
	"gowebserver/app/common/global"
	"gowebserver/app/common/utils/page"
	"gowebserver/app/common/utils/valueCheck"
	"gowebserver/app/system/dao"
	"gowebserver/app/system/define"
	"gowebserver/app/system/model"
	"strings"
)

//ISysBlackList 系统黑名单管理业务接口
type ISysBlackList interface {
	VerifyIpPass(ipAddr string) bool
	GetListSearch(ctx context.Context, req *define.SysBlackListSelectReq) (pageInfo *page.Paging, blackList []*model.SysBlackList, err error)
	Add(ctx context.Context, blackInfo *model.SysBlackList) (err error)
	Del(ctx context.Context, cfgIds []int) (err error)
	GetMode(ctx context.Context) (filterMode int)
	SetMode(ctx context.Context, mode int) (err error)
}

type sysBlackListImpl struct{}

var sysBlackListService = sysBlackListImpl{}

func InterfaceSysBlackList() ISysBlackList {
	return &sysBlackListService
}

func (s *sysBlackListImpl) VerifyIpPass(ipAddr string) bool {
	filterMode := s.GetMode(context.TODO())
	_, ipList, err := s.GetListSearch(context.TODO(), &define.SysBlackListSelectReq{
		Type: uint(filterMode),
		SelectPageReq: defCom.SelectPageReq{
			PageNum:  1,
			PageSize: 0,
		},
	})
	if err != nil {
		g.Log().Error(err)
		return true
	}

	switch filterMode {
	case 0:
		//未启用过滤
		return true

	case 1:
		//黑名单模式
		{
			for _, filterIpCfg := range ipList {
				if filterIpCfg.IpAddr == ipAddr {
					//在黑名单中
					return false
				}
			}
			return true
		}

	case 2:
		{
			//白名单模式
			if len(ipList) == 0 {
				return true
			}
			for _, filterIpCfg := range ipList {
				if filterIpCfg.IpAddr == ipAddr {
					//在白名单中
					return true
				}
			}
			return false
		}
	}
	return true
}

func (s *sysBlackListImpl) CheckIpAddrUnique(cfgId int, filterType uint, filterIp string) (err error) {
	err = g.Try(func() {
		count, e := dao.SysBlackList.Ctx(context.TODO()).
			Where(dao.SysBlackList.Columns.IpAddr, filterIp).
			Where(dao.SysBlackList.Columns.Type, filterType).
			WhereNot(dao.SysBlackList.Columns.Id, cfgId).
			FindCount()
		valueCheck.ErrIsNil(context.TODO(), e, errcode.ErrCommonDbOperationError.Message())
		if count > 0 {
			valueCheck.ErrIsNil(context.TODO(), gerror.NewCode(errcode.ErrCommonInvalidParameter))
		}
	})
	return
}

func (s *sysBlackListImpl) GetListSearch(ctx context.Context, req *define.SysBlackListSelectReq) (pageInfo *page.Paging, blackList []*model.SysBlackList, err error) {
	var total int
	err = g.Try(func() {
		m := dao.SysBlackList.Ctx(ctx)
		if req.Keyword != "" {
			req.Keyword = strings.Replace(req.Keyword, "'", "''", -1)
			m = m.Where(fmt.Sprintf("locate('%s', %s) > 0", req.Keyword, dao.SysBlackList.Columns.IpAddr))
		}
		if req.Type != 0 {
			m = m.Where(dao.SysBlackList.Columns.Type, req.Type)
		}

		total, err = m.Count()
		valueCheck.ErrIsNil(ctx, err)

		if req.PageSize == 0 {
			req.PageSize = total
		}
		if req.OrderByColumn != "" && req.IsAsc != "" {
			req.OrderByColumn = gstr.CaseSnake(req.OrderByColumn)
			m = m.Order(req.OrderByColumn + " " + req.IsAsc)
		} else {
			m = m.OrderAsc(dao.SysBlackList.Columns.Id)
		}

		pageInfo = page.CreatePaging(req.PageNum, req.PageSize, total)
		m = m.Limit(pageInfo.StartNum, pageInfo.PageSize)
		err = m.Scan(&blackList)
		valueCheck.ErrIsNil(ctx, err)
	})
	return
}

func (s *sysBlackListImpl) Add(ctx context.Context, blackInfo *model.SysBlackList) (err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		err = g.Try(func() {
			//校验过滤IP地址是否重复
			err = s.CheckIpAddrUnique(blackInfo.Id, blackInfo.Type, blackInfo.IpAddr)
			valueCheck.ErrIsNil(ctx, err)

			cfgId, e := dao.SysBlackList.Ctx(ctx).TX(tx).InsertAndGetId(blackInfo)
			valueCheck.ErrIsNil(ctx, e, errcode.ErrCommonDbOperationError.Message())
			blackInfo.Id = int(cfgId)
		})
		return err
	})
	return
}

func (s *sysBlackListImpl) Del(ctx context.Context, cfgIds []int) (err error) {
	err = g.Try(func() {
		_, err = dao.SysBlackList.Ctx(ctx).WhereIn(dao.SysBlackList.Columns.Id, cfgIds).Unscoped().Delete()
		valueCheck.ErrIsNil(ctx, err, errcode.ErrCommonDbOperationError.Message())
	})
	return
}

//GetMode
//@summary 查询系统黑名单过滤方式
//         0=未启用 1=黑名单 2=白名单
func (s *sysBlackListImpl) GetMode(ctx context.Context) (filterMode int) {
	sysParam, err := InterfaceSysParamMgr().GetSysParamCfg(ctx, global.SysParamKey_AccessingBlacklistFilterType)
	if err != nil {
		return 0
	}
	return gconv.Int(sysParam.ParamValue)
}

//SetMode
//@summary 设置系统黑白名单过滤方式
//         0=未启用 1=黑名单 2=白名单
func (s *sysBlackListImpl) SetMode(ctx context.Context, filterMode int) (err error) {
	return InterfaceSysParamMgr().SetSysParamCfg(ctx, global.SysParamKey_AccessingBlacklistFilterType, gconv.String(filterMode))
}
