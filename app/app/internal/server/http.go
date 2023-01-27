package server

import (
	"context"
	v1 "dhb/app/app/api"
	"dhb/app/app/internal/conf"
	"dhb/app/app/internal/service"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/handlers"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, app *service.AppService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			selector.Server( // jwt 验证
				jwt.Server(func(token *jwt2.Token) (interface{}, error) {
					return []byte("ebb0a1577eaebbad70115c472602da54"), nil
				}, jwt.WithSigningMethod(jwt2.SigningMethodHS256)),
			).Match(NewWhiteListMatcher()).Build(),
		),
		http.Filter(handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}),
		)),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterAppHTTPServer(srv, app)
	return srv
}

// NewWhiteListMatcher 设置白名单，不需要 token 验证的接口
func NewWhiteListMatcher() selector.MatchFunc {
	whiteList := make(map[string]struct{})
	whiteList["/api.App/EthAuthorize"] = struct{}{}
	//whiteList["/api.App/Deposit"] = struct{}{}
	//whiteList["/api.App/AdminLocationList"] = struct{}{}
	//whiteList["/api.App/AdminRewardList"] = struct{}{}
	//whiteList["/api.App/AdminUserList"] = struct{}{}
	//whiteList["/api.App/AdminWithdrawList"] = struct{}{}
	//whiteList["/api.App/AdminWithdraw"] = struct{}{}
	//whiteList["/api.App/AdminWithdrawEth"] = struct{}{}
	//whiteList["/api.App/AdminFee"] = struct{}{}
	//whiteList["/api.App/AdminAll"] = struct{}{}
	//whiteList["/api.App/AdminConfigUpdate"] = struct{}{}
	//whiteList["/api.App/AdminConfig"] = struct{}{}
	//whiteList["/api.App/AdminUserRecommend"] = struct{}{}
	//whiteList["/api.App/AdminMonthRecommend"] = struct{}{}
	return func(ctx context.Context, operation string) bool {
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}
}
