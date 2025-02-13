package pagination

func New(page *int32, pageSize *int32) (limit int, offset int) {
	if page == nil || *page <= 0 {
		*page = 1
	}

	if pageSize == nil || *pageSize <= 0 {
		*pageSize = 10
	}

	offset = (int(*page) - 1) * int(*pageSize)
	limit = int(*pageSize)

	return
}
