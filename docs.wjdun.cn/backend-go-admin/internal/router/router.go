package router

import (
	"context"
	"net/http"
	"time"

	"base/internal/config"
	"base/internal/db"
	"base/internal/handler"
	authhandler "base/internal/handler/auth"
	systemhandler "base/internal/handler/system"
	"base/internal/middleware"
	rediscache "base/internal/redis"
	authservice "base/internal/service/auth"
	systemservice "base/internal/service/system"

	casbinlib "github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Dependencies struct {
	Database *db.Resource
	Redis    *rediscache.Client
	Enforcer *casbinlib.Enforcer
}

func New(ctx context.Context, cfg config.Config, log *zap.Logger, deps Dependencies) *gin.Engine {
	if cfg.App.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.PublicArticleCORS())
	r.Use(func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.JSON(http.StatusOK, gin.H{})
			c.Abort()
		} else {
			c.Next()
		}
	})
	r.Use(middleware.AppContext(ctx))
	r.Use(middleware.RequestID())
	r.Use(middleware.Resources(deps.Database.Client, deps.Redis, log))
	r.Use(middleware.RequestLogger(log))

	health := handler.NewHealthHandler()
	authSvc := authservice.New(deps.Database.Client, cfg.JWT.Secret)
	authAPI := authhandler.New(authSvc)
	systemSvc := systemservice.New(deps.Database.Client)
	startTime, _ := ctx.Value("app_start_time").(time.Time)
	monitorSvc := systemservice.NewMonitor(deps.Database, deps.Redis, cfg, startTime)
	systemAPI := systemhandler.New(systemSvc, monitorSvc)
	api := r.Group("/api")
	{
		api.GET("/health", health.Health)
		api.POST("/auth/login", authAPI.Login)
		api.GET("/articles", systemAPI.ListPublishedArticles)
		api.GET("/articles/:id", systemAPI.GetPublishedArticle)
		api.GET("/uploads/articles/:name", systemhandler.ServeArticleImage)
		api.POST("/upload", middleware.Authentication(authSvc), middleware.PermissionAny(deps.Enforcer, "content:article:create", "content:article:update"), systemhandler.UploadArticleImage)

		api.GET("/user/info", middleware.Authentication(authSvc), authAPI.UserInfo)
		api.GET("/auth/codes", middleware.Authentication(authSvc), authAPI.AccessCodes)
		api.GET("/menu/all", middleware.Authentication(authSvc), authAPI.Menus)

		api.GET("/private", middleware.Authentication(authSvc), health.Private)

		system := api.Group("/system")
		system.Use(middleware.Authentication(authSvc))
		{
			system.GET("/users", middleware.Permission(deps.Enforcer, "sys:user:list"), systemAPI.ListUsers)
			system.GET("/users/:id", middleware.Permission(deps.Enforcer, "sys:user:list"), systemAPI.GetUser)
			system.POST("/users", middleware.Permission(deps.Enforcer, "sys:user:create"), middleware.Audit(log, "创建用户"), systemAPI.CreateUser)
			system.PUT("/users/:id", middleware.Permission(deps.Enforcer, "sys:user:update"), middleware.Audit(log, "修改用户"), systemAPI.UpdateUser)
			system.DELETE("/users/:id", middleware.Permission(deps.Enforcer, "sys:user:delete"), middleware.Audit(log, "删除用户"), systemAPI.DeleteUser)
			system.PUT("/users/:id/roles", middleware.Permission(deps.Enforcer, "sys:user:update"), middleware.Audit(log, "分配用户角色"), systemAPI.ReplaceUserRoles)
			system.PUT("/users/:id/password", middleware.Permission(deps.Enforcer, "sys:user:update"), middleware.Audit(log, "重置用户密码"), systemAPI.ResetUserPassword)

			system.GET("/roles", middleware.Permission(deps.Enforcer, "sys:role:list"), systemAPI.ListRoles)
			system.GET("/roles/:id", middleware.Permission(deps.Enforcer, "sys:role:list"), systemAPI.GetRole)
			system.POST("/roles", middleware.Permission(deps.Enforcer, "sys:role:create"), middleware.Audit(log, "创建角色"), systemAPI.CreateRole)
			system.PUT("/roles/:id", middleware.Permission(deps.Enforcer, "sys:role:update"), middleware.Audit(log, "修改角色"), systemAPI.UpdateRole)
			system.DELETE("/roles/:id", middleware.Permission(deps.Enforcer, "sys:role:delete"), middleware.Audit(log, "删除角色"), systemAPI.DeleteRole)
			system.PUT("/roles/:id/menus", middleware.Permission(deps.Enforcer, "sys:role:update"), middleware.Audit(log, "分配角色菜单"), systemAPI.ReplaceRoleMenus)

			system.GET("/menus", middleware.Permission(deps.Enforcer, "sys:menu:list"), systemAPI.ListMenus)
			system.GET("/menus/tree", middleware.Permission(deps.Enforcer, "sys:menu:list"), systemAPI.ListMenuTree)
			system.GET("/menus/:id", middleware.Permission(deps.Enforcer, "sys:menu:list"), systemAPI.GetMenu)
			system.POST("/menus", middleware.Permission(deps.Enforcer, "sys:menu:create"), middleware.Audit(log, "创建菜单"), systemAPI.CreateMenu)
			system.PUT("/menus/:id", middleware.Permission(deps.Enforcer, "sys:menu:update"), middleware.Audit(log, "修改菜单"), systemAPI.UpdateMenu)
			system.DELETE("/menus/:id", middleware.Permission(deps.Enforcer, "sys:menu:delete"), middleware.Audit(log, "删除菜单"), systemAPI.DeleteMenu)

			system.GET("/organizations", middleware.Permission(deps.Enforcer, "sys:organization:list"), systemAPI.ListOrganizations)
			system.GET("/organizations/tree", middleware.Permission(deps.Enforcer, "sys:organization:list"), systemAPI.ListOrganizationTree)
			system.GET("/organizations/:id", middleware.Permission(deps.Enforcer, "sys:organization:list"), systemAPI.GetOrganization)
			system.POST("/organizations", middleware.Permission(deps.Enforcer, "sys:organization:create"), middleware.Audit(log, "创建组织"), systemAPI.CreateOrganization)
			system.PUT("/organizations/:id", middleware.Permission(deps.Enforcer, "sys:organization:update"), middleware.Audit(log, "修改组织"), systemAPI.UpdateOrganization)
			system.DELETE("/organizations/:id", middleware.Permission(deps.Enforcer, "sys:organization:delete"), middleware.Audit(log, "删除组织"), systemAPI.DeleteOrganization)

			system.GET("/dict-types", middleware.Permission(deps.Enforcer, "sys:dict:list"), systemAPI.ListDictTypes)
			system.GET("/dict-types/:id", middleware.Permission(deps.Enforcer, "sys:dict:list"), systemAPI.GetDictType)
			system.POST("/dict-types", middleware.Permission(deps.Enforcer, "sys:dict:create"), middleware.Audit(log, "创建字典类型"), systemAPI.CreateDictType)
			system.PUT("/dict-types/:id", middleware.Permission(deps.Enforcer, "sys:dict:update"), middleware.Audit(log, "修改字典类型"), systemAPI.UpdateDictType)
			system.DELETE("/dict-types/:id", middleware.Permission(deps.Enforcer, "sys:dict:delete"), middleware.Audit(log, "删除字典类型"), systemAPI.DeleteDictType)

			system.GET("/dict-data", middleware.Permission(deps.Enforcer, "sys:dict:list"), systemAPI.ListDictData)
			system.GET("/dict-data/:id", middleware.Permission(deps.Enforcer, "sys:dict:list"), systemAPI.GetDictData)
			system.POST("/dict-data", middleware.Permission(deps.Enforcer, "sys:dict:create"), middleware.Audit(log, "创建字典数据"), systemAPI.CreateDictData)
			system.PUT("/dict-data/:id", middleware.Permission(deps.Enforcer, "sys:dict:update"), middleware.Audit(log, "修改字典数据"), systemAPI.UpdateDictData)
			system.DELETE("/dict-data/:id", middleware.Permission(deps.Enforcer, "sys:dict:delete"), middleware.Audit(log, "删除字典数据"), systemAPI.DeleteDictData)

			system.GET("/monitor", middleware.Permission(deps.Enforcer, "sys:monitor:view"), systemAPI.GetMonitor)
			system.GET("/articles", middleware.Permission(deps.Enforcer, "content:article:list"), systemAPI.ListArticles)
			system.GET("/articles/:id", middleware.Permission(deps.Enforcer, "content:article:list"), systemAPI.GetArticle)
			system.POST("/articles", middleware.Permission(deps.Enforcer, "content:article:create"), middleware.Audit(log, "create article"), systemAPI.CreateArticle)
			system.PUT("/articles/:id", middleware.Permission(deps.Enforcer, "content:article:update"), middleware.Audit(log, "update article"), systemAPI.UpdateArticle)
			system.DELETE("/articles/:id", middleware.Permission(deps.Enforcer, "content:article:delete"), middleware.Audit(log, "delete article"), systemAPI.DeleteArticle)
		}
	}

	return r
}
