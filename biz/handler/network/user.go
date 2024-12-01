// Code generated by hertz generator.

package network

import (
	"context"

	domain "go-social-network/biz/domain"
	"go-social-network/biz/logic"
	base "go-social-network/biz/model/base"
	network "go-social-network/biz/model/network"
	"go-social-network/data"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/jinzhu/copier"
)

// Register .
// @router /api/v1/register [POST]
func Register(ctx context.Context, c *app.RequestContext) {
	var err error
	var req network.UserRegisterReq
	resp := new(base.BaseResp)
	err = c.BindAndValidate(&req)
	if err != nil {
		resp.ErrCode = base.ErrCode_ArgumentError
		resp.ErrMsg = err.Error()
		c.JSON(consts.StatusBadRequest, resp)
		return
	}
	var userRegisterReq domain.UserRegisterReq
	_ = copier.Copy(&userRegisterReq, &req)
	err = logic.NewUser(data.Default()).Register(ctx, userRegisterReq)
	if err != nil {
		resp.ErrCode = base.ErrCode_CreateUserError
		resp.ErrMsg = err.Error()
		c.JSON(consts.StatusInternalServerError, resp)
		return
	}
	c.JSON(consts.StatusOK, resp)
}
