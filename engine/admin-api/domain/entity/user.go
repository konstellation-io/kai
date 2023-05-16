package entity

const DefaultAdminRole = "ADMIN"

type User struct {
	ID            string
	Roles         []string
	ProductGrants ProductGrants
}

func (u User) IsAdmin(optAdminRole ...string) bool {
	adminRole := DefaultAdminRole
	if len(optAdminRole) > 0 {
		adminRole = optAdminRole[0]
	}

	for _, role := range u.Roles {
		if role == adminRole {
			return true
		}
	}
	return false
}

type ProductGrants map[string][]string
