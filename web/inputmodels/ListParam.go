package inputmodels

type ListParam struct {
	Page    int    `json:"page" form:"page,default=1" example:"1"`
	PerPage int    `json:"per_page" form:"per_page,default=20" example:"20"`
	Order   string `json:"order" form:"order,default=id" example:"id"`
	Desc    bool   `json:"desc" form:"desc,default=false" example:"false"`
}
