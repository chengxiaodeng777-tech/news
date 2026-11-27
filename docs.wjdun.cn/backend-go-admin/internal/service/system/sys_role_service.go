package system

import (
	"context"

	"base/ent"
	"base/ent/sysrole"
	"base/ent/sysrolemenu"
	"base/ent/sysuserrole"
	internalcasbin "base/internal/casbin"
	"base/internal/constants"
	systemmodel "base/internal/model/system"
	"base/internal/pkg/apperror"
)

func (s *Service) ListRoles(ctx context.Context) ([]systemmodel.RoleResponse, error) {
	result, err := s.ListRolesPage(ctx, systemmodel.PageQuery{Page: 1, PageSize: 1000})
	if err != nil {
		return nil, err
	}
	list, _ := result.List.([]systemmodel.RoleResponse)
	return list, nil
}

func (s *Service) ListRolesPage(ctx context.Context, query systemmodel.PageQuery) (*systemmodel.PageResult, error) {
	page, pageSize := normalizePage(query.Page, query.PageSize)
	builder := s.client.SysRole.Query()
	if query.Keyword != "" {
		builder.Where(
			sysrole.Or(
				sysrole.CodeContainsFold(query.Keyword),
				sysrole.NameContainsFold(query.Keyword),
			),
		)
	}
	if query.Status != nil {
		builder.Where(sysrole.StatusEQ(*query.Status))
	}

	total, err := builder.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}

	roles, err := builder.
		Order(ent.Asc(sysrole.FieldSort), ent.Asc(sysrole.FieldID)).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]systemmodel.RoleResponse, 0, len(roles))
	for _, role := range roles {
		menuIDs, err := s.client.SysRoleMenu.Query().
			Where(sysrolemenu.RoleID(role.ID)).
			Select(sysrolemenu.FieldMenuID).
			Ints(ctx)
		if err != nil {
			return nil, err
		}

		result = append(result, systemmodel.RoleResponse{
			ID:      role.ID,
			Code:    role.Code,
			Name:    role.Name,
			Status:  role.Status,
			Sort:    role.Sort,
			MenuIDs: menuIDs,
			Remark:  role.Remark,
		})
	}
	return &systemmodel.PageResult{
		List:     result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *Service) GetRole(ctx context.Context, id int) (*systemmodel.RoleResponse, error) {
	role, err := s.client.SysRole.Query().
		Where(sysrole.IDEQ(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	menuIDs, err := s.client.SysRoleMenu.Query().
		Where(sysrolemenu.RoleID(role.ID)).
		Select(sysrolemenu.FieldMenuID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}

	return &systemmodel.RoleResponse{
		ID:      role.ID,
		Code:    role.Code,
		Name:    role.Name,
		Status:  role.Status,
		Sort:    role.Sort,
		MenuIDs: menuIDs,
		Remark:  role.Remark,
	}, nil
}

func (s *Service) CreateRole(ctx context.Context, req systemmodel.CreateRoleRequest) error {
	exists, err := s.client.SysRole.Query().
		Where(sysrole.CodeEQ(req.Code)).
		Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return apperror.New(apperror.CodeDuplicateResource, "角色编码已存在")
	}

	return s.client.SysRole.Create().
		SetCode(req.Code).
		SetName(req.Name).
		SetStatus(req.Status).
		SetSort(req.Sort).
		SetRemark(req.Remark).
		Exec(ctx)
}

func (s *Service) UpdateRole(ctx context.Context, id int, req systemmodel.UpdateRoleRequest) error {
	role, err := s.client.SysRole.Get(ctx, id)
	if err != nil {
		return err
	}
	if role.Code == constants.SystemAdminRoleCode && req.Status == constants.StatusDisabled {
		return apperror.New(apperror.CodeBuiltinProtect, "系统管理员角色不允许停用")
	}
	if role.Code == constants.SystemAdminRoleCode && req.Name != role.Name {
		return apperror.New(apperror.CodeBuiltinProtect, "系统管理员角色名称不允许修改")
	}

	return s.client.SysRole.UpdateOneID(id).
		SetName(req.Name).
		SetStatus(req.Status).
		SetSort(req.Sort).
		SetRemark(req.Remark).
		Exec(ctx)
}

func (s *Service) DeleteRole(ctx context.Context, id int) error {
	role, err := s.client.SysRole.Get(ctx, id)
	if err != nil {
		return err
	}
	if role.Code == constants.SystemAdminRoleCode {
		return apperror.New(apperror.CodeBuiltinProtect, "系统管理员角色不允许删除")
	}

	if _, err := s.client.SysUserRole.Delete().Where(sysuserrole.RoleID(id)).Exec(ctx); err != nil {
		return err
	}
	if _, err := s.client.SysRoleMenu.Delete().Where(sysrolemenu.RoleID(id)).Exec(ctx); err != nil {
		return err
	}
	return s.client.SysRole.DeleteOneID(id).Exec(ctx)
}

func (s *Service) AssignRoleMenus(ctx context.Context, roleID int, menuIDs []int) error {
	if _, err := s.client.SysRoleMenu.Delete().Where(sysrolemenu.RoleID(roleID)).Exec(ctx); err != nil {
		return err
	}
	for _, menuID := range menuIDs {
		if err := s.client.SysRoleMenu.Create().
			SetRoleID(roleID).
			SetMenuID(menuID).
			SetRole(&ent.SysRole{ID: roleID}).
			SetMenu(&ent.SysMenu{ID: menuID}).
			Exec(ctx); err != nil {
			return err
		}
	}
	return internalcasbin.SyncRolePolicies(ctx, s.client, roleID)
}
