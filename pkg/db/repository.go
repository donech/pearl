package db

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type Repository struct {
	DB *gorm.DB
}

func (d *Repository) Create(ctx context.Context, entity interface{}) error {
	return d.CreateTx(ctx, d.DB, entity)
}

func (d *Repository) Save(ctx context.Context, entity interface{}) error {
	return d.SaveTx(ctx, d.DB, entity)
}

func (d *Repository) Delete(ctx context.Context, entity interface{}, where ...interface{}) error {
	return d.DeleteTx(ctx, d.DB, entity, where)
}

func (d *Repository) CreateTx(ctx context.Context, tx *gorm.DB, entity interface{}) error {
	return Trace(ctx, tx).Create(entity).Error
}

func (d *Repository) SaveTx(ctx context.Context, tx *gorm.DB, entity interface{}) error {
	err := Trace(ctx, tx).Save(entity).Error
	return errors.WithStack(err)
}

func (d *Repository) DeleteTx(ctx context.Context, tx *gorm.DB, entity interface{}, where ...interface{}) error {
	err := Trace(ctx, tx).Delete(entity, where).Error
	return errors.WithStack(err)
}

func (d *Repository) Find(ctx context.Context, tx *gorm.DB, entity interface{}, out interface{}, limit interface{}, offset interface{}, where ...interface{}) (interface{}, error) {
	return d.FindTx(ctx, d.DB, entity, out, limit, offset, where)
}

func (d *Repository) FindTx(ctx context.Context, tx *gorm.DB, entity interface{}, out interface{}, limit interface{}, offset interface{}, where ...interface{}) (interface{}, error) {
	err := Trace(ctx, tx).Model(entity).Where(where[0], where[1:]).Limit(limit).Offset(offset).Find(out).Error
	return out, errors.WithStack(err)
}
