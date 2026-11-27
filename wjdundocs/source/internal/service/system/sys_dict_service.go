package system

import (
	"context"

	"base/ent"
	"base/ent/sysdictdata"
	"base/ent/sysdicttype"
	systemmodel "base/internal/model/system"
	"base/internal/pkg/apperror"
)

func (s *Service) ListDictTypes(ctx context.Context) ([]systemmodel.DictTypeResponse, error) {
	result, err := s.ListDictTypesPage(ctx, systemmodel.PageQuery{Page: 1, PageSize: 1000})
	if err != nil {
		return nil, err
	}
	list, _ := result.List.([]systemmodel.DictTypeResponse)
	return list, nil
}

func (s *Service) ListDictTypesPage(ctx context.Context, query systemmodel.PageQuery) (*systemmodel.PageResult, error) {
	page, pageSize := normalizePage(query.Page, query.PageSize)
	builder := s.client.SysDictType.Query()
	if query.Keyword != "" {
		builder.Where(
			sysdicttype.Or(
				sysdicttype.CodeContainsFold(query.Keyword),
				sysdicttype.NameContainsFold(query.Keyword),
			),
		)
	}
	if query.Status != nil {
		builder.Where(sysdicttype.StatusEQ(*query.Status))
	}

	total, err := builder.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}

	dictTypes, err := builder.
		Order(ent.Asc(sysdicttype.FieldID)).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]systemmodel.DictTypeResponse, 0, len(dictTypes))
	for _, item := range dictTypes {
		result = append(result, systemmodel.DictTypeResponse{
			ID:     item.ID,
			Code:   item.Code,
			Name:   item.Name,
			Status: item.Status,
			Remark: item.Remark,
		})
	}
	return &systemmodel.PageResult{
		List:     result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *Service) GetDictType(ctx context.Context, id int) (*systemmodel.DictTypeResponse, error) {
	item, err := s.client.SysDictType.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &systemmodel.DictTypeResponse{
		ID:     item.ID,
		Code:   item.Code,
		Name:   item.Name,
		Status: item.Status,
		Remark: item.Remark,
	}, nil
}

func (s *Service) ListDictData(ctx context.Context, dictTypeID int) ([]systemmodel.DictDataResponse, error) {
	result, err := s.ListDictDataPage(ctx, systemmodel.PageQuery{Page: 1, PageSize: 1000}, dictTypeID)
	if err != nil {
		return nil, err
	}
	list, _ := result.List.([]systemmodel.DictDataResponse)
	return list, nil
}

func (s *Service) ListDictDataPage(ctx context.Context, query systemmodel.PageQuery, dictTypeID int) (*systemmodel.PageResult, error) {
	page, pageSize := normalizePage(query.Page, query.PageSize)
	builder := s.client.SysDictData.Query()
	if dictTypeID > 0 {
		builder.Where(sysdictdata.DictTypeID(dictTypeID))
	}
	if query.Keyword != "" {
		builder.Where(
			sysdictdata.Or(
				sysdictdata.LabelContainsFold(query.Keyword),
				sysdictdata.ValueContainsFold(query.Keyword),
			),
		)
	}
	if query.Status != nil {
		builder.Where(sysdictdata.StatusEQ(*query.Status))
	}

	total, err := builder.Clone().Count(ctx)
	if err != nil {
		return nil, err
	}
	list, err := builder.
		Order(ent.Asc(sysdictdata.FieldSort), ent.Asc(sysdictdata.FieldID)).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]systemmodel.DictDataResponse, 0, len(list))
	for _, item := range list {
		result = append(result, systemmodel.DictDataResponse{
			ID:         item.ID,
			DictTypeID: item.DictTypeID,
			Label:      item.Label,
			Value:      item.Value,
			TagType:    item.TagType,
			CSSClass:   item.CSSClass,
			Sort:       item.Sort,
			Status:     item.Status,
			IsDefault:  item.IsDefault,
			Remark:     item.Remark,
		})
	}
	return &systemmodel.PageResult{
		List:     result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *Service) GetDictData(ctx context.Context, id int) (*systemmodel.DictDataResponse, error) {
	item, err := s.client.SysDictData.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &systemmodel.DictDataResponse{
		ID:         item.ID,
		DictTypeID: item.DictTypeID,
		Label:      item.Label,
		Value:      item.Value,
		TagType:    item.TagType,
		CSSClass:   item.CSSClass,
		Sort:       item.Sort,
		Status:     item.Status,
		IsDefault:  item.IsDefault,
		Remark:     item.Remark,
	}, nil
}

func (s *Service) CreateDictType(ctx context.Context, req systemmodel.CreateDictTypeRequest) error {
	exists, err := s.client.SysDictType.Query().
		Where(sysdicttype.CodeEQ(req.Code)).
		Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return apperror.New(apperror.CodeDuplicateResource, "字典类型编码已存在")
	}

	return s.client.SysDictType.Create().
		SetCode(req.Code).
		SetName(req.Name).
		SetStatus(req.Status).
		SetRemark(req.Remark).
		Exec(ctx)
}

func (s *Service) UpdateDictType(ctx context.Context, id int, req systemmodel.UpdateDictTypeRequest) error {
	exists, err := s.client.SysDictType.Query().
		Where(
			sysdicttype.CodeEQ(req.Code),
			sysdicttype.IDNEQ(id),
		).
		Exist(ctx)
	if err != nil {
		return err
	}
	if exists {
		return apperror.New(apperror.CodeDuplicateResource, "字典类型编码已存在")
	}

	return s.client.SysDictType.UpdateOneID(id).
		SetCode(req.Code).
		SetName(req.Name).
		SetStatus(req.Status).
		SetRemark(req.Remark).
		Exec(ctx)
}

func (s *Service) DeleteDictType(ctx context.Context, id int) error {
	count, err := s.client.SysDictData.Query().Where(sysdictdata.DictTypeID(id)).Count(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return apperror.New(apperror.CodeDeleteRestricted, "当前字典类型存在字典数据，不能删除")
	}
	return s.client.SysDictType.DeleteOneID(id).Exec(ctx)
}

func (s *Service) CreateDictData(ctx context.Context, req systemmodel.CreateDictDataRequest) error {
	return s.client.SysDictData.Create().
		SetDictTypeID(req.DictTypeID).
		SetDictType(&ent.SysDictType{ID: req.DictTypeID}).
		SetLabel(req.Label).
		SetValue(req.Value).
		SetTagType(req.TagType).
		SetCSSClass(req.CSSClass).
		SetSort(req.Sort).
		SetStatus(req.Status).
		SetIsDefault(req.IsDefault).
		SetRemark(req.Remark).
		Exec(ctx)
}

func (s *Service) UpdateDictData(ctx context.Context, id int, req systemmodel.UpdateDictDataRequest) error {
	return s.client.SysDictData.UpdateOneID(id).
		SetDictTypeID(req.DictTypeID).
		SetLabel(req.Label).
		SetValue(req.Value).
		SetTagType(req.TagType).
		SetCSSClass(req.CSSClass).
		SetSort(req.Sort).
		SetStatus(req.Status).
		SetIsDefault(req.IsDefault).
		SetRemark(req.Remark).
		Exec(ctx)
}

func (s *Service) DeleteDictData(ctx context.Context, id int) error {
	return s.client.SysDictData.DeleteOneID(id).Exec(ctx)
}
