// Code generated by auto_generate_repo.
// source:
//
// DO NOT EDIT

package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	errors "github.com/rotisserie/eris"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type GameQueryOption func(db *gorm.DB) *gorm.DB

type Game struct {
	Id               *int64          `json:"id,omitempty" gorm:"size:8;primaryKey"`
	Name             *string         `json:"name,omitempty"`
	Owners           *string         `json:"owners,omitempty"`
	NotificationList *string         `json:"notification_list,omitempty"`
	Description      *string         `json:"description,omitempty"`
	Attr             *datatypes.JSON `json:"attr,omitempty" gorm:"type:jsonb;default:'{}'"`
	CreatedBy        *string         `json:"created_by,omitempty"`
	CreatedTs        *time.Time      `json:"created_ts,omitempty"`
	ModifiedBy       *string         `json:"modified_by,omitempty"`
	ModifiedTs       *time.Time      `json:"modified_ts,omitempty"`
}

func (i Game) TableName() string {
	return "game"
}

func (g *Game) BuildQuery() GameQueryOption {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where(g)
		return db
	}
}

func (g *Game) BuildFuzzyQuery() GameQueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if g.Name != nil {
			db = db.Where("name ilike ?", "%"+*g.Name+"%")
		}
		if g.Description != nil {
			db = db.Where("description ilike ?", "%"+*g.Description+"%")
		}

		return db
	}
}

type defaultGameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) GameRepository {
	return &defaultGameRepository{
		db: db,
	}
}

func (d *defaultGameRepository) CreateGame(ctx context.Context, data *Game) (err error) {
	var pgError *pgconn.PgError
	if err := d.db.WithContext(ctx).Table(Game{}.TableName()).Create(data).Error; err != nil {
		if errors.As(err, &pgError) && pgError.Code == PgErrUniqueViolation {
			return errors.Wrap(gorm.ErrDuplicatedKey, pgError.Detail)
		}
		return errors.Wrap(err, "failed to create Game entity")
	}
	return nil
}

func (d *defaultGameRepository) ListGames(ctx context.Context, options map[string]interface{}, query ...GameQueryOption) (data []*Game, err error) {
	exec := d.db.WithContext(ctx).Table(Game{}.TableName())
	for _, v := range query {
		exec = v(exec)
	}
	for key, val := range options {
		switch key {
		case "limit":
			exec = exec.Limit(val.(int))
		case "offset":
			exec = exec.Offset(val.(int))
		case "sort":
			exec = exec.Order(val).Order("id ASC")
		case "select":
			exec = exec.Select(val)
		}
	}

	err = exec.Find(&data).Error
	if err != nil {
		return data, errors.Wrap(err, "failed to list the Game entities")
	}
	return data, err
}

func (d *defaultGameRepository) GetGame(ctx context.Context, query *Game) (data *Game, err error) {
	data = &Game{}
	exec := d.db.WithContext(ctx).Table(Game{}.TableName())
	exec = exec.Where(query)
	err = exec.Last(data).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the Game")
	}
	return data, err
}

func (d *defaultGameRepository) UpdateGame(ctx context.Context, update, query *Game) (data *Game, err error) {
	var pgError *pgconn.PgError
	exec := d.db.WithContext(ctx).Table(Game{}.TableName())
	exec = exec.Where(query)
	out := &Game{}
	result := exec.Updates(update)
	if result.Error != nil {
		if errors.As(err, &pgError) && pgError.Code == PgErrUniqueViolation {
			return nil, errors.Wrap(gorm.ErrDuplicatedKey, pgError.Detail)
		}
		return nil, errors.Wrap(err, "failed to update Game entity")
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	err = result.Row().Scan(
		&out.Id,
		&out.Name,
		&out.Owners,
		&out.NotificationList,
		&out.Description,
		&out.Attr,
		&out.CreatedBy,
		&out.CreatedTs,
		&out.ModifiedBy,
		&out.ModifiedTs,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update Game entity")
	}

	return out, nil
}

type GameRepository interface {
	CreateGame(ctx context.Context, data *Game) (err error)
	ListGames(ctx context.Context, options map[string]interface{}, query ...GameQueryOption) (data []*Game, err error)
	GetGame(ctx context.Context, query *Game) (data *Game, err error)
	UpdateGame(ctx context.Context, update, query *Game) (data *Game, err error)
}
