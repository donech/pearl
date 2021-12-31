package supplier

import (
	"context"
	"fmt"
	"time"

	"github.com/donech/pearl/internal/data"
	"github.com/donech/pearl/internal/data/ent"
	"github.com/donech/pearl/internal/data/ent/predicate"
	"github.com/donech/pearl/internal/data/ent/supplier"
	"github.com/donech/pearl/internal/data/ent/suppliercontact"
	"github.com/donech/pearl/internal/data/ent/suppliercontactrelation"
	"github.com/donech/pearl/internal/data/ent/user"
	"github.com/donech/pearl/internal/service/supplier/vo"
)

type SuppliersParams struct {
	Name    string
	Address string
	Page    int
	Size    int
}

type SuppliersList struct {
	List  []*ent.Supplier
	Total int
}

func Suppliers(ctx context.Context, params SuppliersParams) (*SuppliersList, error) {
	if params.Page <= 0 {
		params.Page = 1
	}

	if params.Size <= 0 {
		params.Size = 10
	}

	where := make([]predicate.Supplier, 0)
	if len(params.Name) > 0 {
		where = append(where, supplier.NameContains(params.Name))
	}
	if len(params.Address) > 0 {
		where = append(where, supplier.AddressContains(params.Address))
	}
	where = append(where, supplier.DeletedTimeIsNil())
	builder := data.DB.Supplier.Query().Unique(false).Where(where...)
	res := new(SuppliersList)
	total, err := builder.Clone().Count(ctx)
	if err != nil {
		return res, err
	}
	res.Total = total
	lists, err := builder.Offset((params.Page - 1) * params.Size).Limit(params.Size).Order(ent.Desc(user.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}
	res.List = lists
	return res, nil
}

func Supplier(ctx context.Context, id int64) (*ent.Supplier, error) {
	return data.DB.Supplier.Query().Where(supplier.IDEQ(id)).WithContacts(func(scq *ent.SupplierContactQuery) {
		scq.WithRelations(func(scrq *ent.SupplierContactRelationQuery) {
			scrq.Where(suppliercontactrelation.DeletedTimeIsNil())
		}).Where(suppliercontact.DeletedTimeIsNil())
	}).First(ctx)
}

func SaveSupplier(ctx context.Context, model *vo.Supplier) (sup *ent.Supplier, err error) {
	tx, err := data.DB.Tx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			if rerr := tx.Rollback(); rerr != nil {
				err = fmt.Errorf("%w: %v", err, rerr)
			}
			return
		}
		err = tx.Commit()
	}()
	if model.ID == 0 {
		return CreateSupplier(ctx, tx, model)
	}
	return UpdateSupplier(ctx, tx, model)
}

//UpdateSupplier 目前只支持添加和移除 constact, 不支持修改 constact 信息
func UpdateSupplier(ctx context.Context, tx *ent.Tx, model *vo.Supplier) (*ent.Supplier, error) {
	builder := tx.Supplier.UpdateOneID(model.ID)
	if model.Address != "" {
		builder.SetAddress(model.Address)
	}
	if model.Name != "" {
		builder.SetName(model.Name)
	}
	if len(model.Constacts) > 0 {
		newConstacts := make([]*vo.Constact, 0)
		for _, v := range model.Constacts {
			if v.ID == 0 {
				newConstacts = append(newConstacts, v)
			}
		}
		cons, err := CreateConstacts(ctx, tx, newConstacts)
		if err != nil {
			return nil, err
		}
		builder.AddContacts(cons...)
	}
	if len(model.RemoveConstacts) > 0 {
		consIDs := make([]int64, len(model.RemoveConstacts))
		for _, v := range model.RemoveConstacts {
			consIDs = append(consIDs, v.ID)
		}
		builder.RemoveContactIDs(consIDs...)
	}
	return builder.Save(ctx)
}

func CreateSupplier(ctx context.Context, tx *ent.Tx, model *vo.Supplier) (*ent.Supplier, error) {
	builder := tx.Supplier.Create().SetAddress(model.Address).SetName(model.Name)
	if len(model.Constacts) != 0 {
		cons, err := CreateConstacts(ctx, tx, model.Constacts)
		if err != nil {
			return nil, err
		}
		builder.AddContacts(cons...)
	}
	return builder.Save(ctx)
}

func DeleteSupplier(ctx context.Context, id int64) error {
	return data.DB.Supplier.Update().SetDeletedTime(time.Now()).Exec(ctx)
}

func CreateConstacts(ctx context.Context, tx *ent.Tx, cons []*vo.Constact) ([]*ent.SupplierContact, error) {
	if len(cons) == 0 {
		return nil, nil
	}
	bulk := make([]*ent.SupplierContactCreate, len(cons))
	for i, c := range cons {
		bulk[i] = tx.SupplierContact.Create().SetFirstname(c.Firstname).SetLastname(c.Lastname).SetTitle(c.Title)
		if len(c.Relations) != 0 {
			rels, err := CreateRelations(ctx, tx, c.Relations)
			if err != nil {
				return nil, err
			}
			bulk[i] = bulk[i].AddRelations(rels...)
		}
	}
	return tx.SupplierContact.CreateBulk(bulk...).Save(ctx)
}

func SaveConstact(ctx context.Context, model *vo.Constact) (*ent.SupplierContact, error) {
	tx, err := data.DB.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			if rerr := tx.Rollback(); rerr != nil {
				err = fmt.Errorf("%w: %v", err, rerr)
			}
			return
		}
		err = tx.Commit()
	}()

	if model.ID == 0 {
		cons, err := CreateConstacts(ctx, tx, []*vo.Constact{model})
		if err != nil {
			return nil, err
		}
		return cons[0], err
	}
	return UpdateConstact(ctx, tx, model)
}

func UpdateConstact(ctx context.Context, tx *ent.Tx, model *vo.Constact) (*ent.SupplierContact, error) {
	builder := tx.SupplierContact.UpdateOneID(model.ID)
	if model.Firstname != "" {
		builder.SetFirstname(model.Firstname)
	}
	if model.Lastname != "" {
		builder.SetLastname(model.Lastname)
	}
	if model.Title != "" {
		builder.SetTitle(model.Title)
	}
	if model.SupplierID > 0 {
		builder.SetSupplierID(model.SupplierID)
	}
	if len(model.Relations) > 0 {
		rels, err := CreateRelations(ctx, tx, model.Relations)
		if err != nil {
			return nil, err
		}
		builder.AddRelations(rels...)
	}
	if len(model.RemoveRelations) > 0 {
		relsIDs := make([]int64, len(model.RemoveRelations))
		for _, v := range model.RemoveRelations {
			relsIDs = append(relsIDs, v.ID)
		}
		builder.RemoveRelationIDs(relsIDs...)
	}
	return builder.Save(ctx)
}

func DeleteConstact(ctx context.Context, id int64) error {
	return data.DB.SupplierContact.UpdateOneID(id).SetDeletedTime(time.Now()).Exec(ctx)
}

func SaveRelation(ctx context.Context, model *vo.Relation) (*ent.SupplierContactRelation, error) {
	if model.ID == 0 {
		return data.DB.SupplierContactRelation.Create().
			SetContactID(model.ContactID).
			SetType(model.Type).SetValue(model.Value).Save(ctx)
	}
	return data.DB.SupplierContactRelation.UpdateOneID(model.ID).SetContactID(model.ContactID).SetType(model.Type).SetValue(model.Value).Save(ctx)
}

func DeleteRelation(ctx context.Context, id int64) error {
	return data.DB.SupplierContactRelation.UpdateOneID(id).SetDeletedTime(time.Now()).Exec(ctx)
}

func CreateRelations(ctx context.Context, tx *ent.Tx, rels []*vo.Relation) ([]*ent.SupplierContactRelation, error) {
	if len(rels) == 0 {
		return nil, nil
	}
	bulk := make([]*ent.SupplierContactRelationCreate, len(rels))
	for i, c := range rels {
		bulk[i] = tx.SupplierContactRelation.Create().SetType(c.Type).SetValue(c.Value)
	}
	return tx.SupplierContactRelation.CreateBulk(bulk...).Save(ctx)
}
