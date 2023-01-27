package biz

import (
	"context"
	v1 "dhb/app/app/api"
	"encoding/base64"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"strings"
	"time"
)

type User struct {
	ID        int64
	Address   string
	CreatedAt time.Time
}

type UserInfo struct {
	ID               int64
	UserId           int64
	Vip              int64
	HistoryRecommend int64
}

type UserRecommend struct {
	ID            int64
	UserId        int64
	RecommendCode string
	CreatedAt     time.Time
}

type UserCurrentMonthRecommend struct {
	ID              int64
	UserId          int64
	RecommendUserId int64
	Date            time.Time
}

type Config struct {
	ID      int64
	KeyName string
	Name    string
	Value   string
}

type UserBalance struct {
	ID          int64
	UserId      int64
	BalanceUsdt int64
	BalanceDhb  int64
}

type Withdraw struct {
	ID              int64
	UserId          int64
	Amount          int64
	RelAmount       int64
	BalanceRecordId int64
	Status          string
	Type            string
	CreatedAt       time.Time
}

type UserSortRecommendReward struct {
	UserId int64
	Total  int64
}

type UserUseCase struct {
	repo                          UserRepo
	urRepo                        UserRecommendRepo
	configRepo                    ConfigRepo
	uiRepo                        UserInfoRepo
	ubRepo                        UserBalanceRepo
	locationRepo                  LocationRepo
	userCurrentMonthRecommendRepo UserCurrentMonthRecommendRepo
	tx                            Transaction
	log                           *log.Helper
}

type Reward struct {
	ID               int64
	UserId           int64
	Amount           int64
	BalanceRecordId  int64
	Type             string
	TypeRecordId     int64
	Reason           string
	ReasonLocationId int64
	LocationType     string
	CreatedAt        time.Time
}

type Pagination struct {
	PageNum  int
	PageSize int
}

type ConfigRepo interface {
	GetConfigByKeys(ctx context.Context, keys ...string) ([]*Config, error)
	GetConfigs(ctx context.Context) ([]*Config, error)
	UpdateConfig(ctx context.Context, id int64, value string) (bool, error)
}

type UserBalanceRepo interface {
	CreateUserBalance(ctx context.Context, u *User) (*UserBalance, error)
	LocationReward(ctx context.Context, userId int64, amount int64, locationId int64, myLocationId int64, locationType string) (int64, error)
	WithdrawReward(ctx context.Context, userId int64, amount int64, locationId int64, myLocationId int64, locationType string) (int64, error)
	RecommendReward(ctx context.Context, userId int64, amount int64, locationId int64) (int64, error)
	SystemWithdrawReward(ctx context.Context, amount int64, locationId int64) error
	SystemReward(ctx context.Context, amount int64, locationId int64) error
	SystemFee(ctx context.Context, amount int64, locationId int64) error
	GetSystemYesterdayDailyReward(ctx context.Context) (*Reward, error)
	UserFee(ctx context.Context, userId int64, amount int64) (int64, error)
	RecommendWithdrawReward(ctx context.Context, userId int64, amount int64, locationId int64) (int64, error)
	NormalRecommendReward(ctx context.Context, userId int64, amount int64, locationId int64) (int64, error)
	NormalWithdrawRecommendReward(ctx context.Context, userId int64, amount int64, locationId int64) (int64, error)
	Deposit(ctx context.Context, userId int64, amount int64) (int64, error)
	DepositLast(ctx context.Context, userId int64, lastAmount int64, locationId int64) (int64, error)
	DepositDhb(ctx context.Context, userId int64, amount int64) (int64, error)
	GetUserBalance(ctx context.Context, userId int64) (*UserBalance, error)
	GetUserRewardByUserId(ctx context.Context, userId int64) ([]*Reward, error)
	GetUserRewardByUserIds(ctx context.Context, userIds ...int64) (map[int64]*UserSortRecommendReward, error)
	GetUserRewards(ctx context.Context, b *Pagination, userId int64) ([]*Reward, error, int64)
	GetUserRewardsLastMonthFee(ctx context.Context) ([]*Reward, error)
	GetUserBalanceByUserIds(ctx context.Context, userIds ...int64) (map[int64]*UserBalance, error)
	GetUserBalanceUsdtTotal(ctx context.Context) (int64, error)
	GreateWithdraw(ctx context.Context, userId int64, amount int64, coinType string) (*Withdraw, error)
	WithdrawUsdt(ctx context.Context, userId int64, amount int64) error
	WithdrawDhb(ctx context.Context, userId int64, amount int64) error
	GetWithdrawByUserId(ctx context.Context, userId int64) ([]*Withdraw, error)
	GetWithdraws(ctx context.Context, b *Pagination, userId int64) ([]*Withdraw, error, int64)
	GetWithdrawPassOrRewarded(ctx context.Context) ([]*Withdraw, error)
	UpdateWithdraw(ctx context.Context, id int64, status string) (*Withdraw, error)
	GetWithdrawById(ctx context.Context, id int64) (*Withdraw, error)
	GetWithdrawNotDeal(ctx context.Context) ([]*Withdraw, error)
	GetUserBalanceRecordUserUsdtTotal(ctx context.Context, userId int64) (int64, error)
	GetUserBalanceRecordUsdtTotal(ctx context.Context) (int64, error)
	GetUserBalanceRecordUsdtTotalToday(ctx context.Context) (int64, error)
	GetUserWithdrawUsdtTotalToday(ctx context.Context) (int64, error)
	GetUserWithdrawUsdtTotal(ctx context.Context) (int64, error)
	GetUserRewardUsdtTotal(ctx context.Context) (int64, error)
	GetSystemRewardUsdtTotal(ctx context.Context) (int64, error)
	UpdateWithdrawAmount(ctx context.Context, id int64, status string, amount int64) (*Withdraw, error)
	GetUserRewardRecommendSort(ctx context.Context) ([]*UserSortRecommendReward, error)
	GetUserRewardTodayTotalByUserId(ctx context.Context, userId int64) (*UserSortRecommendReward, error)
}

type UserRecommendRepo interface {
	GetUserRecommendByUserId(ctx context.Context, userId int64) (*UserRecommend, error)
	CreateUserRecommend(ctx context.Context, u *User, recommendUser *UserRecommend) (*UserRecommend, error)
	GetUserRecommendByCode(ctx context.Context, code string) ([]*UserRecommend, error)
	GetUserRecommendLikeCode(ctx context.Context, code string) ([]*UserRecommend, error)
}

type UserCurrentMonthRecommendRepo interface {
	GetUserCurrentMonthRecommendByUserId(ctx context.Context, userId int64) ([]*UserCurrentMonthRecommend, error)
	GetUserCurrentMonthRecommendGroupByUserId(ctx context.Context, b *Pagination, userId int64) ([]*UserCurrentMonthRecommend, error, int64)
	CreateUserCurrentMonthRecommend(ctx context.Context, u *UserCurrentMonthRecommend) (*UserCurrentMonthRecommend, error)
	GetUserCurrentMonthRecommendCountByUserIds(ctx context.Context, userIds ...int64) (map[int64]int64, error)
	GetUserLastMonthRecommend(ctx context.Context) ([]int64, error)
}

type UserInfoRepo interface {
	CreateUserInfo(ctx context.Context, u *User) (*UserInfo, error)
	GetUserInfoByUserId(ctx context.Context, userId int64) (*UserInfo, error)
	UpdateUserInfo(ctx context.Context, u *UserInfo) (*UserInfo, error)
	GetUserInfoByUserIds(ctx context.Context, userIds ...int64) (map[int64]*UserInfo, error)
}

type UserRepo interface {
	GetUserById(ctx context.Context, Id int64) (*User, error)
	GetUserByAddresses(ctx context.Context, Addresses ...string) (map[string]*User, error)
	GetUserByAddress(ctx context.Context, address string) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByUserIds(ctx context.Context, userIds ...int64) (map[int64]*User, error)
	GetUsers(ctx context.Context, b *Pagination, address string) ([]*User, error, int64)
	GetUserCount(ctx context.Context) (int64, error)
	GetUserCountToday(ctx context.Context) (int64, error)
}

func NewUserUseCase(repo UserRepo, tx Transaction, configRepo ConfigRepo, uiRepo UserInfoRepo, urRepo UserRecommendRepo, locationRepo LocationRepo, userCurrentMonthRecommendRepo UserCurrentMonthRecommendRepo, ubRepo UserBalanceRepo, logger log.Logger) *UserUseCase {
	return &UserUseCase{
		repo:                          repo,
		tx:                            tx,
		configRepo:                    configRepo,
		locationRepo:                  locationRepo,
		userCurrentMonthRecommendRepo: userCurrentMonthRecommendRepo,
		uiRepo:                        uiRepo,
		urRepo:                        urRepo,
		ubRepo:                        ubRepo,
		log:                           log.NewHelper(logger),
	}
}

func (uuc *UserUseCase) GetUserByAddress(ctx context.Context, Addresses ...string) (map[string]*User, error) {
	return uuc.repo.GetUserByAddresses(ctx, Addresses...)
}

func (uuc *UserUseCase) GetDhbConfig(ctx context.Context) ([]*Config, error) {
	return uuc.configRepo.GetConfigByKeys(ctx, "level1Dhb", "level2Dhb", "level3Dhb")
}

func (uuc *UserUseCase) GetExistUserByAddressOrCreate(ctx context.Context, u *User, req *v1.EthAuthorizeRequest) (*User, error) {
	var (
		user          *User
		recommendUser *UserRecommend
		userRecommend *UserRecommend
		userInfo      *UserInfo
		userBalance   *UserBalance
		err           error
		userId        int64
		decodeBytes   []byte
	)

	user, err = uuc.repo.GetUserByAddress(ctx, u.Address) // 查询用户
	if nil == user || nil != err {
		code := req.SendBody.Code // 查询推荐码 abf00dd52c08a9213f225827bc3fb100 md5 dhbmachinefirst
		if "abf00dd52c08a9213f225827bc3fb100" != code {
			decodeBytes, err = base64.StdEncoding.DecodeString(code)
			code = string(decodeBytes)
			if 1 >= len(code) {
				return nil, errors.New(500, "USER_ERROR", "无效的推荐码")
			}
			if userId, err = strconv.ParseInt(code[1:], 10, 64); 0 >= userId || nil != err {
				return nil, errors.New(500, "USER_ERROR", "无效的推荐码")
			}

			// 查询推荐人的相关信息
			recommendUser, err = uuc.urRepo.GetUserRecommendByUserId(ctx, userId)
			if err != nil {
				return nil, errors.New(500, "USER_ERROR", "无效的推荐码")
			}
		}

		if err = uuc.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务
			user, err = uuc.repo.CreateUser(ctx, u) // 用户创建
			if err != nil {
				return err
			}

			userInfo, err = uuc.uiRepo.CreateUserInfo(ctx, user) // 创建用户信息
			if err != nil {
				return err
			}

			userRecommend, err = uuc.urRepo.CreateUserRecommend(ctx, user, recommendUser) // 创建用户信息
			if err != nil {
				return err
			}

			userBalance, err = uuc.ubRepo.CreateUserBalance(ctx, user) // 创建余额信息
			if err != nil {
				return err
			}

			return nil
		}); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (uuc *UserUseCase) UserInfo(ctx context.Context, user *User) (*v1.UserInfoReply, error) {
	var (
		myUser                   *User
		userInfo                 *UserInfo
		locations                []*Location
		userBalance              *UserBalance
		userRecommend            *UserRecommend
		userRecommends           []*UserRecommend
		userRewards              []*Reward
		rewardLocations          []*Location
		userRewardTotal          int64
		encodeString             string
		myUserRecommendUserId    int64
		myRecommendUser          *User
		myRow                    int64
		rowNum                   int64
		colNum                   int64
		myCol                    int64
		recommendTeamNum         int64
		recommendTotal           int64
		recommendVipTotal        int64
		feeDaily                 int64
		locationTotal            int64
		locationTotalCol         int64
		locationTotalRow         int64
		myCode                   string
		inviteUserAddress        string
		amount                   string
		userCount                string
		status                   = "no"
		configs                  []*Config
		myLastStopLocation       *Location
		myLastLocationCurrent    int64
		hasRunningLocation       bool
		myWithdraws              []*Withdraw
		userLocations            []*Location
		totalDepoist             int64
		fee                      int64
		withdrawAmount           int64
		userSortRecommendRewards []*UserSortRecommendReward
		topUsersReply            []*v1.UserInfoReply_List
		topUsers                 map[int64]*User
		topUserIds               []int64
		locationCount            int64
		poolAmount               int64
		systemYesterdayreward    *Reward
		userTodayRewardTotal     *UserSortRecommendReward
		bnbAmount                string
		userTodayReward          int64
		recommendTop             int64
		err                      error
	)

	myUser, err = uuc.repo.GetUserById(ctx, user.ID)
	if nil != err {
		return nil, err
	}

	userInfo, err = uuc.uiRepo.GetUserInfoByUserId(ctx, myUser.ID)
	if nil != err {
		return nil, err
	}

	locations, err = uuc.locationRepo.GetLocationsByUserId(ctx, myUser.ID)
	if nil != locations && 0 < len(locations) {
		status = "stop"
		for _, v := range locations {
			if "running" == v.Status {
				status = "running"
				if 0 == v.Current {
					status = "yes"
				}
				hasRunningLocation = true
				if v.CurrentMax > v.Current {
					amount = fmt.Sprintf("%.2f", float64(v.CurrentMax-v.Current)/float64(10000000000))
				} else {
					amount = fmt.Sprintf("%.2f", float64(v.CurrentMax)/float64(10000000000))
				}
				myCol = v.Col
				myRow = v.Row
				break
			}
		}
	}

	locationCount = int64(len(locations))

	// 提现记录
	myWithdraws, err = uuc.ubRepo.GetWithdrawByUserId(ctx, myUser.ID)
	for _, vMyWithdraw := range myWithdraws {
		withdrawAmount += vMyWithdraw.RelAmount
	}

	// 充值记录
	totalDepoist, err = uuc.ubRepo.GetUserBalanceRecordUserUsdtTotal(ctx, myUser.ID)

	now := time.Now().UTC().Add(8 * time.Hour)
	myLastStopLocation, err = uuc.locationRepo.GetMyStopLocationLast(ctx, myUser.ID) // 冻结
	if nil != myLastStopLocation && now.Before(myLastStopLocation.StopDate.Add(24*time.Hour)) && !hasRunningLocation {
		myLastLocationCurrent = myLastStopLocation.Current - myLastStopLocation.CurrentMax // 补上
	}

	userBalance, err = uuc.ubRepo.GetUserBalance(ctx, myUser.ID)
	if nil != err {
		return nil, err
	}

	userRecommend, err = uuc.urRepo.GetUserRecommendByUserId(ctx, myUser.ID)
	if nil == userRecommend {
		return nil, err
	}

	myCode = "D" + strconv.FormatInt(myUser.ID, 10)
	codeByte := []byte(myCode)
	encodeString = base64.StdEncoding.EncodeToString(codeByte)

	if "" != userRecommend.RecommendCode {
		tmpRecommendUserIds := strings.Split(userRecommend.RecommendCode, "D")
		if 2 <= len(tmpRecommendUserIds) {
			myUserRecommendUserId, _ = strconv.ParseInt(tmpRecommendUserIds[len(tmpRecommendUserIds)-1], 10, 64) // 最后一位是直推人
		}
		myRecommendUser, err = uuc.repo.GetUserById(ctx, myUserRecommendUserId)
		if nil != err {
			return nil, err
		}
		inviteUserAddress = myRecommendUser.Address
		myCode = userRecommend.RecommendCode + myCode
	}

	// 团队
	userRecommends, err = uuc.urRepo.GetUserRecommendLikeCode(ctx, myCode)
	if nil != userRecommends {
		recommendTeamNum = int64(len(userRecommends))
	}

	// 累计奖励
	userRewards, err = uuc.ubRepo.GetUserRewardByUserId(ctx, myUser.ID)
	if nil != userRewards {
		for _, vUserReward := range userRewards {
			userRewardTotal += vUserReward.Amount
			if "recommend" == vUserReward.Reason {
				recommendTotal += vUserReward.Amount
			} else if "location" == vUserReward.Reason {
				locationTotal += vUserReward.Amount
				if "col" == vUserReward.LocationType {
					locationTotalCol += vUserReward.Amount
				} else if "row" == vUserReward.LocationType {
					locationTotalRow += vUserReward.Amount
				}
			} else if "recommend_vip" == vUserReward.Reason {
				recommendVipTotal += vUserReward.Amount
			} else if "fee_daily" == vUserReward.Reason {
				feeDaily += vUserReward.Amount
			} else if "recommend_top" == vUserReward.Reason {
				recommendTop += vUserReward.Amount
			} else if "recommend_vip_top" == vUserReward.Reason {
				recommendVipTotal += vUserReward.Amount
			}
		}
	}

	// 位置
	if 0 < myRow && 0 < myCol {
		rewardLocations, err = uuc.locationRepo.GetRewardLocationByRowOrCol(ctx, myRow, myCol)
		if nil != rewardLocations {
			for _, vRewardLocation := range rewardLocations {
				if myRow == vRewardLocation.Row && myCol == vRewardLocation.Col { // 跳过自己
					continue
				}
				if myRow == vRewardLocation.Row {
					colNum++
				}
				if myCol == vRewardLocation.Col {
					rowNum++
				}
			}
		}
	}

	// 配置
	configs, err = uuc.configRepo.GetConfigByKeys(ctx, "user_count", "bnb_amount")
	if nil != configs {
		for _, vConfig := range configs {
			if "user_count" == vConfig.KeyName {
				userCount = vConfig.Value
			}
			if "bnb_amount" == vConfig.KeyName {
				bnbAmount = vConfig.Value
			}
		}
	}

	// 全网手续费
	userLocations, err = uuc.locationRepo.GetLocationDaily(ctx)
	if nil != err {
		return nil, err
	}

	for _, userLocation := range userLocations {
		fee += userLocation.CurrentMax / 5
	}

	// 昨日剩余全网手续费
	systemYesterdayreward, _ = uuc.ubRepo.GetSystemYesterdayDailyReward(ctx)
	if nil != systemYesterdayreward {
		fee += systemYesterdayreward.Amount
	}

	if fee > 0 {
		poolAmount = fee / 100 * 3
		fee = fee / 10000 * 3 * 70
	}

	// 前四
	userSortRecommendRewards, err = uuc.ubRepo.GetUserRewardRecommendSort(ctx)
	if 0 < len(userSortRecommendRewards) {
		for _, userSortRecommendReward := range userSortRecommendRewards {
			topUserIds = append(topUserIds, userSortRecommendReward.UserId)
		}
		topUsersReply = make([]*v1.UserInfoReply_List, 0)
		if 0 < len(topUserIds) {
			topUsers, err = uuc.repo.GetUserByUserIds(ctx, topUserIds...)
			if nil != topUsers {
				for k, userSortRecommendReward := range userSortRecommendRewards {
					if _, ok := topUsers[userSortRecommendReward.UserId]; !ok {
						continue
					}

					var tmpAmount int64
					if 0 == k {
						tmpAmount = fee / 100 * 40
					} else if 1 == k {
						tmpAmount = fee / 100 * 30
					} else if 2 == k {
						tmpAmount = fee / 100 * 20
					} else if 3 == k {
						tmpAmount = fee / 100 * 10
					}

					topUsersReply = append(topUsersReply, &v1.UserInfoReply_List{
						Account:         topUsers[userSortRecommendReward.UserId].Address,
						RecommendReward: fmt.Sprintf("%.2f", float64(userSortRecommendReward.Total)/float64(10000000000)),
						Reward:          fmt.Sprintf("%.2f", float64(tmpAmount)/float64(10000000000)),
					})
				}
			}
		}
	}

	// 今日收益
	userTodayRewardTotal, err = uuc.ubRepo.GetUserRewardTodayTotalByUserId(ctx, myUser.ID)
	if nil != userTodayRewardTotal {
		userTodayReward = userTodayRewardTotal.Total
	}

	return &v1.UserInfoReply{
		Address:           myUser.Address,
		Level:             userInfo.Vip,
		Status:            status,
		Amount:            amount,
		RecommendVipTotal: fmt.Sprintf("%.2f", float64(recommendVipTotal)/float64(10000000000)),
		FeeDaily:          fmt.Sprintf("%.2f", float64(feeDaily)/float64(10000000000)),
		BalanceUsdt:       fmt.Sprintf("%.2f", float64(userBalance.BalanceUsdt)/float64(10000000000)),
		BalanceDhb:        fmt.Sprintf("%.2f", float64(userBalance.BalanceDhb)/float64(10000000000)),
		InviteUrl:         encodeString,
		InviteUserAddress: inviteUserAddress,
		RecommendNum:      userInfo.HistoryRecommend,
		RecommendTeamNum:  recommendTeamNum,
		Total:             fmt.Sprintf("%.2f", float64(userRewardTotal)/float64(10000000000)),
		WithdrawAmount:    fmt.Sprintf("%.2f", float64(withdrawAmount)/float64(10000000000)),
		Row:               rowNum,
		Col:               colNum,
		RecommendTotal:    fmt.Sprintf("%.2f", float64(recommendTotal)/float64(10000000000)),
		LocationTotal:     fmt.Sprintf("%.2f", float64(locationTotal)/float64(10000000000)),
		Usdt:              "0x55d398326f99059fF775485246999027B3197955",
		Account:           "0x12ad46402226ebf2b5dff48485bd0b1eb84cabce",
		AmountB:           fmt.Sprintf("%.2f", float64(myLastLocationCurrent)/float64(10000000000)),
		UserCount:         userCount,
		TotalDeposit:      fmt.Sprintf("%.2f", float64(totalDepoist)/float64(10000000000)),
		PoolAmount:        fmt.Sprintf("%.2f", float64(poolAmount)/float64(10000000000)),
		TopUser:           topUsersReply,
		LocationCount:     locationCount,
		TodayReward:       fmt.Sprintf("%.2f", float64(userTodayReward)/float64(10000000000)),
		RecommendTop:      fmt.Sprintf("%.2f", float64(recommendTop)/float64(10000000000)),
		LocationTotalCol:  fmt.Sprintf("%.2f", float64(locationTotalCol)/float64(10000000000)),
		LocationTotalRow:  fmt.Sprintf("%.2f", float64(locationTotalRow)/float64(10000000000)),
		BnbAmount:         bnbAmount,
	}, nil
}

func (uuc *UserUseCase) RewardList(ctx context.Context, req *v1.RewardListRequest, user *User) (*v1.RewardListReply, error) {
	var (
		userRewards    []*Reward
		locationIdsMap map[int64]int64
		locations      map[int64]*Location
		err            error
	)
	res := &v1.RewardListReply{
		Rewards: make([]*v1.RewardListReply_List, 0),
	}

	userRewards, err = uuc.ubRepo.GetUserRewardByUserId(ctx, user.ID)
	if nil != err {
		return res, nil
	}

	locationIdsMap = make(map[int64]int64, 0)
	if nil != userRewards {
		for _, vUserReward := range userRewards {
			if "location" == vUserReward.Reason && req.Type == vUserReward.LocationType && 1 <= vUserReward.ReasonLocationId {
				locationIdsMap[vUserReward.ReasonLocationId] = vUserReward.ReasonLocationId
			}
		}

		var tmpLocationIds []int64
		for _, v := range locationIdsMap {
			tmpLocationIds = append(tmpLocationIds, v)
		}
		if 0 >= len(tmpLocationIds) {
			return res, nil
		}

		locations, err = uuc.locationRepo.GetRewardLocationByIds(ctx, tmpLocationIds...)

		for _, vUserReward := range userRewards {
			if "location" == vUserReward.Reason && req.Type == vUserReward.LocationType {
				if _, ok := locations[vUserReward.ReasonLocationId]; !ok {
					continue
				}

				res.Rewards = append(res.Rewards, &v1.RewardListReply_List{
					CreatedAt:      vUserReward.CreatedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
					Amount:         fmt.Sprintf("%.2f", float64(vUserReward.Amount)/float64(10000000000)),
					LocationStatus: locations[vUserReward.ReasonLocationId].Status,
					Type:           vUserReward.Type,
				})
			}
		}
	}

	return res, nil
}

func (uuc *UserUseCase) RecommendRewardList(ctx context.Context, user *User) (*v1.RecommendRewardListReply, error) {
	var (
		userRewards []*Reward
		err         error
	)
	res := &v1.RecommendRewardListReply{
		Rewards: make([]*v1.RecommendRewardListReply_List, 0),
	}

	userRewards, err = uuc.ubRepo.GetUserRewardByUserId(ctx, user.ID)
	if nil != err {
		return res, nil
	}

	for _, vUserReward := range userRewards {
		if "recommend" == vUserReward.Reason || "recommend_vip" == vUserReward.Reason {
			res.Rewards = append(res.Rewards, &v1.RecommendRewardListReply_List{
				CreatedAt: vUserReward.CreatedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
				Amount:    fmt.Sprintf("%.2f", float64(vUserReward.Amount)/float64(10000000000)),
				Type:      vUserReward.Type,
				Reason:    vUserReward.Reason,
			})
		}
	}

	return res, nil
}

func (uuc *UserUseCase) FeeRewardList(ctx context.Context, user *User) (*v1.FeeRewardListReply, error) {
	var (
		userRewards []*Reward
		err         error
	)
	res := &v1.FeeRewardListReply{
		Rewards: make([]*v1.FeeRewardListReply_List, 0),
	}

	userRewards, err = uuc.ubRepo.GetUserRewardByUserId(ctx, user.ID)
	if nil != err {
		return res, nil
	}

	for _, vUserReward := range userRewards {
		if "fee" == vUserReward.Reason {
			res.Rewards = append(res.Rewards, &v1.FeeRewardListReply_List{
				CreatedAt: vUserReward.CreatedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
				Amount:    fmt.Sprintf("%.2f", float64(vUserReward.Amount)/float64(10000000000)),
			})
		}
	}

	return res, nil
}

func (uuc *UserUseCase) WithdrawList(ctx context.Context, user *User) (*v1.WithdrawListReply, error) {

	var (
		withdraws []*Withdraw
		err       error
	)

	res := &v1.WithdrawListReply{
		Withdraw: make([]*v1.WithdrawListReply_List, 0),
	}

	withdraws, err = uuc.ubRepo.GetWithdrawByUserId(ctx, user.ID)
	if nil != err {
		return res, err
	}

	for _, v := range withdraws {
		res.Withdraw = append(res.Withdraw, &v1.WithdrawListReply_List{
			CreatedAt: v.CreatedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
			Amount:    fmt.Sprintf("%.2f", float64(v.RelAmount)/float64(10000000000)),
			Status:    v.Status,
			Type:      v.Type,
		})
	}

	return res, nil
}

func (uuc *UserUseCase) Withdraw(ctx context.Context, req *v1.WithdrawRequest, user *User) (*v1.WithdrawReply, error) {
	var (
		err         error
		userBalance *UserBalance
	)

	if "dhb" != req.SendBody.Type && "usdt" != req.SendBody.Type {
		return &v1.WithdrawReply{
			Status: "fail",
		}, nil
	}

	userBalance, err = uuc.ubRepo.GetUserBalance(ctx, user.ID)
	if nil != err {
		return nil, err
	}

	amountFloat, _ := strconv.ParseFloat(req.SendBody.Amount, 10)
	amountFloat *= 10000000000
	amount, _ := strconv.ParseInt(strconv.FormatFloat(amountFloat, 'f', -1, 64), 10, 64)
	if 0 >= amount {
		return &v1.WithdrawReply{
			Status: "fail",
		}, nil
	}

	if "dhb" == req.SendBody.Type && userBalance.BalanceDhb < amount {
		return &v1.WithdrawReply{
			Status: "fail",
		}, nil
	}

	if "usdt" == req.SendBody.Type && userBalance.BalanceUsdt < amount {
		return &v1.WithdrawReply{
			Status: "fail",
		}, nil
	}
	if err = uuc.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务

		if "usdt" == req.SendBody.Type {
			err = uuc.ubRepo.WithdrawUsdt(ctx, user.ID, amount) // 提现
			if nil != err {
				return err
			}
			_, err = uuc.ubRepo.GreateWithdraw(ctx, user.ID, amount, req.SendBody.Type)
			if nil != err {
				return err
			}

		} else if "dhb" == req.SendBody.Type {
			err = uuc.ubRepo.WithdrawDhb(ctx, user.ID, amount) // 提现
			if nil != err {
				return err
			}
			_, err = uuc.ubRepo.GreateWithdraw(ctx, user.ID, amount, req.SendBody.Type)
			if nil != err {
				return err
			}
		}

		return nil
	}); nil != err {
		return nil, err
	}

	return &v1.WithdrawReply{
		Status: "ok",
	}, nil
}

func (uuc *UserUseCase) AdminRewardList(ctx context.Context, req *v1.AdminRewardListRequest) (*v1.AdminRewardListReply, error) {
	var (
		userSearch  *User
		userId      int64 = 0
		userRewards []*Reward
		users       map[int64]*User
		userIdsMap  map[int64]int64
		userIds     []int64
		err         error
		count       int64
	)
	res := &v1.AdminRewardListReply{
		Rewards: make([]*v1.AdminRewardListReply_List, 0),
	}

	// 地址查询
	if "" != req.Address {
		userSearch, err = uuc.repo.GetUserByAddress(ctx, req.Address)
		if nil != err {
			return res, nil
		}
		userId = userSearch.ID
	}

	userRewards, err, count = uuc.ubRepo.GetUserRewards(ctx, &Pagination{
		PageNum:  int(req.Page),
		PageSize: 10,
	}, userId)
	if nil != err {
		return res, nil
	}
	res.Count = count

	userIdsMap = make(map[int64]int64, 0)
	for _, vUserReward := range userRewards {
		userIdsMap[vUserReward.UserId] = vUserReward.UserId
	}
	for _, v := range userIdsMap {
		userIds = append(userIds, v)
	}

	users, err = uuc.repo.GetUserByUserIds(ctx, userIds...)
	for _, vUserReward := range userRewards {
		tmpUser := ""
		if nil != users {
			if _, ok := users[vUserReward.UserId]; ok {
				tmpUser = users[vUserReward.UserId].Address
			}
		}

		res.Rewards = append(res.Rewards, &v1.AdminRewardListReply_List{
			CreatedAt: vUserReward.CreatedAt.Format("2006-01-02 15:04:05"),
			Amount:    fmt.Sprintf("%.2f", float64(vUserReward.Amount)/float64(10000000000)),
			Type:      vUserReward.Type,
			Address:   tmpUser,
			Reason:    vUserReward.Reason,
		})
	}

	return res, nil
}

func (uuc *UserUseCase) AdminUserList(ctx context.Context, req *v1.AdminUserListRequest) (*v1.AdminUserListReply, error) {
	var (
		users                          []*User
		userIds                        []int64
		userBalances                   map[int64]*UserBalance
		userInfos                      map[int64]*UserInfo
		userCurrentMonthRecommendCount map[int64]int64
		count                          int64
		err                            error
	)

	res := &v1.AdminUserListReply{
		Users: make([]*v1.AdminUserListReply_UserList, 0),
	}

	users, err, count = uuc.repo.GetUsers(ctx, &Pagination{
		PageNum:  int(req.Page),
		PageSize: 10,
	}, req.Address)
	if nil != err {
		return res, nil
	}
	res.Count = count

	for _, vUsers := range users {
		userIds = append(userIds, vUsers.ID)
	}

	userBalances, err = uuc.ubRepo.GetUserBalanceByUserIds(ctx, userIds...)
	if nil != err {
		return res, nil
	}

	userInfos, err = uuc.uiRepo.GetUserInfoByUserIds(ctx, userIds...)
	if nil != err {
		return res, nil
	}

	userCurrentMonthRecommendCount, err = uuc.userCurrentMonthRecommendRepo.GetUserCurrentMonthRecommendCountByUserIds(ctx, userIds...)

	for _, v := range users {
		if _, ok := userBalances[v.ID]; !ok {
			continue
		}
		if _, ok := userInfos[v.ID]; !ok {
			continue
		}

		var tmpCount int64
		if nil != userCurrentMonthRecommendCount {
			if _, ok := userCurrentMonthRecommendCount[v.ID]; ok {
				tmpCount = userCurrentMonthRecommendCount[v.ID]
			}
		}

		res.Users = append(res.Users, &v1.AdminUserListReply_UserList{
			UserId:           v.ID,
			CreatedAt:        v.CreatedAt.Format("2006-01-02 15:04:05"),
			Address:          v.Address,
			BalanceUsdt:      fmt.Sprintf("%.2f", float64(userBalances[v.ID].BalanceUsdt)/float64(10000000000)),
			BalanceDhb:       fmt.Sprintf("%.2f", float64(userBalances[v.ID].BalanceDhb)/float64(10000000000)),
			Vip:              userInfos[v.ID].Vip,
			MonthRecommend:   tmpCount,
			HistoryRecommend: userInfos[v.ID].HistoryRecommend,
		})
	}

	return res, nil
}

func (uuc *UserUseCase) GetUserByUserIds(ctx context.Context, userIds ...int64) (map[int64]*User, error) {
	return uuc.repo.GetUserByUserIds(ctx, userIds...)
}

func (uuc *UserUseCase) AdminLocationList(ctx context.Context, req *v1.AdminLocationListRequest) (*v1.AdminLocationListReply, error) {
	var (
		locations  []*Location
		userSearch *User
		userId     int64
		userIds    []int64
		userIdsMap map[int64]int64
		users      map[int64]*User
		count      int64
		err        error
	)

	res := &v1.AdminLocationListReply{
		Locations: make([]*v1.AdminLocationListReply_LocationList, 0),
	}

	// 地址查询
	if "" != req.Address {
		userSearch, err = uuc.repo.GetUserByAddress(ctx, req.Address)
		if nil != err {
			return res, nil
		}
		userId = userSearch.ID
	}

	locations, err, count = uuc.locationRepo.GetLocations(ctx, &Pagination{
		PageNum:  int(req.Page),
		PageSize: 10,
	}, userId)
	if nil != err {
		return res, nil
	}
	res.Count = count

	userIdsMap = make(map[int64]int64, 0)
	for _, vLocations := range locations {
		userIdsMap[vLocations.UserId] = vLocations.UserId
	}
	for _, v := range userIdsMap {
		userIds = append(userIds, v)
	}

	users, err = uuc.repo.GetUserByUserIds(ctx, userIds...)
	if nil != err {
		return res, nil
	}

	for _, v := range locations {
		if _, ok := users[v.UserId]; !ok {
			continue
		}

		res.Locations = append(res.Locations, &v1.AdminLocationListReply_LocationList{
			CreatedAt:    v.CreatedAt.Format("2006-01-02 15:04:05"),
			Address:      users[v.UserId].Address,
			Row:          v.Row,
			Col:          v.Col,
			Status:       v.Status,
			CurrentLevel: v.CurrentLevel,
			Current:      fmt.Sprintf("%.2f", float64(v.Current)/float64(10000000000)),
			CurrentMax:   fmt.Sprintf("%.2f", float64(v.CurrentMax)/float64(10000000000)),
		})
	}

	return res, nil

}

func (uuc *UserUseCase) AdminRecommendList(ctx context.Context, req *v1.AdminUserRecommendRequest) (*v1.AdminUserRecommendReply, error) {
	var (
		userRecommends []*UserRecommend
		userRecommend  *UserRecommend
		userIdsMap     map[int64]int64
		userIds        []int64
		users          map[int64]*User
		err            error
	)

	res := &v1.AdminUserRecommendReply{
		Users: make([]*v1.AdminUserRecommendReply_List, 0),
	}

	// 地址查询
	if 0 < req.UserId {
		userRecommend, err = uuc.urRepo.GetUserRecommendByUserId(ctx, req.UserId)
		if nil == userRecommend {
			return res, nil
		}

		userRecommends, err = uuc.urRepo.GetUserRecommendByCode(ctx, userRecommend.RecommendCode+"D"+strconv.FormatInt(userRecommend.UserId, 10))
		if nil != err {
			return res, nil
		}
	}

	userIdsMap = make(map[int64]int64, 0)
	for _, vLocations := range userRecommends {
		userIdsMap[vLocations.UserId] = vLocations.UserId
	}
	for _, v := range userIdsMap {
		userIds = append(userIds, v)
	}

	users, err = uuc.repo.GetUserByUserIds(ctx, userIds...)
	if nil != err {
		return res, nil
	}

	for _, v := range userRecommends {
		if _, ok := users[v.UserId]; !ok {
			continue
		}

		res.Users = append(res.Users, &v1.AdminUserRecommendReply_List{
			Address:   users[v.UserId].Address,
			Id:        v.ID,
			UserId:    v.UserId,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return res, nil
}

func (uuc *UserUseCase) AdminMonthRecommend(ctx context.Context, req *v1.AdminMonthRecommendRequest) (*v1.AdminMonthRecommendReply, error) {
	var (
		userCurrentMonthRecommends []*UserCurrentMonthRecommend
		searchUser                 *User
		userIdsMap                 map[int64]int64
		userIds                    []int64
		searchUserId               int64
		users                      map[int64]*User
		count                      int64
		err                        error
	)

	res := &v1.AdminMonthRecommendReply{
		Users: make([]*v1.AdminMonthRecommendReply_List, 0),
	}

	// 地址查询
	if "" != req.Address {
		searchUser, err = uuc.repo.GetUserByAddress(ctx, req.Address)
		if nil == searchUser {
			return res, nil
		}
		searchUserId = searchUser.ID
	}

	userCurrentMonthRecommends, err, count = uuc.userCurrentMonthRecommendRepo.GetUserCurrentMonthRecommendGroupByUserId(ctx, &Pagination{
		PageNum:  int(req.Page),
		PageSize: 10,
	}, searchUserId)
	if nil != err {
		return res, nil
	}
	res.Count = count

	userIdsMap = make(map[int64]int64, 0)
	for _, vRecommend := range userCurrentMonthRecommends {
		userIdsMap[vRecommend.UserId] = vRecommend.UserId
		userIdsMap[vRecommend.RecommendUserId] = vRecommend.RecommendUserId
	}
	for _, v := range userIdsMap {
		userIds = append(userIds, v)
	}

	users, err = uuc.repo.GetUserByUserIds(ctx, userIds...)
	if nil != err {
		return res, nil
	}

	for _, v := range userCurrentMonthRecommends {
		if _, ok := users[v.UserId]; !ok {
			continue
		}

		res.Users = append(res.Users, &v1.AdminMonthRecommendReply_List{
			Address:          users[v.UserId].Address,
			Id:               v.ID,
			RecommendAddress: users[v.RecommendUserId].Address,
			CreatedAt:        v.Date.Format("2006-01-02 15:04:05"),
		})
	}

	return res, nil
}

func (uuc *UserUseCase) AdminConfig(ctx context.Context, req *v1.AdminConfigRequest) (*v1.AdminConfigReply, error) {
	var (
		configs []*Config
	)

	res := &v1.AdminConfigReply{
		Config: make([]*v1.AdminConfigReply_List, 0),
	}

	configs, _ = uuc.configRepo.GetConfigs(ctx)
	if nil == configs {
		return res, nil
	}

	for _, v := range configs {
		res.Config = append(res.Config, &v1.AdminConfigReply_List{
			Id:    v.ID,
			Name:  v.Name,
			Value: v.Value,
		})
	}

	return res, nil
}

func (uuc *UserUseCase) AdminConfigUpdate(ctx context.Context, req *v1.AdminConfigUpdateRequest) (*v1.AdminConfigUpdateReply, error) {
	var (
		err error
	)

	res := &v1.AdminConfigUpdateReply{}

	_, err = uuc.configRepo.UpdateConfig(ctx, req.SendBody.Id, req.SendBody.Value)
	if nil != err {
		return res, err
	}

	return res, nil
}

func (uuc *UserUseCase) GetWithdrawPassOrRewardedList(ctx context.Context) ([]*Withdraw, error) {
	return uuc.ubRepo.GetWithdrawPassOrRewarded(ctx)
}

func (uuc *UserUseCase) UpdateWithdrawDoing(ctx context.Context, id int64) (*Withdraw, error) {
	return uuc.ubRepo.UpdateWithdraw(ctx, id, "doing")
}

func (uuc *UserUseCase) UpdateWithdrawSuccess(ctx context.Context, id int64) (*Withdraw, error) {
	return uuc.ubRepo.UpdateWithdraw(ctx, id, "success")
}

func (uuc *UserUseCase) AdminWithdrawList(ctx context.Context, req *v1.AdminWithdrawListRequest) (*v1.AdminWithdrawListReply, error) {
	var (
		withdraws  []*Withdraw
		userIds    []int64
		userSearch *User
		userId     int64
		userIdsMap map[int64]int64
		users      map[int64]*User
		count      int64
		err        error
	)

	res := &v1.AdminWithdrawListReply{
		Withdraw: make([]*v1.AdminWithdrawListReply_List, 0),
	}

	// 地址查询
	if "" != req.Address {
		userSearch, err = uuc.repo.GetUserByAddress(ctx, req.Address)
		if nil != err {
			return res, nil
		}
		userId = userSearch.ID
	}

	withdraws, err, count = uuc.ubRepo.GetWithdraws(ctx, &Pagination{
		PageNum:  int(req.Page),
		PageSize: 10,
	}, userId)
	if nil != err {
		return res, err
	}
	res.Count = count

	userIdsMap = make(map[int64]int64, 0)
	for _, vWithdraws := range withdraws {
		userIdsMap[vWithdraws.UserId] = vWithdraws.UserId
	}
	for _, v := range userIdsMap {
		userIds = append(userIds, v)
	}

	users, err = uuc.repo.GetUserByUserIds(ctx, userIds...)
	if nil != err {
		return res, nil
	}

	for _, v := range withdraws {
		if _, ok := users[v.UserId]; !ok {
			continue
		}
		res.Withdraw = append(res.Withdraw, &v1.AdminWithdrawListReply_List{
			Id:        v.ID,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			Amount:    fmt.Sprintf("%.2f", float64(v.Amount)/float64(10000000000)),
			Status:    v.Status,
			Type:      v.Type,
			Address:   users[v.UserId].Address,
			RelAmount: fmt.Sprintf("%.2f", float64(v.RelAmount)/float64(10000000000)),
		})
	}

	return res, nil

}

func (uuc *UserUseCase) AdminFee(ctx context.Context, req *v1.AdminFeeRequest) (*v1.AdminFeeReply, error) {

	var (
		userIds        []int64
		userRewardFees []*Reward
		userCount      int64
		fee            int64
		myLocationLast *Location
		err            error
	)

	userIds, err = uuc.userCurrentMonthRecommendRepo.GetUserLastMonthRecommend(ctx)
	if nil != err {
		return nil, err
	}

	if 0 >= len(userIds) {
		return &v1.AdminFeeReply{}, err
	}

	// 全网手续费
	userRewardFees, err = uuc.ubRepo.GetUserRewardsLastMonthFee(ctx)
	if nil != err {
		return nil, err
	}

	for _, vUserRewardFee := range userRewardFees {
		fee += vUserRewardFee.Amount
	}

	if 0 >= fee {
		return &v1.AdminFeeReply{}, err
	}

	userCount = int64(len(userIds))
	fee = fee / 100 / userCount

	for _, v := range userIds {
		// 获取当前用户的占位信息，已经有运行中的跳过
		myLocationLast, err = uuc.locationRepo.GetMyLocationRunningLast(ctx, v)
		if nil == myLocationLast { // 无占位信息
			continue
		}

		if err = uuc.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务
			tmpCurrentStatus := myLocationLast.Status // 现在还在运行中
			tmpCurrent := myLocationLast.Current
			tmpBalanceAmount := fee
			myLocationLast.Status = "running"
			myLocationLast.Current += fee
			if myLocationLast.Current >= myLocationLast.CurrentMax { // 占位分红人分满停止
				if "running" == tmpCurrentStatus {
					myLocationLast.StopDate = time.Now().UTC().Add(8 * time.Hour)
				}
				myLocationLast.Status = "stop"
			}

			if 0 < tmpBalanceAmount {
				err = uuc.locationRepo.UpdateLocation(ctx, myLocationLast.ID, myLocationLast.Status, tmpBalanceAmount, myLocationLast.StopDate) // 分红占位数据修改
				if nil != err {
					return err
				}

				if 0 < tmpBalanceAmount && "running" == tmpCurrentStatus && tmpCurrent < myLocationLast.CurrentMax { // 这次还能分红
					tmpCurrentAmount := myLocationLast.CurrentMax - tmpCurrent // 最大可分红额度
					rewardAmount := tmpBalanceAmount
					if tmpCurrentAmount < tmpBalanceAmount { // 大于最大可分红额度
						rewardAmount = tmpCurrentAmount
					}

					_, err = uuc.ubRepo.UserFee(ctx, v, rewardAmount)
					if nil != err {
						return err
					}
				}
			}

			return nil
		}); nil != err {
			return nil, err
		}
	}

	return &v1.AdminFeeReply{}, err
}

func (uuc *UserUseCase) AdminAll(ctx context.Context, req *v1.AdminAllRequest) (*v1.AdminAllReply, error) {

	var (
		userCount                       int64
		userTodayCount                  int64
		userBalanceUsdtTotal            int64
		userBalanceRecordUsdtTotal      int64
		userBalanceRecordUsdtTotalToday int64
		userWithdrawUsdtTotalToday      int64
		userWithdrawUsdtTotal           int64
		userRewardUsdtTotal             int64
		systemRewardUsdtTotal           int64
	)
	userCount, _ = uuc.repo.GetUserCount(ctx)
	userTodayCount, _ = uuc.repo.GetUserCountToday(ctx)
	userBalanceUsdtTotal, _ = uuc.ubRepo.GetUserBalanceUsdtTotal(ctx)
	userBalanceRecordUsdtTotal, _ = uuc.ubRepo.GetUserBalanceRecordUsdtTotal(ctx)
	userBalanceRecordUsdtTotalToday, _ = uuc.ubRepo.GetUserBalanceRecordUsdtTotalToday(ctx)
	userWithdrawUsdtTotalToday, _ = uuc.ubRepo.GetUserWithdrawUsdtTotalToday(ctx)
	userWithdrawUsdtTotal, _ = uuc.ubRepo.GetUserWithdrawUsdtTotal(ctx)
	userRewardUsdtTotal, _ = uuc.ubRepo.GetUserRewardUsdtTotal(ctx)
	systemRewardUsdtTotal, _ = uuc.ubRepo.GetSystemRewardUsdtTotal(ctx)

	return &v1.AdminAllReply{
		TodayTotalUser:        userTodayCount,
		TotalUser:             userCount,
		AllBalance:            fmt.Sprintf("%.2f", float64(userBalanceUsdtTotal)/float64(10000000000)),
		TodayLocation:         fmt.Sprintf("%.2f", float64(userBalanceRecordUsdtTotalToday)/float64(10000000000)),
		AllLocation:           fmt.Sprintf("%.2f", float64(userBalanceRecordUsdtTotal)/float64(10000000000)),
		TodayWithdraw:         fmt.Sprintf("%.2f", float64(userWithdrawUsdtTotalToday)/float64(10000000000)),
		AllWithdraw:           fmt.Sprintf("%.2f", float64(userWithdrawUsdtTotal)/float64(10000000000)),
		AllReward:             fmt.Sprintf("%.2f", float64(userRewardUsdtTotal)/float64(10000000000)),
		AllSystemRewardAndFee: fmt.Sprintf("%.2f", float64(systemRewardUsdtTotal)/float64(10000000000)),
	}, nil
}

func (uuc *UserUseCase) AdminWithdraw(ctx context.Context, req *v1.AdminWithdrawRequest) (*v1.AdminWithdrawReply, error) {
	var (
		currentValue                    int64
		systemAmount                    int64
		rewardLocations                 []*Location
		userRecommend                   *UserRecommend
		myLocationLast                  *Location
		myUserRecommendUserLocationLast *Location
		myUserRecommendUserId           int64
		myUserRecommendUserInfo         *UserInfo
		withdrawAmount                  int64
		stopLocations                   []*Location
		//lock                            bool
		withdrawNotDeal   []*Withdraw
		configs           []*Config
		recommendNeed     int64
		recommendNeedVip1 int64
		recommendNeedVip2 int64
		recommendNeedVip3 int64
		recommendNeedVip4 int64
		recommendNeedVip5 int64
		err               error
	)
	// 配置
	configs, _ = uuc.configRepo.GetConfigByKeys(ctx, "recommend_need", "recommend_need_vip1", "recommend_need_vip2",
		"recommend_need_vip3", "recommend_need_vip4", "recommend_need_vip5")
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
			}
		}
	}

	time.Sleep(30 * time.Second) // 错开时间和充值

	// todo 全局锁
	//for i := 0; i < 3; i++ {
	//	lock, _ = uuc.locationRepo.LockGlobalWithdraw(ctx)
	//	if !lock {
	//		time.Sleep(12 * time.Second)
	//		continue
	//	}
	//	break
	//}
	//if !lock {
	//	return &v1.AdminWithdrawReply{}, nil
	//}

	withdrawNotDeal, err = uuc.ubRepo.GetWithdrawNotDeal(ctx)
	if nil == withdrawNotDeal {
		//_, _ = uuc.locationRepo.UnLockGlobalWithdraw(ctx)
		return &v1.AdminWithdrawReply{}, nil
	}

	for _, withdraw := range withdrawNotDeal {
		if "" != withdraw.Status {
			continue
		}

		currentValue = withdraw.Amount

		if "dhb" == withdraw.Type { // 提现dhb
			//if err = uuc.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务
			//	_, err = uuc.ubRepo.UpdateWithdraw(ctx, withdraw.ID, "pass")
			//	if nil != err {
			//		return err
			//	}
			//
			//	return nil
			//}); nil != err {
			//
			//	return nil, err
			//}

			continue
		}

		// 先紧缩一次位置
		stopLocations, err = uuc.locationRepo.GetLocationsStopNotUpdate(ctx)
		if nil != stopLocations {
			// 调整位置紧缩
			for _, vStopLocations := range stopLocations {

				if err = uuc.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务
					err = uuc.locationRepo.UpdateLocationRowAndCol(ctx, vStopLocations.ID)
					if nil != err {
						return err
					}
					return nil
				}); nil != err {
					continue
				}
			}
		}

		// 获取当前用户的占位信息，已经有运行中的跳过
		myLocationLast, err = uuc.locationRepo.GetMyLocationLast(ctx, withdraw.UserId)
		if nil == myLocationLast { // 无占位信息
			return nil, err
		}
		// 占位分红人
		rewardLocations, err = uuc.locationRepo.GetRewardLocationByRowOrCol(ctx, myLocationLast.Row, myLocationLast.Col)

		// 推荐人
		userRecommend, err = uuc.urRepo.GetUserRecommendByUserId(ctx, withdraw.UserId)
		if nil != err {
			return nil, err
		}
		if "" != userRecommend.RecommendCode {
			tmpRecommendUserIds := strings.Split(userRecommend.RecommendCode, "D")
			if 2 <= len(tmpRecommendUserIds) {
				myUserRecommendUserId, _ = strconv.ParseInt(tmpRecommendUserIds[len(tmpRecommendUserIds)-1], 10, 64) // 最后一位是直推人
			}
		}
		myUserRecommendUserInfo, err = uuc.uiRepo.GetUserInfoByUserId(ctx, myUserRecommendUserId)

		if err = uuc.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务
			fmt.Println(withdraw.Amount)
			currentValue -= withdraw.Amount / 100 * 5 // 手续费

			// 手续费记录
			err = uuc.ubRepo.SystemFee(ctx, withdraw.Amount/100*5, myLocationLast.ID) // 推荐人奖励
			if nil != err {
				return err
			}

			currentValue = currentValue / 100 * 50 // 百分之50重新分配
			withdrawAmount = currentValue
			systemAmount = currentValue
			fmt.Println(withdrawAmount)
			// 占位分红人分红
			if nil != rewardLocations {
				for _, vRewardLocations := range rewardLocations {
					if "running" != vRewardLocations.Status {
						continue
					}
					if myLocationLast.Row == vRewardLocations.Row && myLocationLast.Col == vRewardLocations.Col { // 跳过自己
						continue
					}

					var locationType string
					var tmpAmount int64
					if myLocationLast.Row == vRewardLocations.Row { // 同行的人
						tmpAmount = currentValue / 100 * 5
						locationType = "row"
					} else if myLocationLast.Col == vRewardLocations.Col { // 同列的人
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
						vRewardLocations.Status = "stop"
						if "running" == tmpCurrentStatus {
							vRewardLocations.StopDate = time.Now().UTC().Add(8 * time.Hour)
						}
					}
					fmt.Println(vRewardLocations.StopDate)
					if 0 < tmpBalanceAmount {
						err = uuc.locationRepo.UpdateLocation(ctx, vRewardLocations.ID, vRewardLocations.Status, tmpBalanceAmount, vRewardLocations.StopDate) // 分红占位数据修改
						if nil != err {
							return err
						}
						systemAmount -= tmpBalanceAmount // 占位分红后剩余金额

						if 0 < tmpBalanceAmount && "running" == tmpCurrentStatus && tmpCurrent < vRewardLocations.CurrentMax { // 这次还能分红
							tmpCurrentAmount := vRewardLocations.CurrentMax - tmpCurrent // 最大可分红额度
							rewardAmount := tmpBalanceAmount
							if tmpCurrentAmount < tmpBalanceAmount { // 大于最大可分红额度
								rewardAmount = tmpCurrentAmount
							}

							_, err = uuc.ubRepo.WithdrawReward(ctx, vRewardLocations.UserId, rewardAmount, myLocationLast.ID, vRewardLocations.ID, locationType) // 分红信息修改
							if nil != err {
								return err
							}
						}
					}
				}
			}

			// 获取当前用户的占位信息，已经有运行中的跳过
			if nil != myUserRecommendUserInfo {
				// 有占位信息
				myUserRecommendUserLocationLast, err = uuc.locationRepo.GetMyLocationLast(ctx, myUserRecommendUserInfo.UserId)
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

					fmt.Println(myUserRecommendUserLocationLast.StopDate)
					if 0 < tmpBalanceAmount {
						err = uuc.locationRepo.UpdateLocation(ctx, myUserRecommendUserLocationLast.ID, myUserRecommendUserLocationLast.Status, tmpBalanceAmount, myUserRecommendUserLocationLast.StopDate) // 分红占位数据修改
						if nil != err {
							return err
						}
					}
					systemAmount -= tmpBalanceAmount // 扣除

					if 0 < tmpBalanceAmount && "running" == tmpStatus && tmpCurrent < myUserRecommendUserLocationLast.CurrentMax { // 这次还能分红
						tmpCurrentAmount := myUserRecommendUserLocationLast.CurrentMax - tmpCurrent // 最大可分红额度
						rewardAmount := tmpBalanceAmount
						if tmpCurrentAmount < tmpBalanceAmount { // 大于最大可分红额度
							rewardAmount = tmpCurrentAmount
						}
						_, err = uuc.ubRepo.NormalWithdrawRecommendReward(ctx, myUserRecommendUserId, rewardAmount, myLocationLast.ID) // 直推人奖励
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
							err = uuc.locationRepo.UpdateLocation(ctx, myUserRecommendUserLocationLast.ID, myUserRecommendUserLocationLast.Status, tmpBalanceAmount, myUserRecommendUserLocationLast.StopDate) // 分红占位数据修改
							if nil != err {
								return err
							}
						}
						systemAmount -= tmpBalanceAmount                                                                               // 扣除                                                                                    // 扣除
						if 0 < tmpBalanceAmount && "running" == tmpStatus && tmpCurrent < myUserRecommendUserLocationLast.CurrentMax { // 这次还能分红
							tmpCurrentAmount := myUserRecommendUserLocationLast.CurrentMax - tmpCurrent // 最大可分红额度
							rewardAmount := tmpBalanceAmount
							if tmpCurrentAmount < tmpBalanceAmount { // 大于最大可分红额度
								rewardAmount = tmpCurrentAmount
							}
							_, err = uuc.ubRepo.RecommendWithdrawReward(ctx, myUserRecommendUserId, rewardAmount, myLocationLast.ID) // 推荐人奖励
							if nil != err {
								return err
							}

						}
					}
				}
			}

			err = uuc.ubRepo.SystemWithdrawReward(ctx, systemAmount, myLocationLast.ID)
			if nil != err {
				return err
			}

			_, err = uuc.ubRepo.UpdateWithdrawAmount(ctx, withdraw.ID, "rewarded", withdrawAmount)
			if nil != err {
				return err
			}

			return nil
		}); nil != err {
			continue
		}

		// 调整位置紧缩
		stopLocations, err = uuc.locationRepo.GetLocationsStopNotUpdate(ctx)
		if nil != stopLocations {
			// 调整位置紧缩
			for _, vStopLocations := range stopLocations {

				if err = uuc.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务
					err = uuc.locationRepo.UpdateLocationRowAndCol(ctx, vStopLocations.ID)
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

	//_, _ = uuc.locationRepo.UnLockGlobalWithdraw(ctx)

	return &v1.AdminWithdrawReply{}, nil
}
