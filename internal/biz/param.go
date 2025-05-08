package biz

type ReplyParam struct {
	ReviewID  int64
	StoreID   int64
	Content   string
	PicInfo   string
	VideoInfo string
}

type DeleteParam struct {
	ReviewID int64
	UserID   int64
}

type ListParam struct {
	UserID   int64
	Page     int32
	PageSize int32
}

type AppealParam struct {
	ReviewID  int64
	StoreID   int64
	Reason    string
	Content   string
	PicInfo   string
	VideoInfo string
}

type AuditParam struct {
	ReviewID  int64
	Status    int32
	OpRemarks string
	OpUser    string
	ExtJSON   string
	CtrlJSON  string
}
