package bootstrap

import (
	"context"
	"fmt"

	"base/ent"
	"base/ent/syscasbinrule"
	"base/ent/sysdictdata"
	"base/ent/sysdicttype"
	"base/ent/sysmenu"
	"base/ent/sysorganization"
	"base/ent/sysrole"
	"base/ent/sysrolemenu"
	"base/ent/sysuser"
	"base/ent/sysuserrole"
	"base/internal/constants"
	"base/internal/pkg/password"
)

type MenuSeed struct {
	Name           string
	Title          string
	Path           string
	Component      string
	Icon           string
	Type           string
	PermissionCode string
	Redirect       string
	Sort           int
	Hidden         bool
	KeepAlive      bool
	Remark         string
	Children       []MenuSeed
}

func InitSystemData(ctx context.Context, client *ent.Client) error {
	rootOrgID, err := initOrganizations(ctx, client)
	if err != nil {
		return err
	}

	roleID, err := initRoles(ctx, client)
	if err != nil {
		return err
	}

	userID, err := initUsers(ctx, client, rootOrgID)
	if err != nil {
		return err
	}

	if err := initUserRoles(ctx, client, userID, roleID); err != nil {
		return err
	}
	if err := initDictData(ctx, client); err != nil {
		return err
	}

	menuIDs, permissionCodes, err := initMenus(ctx, client)
	if err != nil {
		return err
	}
	if err := initRoleMenus(ctx, client, roleID, menuIDs); err != nil {
		return err
	}
	if err := initCasbinRules(ctx, client, roleID, permissionCodes); err != nil {
		return err
	}
	return nil
}

func initOrganizations(ctx context.Context, client *ent.Client) (int, error) {
	exists, err := client.SysOrganization.Query().
		Where(sysorganization.CodeEQ("root")).
		Only(ctx)
	if err == nil {
		return exists.ID, nil
	}
	if !ent.IsNotFound(err) {
		return 0, err
	}

	org, err := client.SysOrganization.Create().
		SetCode("root").
		SetName("平台组织").
		SetType("platform").
		SetLeader("系统").
		SetSort(0).
		SetStatus(constants.StatusEnabled).
		SetRemark("系统初始化根组织").
		Save(ctx)
	if err != nil {
		return 0, err
	}
	return org.ID, nil
}

func initRoles(ctx context.Context, client *ent.Client) (int, error) {
	exists, err := client.SysRole.Query().
		Where(sysrole.CodeEQ("admin")).
		Only(ctx)
	if err == nil {
		return exists.ID, nil
	}
	if !ent.IsNotFound(err) {
		return 0, err
	}

	role, err := client.SysRole.Create().
		SetCode("admin").
		SetName("管理员").
		SetSort(0).
		SetStatus(constants.StatusEnabled).
		SetRemark("系统初始化管理员角色").
		Save(ctx)
	if err != nil {
		return 0, err
	}
	return role.ID, nil
}

func initUsers(ctx context.Context, client *ent.Client, organizationID int) (int, error) {
	exists, err := client.SysUser.Query().
		Where(sysuser.UsernameEQ("admin")).
		Only(ctx)
	if err == nil {
		return exists.ID, nil
	}
	if !ent.IsNotFound(err) {
		return 0, err
	}

	passwordHash, err := password.Hash("admin123456")
	if err != nil {
		return 0, err
	}

	user, err := client.SysUser.Create().
		SetUsername("admin").
		SetPasswordHash(passwordHash).
		SetRealName("超级管理员").
		SetNickname("管理员").
		SetStatus(constants.StatusEnabled).
		SetIsSuperAdmin(true).
		SetOrganizationID(organizationID).
		SetRemark("系统初始化管理员账号").
		Save(ctx)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func initUserRoles(ctx context.Context, client *ent.Client, userID, roleID int) error {
	exists, err := client.SysUserRole.Query().
		Where(
			sysuserrole.UserID(userID),
			sysuserrole.RoleID(roleID),
		).
		Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	return client.SysUserRole.Create().
		SetUserID(userID).
		SetRoleID(roleID).
		SetUser(&ent.SysUser{ID: userID}).
		SetRole(&ent.SysRole{ID: roleID}).
		Exec(ctx)
}

func initDictData(ctx context.Context, client *ent.Client) error {
	type dictItem struct {
		Code string
		Name string
		Data []struct {
			Label     string
			Value     string
			TagType   string
			Sort      int
			IsDefault bool
		}
	}

	items := []dictItem{
		{
			Code: "sys_common_status",
			Name: "通用状态",
			Data: []struct {
				Label     string
				Value     string
				TagType   string
				Sort      int
				IsDefault bool
			}{
				{Label: "启用", Value: "1", TagType: "success", Sort: 1, IsDefault: true},
				{Label: "停用", Value: "0", TagType: "error", Sort: 2},
			},
		},
		{
			Code: "sys_menu_type",
			Name: "菜单类型",
			Data: []struct {
				Label     string
				Value     string
				TagType   string
				Sort      int
				IsDefault bool
			}{
				{Label: "目录", Value: constants.MenuTypeCatalog, TagType: "processing", Sort: 1},
				{Label: "菜单", Value: constants.MenuTypeMenu, TagType: "success", Sort: 2},
				{Label: "按钮", Value: constants.MenuTypeButton, TagType: "warning", Sort: 3},
				{Label: "接口", Value: constants.MenuTypeAPI, TagType: "error", Sort: 4},
			},
		},
		{
			Code: "sys_org_type",
			Name: "组织类型",
			Data: []struct {
				Label     string
				Value     string
				TagType   string
				Sort      int
				IsDefault bool
			}{
				{Label: "平台", Value: "platform", TagType: "processing", Sort: 1, IsDefault: true},
				{Label: "公司", Value: "company", TagType: "success", Sort: 2},
				{Label: "部门", Value: "department", TagType: "warning", Sort: 3},
				{Label: "小组", Value: "group", TagType: "default", Sort: 4},
			},
		},
	}

	for _, item := range items {
		dictType, err := client.SysDictType.Query().
			Where(sysdicttype.CodeEQ(item.Code)).
			Only(ctx)
		if ent.IsNotFound(err) {
			dictType, err = client.SysDictType.Create().
				SetCode(item.Code).
				SetName(item.Name).
				SetStatus(constants.StatusEnabled).
				SetRemark("系统初始化字典类型").
				Save(ctx)
		}
		if err != nil {
			return err
		}

		for _, data := range item.Data {
			exists, err := client.SysDictData.Query().
				Where(
					sysdictdata.DictTypeID(dictType.ID),
					sysdictdata.ValueEQ(data.Value),
				).
				Exist(ctx)
			if err != nil {
				return err
			}
			if exists {
				continue
			}

			err = client.SysDictData.Create().
				SetDictTypeID(dictType.ID).
				SetDictType(dictType).
				SetLabel(data.Label).
				SetValue(data.Value).
				SetTagType(data.TagType).
				SetSort(data.Sort).
				SetStatus(constants.StatusEnabled).
				SetIsDefault(data.IsDefault).
				SetRemark("系统初始化字典数据").
				Exec(ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func initMenus(ctx context.Context, client *ent.Client) ([]int, []string, error) {
	seeds := []MenuSeed{
		{
			Name:      constants.SystemDashboardName,
			Title:     "控制台",
			Path:      "/dashboard",
			Component: "LAYOUT",
			Icon:      "lucide:layout-dashboard",
			Type:      constants.MenuTypeCatalog,
			Redirect:  "/dashboard/monitor",
			Sort:      1,
			Remark:    "系统控制台目录",
			Children: []MenuSeed{
				{
					Name:      "SystemMonitor",
					Title:     "系统监控",
					Path:      "/dashboard/monitor",
					Component: "/dashboard/monitor/index",
					Icon:      "carbon:chart-line-data",
					Type:      constants.MenuTypeMenu,
					Sort:      1,
					KeepAlive: true,
					Remark:    "系统监控首页",
					Children: []MenuSeed{
						{
							Name:           "SystemMonitorViewAPI",
							Title:          "系统监控接口",
							Type:           constants.MenuTypeAPI,
							PermissionCode: "sys:monitor:view",
							Sort:           101,
							Hidden:         true,
							Remark:         "系统监控查看接口权限",
						},
					},
				},
			},
		},
		{
			Name:      "System",
			Title:     "系统管理",
			Path:      "/system",
			Component: "LAYOUT",
			Icon:      "carbon:settings",
			Type:      constants.MenuTypeCatalog,
			Redirect:  "/system/user",
			Sort:      100,
			Remark:    "系统管理目录",
			Children: []MenuSeed{
				buildEntityMenuSeed("User", "用户管理", "/system/user", "/system/user/index", "carbon:user-multiple", "sys:user"),
				buildEntityMenuSeed("Role", "角色管理", "/system/role", "/system/role/index", "carbon:user-role", "sys:role"),
				buildEntityMenuSeed("Menu", "菜单管理", "/system/menu", "/system/menu/index", "carbon:menu", "sys:menu"),
				buildEntityMenuSeed("Organization", "组织管理", "/system/organization", "/system/organization/index", "carbon:tree-view-alt", "sys:organization"),
				buildEntityMenuSeed("Dict", "字典管理", "/system/dict", "/system/dict/index", "carbon:data-structured", "sys:dict"),
				buildEntityMenuSeed("Article", "文章管理", "/system/article", "/system/article/index", "carbon:document", "content:article"),
			},
		},
	}

	var menuIDs []int
	var permissionCodes []string
	for _, seed := range seeds {
		ids, codes, err := saveMenuTree(ctx, client, nil, seed)
		if err != nil {
			return nil, nil, err
		}
		menuIDs = append(menuIDs, ids...)
		permissionCodes = append(permissionCodes, codes...)
	}
	return menuIDs, permissionCodes, nil
}

func saveMenuTree(ctx context.Context, client *ent.Client, parentID *int, seed MenuSeed) ([]int, []string, error) {
	menu, err := client.SysMenu.Query().
		Where(sysmenu.NameEQ(seed.Name)).
		Only(ctx)
	if ent.IsNotFound(err) {
		create := client.SysMenu.Create().
			SetName(seed.Name).
			SetTitle(seed.Title).
			SetPath(seed.Path).
			SetComponent(seed.Component).
			SetIcon(seed.Icon).
			SetType(seed.Type).
			SetPermissionCode(seed.PermissionCode).
			SetRedirect(seed.Redirect).
			SetSort(seed.Sort).
			SetHidden(seed.Hidden).
			SetKeepAlive(seed.KeepAlive).
			SetStatus(constants.StatusEnabled).
			SetRemark(seed.Remark)
		if parentID != nil {
			create.SetParentID(*parentID)
		}
		menu, err = create.Save(ctx)
	} else if err == nil {
		update := client.SysMenu.UpdateOneID(menu.ID).
			SetTitle(seed.Title).
			SetPath(seed.Path).
			SetComponent(seed.Component).
			SetIcon(seed.Icon).
			SetType(seed.Type).
			SetPermissionCode(seed.PermissionCode).
			SetRedirect(seed.Redirect).
			SetSort(seed.Sort).
			SetHidden(seed.Hidden).
			SetKeepAlive(seed.KeepAlive).
			SetStatus(constants.StatusEnabled).
			SetRemark(seed.Remark)
		if parentID == nil {
			update.ClearParentID()
		} else {
			update.SetParentID(*parentID)
		}
		menu, err = update.Save(ctx)
	}
	if err != nil {
		return nil, nil, err
	}

	ids := []int{menu.ID}
	var codes []string
	if seed.PermissionCode != "" {
		codes = append(codes, seed.PermissionCode)
	}

	for _, child := range seed.Children {
		childIDs, childCodes, err := saveMenuTree(ctx, client, &menu.ID, child)
		if err != nil {
			return nil, nil, err
		}
		ids = append(ids, childIDs...)
		codes = append(codes, childCodes...)
	}
	return ids, codes, nil
}

func buildEntityMenuSeed(name, title, path, component, icon, permissionPrefix string) MenuSeed {
	return MenuSeed{
		Name:      name,
		Title:     title,
		Path:      path,
		Component: component,
		Icon:      icon,
		Type:      constants.MenuTypeMenu,
		Sort:      10,
		Remark:    fmt.Sprintf("%s菜单", title),
		Children: []MenuSeed{
			{
				Name:           name + "ListButton",
				Title:          "查询",
				Type:           constants.MenuTypeButton,
				PermissionCode: permissionPrefix + ":list",
				Sort:           1,
				Hidden:         true,
				Remark:         title + "查询按钮权限",
			},
			{
				Name:           name + "CreateButton",
				Title:          "新增",
				Type:           constants.MenuTypeButton,
				PermissionCode: permissionPrefix + ":create",
				Sort:           2,
				Hidden:         true,
				Remark:         title + "新增按钮权限",
			},
			{
				Name:           name + "UpdateButton",
				Title:          "修改",
				Type:           constants.MenuTypeButton,
				PermissionCode: permissionPrefix + ":update",
				Sort:           3,
				Hidden:         true,
				Remark:         title + "修改按钮权限",
			},
			{
				Name:           name + "DeleteButton",
				Title:          "删除",
				Type:           constants.MenuTypeButton,
				PermissionCode: permissionPrefix + ":delete",
				Sort:           4,
				Hidden:         true,
				Remark:         title + "删除按钮权限",
			},
			{
				Name:           name + "ListAPI",
				Title:          "列表接口",
				Type:           constants.MenuTypeAPI,
				PermissionCode: permissionPrefix + ":list",
				Sort:           101,
				Hidden:         true,
				Remark:         title + "列表接口权限",
			},
			{
				Name:           name + "CreateAPI",
				Title:          "新增接口",
				Type:           constants.MenuTypeAPI,
				PermissionCode: permissionPrefix + ":create",
				Sort:           102,
				Hidden:         true,
				Remark:         title + "新增接口权限",
			},
			{
				Name:           name + "UpdateAPI",
				Title:          "修改接口",
				Type:           constants.MenuTypeAPI,
				PermissionCode: permissionPrefix + ":update",
				Sort:           103,
				Hidden:         true,
				Remark:         title + "修改接口权限",
			},
			{
				Name:           name + "DeleteAPI",
				Title:          "删除接口",
				Type:           constants.MenuTypeAPI,
				PermissionCode: permissionPrefix + ":delete",
				Sort:           104,
				Hidden:         true,
				Remark:         title + "删除接口权限",
			},
		},
	}
}

func initRoleMenus(ctx context.Context, client *ent.Client, roleID int, menuIDs []int) error {
	for _, menuID := range menuIDs {
		exists, err := client.SysRoleMenu.Query().
			Where(
				sysrolemenu.RoleID(roleID),
				sysrolemenu.MenuID(menuID),
			).
			Exist(ctx)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		err = client.SysRoleMenu.Create().
			SetRoleID(roleID).
			SetMenuID(menuID).
			SetRole(&ent.SysRole{ID: roleID}).
			SetMenu(&ent.SysMenu{ID: menuID}).
			Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func initCasbinRules(ctx context.Context, client *ent.Client, roleID int, permissionCodes []string) error {
	roleSubject := fmt.Sprintf("role:%d", roleID)
	userSubject := "user:admin"

	if err := ensureCasbinRule(ctx, client, "g", userSubject, roleSubject, "", "", "", ""); err != nil {
		return err
	}

	uniqueCodes := map[string]struct{}{}
	for _, code := range permissionCodes {
		if code == "" {
			continue
		}
		uniqueCodes[code] = struct{}{}
	}
	for code := range uniqueCodes {
		if err := ensureCasbinRule(ctx, client, "p", roleSubject, code, "access", "", "", ""); err != nil {
			return err
		}
	}
	return nil
}

func ensureCasbinRule(ctx context.Context, client *ent.Client, ptype, v0, v1, v2, v3, v4, v5 string) error {
	exists, err := client.SysCasbinRule.Query().
		Where(
			syscasbinrule.PtypeEQ(ptype),
			syscasbinrule.V0EQ(v0),
			syscasbinrule.V1EQ(v1),
			syscasbinrule.V2EQ(v2),
			syscasbinrule.V3EQ(v3),
			syscasbinrule.V4EQ(v4),
			syscasbinrule.V5EQ(v5),
		).
		Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	return client.SysCasbinRule.Create().
		SetPtype(ptype).
		SetV0(v0).
		SetV1(v1).
		SetV2(v2).
		SetV3(v3).
		SetV4(v4).
		SetV5(v5).
		Exec(ctx)
}
