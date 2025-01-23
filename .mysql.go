// Code generated by auto_generate_repo.
// source:
//
// DO NOT EDIT

package 

import (
    "time"
	"context"

	"gorm.io/gorm"
	errors "github.com/rotisserie/eris"
)

type QueryOption func(db *gorm.DB) *gorm.DB

type  struct {
    CreatedBy *string `json:"created_by,omitempty" gorm:"size:100;comment:创建者"`
    CreatedAt *time.Time `json:"created_at,omitempty" gorm:"index;comment:创建时间"`
    ModifiedBy *string `json:"modified_by,omitempty" gorm:"size:100;comment:修改者"`
    ModifiedAt *time.Time `json:"modified_at,omitempty" gorm:"index;comment:修改时间"`
    DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index;comment:删除时间"`
}

func ( ) TableName() string {
    return ""
}

func ( *) BuildFuzzyQuery() QueryOption {
	g := 
	return func(db *gorm.DB) *gorm.DB {

		return db
	}
}

func ( *) BuildQuery() QueryOption {
	g := 
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where(g)
		return db
	}
}

type defaultRepository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &defaultRepository{
		db: db,
	}
}

func (d *defaultRepository) Create(ctx context.Context, data *) (err error) {
	if err := d.db.WithContext(ctx).Model(&{}).Create(data).Error; err != nil {
        return errors.Wrap(err, "failed to create  entity")
	}
	return nil
}


func (d *defaultRepository) Lists(ctx context.Context,  options map[string]interface{}, query ...QueryOption) (data []*,count int64, err error) {
    exec := d.db.WithContext(ctx).Model(&{})
	for _, v := range query {
    	exec = v(exec)
    }
    exec = exec.Count(&count)
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
		return data,count, errors.Wrap(err, "failed to list the  entities")
	}
	return data,count,err
}


func (d *defaultRepository) Get(ctx context.Context,  query *) (data *, err error) {
    data = &{}
    exec := d.db.WithContext(ctx).Model(&{})
	exec = exec.Where(query)
	err = exec.Last(data).Error
    if err != nil {
		return nil, errors.Wrap(err, "failed to get the ")
	}
	return data, err
}

func (d *defaultRepository) Update(ctx context.Context,update, query *) (data *,err error) {
	exec := d.db.WithContext(ctx).Model(&{})
    exec = exec.Where(query)
	out := &{}
	result := exec.Updates(update)
	if result.Error != nil {
		return nil, errors.Wrap(err, "failed to update  entity")
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	err = result.Scan(&out).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to update  entity")
	}

	return out, nil
}

func (d *defaultRepository) Delete(ctx context.Context,  query *) (err error){
    exec := d.db.WithContext(ctx).Model(&{})
	exec = exec.Where(query).Delete(&{})
	err = exec.Error
	if err != nil {
		return errors.Wrap(err, "failed to delete the ")
	}
	return nil
}

func (d *defaultRepository) Deletes(ctx context.Context,  ids []any) (err error){
    exec := d.db.WithContext(ctx).Model(&{})
	exec = exec.Where("id in ?", ids)
	err = exec.Delete(&{}).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete the s")
	}
	return nil
}

type Repository interface {
	Create(ctx context.Context, data *) (err error)
	Lists(ctx context.Context, options map[string]interface{}, query ...QueryOption) (data []*,count int64,err error)
	Get(ctx context.Context, query *) (data *, err error)
	Update(ctx context.Context, update , query *) (data *,err error)
	Delete(ctx context.Context,  query *) (err error)
	Deletes(ctx context.Context,  ids []any) (err error)
}
