package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"strings"
	"time"
)

type EthUserRecord struct {
	ID       int64
	UserId   int64
	Hash     string
	Status   string
	Type     string
	Amount   string
	CoinType string
}

type Location struct {
	ID           int64
	UserId       int64
	Status       string
	CurrentLevel int64
	Current      int64
	CurrentMax   int64
	Row          int64
	Col          int64
	StopDate     time.Time
	CreatedAt    time.Time
}

type GlobalLock struct {
	ID     int64
	Status int64
}

type RecordUseCase struct {
	ethUserRecordRepo             EthUserRecordRepo
	userRecommendRepo             UserRecommendRepo
	configRepo                    ConfigRepo
	locationRepo                  LocationRepo
	userBalanceRepo               UserBalanceRepo
	userInfoRepo                  UserInfoRepo
	userCurrentMonthRecommendRepo UserCurrentMonthRecommendRepo
	tx                            Transaction
	log                           *log.Helper
}

type EthUserRecordRepo interface {
	GetEthUserRecordListByHash(ctx context.Context, hash ...string) (map[string]*EthUserRecord, error)
	CreateEthUserRecordListByHash(ctx context.Context, r *EthUserRecord) (*EthUserRecord, error)
}

type LocationRepo interface {
	CreateLocation(ctx context.Context, rel *Location) (*Location, error)
	GetLocationLast(ctx context.Context) (*Location, error)
	GetLocationDaily(ctx context.Context) ([]*Location, error)
	GetMyLocationLast(ctx context.Context, userId int64) (*Location, error)
	GetMyStopLocationLast(ctx context.Context, userId int64) (*Location, error)
	GetMyLocationRunningLast(ctx context.Context, userId int64) (*Location, error)
	GetLocationsByUserId(ctx context.Context, userId int64) ([]*Location, error)
	GetRewardLocationByRowOrCol(ctx context.Context, row int64, col int64) ([]*Location, error)
	GetRewardLocationByIds(ctx context.Context, ids ...int64) (map[int64]*Location, error)
	UpdateLocation(ctx context.Context, id int64, status string, current int64, stopDate time.Time) error
	GetLocations(ctx context.Context, b *Pagination, userId int64) ([]*Location, error, int64)
	UpdateLocationRowAndCol(ctx context.Context, id int64) error
	GetLocationsStopNotUpdate(ctx context.Context) ([]*Location, error)
	LockGlobalLocation(ctx context.Context) (bool, error)
	UnLockGlobalLocation(ctx context.Context) (bool, error)
	LockGlobalWithdraw(ctx context.Context) (bool, error)
	UnLockGlobalWithdraw(ctx context.Context) (bool, error)
	GetLockGlobalLocation(ctx context.Context) (*GlobalLock, error)
}

func NewRecordUseCase(
	ethUserRecordRepo EthUserRecordRepo,
	locationRepo LocationRepo,
	userBalanceRepo UserBalanceRepo,
	userRecommendRepo UserRecommendRepo,
	userInfoRepo UserInfoRepo,
	configRepo ConfigRepo,
	userCurrentMonthRecommendRepo UserCurrentMonthRecommendRepo,
	tx Transaction,
	logger log.Logger) *RecordUseCase {
	return &RecordUseCase{
		ethUserRecordRepo:             ethUserRecordRepo,
		locationRepo:                  locationRepo,
		configRepo:                    configRepo,
		userRecommendRepo:             userRecommendRepo,
		userBalanceRepo:               userBalanceRepo,
		userCurrentMonthRecommendRepo: userCurrentMonthRecommendRepo,
		userInfoRepo:                  userInfoRepo,
		tx:                            tx,
		log:                           log.NewHelper(logger),
	}
}

func (ruc *RecordUseCase) GetEthUserRecordByTxHash(ctx context.Context, txHash ...string) (map[string]*EthUserRecord, error) {
	return ruc.ethUserRecordRepo.GetEthUserRecordListByHash(ctx, txHash...)
}

func (ruc *RecordUseCase) EthUserRecordHandle(ctx context.Context, ethUserRecord ...*EthUserRecord) (bool, error) {

	var (
		configs           []*Config
		recommendNeed     int64
		recommendNeedVip1 int64
		recommendNeedVip2 int64
		recommendNeedVip3 int64
		recommendNeedVip4 int64
		recommendNeedVip5 int64
		timeAgain         int64
	)
	// 配置
	configs, _ = ruc.configRepo.GetConfigByKeys(ctx, "recommend_need", "recommend_need_vip1", "recommend_need_vip2",
		"recommend_need_vip3", "recommend_need_vip4", "recommend_need_vip5", "time_again")
	if nil != configs {
		for _, vConfig := range configs {
			if "recommend_need" == vConfig.KeyName {
				recommendNeed, _ = strconv.ParseInt(vConfig.Value, 10, 64)
			} else if "recommend_need_vip1" == vConfig.KeyName {
				recommendNeedVip1, _ = strconv.ParseInt(vConfig.Value, 10, 64)
			} else if "recommend_need_vip2" == vConfig.KeyName {
				recommendNeedVip2, _ = strconv.ParseInt(vConfig.Value, 10, 64)
			} else if "recommend_need_vip3" == vConfig.KeyName {
				recommendNeedVip3, _ = strconv.ParseInt(vConfig.Value, 10, 64)
			} else if "recommend_need_vip4" == vConfig.KeyName {
				recommendNeedVip4, _ = strconv.ParseInt(vConfig.Value, 10, 64)
			} else if "recommend_need_vip5" == vConfig.KeyName {
				recommendNeedVip5, _ = strconv.ParseInt(vConfig.Value, 10, 64)
			} else if "time_again" == vConfig.KeyName {
				timeAgain, _ = strconv.ParseInt(vConfig.Value, 10, 64)
			}
		}
	}
	fmt.Println(recommendNeed, recommendNeedVip1, recommendNeedVip2, recommendNeedVip3, recommendNeedVip4, recommendNeedVip5, timeAgain)

	for _, v := range ethUserRecord {
		var (
			lastLocation                    *Location
			myLocations                     []*Location
			currentValue                    int64
			amount                          int64
			locationCurrentLevel            int64
			locationCurrent                 int64
			locationCurrentMax              int64
			locationRow                     int64
			locationCol                     int64
			currentLocation                 *Location
			rewardLocations                 []*Location
			userRecommend                   *UserRecommend
			myUserRecommendUserId           int64
			myUserRecommendUserInfo         *UserInfo
			myUserRecommendUserLocationLast *Location
			stopLocations                   []*Location
			myLastStopLocation              *Location
			err                             error
		)

		//if "DHB" == v.CoinType {
		//	continue
		//}

		// 获取当前用户的占位信息，已经有运行中的跳过
		myLocations, err = ruc.locationRepo.GetLocationsByUserId(ctx, v.UserId)
		if nil == myLocations { // 查询异常跳过本次循环
			continue
		}
		if 0 < len(myLocations) { // 也代表复投
			tmpStatusRunning := false
			for _, vMyLocations := range myLocations {
				if "running" == vMyLocations.Status {
					tmpStatusRunning = true
					break
				}
			}

			if tmpStatusRunning { // 有运行中直接跳过本次循环
				continue
			}
		}

		// 先紧缩一次位置
		stopLocations, err = ruc.locationRepo.GetLocationsStopNotUpdate(ctx)
		if nil != stopLocations {
			// 调整位置紧缩
			for _, vStopLocations := range stopLocations {

				if err = ruc.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务
					err = ruc.locationRepo.UpdateLocationRowAndCol(ctx, vStopLocations.ID)
					if nil != err {
						return err
					}
					return nil
				}); nil != err {
					continue
				}
			}
		}

		// 获取最后一行数据
		lastLocation, err = ruc.locationRepo.GetLocationLast(ctx)
		if nil == lastLocation {
			locationRow = 1
			locationCol = 1
			fmt.Println(25, locationRow, locationRow)
		} else {
			if 3 > lastLocation.Col {
				locationCol = lastLocation.Col + 1
				locationRow = lastLocation.Row
				fmt.Println(33, locationCol, locationRow)
			} else {
				locationCol = 1
				locationRow = lastLocation.Row + 1
				fmt.Println(22, locationRow, locationRow)
			}
		}

		// todo
		if "100000000000000000000" == v.Amount {
			locationCurrentLevel = 1
			locationCurrentMax = 5000000000000
			currentValue = 1000000000000
		} else if "200000000000000000000" == v.Amount {
			locationCurrentLevel = 2
			locationCurrentMax = 10000000000000
			currentValue = 2000000000000
		} else if "500000000000000000000" == v.Amount {
			locationCurrentLevel = 3
			locationCurrentMax = 25000000000000
			currentValue = 5000000000000
		} else {
			continue
		}
		amount = currentValue

		// 占位分红人
		rewardLocations, err = ruc.locationRepo.GetRewardLocationByRowOrCol(ctx, locationRow, locationCol)

		// 推荐人
		userRecommend, err = ruc.userRecommendRepo.GetUserRecommendByUserId(ctx, v.UserId)
		if nil != err {
			continue
		}
		if "" != userRecommend.RecommendCode {
			tmpRecommendUserIds := strings.Split(userRecommend.RecommendCode, "D")
			if 2 <= len(tmpRecommendUserIds) {
				myUserRecommendUserId, _ = strconv.ParseInt(tmpRecommendUserIds[len(tmpRecommendUserIds)-1], 10, 64) // 最后一位是直推人
			}
		}
		if 0 < myUserRecommendUserId {
			myUserRecommendUserInfo, err = ruc.userInfoRepo.GetUserInfoByUserId(ctx, myUserRecommendUserId)
		}

		myLastStopLocation, err = ruc.locationRepo.GetMyStopLocationLast(ctx, v.UserId)
		now := time.Now().UTC().Add(8 * time.Hour)
		if nil != myLastStopLocation && now.Before(myLastStopLocation.StopDate.Add(time.Duration(timeAgain)*time.Minute)) {
			locationCurrent = myLastStopLocation.Current - myLastStopLocation.CurrentMax // 补上
		}

		if err = ruc.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务
			currentLocation, err = ruc.locationRepo.CreateLocation(ctx, &Location{ // 占位
				UserId:       v.UserId,
				Status:       "running",
				CurrentLevel: locationCurrentLevel,
				Current:      locationCurrent,
				CurrentMax:   locationCurrentMax,
				Row:          locationRow,
				Col:          locationCol,
			})
			if nil != err {
				return err
			}

			// 占位分红人分红
			if nil != rewardLocations {
				for _, vRewardLocations := range rewardLocations {
					if "running" != vRewardLocations.Status {
						continue
					}
					if locationRow == vRewardLocations.Row && locationCol == vRewardLocations.Col { // 跳过自己
						continue
					}

					var locationType string
					var tmpAmount int64
					if locationRow == vRewardLocations.Row { // 同行的人
						tmpAmount = currentValue / 100 * 5
						locationType = "row"
					} else if locationCol == vRewardLocations.Col { // 同列的人
						tmpAmount = currentValue / 100
						locationType = "col"
					} else {
						continue
					}

					tmpCurrentStatus := vRewardLocations.Status // 现在还在运行中
					tmpCurrent := vRewardLocations.Current

					tmpBalanceAmount := tmpAmount
					vRewardLocations.Status = "running"
					vRewardLocations.Current += tmpAmount
					if vRewardLocations.Current >= vRewardLocations.CurrentMax { // 占位分红人分满停止
						if "running" == tmpCurrentStatus {
							vRewardLocations.StopDate = time.Now().UTC().Add(8 * time.Hour)
						}
						vRewardLocations.Status = "stop"
					}

					if 0 < tmpBalanceAmount {
						err = ruc.locationRepo.UpdateLocation(ctx, vRewardLocations.ID, vRewardLocations.Status, tmpBalanceAmount, vRewardLocations.StopDate) // 分红占位数据修改
						if nil != err {
							return err
						}
						amount -= tmpBalanceAmount // 扣除

						if 0 < tmpBalanceAmount && "running" == tmpCurrentStatus && tmpCurrent < vRewardLocations.CurrentMax { // 这次还能分红
							tmpCurrentAmount := vRewardLocations.CurrentMax - tmpCurrent // 最大可分红额度
							rewardAmount := tmpBalanceAmount
							if tmpCurrentAmount < tmpBalanceAmount { // 大于最大可分红额度
								rewardAmount = tmpCurrentAmount
							}

							_, err = ruc.userBalanceRepo.LocationReward(ctx, vRewardLocations.UserId, rewardAmount, currentLocation.ID, vRewardLocations.ID, locationType) // 分红信息修改
							if nil != err {
								return err
							}
						}
					}
				}

			}

			// 推荐人
			if nil != myUserRecommendUserInfo {
				if 0 == len(myLocations) { // vip 等级调整，被推荐人首次入单
					myUserRecommendUserInfo.HistoryRecommend += 1
					if myUserRecommendUserInfo.HistoryRecommend >= 10 {
						myUserRecommendUserInfo.Vip = 5
					} else if myUserRecommendUserInfo.HistoryRecommend >= 8 {
						myUserRecommendUserInfo.Vip = 4
					} else if myUserRecommendUserInfo.HistoryRecommend >= 6 {
						myUserRecommendUserInfo.Vip = 3
					} else if myUserRecommendUserInfo.HistoryRecommend >= 4 {
						myUserRecommendUserInfo.Vip = 2
					} else if myUserRecommendUserInfo.HistoryRecommend >= 2 {
						myUserRecommendUserInfo.Vip = 1
					}

					_, err = ruc.userInfoRepo.UpdateUserInfo(ctx, myUserRecommendUserInfo) // 推荐人信息修改
					if nil != err {
						return err
					}

					_, err = ruc.userCurrentMonthRecommendRepo.CreateUserCurrentMonthRecommend(ctx, &UserCurrentMonthRecommend{ // 直推人本月推荐人数
						UserId:          myUserRecommendUserId,
						RecommendUserId: v.UserId,
						Date:            time.Now().UTC().Add(8 * time.Hour),
					})
					if nil != err {
						return err
					}
				}

				// 有占位信息
				myUserRecommendUserLocationLast, err = ruc.locationRepo.GetMyLocationLast(ctx, myUserRecommendUserInfo.UserId)
				if nil != myUserRecommendUserLocationLast {
					tmpStatus := myUserRecommendUserLocationLast.Status // 现在还在运行中
					tmpCurrent := myUserRecommendUserLocationLast.Current

					tmpBalanceAmount := currentValue / 100 * recommendNeed // 记录下一次
					myUserRecommendUserLocationLast.Status = "running"
					myUserRecommendUserLocationLast.Current += tmpBalanceAmount
					if myUserRecommendUserLocationLast.Current >= myUserRecommendUserLocationLast.CurrentMax { // 占位分红人分满停止
						myUserRecommendUserLocationLast.Status = "stop"
						if "running" == tmpStatus {
							myUserRecommendUserLocationLast.StopDate = time.Now().UTC().Add(8 * time.Hour)
						}
					}
					if 0 < tmpBalanceAmount {
						err = ruc.locationRepo.UpdateLocation(ctx, myUserRecommendUserLocationLast.ID, myUserRecommendUserLocationLast.Status, tmpBalanceAmount, myUserRecommendUserLocationLast.StopDate) // 分红占位数据修改
						if nil != err {
							return err
						}
					}
					amount -= tmpBalanceAmount // 扣除

					if 0 < tmpBalanceAmount && "running" == tmpStatus && tmpCurrent < myUserRecommendUserLocationLast.CurrentMax { // 这次还能分红
						tmpCurrentAmount := myUserRecommendUserLocationLast.CurrentMax - tmpCurrent // 最大可分红额度
						rewardAmount := tmpBalanceAmount
						if tmpCurrentAmount < tmpBalanceAmount { // 大于最大可分红额度
							rewardAmount = tmpCurrentAmount
						}
						_, err = ruc.userBalanceRepo.NormalRecommendReward(ctx, myUserRecommendUserId, rewardAmount, currentLocation.ID) // 直推人奖励
						if nil != err {
							return err
						}

					}
				}

				if nil != myUserRecommendUserLocationLast {
					var tmpMyRecommendAmount int64
					if 5 == myUserRecommendUserInfo.Vip { // 会员等级分红
						tmpMyRecommendAmount = currentValue / 100 * recommendNeedVip5
					} else if 4 == myUserRecommendUserInfo.Vip {
						tmpMyRecommendAmount = currentValue / 100 * recommendNeedVip4
					} else if 3 == myUserRecommendUserInfo.Vip {
						tmpMyRecommendAmount = currentValue / 100 * recommendNeedVip3
					} else if 2 == myUserRecommendUserInfo.Vip {
						tmpMyRecommendAmount = currentValue / 100 * recommendNeedVip2
					} else if 1 == myUserRecommendUserInfo.Vip {
						tmpMyRecommendAmount = currentValue / 100 * recommendNeedVip1
					}
					if 0 < tmpMyRecommendAmount { // 扣除推荐人分红
						tmpStatus := myUserRecommendUserLocationLast.Status // 现在还在运行中
						tmpCurrent := myUserRecommendUserLocationLast.Current

						tmpBalanceAmount := tmpMyRecommendAmount // 记录下一次
						myUserRecommendUserLocationLast.Status = "running"
						myUserRecommendUserLocationLast.Current += tmpBalanceAmount
						if myUserRecommendUserLocationLast.Current >= myUserRecommendUserLocationLast.CurrentMax { // 占位分红人分满停止
							myUserRecommendUserLocationLast.Status = "stop"
							if "running" == tmpStatus {
								myUserRecommendUserLocationLast.StopDate = time.Now().UTC().Add(8 * time.Hour)
							}
						}
						if 0 < tmpBalanceAmount {
							err = ruc.locationRepo.UpdateLocation(ctx, myUserRecommendUserLocationLast.ID, myUserRecommendUserLocationLast.Status, tmpBalanceAmount, myUserRecommendUserLocationLast.StopDate) // 分红占位数据修改
							if nil != err {
								return err
							}
						}
						amount -= tmpBalanceAmount                                                                                     // 扣除
						if 0 < tmpBalanceAmount && "running" == tmpStatus && tmpCurrent < myUserRecommendUserLocationLast.CurrentMax { // 这次还能分红
							tmpCurrentAmount := myUserRecommendUserLocationLast.CurrentMax - tmpCurrent // 最大可分红额度
							rewardAmount := tmpBalanceAmount
							if tmpCurrentAmount < tmpBalanceAmount { // 大于最大可分红额度
								rewardAmount = tmpCurrentAmount
							}
							_, err = ruc.userBalanceRepo.RecommendReward(ctx, myUserRecommendUserId, rewardAmount, currentLocation.ID) // 推荐人奖励
							if nil != err {
								return err
							}

						}
					}
				}
			}

			_, err = ruc.userBalanceRepo.Deposit(ctx, v.UserId, currentValue) // 充值
			if nil != err {
				return err
			}

			if 0 < locationCurrent && nil != myLastStopLocation {
				_, err = ruc.userBalanceRepo.DepositLast(ctx, v.UserId, locationCurrent, myLastStopLocation.ID) // 充值
				if nil != err {
					return err
				}
			}

			err = ruc.userBalanceRepo.SystemReward(ctx, amount, currentLocation.ID)
			if nil != err {
				return err
			}

			_, err = ruc.ethUserRecordRepo.CreateEthUserRecordListByHash(ctx, &EthUserRecord{
				Hash:     v.Hash,
				UserId:   v.UserId,
				Status:   v.Status,
				Type:     v.Type,
				Amount:   v.Amount,
				CoinType: v.CoinType,
			})
			if nil != err {
				return err
			}

			//dhbAmount, _ := strconv.ParseInt(ethUserRecord[k+1].Amount, 10, 64)
			//dhbAmount /= 100000000                                                             // 转换为系统精度
			//_, err = ruc.userBalanceRepo.DepositDhb(ctx, ethUserRecord[k+1].UserId, dhbAmount) // 充值
			//if nil != err {
			//	return err
			//}
			//
			//_, err = ruc.ethUserRecordRepo.CreateEthUserRecordListByHash(ctx, &EthUserRecord{
			//	Hash:     ethUserRecord[k+1].Hash,
			//	UserId:   ethUserRecord[k+1].UserId,
			//	Status:   ethUserRecord[k+1].Status,
			//	Type:     ethUserRecord[k+1].Type,
			//	Amount:   ethUserRecord[k+1].Amount,
			//	CoinType: ethUserRecord[k+1].CoinType,
			//})
			//if nil != err {
			//	return err
			//}

			return nil
		}); nil != err {
			continue
		}

		// 调整位置紧缩
		stopLocations, err = ruc.locationRepo.GetLocationsStopNotUpdate(ctx)
		if nil != stopLocations {
			// 调整位置紧缩
			for _, vStopLocations := range stopLocations {

				if err = ruc.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务
					err = ruc.locationRepo.UpdateLocationRowAndCol(ctx, vStopLocations.ID)
					if nil != err {
						return err
					}
					return nil
				}); nil != err {
					continue
				}
			}
		}
	}

	return true, nil
}

func (ruc *RecordUseCase) LockEthUserRecordHandle(ctx context.Context, ethUserRecord ...*EthUserRecord) (bool, error) {
	var (
		lock bool
	)
	// todo 全局锁
	for i := 0; i < 3; i++ {
		lock, _ = ruc.locationRepo.LockGlobalLocation(ctx)
		if lock {
			return true, nil
		}
		time.Sleep(5 * time.Second)
	}

	return false, nil
}

func (ruc *RecordUseCase) UnLockEthUserRecordHandle(ctx context.Context, ethUserRecord ...*EthUserRecord) (bool, error) {
	return ruc.locationRepo.UnLockGlobalLocation(ctx)
}
