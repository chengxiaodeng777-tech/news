package system

import (
	"context"

	"base/ent"
	"base/ent/sysuser"
	"base/ent/sysuserrole"
	internalcasbin "base/internal/casbin"
	"base/internal/constants"
	systemmodel "base/internal/model/system"
	"base/internal/pkg/apperror"
	"base/internal/pkg/password"
)

func organizationID(id *int) int {
	if id == nil {
		return 0
	}
	return *id
}

func (s *Service) ListUsers(ctx context.Context) ([]systemmodel.UserResponse, error) {
	result, err := s.ListUsersPage(ctx, systemmodel.PageQuery{Page: 1, PageSize: 1000})
	if err != nil {
		return nil, err
	}
	list, _ := result.List.([]systemmodel.UserResponse)
	return list, nil
}

func (s *Service) ListUsersPage(ctx context.Context, query systemmodel.PageQuery) (*systemmodel.PageResult, error) {
	page, pageSize := normalizePage(query.Page, query.PageSize)
	builder := s.client.SysUser.Query().WithOrganization()
	if query.Keyword != "" {
		builder.Where(
			sysuser.Or(
				sysuser.UsernameContainsFold(query.Keyword),
				sysuser.RealNameContainsFold(query.Keyword),
				sysuser.NicknameContainsFold(query.Keyword),
			),
		)
	}
	if query.Status != nil {
		builder.Where(sysuser.StatusEQ(*query.Status))
	}

	total, err := builder.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}

	users, err := builder.
		Order(ent.Asc(sysuser.FieldID)).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]systemmodel.UserResponse, 0, len(users))
	for _, user := range users {
		roleIDs, err := s.client.SysUserRole.Query().
			Where(sysuserrole.UserID(user.ID)).
			Select(sysuserrole.FieldRoleID).
			Ints(ctx)
		if err != nil {
			return nil, err
		}

		orgName := ""
		if user.Edges.Organization != nil {
			orgName = user.Edges.Organization.Name
		}

		result = append(result, systemmodel.UserResponse{
			ID:               user.ID,
			Username:         user.Username,
			RealName:         user.RealName,
			Nickname:         user.Nickname,
			Mobile:           user.Mobile,
			Email:            user.Email,
			Avatar:           user.Avatar,
			Status:           user.Status,
			OrganizationID:   organizationID(user.OrganizationID),
			OrganizationName: orgName,
			RoleIDs:          roleIDs,
			Remark:           user.Remark,
		})
	}
	return &systemmodel.PageResult{
		List:     result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *Service) GetUser(ctx context.Context, id int) (*systemmodel.UserResponse, error) {
	user, err := s.client.SysUser.Query().
		Where(sysuser.IDEQ(id)).
		WithOrganization().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	roleIDs, err := s.client.SysUserRole.Query().
		Where(sysuserrole.UserID(user.ID)).
		Select(sysuserrole.FieldRoleID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}

	orgName := ""
	if user.Edges.Organization != nil {
		orgName = user.Edges.Organization.Name
	}

	return &systemmodel.UserResponse{
		ID:               user.ID,
		Username:         user.Username,
		RealName:         user.RealName,
		Nickname:         user.Nickname,
		Mobile:           user.Mobile,
		Email:            user.Email,
		Avatar:           user.Avatar,
		Status:           user.Status,
		OrganizationID:   organizationID(user.OrganizationID),
		OrganizationName: orgName,
		RoleIDs:          roleIDs,
		Remark:           user.Remark,
	}, nil
}

func (s *Service) CreateUser(ctx context.Context, req systemmodel.CreateUserRequest) error {
	exists, err := s.client.SysUser.Query().
		Where(sysuser.UsernameEQ(req.Username)).
		Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return apperror.New(apperror.CodeDuplicateResource, "用户名已存在")
	}

	passwordHash, err := password.Hash(req.Password)
	if err != nil {
		return err
	}
	create := s.client.SysUser.Create().
		SetUsername(req.Username).
		SetPasswordHash(passwordHash).
		SetRealName(req.RealName).
		SetNickname(req.Nickname).
		SetMobile(req.Mobile).
		SetEmail(req.Email).
		SetAvatar(req.Avatar).
		SetStatus(req.Status).
		SetRemark(req.Remark)
	if req.OrganizationID > 0 {
		create.SetOrganizationID(req.OrganizationID)
	}
	return create.Exec(ctx)
}

func (s *Service) UpdateUser(ctx context.Context, id int, req systemmodel.UpdateUserRequest) error {
	user, err := s.client.SysUser.Get(ctx, id)
	if err != nil {
		return err
	}
	if user.Username == constants.SystemAdminUsername && req.Status == constants.StatusDisabled {
		return apperror.New(apperror.CodeBuiltinProtect, "系统管理员账号不允许停用")
	}

	update := s.client.SysUser.UpdateOneID(id).
		SetRealName(req.RealName).
		SetNickname(req.Nickname).
		SetMobile(req.Mobile).
		SetEmail(req.Email).
		SetAvatar(req.Avatar).
		SetStatus(req.Status).
		SetRemark(req.Remark)
	if req.OrganizationID > 0 {
		update.SetOrganizationID(req.OrganizationID)
	}
	return update.Exec(ctx)
}

func (s *Service) DeleteUser(ctx context.Context, id int) error {
	user, err := s.client.SysUser.Get(ctx, id)
	if err != nil {
		return err
	}
	if user.Username == constants.SystemAdminUsername {
		return apperror.New(apperror.CodeBuiltinProtect, "系统管理员账号不允许删除")
	}

	if _, err := s.client.SysUserRole.Delete().Where(sysuserrole.UserID(id)).Exec(ctx); err != nil {
		return err
	}
	return s.client.SysUser.DeleteOneID(id).Exec(ctx)
}

func (s *Service) ResetUserPassword(ctx context.Context, id int, newPassword string) error {
	passwordHash, err := password.Hash(newPassword)
	if err != nil {
		return err
	}
	return s.client.SysUser.UpdateOneID(id).SetPasswordHash(passwordHash).Exec(ctx)
}

func (s *Service) AssignUserRoles(ctx context.Context, userID int, roleIDs []int) error {
	if _, err := s.client.SysUserRole.Delete().Where(sysuserrole.UserID(userID)).Exec(ctx); err != nil {
		return err
	}
	for _, roleID := range roleIDs {
		if err := s.client.SysUserRole.Create().
			SetUserID(userID).
			SetRoleID(roleID).
			SetUser(&ent.SysUser{ID: userID}).
			SetRole(&ent.SysRole{ID: roleID}).
			Exec(ctx); err != nil {
			return err
		}
	}
	return internalcasbin.SyncUserRoles(ctx, s.client, userID)
}
