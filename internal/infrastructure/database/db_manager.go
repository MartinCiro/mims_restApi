package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// DBManager simula las operaciones de Prisma
type DBManager struct {
	db *gorm.DB
}

func NewDBManager(db *gorm.DB) *DBManager {
	return &DBManager{
		db: db,
	}
}

// FindUnique simula prisma.findUnique()
func (dm *DBManager) FindUnique(ctx context.Context, model interface{}, where map[string]interface{}) error {
	query := dm.db.WithContext(ctx).Model(model)

	for field, value := range where {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	return query.First(model).Error
}

// FindMany simula prisma.findMany()
func (dm *DBManager) FindMany(ctx context.Context, models interface{}, where map[string]interface{}, options ...QueryOption) error {
	query := dm.db.WithContext(ctx).Model(models)

	for field, value := range where {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	// Aplicar opciones
	for _, option := range options {
		query = option.Apply(query)
	}

	return query.Find(models).Error
}

// Create simula prisma.create()
func (dm *DBManager) Create(ctx context.Context, model interface{}, data map[string]interface{}) error {
	return dm.db.WithContext(ctx).Model(model).Create(data).Error
}

// Update simula prisma.update()
func (dm *DBManager) Update(ctx context.Context, model interface{}, where map[string]interface{}, data map[string]interface{}) error {
	query := dm.db.WithContext(ctx).Model(model)

	for field, value := range where {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	return query.Updates(data).Error
}

// Delete simula prisma.delete()
func (dm *DBManager) Delete(ctx context.Context, model interface{}, where map[string]interface{}) error {
	query := dm.db.WithContext(ctx).Model(model)

	for field, value := range where {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	return query.Delete(model).Error
}

// Include simula el include de Prisma
func (dm *DBManager) Include(ctx context.Context, model interface{}, where map[string]interface{}, includes map[string]interface{}) error {
	query := dm.db.WithContext(ctx).Model(model)

	for field, value := range where {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	// Aplicar includes (preloads)
	for relation, conditions := range includes {
		if conditions == nil {
			query = query.Preload(relation)
		} else {
			query = query.Preload(relation, conditions)
		}
	}

	return query.First(model).Error
}

// QueryOptions para operaciones avanzadas
type QueryOption interface {
	Apply(*gorm.DB) *gorm.DB
}

type WithSelect struct {
	Fields []string
}

func (ws WithSelect) Apply(query *gorm.DB) *gorm.DB {
	return query.Select(ws.Fields)
}

type WithOrder struct {
	Field string
	Desc  bool
}

func (wo WithOrder) Apply(query *gorm.DB) *gorm.DB {
	order := wo.Field
	if wo.Desc {
		order += " DESC"
	}
	return query.Order(order)
}

type WithLimit struct {
	Limit int
}

func (wl WithLimit) Apply(query *gorm.DB) *gorm.DB {
	return query.Limit(wl.Limit)
}
