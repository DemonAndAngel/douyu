package db

type PageInfo struct {
	Page      int64 `json:"page"`
	PageSize  int64 `json:"pageSize"`
	TotalPage int64 `json:"totalPage"`
	Total     int64 `json:"total"`
}

func BuildPage(total, page, pageSize int64) (pageInfo PageInfo, offset int64) {
	var pageMax int64

	if pageSize == 0 {
		pageSize = 15
	}

	if page == 0 {
		page = 1
	}

	if total%pageSize == 0 {
		pageMax = total / pageSize
	} else {
		pageMax = total/pageSize + 1
	}

	offset = pageSize * (page - 1)

	pageInfo.Page = page
	pageInfo.PageSize = pageSize
	pageInfo.Total = total
	pageInfo.TotalPage = pageMax

	return
}
