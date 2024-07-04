package consts

type OrderBy string
type Sort string

const (
	ASC  OrderBy = "asc"
	DESC OrderBy = "desc"
)
const (
	CreatedAt Sort = "created_at"
	UpdatedAt Sort = "updated_at"
)
