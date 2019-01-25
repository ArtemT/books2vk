package books

type Book struct {
	Author		string	`xcol:"1"`
	Title		string	`xcol:"2"`
	Description	string	`xcol:"3"`
	Price		int		`xcol:"11"`
}
