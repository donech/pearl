package system

import (
	"context"

	"github.com/donech/pearl/internal/common"
	v1 "github.com/donech/pearl/internal/proto/gen/go/system/v1"
	"github.com/donech/pearl/internal/service/supplier"
	"github.com/donech/pearl/internal/service/supplier/vo"
)

//SupplierList 供货商列表.
func (s *SystemAPI) SupplierList(ctx context.Context, req *v1.SupplierListRequest) (*v1.SupplierListResponse, error) {
	params := supplier.SuppliersParams{
		Name:    req.GetName(),
		Address: req.GetAddress(),
	}
	if req.GetPager() != nil {
		params.Page = int(req.Pager.Page)
		params.Size = int(req.Pager.Size)
	}
	list, err := supplier.Suppliers(ctx, params)
	if err != nil {
		return nil, err
	}
	respList := make([]*v1.Supplier, len(list.List))
	for i, v := range list.List {
		respList[i] = &v1.Supplier{
			Id:          v.ID,
			Name:        v.Name,
			Address:     v.Address,
			CreatedTime: v.CreatedTime.Local().Format("2006-01-02 15:04:05"),
			UpdatedTime: v.UpdatedTime.Local().Format("2006-01-02 15:04:05"),
		}
	}
	return &v1.SupplierListResponse{
		Code: common.SuccessCode,
		Msg:  common.ResponseMsg(common.SuccessCode),
		Data: &v1.SupplierListResponse_Data{
			Total: int32(list.Total),
			List:  respList,
		},
	}, nil
}

//SaveSupplier 创建/更新供货商.
func (s *SystemAPI) SaveSupplier(ctx context.Context, req *v1.SaveSupplierRequest) (*v1.SaveSupplierResponse, error) {
	if req.Name == "" && req.Address == "" {
		return &v1.SaveSupplierResponse{
			Code: common.ErrorCode,
			Msg:  "参数校验失败",
		}, nil
	}
	model := vo.Supplier{
		ID:        req.Id,
		Name:      req.Name,
		Address:   req.Address,
		Constacts: []*vo.Constact{},
	}
	if len(req.Constacts) > 0 {
		cons := make([]*vo.Constact, len(req.Constacts))
		for i, v := range req.Constacts {
			_temp := vo.Constact{
				ID:         v.Id,
				SupplierID: v.SupplierId,
				Firstname:  v.Firstname,
				Lastname:   v.Lastname,
				Title:      v.Title,
				Relations:  []*vo.Relation{},
			}
			if len(v.Relations) > 0 {
				_tempRelation := make([]*vo.Relation, len(v.Relations))
				for j, r := range v.Relations {
					_tempRelation[j] = &vo.Relation{
						ID:        r.Id,
						ContactID: r.ConstactId,
						Type:      int8(r.Type),
						Value:     r.Value,
					}
				}
			}
			cons[i] = &_temp
		}
		model.Constacts = cons
	}
	supplier, err := supplier.SaveSupplier(ctx, &model)
	if err != nil {
		return nil, err
	}
	return &v1.SaveSupplierResponse{
		Code: common.SuccessCode,
		Msg:  common.ResponseMsg(common.SuccessCode),
		Data: &v1.Supplier{
			Id:          supplier.ID,
			Name:        supplier.Name,
			Address:     supplier.Address,
			CreatedTime: supplier.CreatedTime.Local().Format("2006-01-02 15:04:05"),
			UpdatedTime: supplier.UpdatedTime.Local().Format("2006-01-02 15:04:05"),
		},
	}, nil
}

//Supplier 供货商详情.
func (s *SystemAPI) Supplier(ctx context.Context, req *v1.SupplierRequest) (*v1.SupplierResponse, error) {
	model, err := supplier.Supplier(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	output := v1.Supplier{
		Id:          model.ID,
		Name:        model.Name,
		Address:     model.Address,
		CreatedTime: model.CreatedTime.Local().Format("2006-01-02 15:04:05"),
		UpdatedTime: model.UpdatedTime.Local().Format("2006-01-02 15:04:05"),
	}
	if model.DeletedTime != nil {
		output.DeletedTime = model.DeletedTime.Local().Format("2006-01-02 15:04:05")
	}
	if len(model.Edges.Contacts) > 0 {
		cons := make([]*v1.Constact, len(model.Edges.Contacts))
		for i, v := range model.Edges.Contacts {
			cons[i] = &v1.Constact{
				Id:          v.ID,
				SupplierId:  v.SupplierID,
				Title:       v.Title,
				Firstname:   v.Firstname,
				Lastname:    v.Lastname,
				Relations:   []*v1.Relation{},
				CreatedTime: v.CreatedTime.Local().Format("2006-01-02 15:04:05"),
				UpdatedTime: v.CreatedTime.Local().Format("2006-01-02 15:04:05"),
			}
			if len(v.Edges.Relations) > 0 {
				rels := make([]*v1.Relation, len(v.Edges.Relations))
				for i, rel := range v.Edges.Relations {
					rels[i] = &v1.Relation{
						Id:          rel.ID,
						ConstactId:  *rel.ContactID,
						Type:        v1.RelationType(rel.Type),
						Value:       rel.Value,
						CreatedTime: rel.CreatedTime.Local().Format("2006-01-02 15:04:05"),
						UpdatedTime: rel.UpdatedTime.Local().Format("2006-01-02 15:04:05"),
					}
				}
				cons[i].Relations = rels
			}
		}
		output.Constacts = cons
	}
	return &v1.SupplierResponse{
		Code: common.SuccessCode,
		Msg:  common.ResponseMsg(common.SuccessCode),
		Data: &output,
	}, nil
}

//DeleteSupplier 删除供货商.
func (s *SystemAPI) DeleteSupplier(ctx context.Context, req *v1.DeleteSupplierRequest) (*v1.DeleteSupplierResponse, error) {
	err := supplier.DeleteSupplier(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteSupplierResponse{
		Code: common.SuccessCode,
		Msg:  common.ResponseMsg(common.SuccessCode),
	}, nil
}

//SaveSupplierConstact 创建/更新供货商联系人.
func (s *SystemAPI) SaveSupplierConstact(ctx context.Context, req *v1.SaveSupplierConstactRequest) (*v1.SaveSupplierConstactResponse, error) {
	model := vo.Constact{
		ID:         req.Id,
		SupplierID: req.SupplierId,
		Firstname:  req.Firstname,
		Lastname:   req.Lastname,
		Title:      req.Title,
		Relations:  []*vo.Relation{},
	}
	if len(req.Relations) > 0 {
		for _, r := range req.Relations {
			model.Relations = append(model.Relations, &vo.Relation{
				Type:  int8(r.Type),
				Value: r.Value,
			})
		}
	}
	cons, err := supplier.SaveConstact(ctx, &model)
	if err != nil {
		return nil, err
	}
	return &v1.SaveSupplierConstactResponse{
		Code: 0,
		Msg:  "",
		Data: &v1.Constact{
			Id:          cons.ID,
			SupplierId:  cons.SupplierID,
			Title:       cons.Title,
			Firstname:   cons.Firstname,
			Lastname:    cons.Lastname,
			Relations:   []*v1.Relation{},
			CreatedTime: cons.CreatedTime.Local().Format("2006-01-02 15:04:05"),
			UpdatedTime: cons.UpdatedTime.Local().Format("2006-01-02 15:04:05"),
		},
	}, nil
}

//DeleteSupplierConstact 删除供货商联系人.
func (s *SystemAPI) DeleteSupplierConstact(ctx context.Context, req *v1.DeleteSupplierConstactRequest) (*v1.DeleteSupplierConstactResponse, error) {
	err := supplier.DeleteConstact(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteSupplierConstactResponse{
		Code: common.SuccessCode,
		Msg:  common.ResponseMsg(common.SuccessCode),
	}, nil
}

//SaveSupplierConstactRelation 创建/更新供货商联系人联系方式.
func (s *SystemAPI) SaveSupplierConstactRelation(ctx context.Context, req *v1.SaveSupplierConstactRelationRequest) (*v1.SaveSupplierConstactRelationResponse, error) {
	rel, err := supplier.SaveRelation(ctx, &vo.Relation{
		ID:        req.Id,
		ContactID: req.ConstactId,
		Type:      int8(req.Type),
		Value:     req.Value,
	})
	if err != nil {
		return nil, err
	}
	return &v1.SaveSupplierConstactRelationResponse{
		Code: common.SuccessCode,
		Msg:  common.ResponseMsg(common.SuccessCode),
		Data: &v1.Relation{
			Id:          rel.ID,
			ConstactId:  *rel.ContactID,
			Type:        v1.RelationType(rel.Type),
			Value:       rel.Value,
			CreatedTime: rel.CreatedTime.Local().Format("2006-01-02 15:04:05"),
			UpdatedTime: rel.UpdatedTime.Local().Format("2006-01-02 15:04:05"),
		},
	}, nil
}

//DeleteSupplierConstactRelation 删除供货商联系人联系方式.
func (s *SystemAPI) DeleteSupplierConstactRelation(ctx context.Context, req *v1.DeleteSupplierConstactRelationRequest) (*v1.DeleteSupplierConstactRelationResponse, error) {
	err := supplier.DeleteRelation(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteSupplierConstactRelationResponse{
		Code: common.SuccessCode,
		Msg:  common.ResponseMsg(common.SuccessCode),
	}, nil
}
