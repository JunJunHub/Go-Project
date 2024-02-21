package app

import (
	_ "gowebserver/app/common/task"
	_ "gowebserver/packed"

	"gowebserver/app/common/middleware"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gmode"
	"github.com/gogf/swagger"
)

// app加载的模块列表
import (
	_ "gowebserver/app/common"

	_ "gowebserver/app/system"

	_ "gowebserver/app/customization"
)

// Run 应用启动入口
func Run() {
	httpServer := g.Server()
	httpServer.Use(middleware.CORS)
	httpServer.Use(middleware.CheckClientOS)
	httpServer.Plugin(&swagger.Swagger{
		BasicPath:     g.Cfg().GetString("server.Prefix"),
		BasicAuthUser: "admin",
		BasicAuthPass: "admin123",
	})

	/*
		//前端管理界面默认路由前缀: client 例: http://127.0.0.1:9180/client || http://127.0.0.1:9180
		//管理类前端路由默认指向    ./public/web/manage_client/index.html
		httpServer.SetRewrite("/", "/client")
		httpServer.AddStaticPath("/client", "./public/web/manage_client")
		httpServer.BindHandler("/client/*any", func(r *ghttp.Request) {
			r.Response.ServeFile("./public/web/manage_client/index.html")
		})

		//前端配置界面默认路由前缀: config 例:http://127.0.0.1:9180/config
		//配置类前端路由默认指向    ./public/web/config_client/index.html
		httpServer.AddStaticPath("/config", "./public/web/config_client")
		httpServer.BindHandler("/config/*any", func(r *ghttp.Request) {
			r.Response.ServeFile("./public/web/config_client/index.html")
		})
	*/

	/*
		//Pad管理界面默认路由前缀：
		httpServer.AddStaticPath("/pad", "./public/web/pad/dist")
		httpServer.BindHandler("/pad/*any", func(r *ghttp.Request) {
			r.Response.ServeFile("./public/web/pad/dist/index.html")
		})
		httpServer.BindHandler("/assets/*any", func(r *ghttp.Request) {
			urls := gstr.Split(r.GetUrl(), ".")
			// 如果请求js文件
			if len(urls) > 0 && urls[len(urls)-1] == "js" {
				r.Response.ResponseWriter.Header().Set("Content-Type", "application/javascript; charset=utf-8")
			}
			params := gstr.Split(r.GetUrl(), "assets")
			r.Response.ServeFile("./public/web/pad/dist/assets" + params[1])
		})
		httpServer.BindHandler("/global-config.json*", func(r *ghttp.Request) {
			r.Response.ServeFile("./public/web/pad/dist/global-config.json")
		})
		httpServer.BindHandler("/favicon.ico*", func(r *ghttp.Request) {
			r.Response.ServeFile("./public/web/pad/dist/favicon.ico")
		})
	*/

	//开发阶段禁止浏览器缓存, 方便调试
	if gmode.IsDevelop() {
		httpServer.BindHookHandler("/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
			r.Response.Header().Set("Cache-Control", "no-store")
		})
	}

	httpServer.Run()
}
