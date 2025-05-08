package data

import (
	"context"
	"errors"

	"review-service/internal/biz"
	"review-service/internal/data/model"
	"review-service/internal/data/query"

	"github.com/go-kratos/kratos/v2/log"
)

type reviewRepo struct {
	data *Data
	log  *log.Helper
}

func NewReviewRepo(data *Data, logger log.Logger) biz.ReviewRepo {
	return &reviewRepo{
		data: data,
		log:  log.NewHelper(logger),
	}

}

func (r *reviewRepo) SaveReview(ctx context.Context, g *model.ReviewInfo) (*model.ReviewInfo, error) {
	err := r.data.query.ReviewInfo.WithContext(ctx).Save(g)
	return g, err
}

func (r *reviewRepo) GetReviewByOrderID(ctx context.Context, orderID int64) ([]*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.OrderID.Eq(orderID)).Find()
}

func (r *reviewRepo) GetReviewByReviewID(ctx context.Context, reviewID int64) (*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.ReviewID.Eq(reviewID)).First()
}

func (r *reviewRepo) DeleteReview(ctx context.Context, reply *model.ReviewInfo) (*model.ReviewInfo, error) {
	// 更新评论表和回复表
	err := r.data.query.Transaction(func(tx *query.Query) error {
		if _, err := tx.ReviewInfo.WithContext(ctx).Delete(reply); err != nil {
			return err
		}
		if _, err := tx.ReviewReplyInfo.WithContext(ctx).Where(tx.ReviewReplyInfo.ReviewID.Eq(reply.ReviewID)).Delete(); err != nil {
			return err
		}
		return nil
	})
	return reply, err
}

func (r *reviewRepo) GetReview(ctx context.Context, reviewID int64) (*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.ReviewID.Eq(reviewID)).First()
}

func (r *reviewRepo) ListReview(ctx context.Context, param *biz.ListParam) ([]*model.ReviewInfo, int64, error) {
	var reviews []*model.ReviewInfo
	// 获取总数
	total, err := r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.UserID.Eq(param.UserID)).Count()
	if err != nil {
		return nil, 0, err
	}
	// 分页查询
	offset := (param.Page - 1) * param.PageSize
	if offset > int32(total) {
		return nil, 0, errors.New("查询开始数超过记录数")
	}
	reviews, count, err := r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.UserID.Eq(param.UserID)).FindByPage(int(offset), int(param.Page))
	if err != nil {
		return nil, 0, err
	}
	return reviews, count, nil
}

func (r *reviewRepo) SaveReply(ctx context.Context, reply *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error) {
	// 数据合法性
	review, err := r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.ReviewID.Eq(reply.ReviewID)).First()
	if err != nil {
		return nil, err
	}

	if review.HasReply == 1 {
		return nil, errors.New("改评价已回复")
	}
	// 水平越权
	if review.StoreID != reply.StoreID {
		return nil, errors.New("不能回复非自己店铺的评价")
	}
	// 更新回复表和评论表
	err = r.data.query.Transaction(func(tx *query.Query) error {
		if err := tx.ReviewReplyInfo.WithContext(ctx).Save(reply); err != nil {
			return err
		}
		if _, err := tx.ReviewInfo.WithContext(ctx).Where(tx.ReviewInfo.ReviewID.Eq(reply.ReviewID)).Update(tx.ReviewInfo.HasReply, 1); err != nil {
			return err
		}
		return nil
	})
	return reply, err
}

func (r *reviewRepo) GetStoreID(ctx context.Context, reviewID int64) (*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.ReviewID.Eq(reviewID)).First()
}

func (r *reviewRepo) SaveAppeal(ctx context.Context, appeal *model.ReviewAppealInfo) (*model.ReviewAppealInfo, error) {
	err := r.data.query.ReviewAppealInfo.WithContext(ctx).Save(appeal)
	return appeal, err
}

func (r *reviewRepo) GetAppeal(ctx context.Context, reviewID int64) ([]*model.ReviewAppealInfo, error) {
	return r.data.query.ReviewAppealInfo.WithContext(ctx).Where(r.data.query.ReviewAppealInfo.ReviewID.Eq(reviewID)).Find()
}

func (r *reviewRepo) UpdateAppeal(ctx context.Context, appeal *model.ReviewAppealInfo) (*model.ReviewAppealInfo, error) {
	err := r.data.query.Transaction(func(tx *query.Query) error {
		// 更新申诉表状态
		if _, err := r.data.query.ReviewAppealInfo.WithContext(ctx).Where(r.data.query.ReviewAppealInfo.ReviewID.Eq(appeal.ReviewID)).Updates(appeal); err != nil {
			return err
		}
		if _, err := r.data.query.ReviewInfo.WithContext(ctx).Where(r.data.query.ReviewInfo.ReviewID.Eq(appeal.ReviewID)).Update(r.data.query.ReviewInfo.Status, appeal.Status); err != nil {
			return err
		}
		return nil
	})
	return appeal, err
}
