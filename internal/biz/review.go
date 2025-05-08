package biz

import (
	"context"
	v1 "review-service/api/review/v1"
	"review-service/internal/data/model"
	"review-service/pkg/snowflake"

	"errors"

	"github.com/go-kratos/kratos/v2/log"
)

type ReviewRepo interface {
	SaveReview(context.Context, *model.ReviewInfo) (*model.ReviewInfo, error)
	GetReviewByOrderID(context.Context, int64) ([]*model.ReviewInfo, error)
	SaveReply(context.Context, *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error)
	GetReviewByReviewID(context.Context, int64) (*model.ReviewInfo, error)
	DeleteReview(context.Context, *model.ReviewInfo) (*model.ReviewInfo, error)
	GetReview(context.Context, int64) (*model.ReviewInfo, error)
	ListReview(context.Context, *ListParam) ([]*model.ReviewInfo, int64, error)
	GetStoreID(context.Context, int64) (*model.ReviewInfo, error)
	SaveAppeal(context.Context, *model.ReviewAppealInfo) (*model.ReviewAppealInfo, error)
	GetAppeal(context.Context, int64) ([]*model.ReviewAppealInfo, error)
	UpdateAppeal(context.Context, *model.ReviewAppealInfo) (*model.ReviewAppealInfo, error)
}

type ReviewerUsecase struct {
	repo ReviewRepo
	log  *log.Helper
}

func NewReviewerUsecase(repo ReviewRepo, logger log.Logger) *ReviewerUsecase {
	return &ReviewerUsecase{repo: repo, log: log.NewHelper(logger)}
}

// 创建评价逻辑处理
func (uc *ReviewerUsecase) CreateReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	// 查询是否已存在该评价
	reviewInfo, err := uc.repo.GetReviewByOrderID(ctx, review.OrderID)
	if err != nil {
		return nil, v1.ErrorDbError("不存在该订单")
	}
	if len(reviewInfo) > 0 {
		return nil, v1.ErrorOrderReviewed("订单%d已评价", review.OrderID)
	}
	// 生成reviewID
	review.ReviewID = snowflake.GenID()
	//查询订单服务和商家服务

	// 入库
	return uc.repo.SaveReview(ctx, review)
}

// 删除评价逻辑处理
func (uc *ReviewerUsecase) DeleteReview(ctx context.Context, param *DeleteParam) (*model.ReviewInfo, error) {
	// 查询是否存在该评价
	reviewInfo, err := uc.repo.GetReviewByReviewID(ctx, param.ReviewID)
	if err != nil {
		return nil, v1.ErrorDbError("不存在该评价")
	}
	// 判断该评价是否是该用户评价,防止越权
	if reviewInfo.UserID != param.UserID {
		return nil, errors.New("水平越权")
	}
	return uc.repo.DeleteReview(ctx, reviewInfo)
}

// 获取评价逻辑处理
func (uc *ReviewerUsecase) GetReview(ctx context.Context, reviewID int64) (*v1.GetReviewReply, error) {
	res, err := uc.repo.GetReview(ctx, reviewID)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, v1.ErrorDbError("评价不存在")
	}

	reply := &v1.GetReviewReply{
		ReviewID:     res.ReviewID,
		UserID:       res.UserID,
		OrderID:      res.OrderID,
		Score:        int64(res.Score),
		ServiceScore: int64(res.ServiceScore),
		ExpressScore: int64(res.ExpressScore),
		Content:      res.Content,
		PicInfo:      res.PicInfo,
		VideoInfo:    res.VideoInfo,
		Annoymous:    res.Anonymous == 1,
		Status:       int64(res.Status),
		CreateTime:   res.CreateAt.Format("2006-01-02 15:04:05"),
		UpdateTime:   res.UpdateAt.Format("2006-01-02 15:04:05"),
	}
	return reply, nil
}

// 获取评价列表逻辑处理
func (uc *ReviewerUsecase) ListReview(ctx context.Context, param *ListParam) (*v1.ListReviewReply, error) {
	// 参数校验,默认第一条开始,一页五条
	if param.Page <= 0 {
		param.Page = 1
	}
	if param.PageSize <= 0 {
		param.PageSize = 5
	}

	// 查询评价列表
	reviews, total, err := uc.repo.ListReview(ctx, param)
	if err != nil {
		return nil, err
	}

	// 转换响应数据
	items := make([]*v1.GetReviewReply, 0, len(reviews))
	for _, review := range reviews {
		items = append(items, &v1.GetReviewReply{
			ReviewID:     review.ReviewID,
			UserID:       review.UserID,
			OrderID:      review.OrderID,
			Score:        int64(review.Score),
			ServiceScore: int64(review.ServiceScore),
			ExpressScore: int64(review.ExpressScore),
			Content:      review.Content,
			PicInfo:      review.PicInfo,
			VideoInfo:    review.VideoInfo,
			Annoymous:    review.Anonymous == 1,
			Status:       int64(review.Status),
			CreateTime:   review.CreateAt.Format("2006-01-02 15:04:05"),
			UpdateTime:   review.UpdateAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &v1.ListReviewReply{
		Reviews: items,
		Total:   total,
	}, nil
}

// 创建评价回复逻辑处理
func (uc *ReviewerUsecase) CreateReply(ctx context.Context, param *ReplyParam) (*model.ReviewReplyInfo, error) {
	reply := &model.ReviewReplyInfo{
		ReplyID:   snowflake.GenID(),
		ReviewID:  param.ReviewID,
		StoreID:   param.StoreID,
		Content:   param.Content,
		PicInfo:   param.PicInfo,
		VideoInfo: param.VideoInfo,
	}
	return uc.repo.SaveReply(ctx, reply)
}

// 创建评价申诉逻辑处理
func (uc *ReviewerUsecase) AppealReview(ctx context.Context, param *AppealParam) (*model.ReviewAppealInfo, error) {
	// 判断评价是否属于此商家且是否唯一
	res, err := uc.repo.GetStoreID(ctx, param.ReviewID)
	if err != nil {
		return nil, err
	}
	if res.StoreID != param.StoreID {
		return nil, errors.New("水平越权")
	}
	// 判断是否已经申诉
	appeals, err := uc.repo.GetAppeal(ctx, param.ReviewID)
	if err != nil {
		return nil, err
	}
	if len(appeals) >= 1 {
		return nil, errors.New("申诉已存在")
	}
	// 调用data层记录申诉
	return uc.repo.SaveAppeal(ctx, &model.ReviewAppealInfo{
		AppealID:  snowflake.GenID(),
		ReviewID:  param.ReviewID,
		StoreID:   param.StoreID,
		Reason:    param.Reason,
		Content:   param.Content,
		PicInfo:   param.PicInfo,
		VideoInfo: param.VideoInfo,
	})
}

// 申诉审核逻辑处理
func (uc *ReviewerUsecase) AuditReview(ctx context.Context, param *AuditParam) (*model.ReviewAppealInfo, error) {
	// 查询申诉信息是否存在
	var appeals []*model.ReviewAppealInfo
	appeals, err := uc.repo.GetAppeal(ctx, param.ReviewID)
	if err != nil {
		return nil, err
	}
	if len(appeals) > 1 || len(appeals) == 0 {
		return nil, errors.New("申诉不存在或者超过数量")
	}
	appeal := appeals[0]

	// 更新申诉状态
	appeal.Status = param.Status
	appeal.OpRemarks = param.OpRemarks
	appeal.OpUser = param.OpUser
	appeal.ExtJSON = param.ExtJSON
	appeal.CtrlJSON = param.CtrlJSON
	// 保存申诉审核结果
	return uc.repo.UpdateAppeal(ctx, appeal)
}
