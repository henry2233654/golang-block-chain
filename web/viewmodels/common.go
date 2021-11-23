package viewmodels

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
}

type Pagination struct {
	Total       int64       `json:"total"`
	LastPage    int         `json:"last_page"`
	PerPage     int         `json:"per_page"`
	CurrentPage int         `json:"current_page"`
	PrevPage    int         `json:"prev_page"`
	NextPage    int         `json:"next_page"`
	Data        interface{} `json:"data"`
}

type Paging struct {
	Page int         `json:"page"`
	Data interface{} `json:"data"`
}
