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

// FindUniqueOptions configuración para FindUnique
type FindUniqueOptions struct {
	UseOrderByID bool // Si es true, agrega ORDER BY id, si es false, no agrega ORDER BY
	OrderBy      string
	SelectFields []string
}

// FindUnique simula prisma.findUnique() con opciones mejoradas
func (dm *DBManager) FindUnique(ctx context.Context, table string, result interface{}, where map[string]interface{}, options ...*FindUniqueOptions) error {
	query := dm.db.WithContext(ctx).Table(table)

	// Aplicar condiciones WHERE
	for field, value := range where {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	// Aplicar opciones si se proporcionan
	if len(options) > 0 && options[0] != nil {
		opts := options[0]

		// SELECT fields específicos
		if len(opts.SelectFields) > 0 {
			query = query.Select(opts.SelectFields)
		}

		// ORDER BY
		if opts.OrderBy != "" {
			query = query.Order(opts.OrderBy)
		} else if opts.UseOrderByID {
			// ORDER BY id solo si se solicita explícitamente
			query = query.Order("id")
		}
		// Si no hay opciones de ORDER BY, no se agrega nada (evita el ORDER BY automático)
	}

	return query.First(result).Error
}

// FindUniqueByField busca por un campo específico sin ORDER BY automático
func (dm *DBManager) FindUniqueByField(ctx context.Context, table string, result interface{}, field string, value interface{}) error {
	return dm.FindUnique(ctx, table, result, map[string]interface{}{field: value}, &FindUniqueOptions{
		UseOrderByID: false,
	})
}

// FindUniqueByID busca por ID con ORDER BY (comportamiento tradicional)
func (dm *DBManager) FindUniqueByID(ctx context.Context, table string, result interface{}, id interface{}) error {
	return dm.FindUnique(ctx, table, result, map[string]interface{}{"id": id}, &FindUniqueOptions{
		UseOrderByID: true, // Con ORDER BY id
	})
}

// FindMany simula prisma.findMany()
func (dm *DBManager) FindMany(ctx context.Context, table string, results interface{}, where map[string]interface{}) error {
	query := dm.db.WithContext(ctx).Table(table)

	for field, value := range where {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	return query.Find(results).Error
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

func (dm *DBManager) FindWithJoin(ctx context.Context, results interface{}, query string, args ...interface{}) error {
	return dm.db.WithContext(ctx).Raw(query, args...).Scan(results).Error
}

// GetDB retorna la instancia de GORM para operaciones directas
func (dm *DBManager) GetDB() *gorm.DB {
	return dm.db
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
