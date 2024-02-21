package api

import (
	"context"
	"fmt"
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/mssola/user_agent"
	"gowebserver/app/common/api"
	defCom "gowebserver/app/common/define"
	"gowebserver/app/common/errcode"
	"gowebserver/app/common/global"
	"gowebserver/app/common/service/appstate"
	"gowebserver/app/common/utils"
	"gowebserver/app/common/utils/ip"
	"gowebserver/app/system/define"
	"gowebserver/app/system/model"
	"gowebserver/app/system/service"
	"gowebserver/app/system/service/publish"
	"os"
	"strings"
	"sync"
	"time"
)

type auth struct {
	api.BaseController
}

var (
	Auth                             = new(auth)
	checkOnlineUserAliveTaskMap      = make(map[string]model.SysUserOnline) //key=token value=model.SysUserOnline
	mutexCheckOnlineUserAliveTaskMap sync.Mutex                             //checkOnlineUserAliveTaskMap 读写锁
	GfToken                          = &gtoken.GfToken{
		ServerName:       "mpuaps",
		CacheKey:         global.CacheLoginUserPrefix,
		CacheMode:        g.Cfg().GetInt8("gtoken.CacheMode"),
		Timeout:          g.Cfg().GetInt("gtoken.Timeout"),
		MaxRefresh:       g.Cfg().GetInt("gtoken.MaxRefresh"),
		EncryptKey:       g.Cfg().GetBytes("gtoken.EncryptKey"),
		AuthFailMsg:      g.Cfg().GetString("gtoken.AuthFailMsg"),
		MultiLogin:       g.Cfg().GetBool("gtoken.MultiLogin"),
		LoginPath:        "/auth/login",
		LoginBeforeFunc:  Auth.Login,
		LoginAfterFunc:   Auth.LoginAfter,
		LogoutPath:       "/auth/logout",
		LogoutBeforeFunc: Auth.Logout,
		LogoutAfterFunc:  Auth.LogoutAfter,
		AuthAfterFunc:    Auth.AuthAfter,
	}
)

// AuthAfter
// @summary 用户登录请求认证成功
//          tips: 用户登录请求认证成功后,注册到gtoken执行的方法
//                1、登录用户信息设置到请求上下文
// @param1  r *ghttp.Request "http请求对象"
// @param2  respData gtoken.Resp ""
// @return1 nil
func (a *auth) AuthAfter(r *ghttp.Request, respData gtoken.Resp) {
	if r.Method == "OPTIONS" || respData.Success() {
		var user *model.SysUserExtend
		if err := gjson.DecodeTo(respData.GetString("data"), &user); err == nil {
			//请求上下文设置用户信息
			service.Context.SetUser(r.Context(), user)

			//解析请求头中的token
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				g.Log().Noticef("request unAuthorized: %s %s %v", r.RemoteAddr, r.RequestURI, r.Header)
				a.RespJsonExit(r, gerror.WrapCode(errcode.ErrHttpUnauthorized, err))
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" && parts[1] != "" {
				service.Context.SetToken(r.Context(), parts[1])
			} else {
				a.RespJsonExit(r, gerror.WrapCode(errcode.ErrHttpUnauthorized, err))
			}
		} else {
			a.RespJsonExit(r, gerror.WrapCode(errcode.ErrHttpUnauthorized, err))
			return
		}

		r.Middleware.Next()
	} else {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrHttpUnauthorized))
	}
}

// Login
// @summary 登录校验
//          1、登录用户信息校验
//          2、登录历史日志记录
//          3、接口返回登录用户详细信息,gtoken会存缓存用户详细信息与token对应,之后通过请求中的token即可获取对应的用户详细信息
// @tags 	Auth
// @Accept  json
// @Param   UserApiLoginReq body define.UserApiLoginReq true "登录信息"
// @Produce json
// @Success 200 {object} response.Response{data=define.UserApiLoginRsp} "登录成功"
// @Failure 500 {object} response.Response "登陆失败,msg描述错误原因"
// @Router 	/auth/login [POST]
func (a *auth) Login(r *ghttp.Request) (string, interface{}) {
	var (
		req  *define.UserApiLoginReq
		user *model.SysUser
		err  error
	)
	defer func() {
		var msg string
		if err == nil {
			msg = "登录成功"
		} else {
			msg = err.Error()
		}
		loginRetCode := gerror.Code(err)
		if loginRetCode == gcode.CodeNil {
			loginRetCode = errcode.CodeOk
		}
		if loginRetCode == errcode.ErrSysUserLoginAuthFailed {
			//账户密码验证错误
			errTimes := service.User.SetPasswordCounts(req.UserName)
			havingTime := global.UserPasswdErrCountMax - errTimes

			language := utils.ParseAcceptLanguage(r.GetHeader("Accept-Language"))
			g.I18n().SetLanguage(language)
			msg = g.I18n().Tf(r.GetCtx(), `{#LoginAuthFailedTipsMsg}`, havingTime)

			err = gerror.NewCode(errcode.ErrSysUserLoginAuthFailed, msg)
		}

		//登录方式
		msg = fmt.Sprintf("%s; 登录方式:%s", msg, global.ELoginType(req.LoginType).AliasString())

		//录登录历史日志
		go service.LoginHistory.Create(context.Background(), gconv.String(loginRetCode.Code()), req.UserName, r.GetClientIp(), r.Header.Get("User-Agent"), msg)

		if loginRetCode != errcode.CodeOk {
			a.RespJsonExit(r, err, nil)
		}
	}()

	//解析登录请求参数
	if err = r.Parse(&req); err != nil {
		err = gerror.WrapCode(errcode.ErrCommonInvalidParameter, err)
		return "", nil
	}

	//登录地IP过滤
	if !service.InterfaceSysBlackList().VerifyIpPass(r.GetRemoteIp()) {
		err = gerror.NewCode(errcode.ErrSysUserLoginBlackForbid)
		return "", nil
	}

	//非调试模式并且是Web客户端登录,校验图形验证码
	//if !g.Cfg().GetBool("debug.enable") &&
	//	(req.LoginType == int(global.ELoginTypeManagerClient) || req.LoginType == int(global.ELoginTypeConfigClient)) {
	//	if !comSrv.Captcha.VerifyString(req.IdKey, req.ValidateCode) {
	//		err = gerror.NewCode(errcode.ErrSysUserLoginCaptchaVerifyFailed)
	//		return "", nil
	//	}
	//}

	//检查账户是否已锁定(输入错误密码次数达到上限)
	if service.User.CheckLoginLockState(req.UserName) {
		err = gerror.NewCode(errcode.ErrSysUserLoginLock)
		return "", nil
	}

	//查询登录用户信息
	user, err = service.User.GetUserByLoginName(r.GetCtx(), req.UserName)
	if err != nil {
		err = gerror.WrapCode(errcode.ErrCommonDbOperationError, err)
		return "", nil
	}
	if user == nil {
		err = gerror.NewCode(errcode.ErrSysUserLoginAuthFailed)
		return "", nil
	}

	//校验密码
	if !service.User.CheckPasswordRight(user.LoginName, user.Password, user.Salt, req.Password) {
		err = gerror.NewCode(errcode.ErrSysUserLoginAuthFailed)
		return "", nil
	}

	//判断用户是否在其它地点登录
	//只有web客户端登录限制只能在一个地方登录一次
	//其它方式登录允许登录多次(冗余、级联、第三方集成)
	if req.LoginType == int(global.ELoginTypeManagerClient) || req.LoginType == int(global.ELoginTypeConfigClient) || req.LoginType == int(global.ELoginTypeUnknow) {
		onlineUsers := service.OnlineUser.GetOnlineUserInfoByLoginName(req.UserName)
		for _, loginInfo := range onlineUsers {
			if strings.Contains(loginInfo.Browser, global.ELoginTypeManagerClient.AliasString()) ||
				strings.Contains(loginInfo.Browser, global.ELoginTypeConfigClient.AliasString()) ||
				strings.Contains(loginInfo.Browser, global.ELoginTypeUnknow.AliasString()) {
				err = gerror.NewCode(errcode.ErrSysUserLoginAlready, fmt.Sprintf("%s[%s]", errcode.ErrSysUserLoginAlready.Message(), loginInfo.LoginLocation))
				return "", nil
			}
		}
	}

	//未绑定角色不可登录
	roles, _ := service.User.GetUserRoles(r.GetCtx(), user.UserId)
	if len(roles) == 0 {
		err = gerror.NewCode(errcode.ErrSysUserUnBindRole)
		return "", nil
	}

	//非系统管理员不可访问配置管理页面
	if !service.User.IsSysAdmin(user.UserId) && req.LoginType == int(global.ELoginTypeConfigClient) {
		err = gerror.NewCode(errcode.ErrSysUserLoginPermissionDenied)
		return "", nil
	}

	//校验并发访账户是否达到设定上限
	if service.User.CheckAccessConnectionsNumExceedMax() {
		err = gerror.NewCode(errcode.ErrSysUserLoginAccessConnectionsNumExceedMax)
		return "", nil
	}

	//校验通过
	//获取用户详细信息,并缓存登录用户信息
	service.Context.GetCtx(r.GetCtx()).Uid = user.UserId
	userInfo, getUserProfileErr := service.User.GetProfile(r.Context())
	if getUserProfileErr != nil || userInfo == nil {
		err = gerror.WrapCode(errcode.ErrSysUserLoginAuthFailed, getUserProfileErr)
		return "", nil
	}

	//登录用户信息填充到请求上下文对象中(LoginAfter中使用)
	service.Context.SetUser(r.Context(), userInfo)

	//gtoken缓存key = GfToken.CacheKey + userKey
	var userKey string
	if !g.Cfg().GetBool("gtoken.MultiLogin") || req.LoginType == int(global.ELoginTypeThirdPartyProgram) || req.LoginType == int(global.ELoginTypeBoxRedundancy) {
		//userKey = loginName + loginIP + time
		userKey = fmt.Sprintf("%s%s%s", user.LoginName, r.GetClientIp(), gtime.Now().String())
	} else {
		//userKey = loginName + userId
		userKey = fmt.Sprintf("%s%d", user.LoginName, user.UserId)
	}
	return userKey, userInfo
}

// LoginAfter
// @summary 用户登录之后的业务处理
//          1、移除用户登录尝试次数
//          2、登录用户信息设置到请求上下文
//          3、记录在线用户信息
// @param1  r *ghttp.Request "http请求"
// @param2  respData gtoken.Resp
// @return1 nil
func (a *auth) LoginAfter(r *ghttp.Request, respData gtoken.Resp) {
	if !respData.Success() {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysUserLoginAuthFailed))
	}

	//返回消息头中填写token信息
	token := respData.GetString("token")
	r.Header.Set("Authorization", "Bearer "+token)
	service.Context.SetToken(r.Context(), token)

	//登录用户信息
	userInfo := service.Context.GetUser(r.Context())
	if userInfo == nil {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysUserLoginAuthFailed))
	}

	//登录成功 移除登陆尝试次数记录
	service.User.RemovePasswordCounts(userInfo.LoginName)

	//更新用户最后登录点、时间
	userInfo.LoginIp = r.GetClientIp()
	userInfo.LoginDate = gtime.Now()
	service.User.UpdateLastLoginInfo(userInfo.UserId, r.GetClientIp())

	//下线其他用户(如果设置的同一个用户不准多端登录)
	/* 用户登录时不允许直接抢占,登陆校验时会返回提示信息
	if !g.Cfg().GetBool("gtoken.MultiLogin") {
		other := service.OnlineUser.GetOnlineUserInfoByLoginName(userInfo.UserExtend.LoginName)
		for _, loginInfo := range other {
			//更新用户在线状态
			err := service.OnlineUser.UpdateOnlineUserState(loginInfo.Token)
			if err != nil {
				g.Log().Debugf("更新在线用户状态失败: %s", loginInfo.Token)
				continue
			}

			//强制登出通知
			publish.WsPublish.PublishMsgToAllClient(defCom.ENotifyForceLogout.String(), loginInfo.Token)
		}
	}
	*/

	var loginReq define.UserApiLoginReq
	_ = r.Parse(&loginReq)
	if global.ELoginType(loginReq.LoginType) != global.ELoginTypeBoxRedundancy {
		//与冗余机箱建立连接
		//mpuSrv.InterfaceMpuBoxRedundancy().AddRedundancyBoxLink(loginReq.UserName, loginReq.Password, token)
	}

	//记录在线用户信息
	userAgent := r.Header.Get("User-Agent")
	ua := user_agent.New(userAgent)
	browser, _ := ua.Browser()
	var userOnline model.SysUserOnline
	userOnline.Token = token
	userOnline.UserId = userInfo.UserId
	userOnline.LoginName = userInfo.LoginName
	userOnline.Browser = fmt.Sprintf("%s; 登录方式:%s", browser, global.ELoginType(loginReq.LoginType).AliasString())
	userOnline.Os = fmt.Sprintf("%s Accept-Language:%s", ua.OS(), utils.ParseAcceptLanguage(r.GetHeader("Accept-Language")))
	userOnline.Ipaddr = r.GetClientIp()
	userOnline.ExpireTime = g.Cfg().GetInt("gtoken.MaxRefresh")
	userOnline.StartTimestamp = gtime.Now()
	userOnline.LastAccessTime = gtime.Now()
	userOnline.Status = global.UserStatusOnline
	userOnline.LoginLocation = ip.GetCityByIp(r.GetClientIp())
	service.OnlineUser.AddOnlineUser(userOnline)

	g.Log().Infof("User:[%d]%s LoginLocation:%s loginTime:%s Authorization: Bearer %s",
		userInfo.UserId, userInfo.LoginName, userOnline.LoginLocation, userOnline.StartTimestamp, userOnline.Token)

	//登录成功返回Token、用户信息、用户功能菜单权限
	menuList, _ := service.InterfaceSysMenu().GetUserMenus(r.GetCtx())
	a.RespJsonExit(r, nil, define.UserApiLoginRsp{
		Token:    token,
		UserInfo: userInfo,
		MenuList: menuList,
		RedundancyState: func() int {
			//0未配置 1主控冗余 2机箱冗余 3主控冗余+机箱冗余
			redundancyState := 0

			//查机箱冗余状态
			//if model2.EBoxRedundancyNotCfg != mpuSrv.InterfaceMpuBoxRedundancy().GetBoxRedundancyState() {
			//	redundancyState += 2
			//}

			//查主控冗余状态
			_, err := os.Stat("/msp/cfg/MgrVipCfg.json")
			if err == nil || os.IsExist(err) {
				redundancyState += 1
			}

			return redundancyState
		}(),
	})
}

// Logout
// @summary 登出之前的业务流程
//          1、nothing to do
// @tags 	Auth
// @Accept  json
// @Param   Authorization header string true "Bearer Token"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/auth/logout [POST]
func (a *auth) Logout(r *ghttp.Request) bool {
	return true
}

// LogoutAfter
// @summary 登出之后的业务流程
//          1、删除在线用户信息
//          2、自定义登出请求返回数据
// @return1 nil
func (a *auth) LogoutAfter(r *ghttp.Request, respData gtoken.Resp) {
	if !respData.Success() {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonOperationFailed))
	}
	g.Log().Noticef("客户端主动登出 %s", respData.DataString())
	_ = a.UserLogOutHandle(respData.DataString())

	a.RespJsonExit(r, nil, respData.DataString())
}

// RefreshToken
// @Summary 刷新token有效期
//          tips: 用于级联登录用户,定时发送http请求,刷新token有效期
// @Tags    公共接口
// @Accept  json
// @Param   Authorization header string true "Bearer Token"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router  /auth/refreshToken [PUT]
func (a *auth) RefreshToken(r *ghttp.Request) {
	//noting to do
	a.RespJsonExit(r, nil)
}

// Subscribe
// @Summary 订阅WS通知API
// @Param   IsCascadeLink header string false "级联连接标识,上级级联连接本机时 IsCascadeLink = true"
// @Param   CascadeId header string false "上级级联连接本级时,需要将上级的级联ID发给本级"
// @Description 客户端订阅代理通知消息,url：ws://127.0.0.1:9180/mpuaps/v1/ws/subscribe/*token
// @Tags 公共接口
// @Router /ws/subscribe [GET]
// @Security
func (a *auth) Subscribe(r *ghttp.Request) {
	//token := r.GetString("token") //接口有Bug,值不能有`//`,会当做转义字符处理
	data := strings.Split(r.URL.Path, "subscribe/")
	if len(data) < 2 {
		g.Log().Errorf("url[%s]解析token为空", r.GetUrl())
		return
	}
	token := data[1]

	if !Auth.CheckGTokenIsValid(token) {
		g.Log().Warningf("未登录,拒绝订阅请求[remoteAddr: %s][token: %s][r.URL.Path: %s][r.RequestURI: %s]", r.RemoteAddr, token, r.URL.Path, r.RequestURI)
		return
	}

	//升级为ws连接
	ws, err := r.WebSocket()
	if err != nil {
		g.Log().Error(err)
		return
	}

	//给对应客户端回复代理资源同步状态
	if appstate.CheckSyncMPUResourcesIsOk() {
		pubMsg := publish.PubCommMsg{
			Topic:     "", //mpurpc.ENotifyType_emNotifyResourcesReady.String(),
			TimeStamp: time.Now().UnixNano() / 1e6,
			Data:      "",
		}
		if err = ws.WriteJSON(pubMsg); err != nil {
			g.Log().Error(err)
			return
		}
		g.Log().Notice("notify client: ", ws.RemoteAddr(), pubMsg)
	}

	//如果有网络点位同步异常告警,给客户端发通知
	if state := appstate.GetSyncUmtChnWarnState(); state != 0 {
		pubMsg := publish.PubCommMsg{
			Topic:     defCom.ENotifyUmtChnSyncWarnState.String(),
			TimeStamp: time.Now().UnixNano() / 1e6,
			Data: g.Map{
				"type":  defCom.EOptModify.String(),
				"state": state,
			},
		}
		if err = ws.WriteJSON(pubMsg); err != nil {
			g.Log().Error(err)
			return
		}
		g.Log().Notice("notify client: ", ws.RemoteAddr(), pubMsg)
	}

	//添加到ws订阅客户端列表
	var bIsCascadeLink bool
	if r.GetHeader("IsCascadeLink") == "true" {
		bIsCascadeLink = true
	}
	publish.WsPublish.AddWsLink(token, bIsCascadeLink, r.GetHeader("CascadeId"), ws)
	g.Log().Infof("订阅通知成功[remoteAddr: %s][token: %s]", r.RemoteAddr, token)
}

// UserLogOutHandle
// @summary 客户端登出时的业务处理
//          可能的登出方式
//             1、正常登出
//             2、管理员强制登出
//             3、token超时强制登出
//             4、ws订阅链路长时间断开强制登出
// @param1  token string "token"
// @return1 nil
func (a *auth) UserLogOutHandle(token string) error {
	//删除所有订阅链路
	publish.WsPublish.DelClient(token, "client logout, ws close")

	//删除所有透明通道链路
	//mpuMsgForwarding.DelInstance(token, "")

	//删除机箱冗余相关连接
	//go mpuSrv.InterfaceMpuBoxRedundancy().DelRedundancyBoxLink(token)

	//释放预览资源
	//mpuSrv.MpuChannel.FreeClientMediaPlayblastRecord(token)

	//更新用户在线状态
	_ = service.OnlineUser.DeleteOnlineByToken(token)
	return nil
}

// CheckOnlineUserWsSubLinkAlive
// @summary 校验登录用户是否创建ws订阅链路,未创建则强制登出
// @param1  onlineUser model.SysUserOnline "在线用户信息"
// @return1 nil
func (a *auth) CheckOnlineUserWsSubLinkAlive(onlineUser model.SysUserOnline) {
	mutexCheckOnlineUserAliveTaskMap.Lock()
	defer mutexCheckOnlineUserAliveTaskMap.Unlock()
	if _, ok := checkOnlineUserAliveTaskMap[onlineUser.Token]; ok {
		return
	}
	checkOnlineUserAliveTaskMap[onlineUser.Token] = onlineUser

	userAliveTaskFuc := func(onlineUser model.SysUserOnline) {
		defer func() {
			mutexCheckOnlineUserAliveTaskMap.Lock()
			delete(checkOnlineUserAliveTaskMap, onlineUser.Token)
			mutexCheckOnlineUserAliveTaskMap.Unlock()
		}()

		//web刷新操作每次会重新建立ws链路
		//不能直接认为客户端已登出
		//多次判断客户端订阅链路是否已重连
		//注：CheckClient存在锁竞争问题(AddWsLink添加链路)
		//   采用短时间多次检测规避锁竞争带来的判断错误
		//400ms检测一次,总计2s内未恢复连接,认为客户端登出
		for i := 0; i < 5; i++ {
			time.Sleep(400 * time.Millisecond)
			if publish.WsPublish.CheckClient(onlineUser.Token) {
				return
			}
		}

		//客户端已登出(ws订阅链路未恢复)
		a.DelGTokenCacheUser(onlineUser.Token)
		//强制登出通知(订阅链路未连接此处无法发送通知)
		publish.WsPublish.PublishMsgToSpecifiedClient(onlineUser.Token, defCom.ENotifyForceLogout.String(), define.OnlineUserForceLogoutNotify{
			Token:  onlineUser.Token,
			Reason: define.EForceLogOut_UnConnWsSunLink,
		})
		err := a.UserLogOutHandle(onlineUser.Token)
		if err != nil {
			g.Log().Error("客户端ws订阅通道已断开,登出失败:", err)
			return
		}

	}
	go userAliveTaskFuc(onlineUser)
}

// CleanOnlineUserTask
// @summary 清理无效的在线用户
//          tips: 接口已注册到系统定时任务方法列表
// @return1 interface{}
func (a *auth) CleanOnlineUserTask() interface{} {
	offLineUsers := service.OnlineUser.GetAllOfflineUser()
	for _, userInfo := range offLineUsers {
		GfToken.RemoveToken(userInfo.Token)
		_ = service.OnlineUser.DeleteOnlineByToken(userInfo.Token)
		g.Log().Info("delete offline user: ", userInfo)
	}

	onlineUsers := service.OnlineUser.GetAllOnlineUser()
	for _, userInfo := range onlineUsers {
		//检查登录用户token是否已经过期,删除token过期的在线用户信息
		if !a.CheckGTokenIsValid(userInfo.Token) {
			g.Log().Info("token timeout. will delete invalid online user: ", userInfo)
			//发送强制登出通知
			publish.WsPublish.PublishMsgToSpecifiedClient(userInfo.Token, defCom.ENotifyForceLogout.String(), define.OnlineUserForceLogoutNotify{
				Token:  userInfo.Token,
				Reason: define.EForceLogOut_TokenInValid,
			})
			time.Sleep(100 * time.Millisecond)
			_ = a.UserLogOutHandle(userInfo.Token)
			continue
		}

		//todo 注释掉此处逻辑,可以方便postman等工具调试
		//检查登录客户端是否创建订阅链路,未创建则强制登出
		a.CheckOnlineUserWsSubLinkAlive(userInfo)
	}
	return nil
}

// CheckGTokenIsValid
// @summary 判断用户登录令牌是否有效
//          根据token找到gtoken模块缓存的在线用户信息.
// @param1  token string "登录令牌"
// @return1 bool "true有效 false无效"
func (a *auth) CheckGTokenIsValid(token string) bool {
	uuid, userKey := a.GetGTokenUuidAndUserKey(token)
	cacheKey := GfToken.CacheKey + userKey
	switch GfToken.CacheMode {
	case gtoken.CacheModeCache:
		userCacheValue, _ := gcache.Get(cacheKey)
		if userCacheValue == nil {
			return false
		}
		return true
	case gtoken.CacheModeRedis:
		var userCache g.Map
		userCacheJson, err := g.Redis().Do("GET", cacheKey)
		if err != nil {
			g.Log().Error("[GToken]cache get error", err)
			return false
		}
		if userCacheJson == nil {
			return false
		}
		err = gjson.DecodeTo(userCacheJson, &userCache)
		if err != nil {
			g.Log().Error("[GToken]cache get json error", err)
			return false
		}
		if uuid != userCache["uuid"] {
			return false
		}
		return true
	}
	return false
}

// GetGTokenUuidAndUserKey
// @summary 根据token获取uuid和userKey
//          tips: 获取gtoken模块token对应用户登录信息缓存key; 用户uuid;
// @param1  token string "登陆令牌"
// @return1 uuid string "登录用户uuid"
// @return2 userKey string "用户缓存key"
func (a *auth) GetGTokenUuidAndUserKey(token string) (uuid, userKey string) {
	decryptToken := GfToken.DecryptToken(token)
	if !decryptToken.Success() {
		return
	}
	userKey = decryptToken.GetString("userKey")
	uuid = decryptToken.GetString("uuid")
	return
}

// GetGTokenCacheUser
// @summary 获取请求携带token对应的用户信息
//          返回空值,说明token无效
// @param1  r *ghttp.Request "请求对象"
// @return1 userInfo *model.SysUserExtend "返回登录用户信息"
func (a *auth) GetGTokenCacheUser(r *ghttp.Request) (userInfo *model.SysUserExtend) {
	tokenInfo := GfToken.GetTokenData(r)
	if !tokenInfo.Success() {
		return nil
	}
	if err := gjson.DecodeTo(tokenInfo.GetString("data"), &userInfo); err == nil {
		return
	}
	return nil
}

// DelGTokenCacheUser
// @summary 删除GToken缓存
func (a *auth) DelGTokenCacheUser(token string) {
	//清除token缓存
	if resp := GfToken.RemoveToken(token); !resp.Success() {
		g.Log().Error("在线用户登出,清除删除gtoken缓存出错:", resp)
	}
}
