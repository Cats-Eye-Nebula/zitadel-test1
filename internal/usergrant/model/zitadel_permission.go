package model

type Permissions struct {
	Permissions []string
}

func (p *Permissions) appendPermissions(ctxID string, permissions ...string) {
	for _, permission := range permissions {
		p.appendPermission(ctxID, permission)
	}
}

func (p *Permissions) appendPermission(ctxID, permission string) {
	if ctxID != "" {
		permission = permission + ":" + ctxID
	}
	for _, existingPermission := range p.Permissions {
		if existingPermission == permission {
			return
		}
	}
	p.Permissions = append(p.Permissions, permission)
}
