package repository

const (
	GetArticlesQuery = `
		SELECT 
			id
			, author
			, title
			, body
			, created_at as created
		FROM 
			articles
		WHERE 
			(
				title LIKE ?
				OR body LIKE ?
			)
			AND ((? = '') IS NOT FALSE OR author = ?)
	`

	CreateArticleQuery = `
		INSERT INTO articles (
			author, 
			title, 
			body, 
			created_at
		)
		VALUES(
			?, 
			?, 
			?, 
			?
		)
	`
)
