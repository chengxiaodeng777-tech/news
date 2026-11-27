package system

import (
	"context"

	"base/ent"
	"base/ent/sysorganization"
	"base/ent/sysuser"
	"base/internal/constants"
	systemmodel "base/internal/model/system"
	"base/internal/pkg/apperror"
)

func organizationParentID(id *int) int {
	if id == nil {
		return 0
	}
	return *id
}

func (s *Service) ListOrganizations(ctx context.Context) ([]systemmodel.OrganizationResponse, error) {
	orgs, err := s.client.SysOrganization.Query().
		Order(ent.Asc(sysorganization.FieldSort), ent.Asc(sysorganization.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]systemmodel.OrganizationResponse, 0, len(orgs))
	for _, org := range orgs {
		result = append(result, systemmodel.OrganizationResponse{
			ID:       org.ID,
			ParentID: organizationParentID(org.ParentID),
			Name:     org.Name,
			Code:     org.Code,
			Type:     org.Type,
			Leader:   org.Leader,
			Phone:    org.Phone,
			Email:    org.Email,
			Sort:     org.Sort,
			Status:   org.Status,
			Remark:   org.Remark,
		})
	}
	return result, nil
}

func (s *Service) ListOrganizationTree(ctx context.Context) ([]systemmodel.OrganizationResponse, error) {
	orgs, err := s.client.SysOrganization.Query().
		Order(ent.Asc(sysorganization.FieldSort), ent.Asc(sysorganization.FieldID)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	orgMap := make(map[int]systemmodel.OrganizationResponse, len(orgs))
	childrenMap := make(map[int][]systemmodel.OrganizationResponse)
	var roots []systemmodel.OrganizationResponse

	for _, org := range orgs {
		orgMap[org.ID] = systemmodel.OrganizationResponse{
			ID:       org.ID,
			ParentID: organizationParentID(org.ParentID),
			Name:     org.Name,
			Code:     org.Code,
			Type:     org.Type,
			Leader:   org.Leader,
			Phone:    org.Phone,
			Email:    org.Email,
			Sort:     org.Sort,
			Status:   org.Status,
			Remark:   org.Remark,
		}
	}

	for _, org := range orgs {
		node := orgMap[org.ID]
		if org.ParentID == nil || *org.ParentID == 0 {
			roots = append(roots, node)
			continue
		}
		childrenMap[*org.ParentID] = append(childrenMap[*org.ParentID], node)
	}

	var build func([]systemmodel.OrganizationResponse) []systemmodel.OrganizationResponse
	build = func(nodes []systemmodel.OrganizationResponse) []systemmodel.OrganizationResponse {
		result := make([]systemmodel.OrganizationResponse, 0, len(nodes))
		for _, node := range nodes {
			node.Children = build(childrenMap[node.ID])
			result = append(result, node)
		}
		return result
	}

	return build(roots), nil
}

func (s *Service) GetOrganization(ctx context.Context, id int) (*systemmodel.OrganizationResponse, error) {
	org, err := s.client.SysOrganization.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &systemmodel.OrganizationResponse{
		ID:       org.ID,
		ParentID: organizationParentID(org.ParentID),
		Name:     org.Name,
		Code:     org.Code,
		Type:     org.Type,
		Leader:   org.Leader,
		Phone:    org.Phone,
		Email:    org.Email,
		Sort:     org.Sort,
		Status:   org.Status,
		Remark:   org.Remark,
	}, nil
}

func (s *Service) CreateOrganization(ctx context.Context, req systemmodel.CreateOrganizationRequest) error {
	exists, err := s.client.SysOrganization.Query().
		Where(sysorganization.CodeEQ(req.Code)).
		Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return apperror.New(apperror.CodeDuplicateResource, "组织编码已存在")
	}

	create := s.client.SysOrganization.Create().
		SetName(req.Name).
		SetCode(req.Code).
		SetType(req.Type).
		SetLeader(req.Leader).
		SetPhone(req.Phone).
		SetEmail(req.Email).
		SetSort(req.Sort).
		SetStatus(req.Status).
		SetRemark(req.Remark)
	if req.ParentID > 0 {
		create.SetParentID(req.ParentID)
	}
	return create.Exec(ctx)
}

func (s *Service) UpdateOrganization(ctx context.Context, id int, req systemmodel.UpdateOrganizationRequest) error {
	org, err := s.client.SysOrganization.Get(ctx, id)
	if err != nil {
		return err
	}
	if org.Code == constants.SystemRootOrganizationCode && req.Status == constants.StatusDisabled {
		return apperror.New(apperror.CodeBuiltinProtect, "根组织不允许停用")
	}
	if org.Code == constants.SystemRootOrganizationCode && req.Code != org.Code {
		return apperror.New(apperror.CodeBuiltinProtect, "根组织编码不允许修改")
	}
	if org.Code == constants.SystemRootOrganizationCode && req.Type != org.Type {
		return apperror.New(apperror.CodeBuiltinProtect, "根组织类型不允许修改")
	}
	if org.Code == constants.SystemRootOrganizationCode && req.Name != org.Name {
		return apperror.New(apperror.CodeBuiltinProtect, "根组织名称不允许修改")
	}

	exists, err := s.client.SysOrganization.Query().
		Where(
			sysorganization.CodeEQ(req.Code),
			sysorganization.IDNEQ(id),
		).
		Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return apperror.New(apperror.CodeDuplicateResource, "组织编码已存在")
	}

	update := s.client.SysOrganization.UpdateOneID(id).
		SetName(req.Name).
		SetCode(req.Code).
		SetType(req.Type).
		SetLeader(req.Leader).
		SetPhone(req.Phone).
		SetEmail(req.Email).
		SetSort(req.Sort).
		SetStatus(req.Status).
		SetRemark(req.Remark)
	if req.ParentID > 0 {
		update.SetParentID(req.ParentID)
	}
	return update.Exec(ctx)
}

func (s *Service) DeleteOrganization(ctx context.Context, id int) error {
	org, err := s.client.SysOrganization.Get(ctx, id)
	if err != nil {
		return err
	}
	if org.Code == constants.SystemRootOrganizationCode {
		return apperror.New(apperror.CodeBuiltinProtect, "根组织不允许删除")
	}

	childCount, err := s.client.SysOrganization.Query().Where(sysorganization.ParentID(id)).Count(ctx)
	if err != nil {
		return err
	}
	if childCount > 0 {
		return apperror.New(apperror.CodeDeleteRestricted, "当前组织存在子组织，不能删除")
	}
	userCount, err := s.client.SysUser.Query().Where(sysuser.OrganizationID(id)).Count(ctx)
	if err != nil {
		return err
	}
	if userCount > 0 {
		return apperror.New(apperror.CodeDeleteRestricted, "当前组织下存在用户，不能删除")
	}
	return s.client.SysOrganization.DeleteOneID(id).Exec(ctx)
}
