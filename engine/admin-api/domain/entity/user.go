package entity

type User struct {
	ID            string
	Roles         []string
	ProductGrants ProductGrants
}

func (u User) IsAdmin() bool {
	for _, role := range u.Roles {
		if role == "ADMIN" {
			return true
		}
	}
	return false
}

type ProductGrants map[string][]string
