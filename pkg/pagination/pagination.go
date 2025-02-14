package pagination

func New(page *int32, pageSize *int32) (limit int, offset int) {
	if page == nil || *page <= 0 {
		defaultPage := int32(1)
		page = &defaultPage
	}

	if pageSize == nil || *pageSize <= 0 {
		defaultPageSize := int32(10)
		pageSize = &defaultPageSize
	}

	offset = (int(*page) - 1) * int(*pageSize)
	limit = int(*pageSize)

	return
}
