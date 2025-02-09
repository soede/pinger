package utils

const (
	defaultSize = 10
)

type PaginationQuery struct {
	Size    int `json:"size,omitempty"`
	Page    int `json:"page,omitempty"`
	OrderBy int `json:"order_by,omitempty"`
}

/*func GetPagination() (*PaginationQuery, error) {
	q := &PaginationQuery{}
	if err := q.SetPage(c.QueryParam("page")); err != nil {
		return nil, err
	}
	if err := q.SetSize(c.QueryParam("size")); err != nil {
		return nil, err
	}
	q.SetOrderBy(c.QueryParam("orderBy"))

	return q, nil
}
*/
