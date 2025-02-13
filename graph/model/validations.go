package model

func (p PostInput) ValidatePostInput() map[string]string {
	errors := make(map[string]string)

	if len(p.Title) == 0 {
		errors["title"] = "title length cannot be zero"
	}

	if len(p.Title) > 100 {
		errors["title"] = "title length cannot be more than 100 symbols"
	}

	if len(p.Author) == 0 {
		errors["author"] = "author name length cannot be zero"
	}

	if len(p.Author) > 100 {
		errors["author"] = "author name length cannot be more than 16 symbols"
	}

	if len(p.Content) == 0 {
		errors["content"] = "content length cannot be zero"
	}

	return errors
}

func (c CommentInput) ValidateCommentInput() map[string]string {
	errors := make(map[string]string)

	if len(c.Author) == 0 {
		errors["author"] = "author name length cannot be zero"
	}

	if len(c.Author) > 16 {
		errors["author"] = "author name length cannot be more than 16 symbols"
	}

	if len(c.Content) == 0 {
		errors["content"] = "content length cannot be zero"
	}

	if len(c.Content) > 2000 {
		errors["content"] = "content length cannot be more than 2000 symbols"
	}

	return errors
}
