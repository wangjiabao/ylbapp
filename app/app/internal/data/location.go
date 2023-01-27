package data

import (
	"context"
	"dhb/app/app/internal/biz"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"time"
)

type Location struct {
	ID           int64     `gorm:"primarykey;type:int"`
	UserId       int64     `gorm:"type:int;not null"`
	Row          int64     `gorm:"type:int;not null"`
	Col          int64     `gorm:"type:int;not null"`
	Status       string    `gorm:"type:varchar(45);not null"`
	CurrentLevel int64     `gorm:"type:int;not null"`
	Current      int64     `gorm:"type:bigint;not null"`
	CurrentMax   int64     `gorm:"type:bigint;not null"`
	StopDate     time.Time `gorm:"type:datetime;not null"`
	CreatedAt    time.Time `gorm:"type:datetime;not null"`
	UpdatedAt    time.Time `gorm:"type:datetime;not null"`
}

type GlobalLock struct {
	ID     int64 `gorm:"primarykey;type:int"`
	Status int64 `gorm:"type:int;not null"`
}

type LocationRepo struct {
	data *Data
	log  *log.Helper
}

func NewLocationRepo(data *Data, logger log.Logger) biz.LocationRepo {
	return &LocationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// CreateLocation .
func (lr *LocationRepo) CreateLocation(ctx context.Context, rel *biz.Location) (*biz.Location, error) {
	var location Location
	location.Col = rel.Col
	location.Row = rel.Row
	location.Status = rel.Status
	location.Current = rel.Current
	location.CurrentMax = rel.CurrentMax
	location.CurrentLevel = rel.CurrentLevel
	location.UserId = rel.UserId
	res := lr.data.DB(ctx).Table("location").Create(&location)
	if res.Error != nil {
		return nil, errors.New(500, "CREATE_LOCATION_ERROR", "占位信息创建失败")
	}

	return &biz.Location{
		ID:           location.ID,
		UserId:       location.UserId,
		Status:       location.Status,
		CurrentLevel: location.CurrentLevel,
		Current:      location.Current,
		CurrentMax:   location.CurrentMax,
		Row:          location.Row,
		Col:          location.Col,
	}, nil
}

// GetLocationLast .
func (lr *LocationRepo) GetLocationLast(ctx context.Context) (*biz.Location, error) {
	var location Location
	if err := lr.data.db.Table("location").Order("id desc").First(&location).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("LOCATION_NOT_FOUND", "location not found")
		}

		return nil, errors.New(500, "LOCATION ERROR", err.Error())
	}

	return &biz.Location{
		ID:           location.ID,
		UserId:       location.UserId,
		Status:       location.Status,
		CurrentLevel: location.CurrentLevel,
		Current:      location.Current,
		CurrentMax:   location.CurrentMax,
		Row:          location.Row,
		Col:          location.Col,
	}, nil
}

// GetMyLocationLast .
func (lr *LocationRepo) GetMyLocationLast(ctx context.Context, userId int64) (*biz.Location, error) {
	var location Location
	if err := lr.data.db.Table("location").Where("user_id", userId).Order("id desc").First(&location).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("LOCATION_NOT_FOUND", "location not found")
		}

		return nil, errors.New(500, "LOCATION ERROR", err.Error())
	}

	return &biz.Location{
		ID:           location.ID,
		UserId:       location.UserId,
		Status:       location.Status,
		CurrentLevel: location.CurrentLevel,
		Current:      location.Current,
		CurrentMax:   location.CurrentMax,
		Row:          location.Row,
		Col:          location.Col,
		StopDate:     location.StopDate,
	}, nil
}

// GetMyStopLocationLast .
func (lr *LocationRepo) GetMyStopLocationLast(ctx context.Context, userId int64) (*biz.Location, error) {
	var location Location
	if err := lr.data.db.Table("location").
		Where("status=?", "stop").
		Where("user_id", userId).Order("id desc").First(&location).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("LOCATION_NOT_FOUND", "location not found")
		}

		return nil, errors.New(500, "LOCATION ERROR", err.Error())
	}

	return &biz.Location{
		ID:           location.ID,
		UserId:       location.UserId,
		Status:       location.Status,
		CurrentLevel: location.CurrentLevel,
		Current:      location.Current,
		CurrentMax:   location.CurrentMax,
		Row:          location.Row,
		Col:          location.Col,
		StopDate:     location.StopDate,
	}, nil
}

// GetMyLocationRunningLast .
func (lr *LocationRepo) GetMyLocationRunningLast(ctx context.Context, userId int64) (*biz.Location, error) {
	var location Location
	if err := lr.data.db.Table("location").Where("user_id", userId).
		Where("status=?", "running").
		Order("id desc").First(&location).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("LOCATION_NOT_FOUND", "location not found")
		}

		return nil, errors.New(500, "LOCATION ERROR", err.Error())
	}

	return &biz.Location{
		ID:           location.ID,
		UserId:       location.UserId,
		Status:       location.Status,
		CurrentLevel: location.CurrentLevel,
		Current:      location.Current,
		CurrentMax:   location.CurrentMax,
		Row:          location.Row,
		Col:          location.Col,
	}, nil
}

// GetLocationsByUserId .
func (lr *LocationRepo) GetLocationsByUserId(ctx context.Context, userId int64) ([]*biz.Location, error) {
	var locations []*Location
	res := make([]*biz.Location, 0)
	if err := lr.data.db.Table("location").
		Where("user_id=?", userId).
		Order("id desc").Find(&locations).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res, errors.NotFound("LOCATION_NOT_FOUND", "location not found")
		}

		return nil, errors.New(500, "LOCATION ERROR", err.Error())
	}

	for _, location := range locations {
		res = append(res, &biz.Location{
			ID:           location.ID,
			UserId:       location.UserId,
			Status:       location.Status,
			CurrentLevel: location.CurrentLevel,
			Current:      location.Current,
			CurrentMax:   location.CurrentMax,
			Row:          location.Row,
			Col:          location.Col,
		})
	}

	return res, nil
}

// GetLocationsStopNotUpdate .
func (lr *LocationRepo) GetLocationsStopNotUpdate(ctx context.Context) ([]*biz.Location, error) {
	var locations []*Location
	res := make([]*biz.Location, 0)
	if err := lr.data.db.Table("location").
		Where("status=?", "stop").
		Where("stop_is_update=?", 0).
		Find(&locations).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res, errors.NotFound("LOCATION_NOT_FOUND", "location not found")
		}

		return nil, errors.New(500, "LOCATION ERROR", err.Error())
	}

	for _, location := range locations {
		res = append(res, &biz.Location{
			ID:           location.ID,
			UserId:       location.UserId,
			Status:       location.Status,
			CurrentLevel: location.CurrentLevel,
			Current:      location.Current,
			CurrentMax:   location.CurrentMax,
			Row:          location.Row,
			Col:          location.Col,
		})
	}

	return res, nil
}

// LockGlobalLocation .
func (lr *LocationRepo) LockGlobalLocation(ctx context.Context) (bool, error) {
	res := lr.data.DB(ctx).Where("id=? and status<=?", 1, 2).
		Table("global_lock").
		Updates(map[string]interface{}{"status": 1})

	if 0 <= res.RowsAffected {
		return true, nil
	}

	return false, res.Error

}

// UnLockGlobalLocation .
func (lr *LocationRepo) UnLockGlobalLocation(ctx context.Context) (bool, error) {
	res := lr.data.DB(ctx).Where("id=? and status=?", 1, 1).
		Table("global_lock").
		Updates(map[string]interface{}{"status": 2})

	if 0 <= res.RowsAffected {
		return true, nil
	}

	return false, res.Error
}

// LockGlobalWithdraw .
func (lr *LocationRepo) LockGlobalWithdraw(ctx context.Context) (bool, error) {
	res := lr.data.DB(ctx).Where("id=? and status>=?", 1, 2).
		Table("global_lock").
		Updates(map[string]interface{}{"status": 3})

	if 0 <= res.RowsAffected {
		return true, nil
	}

	return false, res.Error
}

// GetLockGlobalLocation .
func (lr *LocationRepo) GetLockGlobalLocation(ctx context.Context) (*biz.GlobalLock, error) {
	var globalLock GlobalLock
	if res := lr.data.DB(ctx).Where("id=?", 1).
		Table("global_lock").
		First(&globalLock); res.Error != nil {
		return nil, res.Error
	}

	return &biz.GlobalLock{
		ID:     globalLock.ID,
		Status: globalLock.Status,
	}, nil
}

// UnLockGlobalWithdraw .
func (lr *LocationRepo) UnLockGlobalWithdraw(ctx context.Context) (bool, error) {
	res := lr.data.DB(ctx).Where("id=? and status=?", 1, 3).
		Table("global_lock").
		Updates(map[string]interface{}{"status": 2})

	if 0 <= res.RowsAffected {
		return true, nil
	}

	return false, res.Error
}

// UpdateLocation .
func (lr *LocationRepo) UpdateLocation(ctx context.Context, id int64, status string, current int64, stopDate time.Time) error {

	if "stop" == status {
		res := lr.data.db.Table("location").
			Where("id=?", id).
			Updates(map[string]interface{}{"current": gorm.Expr("current + ?", current), "status": "stop", "stop_date": stopDate})
		if 0 == res.RowsAffected || res.Error != nil {
			return res.Error
		}
	} else {
		res := lr.data.db.Table("location").
			Where("id=?", id).
			Where("status=?", "running").
			Updates(map[string]interface{}{"current": gorm.Expr("current + ?", current), "status": status})
		if 0 == res.RowsAffected || res.Error != nil {
			return res.Error
		}
	}

	return nil
}

// UpdateLocationRowAndCol 事务中使用 .
func (lr *LocationRepo) UpdateLocationRowAndCol(ctx context.Context, id int64) error {

	if res := lr.data.db.Table("location").
		Where("id>?", id).
		Where("col > 1").
		Where("update_status=?", 0).
		Updates(map[string]interface{}{"col": gorm.Expr("col - ?", 1), "update_status": 1}); res.Error != nil {
		return res.Error
	}

	if res := lr.data.db.Table("location").
		Where("id>?", id).
		Where("col = 1").
		Where("update_status=?", 0).
		Updates(map[string]interface{}{"row": gorm.Expr("row - ?", 1), "col": 3, "update_status": 1}); res.Error != nil {
		return res.Error
	}

	if res := lr.data.db.Table("location").
		Where("id>?", id).
		Updates(map[string]interface{}{"update_status": 0}); res.Error != nil {
		return res.Error
	}

	if res := lr.data.db.Table("location").
		Where("id=?", id).
		Updates(map[string]interface{}{"stop_is_update": 1}); res.Error != nil {
		return res.Error
	}
	return nil
}

// GetLocationDaily .
func (lr *LocationRepo) GetLocationDaily(ctx context.Context) ([]*biz.Location, error) {
	var locations []*Location
	res := make([]*biz.Location, 0)
	instance := lr.data.db.Table("location")

	now := time.Now().UTC()
	var startDate time.Time
	var endDate time.Time
	if 14 <= now.Hour() {
		startDate = now
		endDate = now.AddDate(0, 0, 1)
	} else {
		startDate = now.AddDate(0, 0, -1)
		endDate = now
	}
	todayStart := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 14, 0, 0, 0, time.UTC)
	todayEnd := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 14, 0, 0, 0, time.UTC)

	instance = instance.Where("created_at>=?", todayStart)
	instance = instance.Where("created_at<?", todayEnd)
	if err := instance.Order("id desc").Find(&locations).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res, errors.NotFound("LOCATION_NOT_FOUND", "location not found")
		}

		return res, errors.New(500, "LOCATION ERROR", err.Error())
	}

	for _, v := range locations {
		res = append(res, &biz.Location{
			ID:           v.ID,
			UserId:       v.UserId,
			Status:       v.Status,
			CurrentLevel: v.CurrentLevel,
			Current:      v.Current,
			CurrentMax:   v.CurrentMax,
		})
	}

	return res, nil
}

// GetRewardLocationByRowOrCol .
func (lr *LocationRepo) GetRewardLocationByRowOrCol(ctx context.Context, row int64, col int64) ([]*biz.Location, error) {
	var (
		rowMin    int64 = 1
		rowMax    int64
		locations []*Location
	)
	if row > 25 {
		rowMin = row - 25
	}
	rowMax = row + 25

	if err := lr.data.db.Table("location").
		Where("status=?", "running").
		Where("row=? or (col=? and row>=? and row<=?)", row, col, rowMin, rowMax).
		Find(&locations).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("LOCATION_NOT_FOUND", "location not found")
		}

		return nil, errors.New(500, "LOCATION ERROR", err.Error())
	}

	res := make([]*biz.Location, 0)
	for _, location := range locations {
		res = append(res, &biz.Location{
			ID:           location.ID,
			UserId:       location.UserId,
			Status:       location.Status,
			CurrentLevel: location.CurrentLevel,
			Current:      location.Current,
			CurrentMax:   location.CurrentMax,
			Row:          location.Row,
			Col:          location.Col,
			StopDate:     location.StopDate,
		})
	}

	return res, nil
}

// GetRewardLocationByIds .
func (lr *LocationRepo) GetRewardLocationByIds(ctx context.Context, ids ...int64) (map[int64]*biz.Location, error) {
	var locations []*Location
	if err := lr.data.db.Table("location").
		Where("status=?", "running").
		Where("id IN (?)", ids).
		Find(&locations).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("LOCATION_NOT_FOUND", "location not found")
		}

		return nil, errors.New(500, "LOCATION ERROR", err.Error())
	}

	res := make(map[int64]*biz.Location, 0)
	for _, location := range locations {
		res[location.ID] = &biz.Location{
			ID:           location.ID,
			UserId:       location.UserId,
			Status:       location.Status,
			CurrentLevel: location.CurrentLevel,
			Current:      location.Current,
			CurrentMax:   location.CurrentMax,
			Row:          location.Row,
			Col:          location.Col,
		}
	}

	return res, nil
}

// GetLocations .
func (lr *LocationRepo) GetLocations(ctx context.Context, b *biz.Pagination, userId int64) ([]*biz.Location, error, int64) {
	var (
		locations []*Location
		count     int64
	)
	instance := lr.data.db.Table("location").Where("status=?", "running")

	if 0 < userId {
		instance = instance.Where("user_id=?", userId)
	}

	instance = instance.Count(&count)
	if err := instance.Scopes(Paginate(b.PageNum, b.PageSize)).Order("id desc").Find(&locations).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("LOCATION_NOT_FOUND", "location not found"), 0
		}

		return nil, errors.New(500, "LOCATION ERROR", err.Error()), 0
	}

	res := make([]*biz.Location, 0)
	for _, location := range locations {
		res = append(res, &biz.Location{
			ID:           location.ID,
			UserId:       location.UserId,
			Status:       location.Status,
			CurrentLevel: location.CurrentLevel,
			Current:      location.Current,
			CurrentMax:   location.CurrentMax,
			Row:          location.Row,
			Col:          location.Col,
			CreatedAt:    location.CreatedAt,
		})
	}

	return res, nil, count
}
