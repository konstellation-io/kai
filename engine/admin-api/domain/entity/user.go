package entity

type User struct {
	ID            string
	Roles         []string
	ProductGrants ProductGrants
}

type ProductGrants map[string][]string
