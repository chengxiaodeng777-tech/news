package system

import (
	"context"

	"base/ent"
	"base/ent/sysmenu"
	"base/ent/sysrolemenu"
	"base/internal/constants"
	systemmodel "base/internal/model/system"
	"base/internal/pkg/apperror"
)

func menuParentID(id *int) int {
	if id == nil {
		return 0
	}
	return *id
}

func (s *Service) ListMenus(ctx context.Context) ([]systemmodel.MenuResponse, error) {
	menus, err := s.client.SysMenu.Query().
		Order(ent.Asc(sysmenu.FieldSort), ent.Asc(sysmenu.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]systemmodel.MenuResponse, 0, len(menus))
	for _, menu := range menus {
		result = append(result, systemmodel.MenuResponse{
			ID:             menu.ID,
			ParentID:       menuParentID(menu.ParentID),
			Name:           menu.Name,
			Title:          menu.Title,
			Path:           menu.Path,
			Component:      menu.Component,
			Icon:           menu.Icon,
			Type:           menu.Type,
			PermissionCode: menu.PermissionCode,
			Redirect:       menu.Redirect,
			Sort:           menu.Sort,
			Status:         menu.Status,
			Hidden:         menu.Hidden,
			KeepAlive:      menu.KeepAlive,
			ExternalLink:   menu.ExternalLink,
			Remark:         menu.Remark,
		})
	}
	return result, nil
}

func (s *Service) ListMenuTree(ctx context.Context) ([]systemmodel.MenuResponse, error) {
	menus, err := s.client.SysMenu.Query().
		Order(ent.Asc(sysmenu.FieldSort), ent.Asc(sysmenu.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	menuMap := make(map[int]systemmodel.MenuResponse, len(menus))
	childrenMap := make(map[int][]systemmodel.MenuResponse)
	var roots []systemmodel.MenuResponse

	for _, menu := range menus {
		menuMap[menu.ID] = systemmodel.MenuResponse{
			ID:             menu.ID,
			ParentID:       menuParentID(menu.ParentID),
			Name:           menu.Name,
			Title:          menu.Title,
			Path:           menu.Path,
			Component:      menu.Component,
			Icon:           menu.Icon,
			Type:           menu.Type,
			PermissionCode: menu.PermissionCode,
			Redirect:       menu.Redirect,
			Sort:           menu.Sort,
			Status:         menu.Status,
			Hidden:         menu.Hidden,
			KeepAlive:      menu.KeepAlive,
			ExternalLink:   menu.ExternalLink,
			Remark:         menu.Remark,
		}
	}

	for _, menu := range menus {
		node := menuMap[menu.ID]
		if menu.ParentID == nil || *menu.ParentID == 0 {
			roots = append(roots, node)
			continue
		}
		childrenMap[*menu.ParentID] = append(childrenMap[*menu.ParentID], node)
	}

	var build func([]systemmodel.MenuResponse) []systemmodel.MenuResponse
	build = func(nodes []systemmodel.MenuResponse) []systemmodel.MenuResponse {
		result := make([]systemmodel.MenuResponse, 0, len(nodes))
		for _, node := range nodes {
			node.Children = build(childrenMap[node.ID])
			result = append(result, node)
		}
		return result
	}

	return build(roots), nil
}

func (s *Service) GetMenu(ctx context.Context, id int) (*systemmodel.MenuResponse, error) {
	menu, err := s.client.SysMenu.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &systemmodel.MenuResponse{
		ID:             menu.ID,
		ParentID:       menuParentID(menu.ParentID),
		Name:           menu.Name,
		Title:          menu.Title,
		Path:           menu.Path,
		Component:      menu.Component,
		Icon:           menu.Icon,
		Type:           menu.Type,
		PermissionCode: menu.PermissionCode,
		Redirect:       menu.Redirect,
		Sort:           menu.Sort,
		Status:         menu.Status,
		Hidden:         menu.Hidden,
		KeepAlive:      menu.KeepAlive,
		ExternalLink:   menu.ExternalLink,
		Remark:         menu.Remark,
	}, nil
}

func (s *Service) CreateMenu(ctx context.Context, req systemmodel.CreateMenuRequest) error {
	if err := s.validateMenuRequest(ctx, req, 0); err != nil {
		return err
	}

	create := s.client.SysMenu.Create().
		SetName(req.Name).
		SetTitle(req.Title).
		SetPath(req.Path).
		SetComponent(req.Component).
		SetIcon(req.Icon).
		SetType(req.Type).
		SetPermissionCode(req.PermissionCode).
		SetRedirect(req.Redirect).
		SetSort(req.Sort).
		SetStatus(req.Status).
		SetHidden(req.Hidden).
		SetKeepAlive(req.KeepAlive).
		SetExternalLink(req.ExternalLink).
		SetRemark(req.Remark)
	if req.ParentID > 0 {
		create.SetParentID(req.ParentID)
	}
	return create.Exec(ctx)
}

func (s *Service) UpdateMenu(ctx context.Context, id int, req systemmodel.UpdateMenuRequest) error {
	menu, err := s.client.SysMenu.Get(ctx, id)
	if err != nil {
		return err
	}
	if menu.Name == constants.SystemMenuRootName && req.Type != menu.Type {
		return apperror.New(apperror.CodeBuiltinProtect, "系统根菜单类型不允许修改")
	}
	if menu.Name == constants.SystemMenuRootName && req.Name != menu.Name {
		return apperror.New(apperror.CodeBuiltinProtect, "系统根菜单名称不允许修改")
	}
	if menu.Name == constants.SystemMenuRootName && req.Path != menu.Path {
		return apperror.New(apperror.CodeBuiltinProtect, "系统根菜单路由不允许修改")
	}
	if err := s.validateMenuRequest(ctx, req, id); err != nil {
		return err
	}

	update := s.client.SysMenu.UpdateOneID(id).
		SetName(req.Name).
		SetTitle(req.Title).
		SetPath(req.Path).
		SetComponent(req.Component).
		SetIcon(req.Icon).
		SetType(req.Type).
		SetPermissionCode(req.PermissionCode).
		SetRedirect(req.Redirect).
		SetSort(req.Sort).
		SetStatus(req.Status).
		SetHidden(req.Hidden).
		SetKeepAlive(req.KeepAlive).
		SetExternalLink(req.ExternalLink).
		SetRemark(req.Remark)
	if req.ParentID > 0 {
		update.SetParentID(req.ParentID)
	}
	return update.Exec(ctx)
}

func (s *Service) DeleteMenu(ctx context.Context, id int) error {
	menu, err := s.client.SysMenu.Get(ctx, id)
	if err != nil {
		return err
	}
	if menu.Name == constants.SystemMenuRootName {
		return apperror.New(apperror.CodeBuiltinProtect, "系统根菜单不允许删除")
	}

	childCount, err := s.client.SysMenu.Query().Where(sysmenu.ParentID(id)).Count(ctx)
	if err != nil {
		return err
	}
	if childCount > 0 {
		return apperror.New(apperror.CodeDeleteRestricted, "当前菜单存在子节点，不能删除")
	}
	if _, err := s.client.SysRoleMenu.Delete().Where(sysrolemenu.MenuID(id)).Exec(ctx); err != nil {
		return err
	}
	return s.client.SysMenu.DeleteOneID(id).Exec(ctx)
}

func (s *Service) validateMenuRequest(ctx context.Context, req systemmodel.CreateMenuRequest, currentID int) error {
	switch req.Type {
	case constants.MenuTypeCatalog:
		if req.Name == "" || req.Title == "" {
			return apperror.New(apperror.CodeInvalidMenuRule, "目录必须填写名称和标题")
		}
		if req.PermissionCode != "" {
			return apperror.New(apperror.CodeInvalidMenuRule, "目录不允许设置权限标识")
		}
	case constants.MenuTypeMenu:
		if req.Name == "" || req.Title == "" {
			return apperror.New(apperror.CodeInvalidMenuRule, "菜单必须填写名称和标题")
		}
		if req.Path == "" {
			return apperror.New(apperror.CodeInvalidMenuRule, "菜单必须填写路由路径")
		}
		if req.Component == "" {
			return apperror.New(apperror.CodeInvalidMenuRule, "菜单必须填写组件路径")
		}
	case constants.MenuTypeButton:
		if req.Name == "" || req.Title == "" {
			return apperror.New(apperror.CodeInvalidMenuRule, "按钮必须填写名称和标题")
		}
		if req.PermissionCode == "" {
			return apperror.New(apperror.CodeInvalidMenuRule, "按钮必须填写权限标识")
		}
	case constants.MenuTypeAPI:
		if req.Name == "" || req.Title == "" {
			return apperror.New(apperror.CodeInvalidMenuRule, "接口权限必须填写名称和标题")
		}
		if req.PermissionCode == "" {
			return apperror.New(apperror.CodeInvalidMenuRule, "接口权限必须填写权限标识")
		}
	default:
		return apperror.New(apperror.CodeInvalidMenuRule, "菜单类型不合法")
	}

	if req.ParentID > 0 {
		parent, err := s.client.SysMenu.Get(ctx, req.ParentID)
		if err != nil {
			return err
		}
		switch req.Type {
		case constants.MenuTypeCatalog, constants.MenuTypeMenu:
			if parent.Type != constants.MenuTypeCatalog {
				return apperror.New(apperror.CodeInvalidMenuTree, "目录或菜单只能挂在目录下")
			}
		case constants.MenuTypeButton, constants.MenuTypeAPI:
			if parent.Type != constants.MenuTypeMenu {
				return apperror.New(apperror.CodeInvalidMenuTree, "按钮或接口权限只能挂在菜单下")
			}
		}
	}

	if currentID > 0 && req.ParentID == currentID {
		return apperror.New(apperror.CodeInvalidMenuTree, "菜单不允许设置自己为父节点")
	}

	return nil
}
