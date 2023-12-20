package entity

type User struct {
	ID            string
	Name          string
	Email         string
	Roles         []string
	ProductGrants ProductGrants
}

type ProductGrants map[string][]string
