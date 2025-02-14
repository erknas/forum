package model

func (p PostInput) ValidatePostInput() map[string]interface{} {
	errors := make(map[string]interface{})

	if len(p.Title) == 0 {
		errors["title"] = "title length cannot be zero"
	}

	if len(p.Title) > 100 {
		errors["title"] = "title length cannot be more than 100 symbols"
	}

	if len(p.Author) == 0 {
		errors["author"] = "author name length cannot be zero"
	}

	if len(p.Content) == 0 {
		errors["content"] = "content length cannot be zero"
	}

	return errors
}

func (c CommentInput) ValidateCommentInput() map[string]interface{} {
	errors := make(map[string]interface{})

	if len(c.PostID) == 0 {
		errors["postID"] = "postID cannot be empty"
	}

	if len(c.Author) == 0 {
		errors["author"] = "author name length cannot be zero"
	}

	if len(c.Content) == 0 {
		errors["content"] = "content length cannot be zero"
	}

	if len(c.Content) > 2000 {
		errors["content"] = "content length cannot be more than 2000 symbols"
	}

	return errors
}
