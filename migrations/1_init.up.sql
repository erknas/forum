CREATE TABLE IF NOT EXISTS post_author (
	id SERIAL PRIMARY KEY,
	author_name VARCHAR(16) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS comment_author (
	id SERIAL PRIMARY KEY,
	author_name VARCHAR(16) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS post (
	id SERIAL PRIMARY KEY,
	title VARCHAR(100) NOT NULL,
	content TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	comments_allowed BOOLEAN,
	author_id INT NOT NULL REFERENCES post_author(id) 
);

CREATE TABLE IF NOT EXISTS comment (
	id SERIAL PRIMARY KEY,
	content VARCHAR(2000) NOT NULL CHECK (LENGTH(content) <= 2000),
	created_at TIMESTAMPTZ DEFAULT NOW(),
	post_id INT NOT NULL REFERENCES post(id),
	parent_id INT,
	author_id INT NOT NULL REFERENCES comment_author(id)
);
