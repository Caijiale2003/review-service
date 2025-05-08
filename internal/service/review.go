package service

import (
	"context"

	pb "review-service/api/review/v1"
	"review-service/internal/biz"
	"review-service/internal/data/model"
)

type ReviewService struct {
	pb.UnimplementedReviewServer
	uc *biz.ReviewerUsecase
}

func NewReviewService(uc *biz.ReviewerUsecase) *ReviewService {
	return &ReviewService{uc: uc}
}

// CreateReview 创建评价服务
func (s *ReviewService) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.CreateReviewReply, error) {
	// 基础参数校验
	var annoymous int32 = 0
	if req.Annoymous {
		annoymous = 1
	}
	// 格式转换
	review := &model.ReviewInfo{
		UserID:       req.UserID,
		OrderID:      req.OrderID,
		Score:        int32(req.Score),
		ServiceScore: int32(req.ServiceScore),
		ExpressScore: int32(req.ExpressScore),
		Content:      req.Content,
		PicInfo:      req.PicInfo,
		VideoInfo:    req.VideoInfo,
		Anonymous:    annoymous,
	}
	// 调用biz层
	res, err := s.uc.CreateReview(ctx, review)
	if err != nil {
		return nil, err
	}
	return &pb.CreateReviewReply{
		ReviewID: res.ReviewID,
	}, nil
}

// DeleteReview 删除评价服务
func (s *ReviewService) DeleteReview(ctx context.Context, req *pb.DeleteReviewRequest) (*pb.DeleteReviewReply, error) {
	// 调用biz层
	res, err := s.uc.DeleteReview(ctx, &biz.DeleteParam{
		ReviewID: req.ReviewID,
		UserID:   req.UserID,
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteReviewReply{ReviewID: res.ReviewID}, nil
}

// GetReview 获取评价服务
func (s *ReviewService) GetReview(ctx context.Context, req *pb.GetReviewRequest) (*pb.GetReviewReply, error) {
	res, err := s.uc.GetReview(ctx, req.ReviewID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// ListReviewByUserID 列出评价列表服务
func (s *ReviewService) ListReviewByUserID(ctx context.Context, req *pb.ListReviewRequest) (*pb.ListReviewReply, error) {
	res, err := s.uc.ListReview(ctx, &biz.ListParam{
		UserID:   req.UserID,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	return &pb.ListReviewReply{Reviews: res.Reviews, Total: res.Total}, nil
}

// ReplyReview 回复评价服务
func (s *ReviewService) ReplyReview(ctx context.Context, req *pb.ReplyReviewRequest) (*pb.ReplyReviewReply, error) {
	res, err := s.uc.CreateReply(ctx, &biz.ReplyParam{
		ReviewID:  req.GetReviewID(),
		StoreID:   req.GetStoreID(),
		Content:   req.GetContent(),
		PicInfo:   req.GetPicInfo(),
		VideoInfo: req.GetVideoInfo(),
	})
	if err != nil {
		return nil, err
	}
	return &pb.ReplyReviewReply{ReplyID: res.ReviewID}, nil
}

// AppealReview 申诉评价服务
func (s *ReviewService) AppealReview(ctx context.Context, req *pb.AppealReviewRequest) (*pb.AppealReviewReply, error) {
	res, err := s.uc.AppealReview(ctx, &biz.AppealParam{
		ReviewID:  req.ReviewID,
		StoreID:   req.StoreID,
		Reason:    req.Reason,
		Content:   req.Content,
		PicInfo:   req.PicInfo,
		VideoInfo: req.VideoInfo,
	})
	if err != nil {
		return nil, err
	}
	return &pb.AppealReviewReply{ReviewID: res.ReviewID}, nil
}

// AuditReview 审核申诉服务
func (s *ReviewService) AuditReview(ctx context.Context, req *pb.AuditReviewRequest) (*pb.AuditReviewReply, error) {
	res, err := s.uc.AuditReview(ctx, &biz.AuditParam{
		ReviewID:  req.ReviewID,
		Status:    req.Status,
		OpRemarks: req.OpRemarks,
		OpUser:    req.OpUser,
		ExtJSON:   req.ExtJSON,
		CtrlJSON:  req.CtrlJSON,
	})
	if err != nil {
		return nil, err
	}
	return &pb.AuditReviewReply{ReviewID: res.ReviewID}, nil
}
