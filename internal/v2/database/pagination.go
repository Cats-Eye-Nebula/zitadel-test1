package database

type Pagination struct {
	Limit  uint32
	Offset uint32
}

func (p *Pagination) Write(stmt *Statement) {
	if p.Limit > 0 {
		stmt.Builder.WriteString(" LIMIT ")
		stmt.WriteArg(p.Limit)
	}
	if p.Offset > 0 {
		stmt.Builder.WriteString(" OFFSET ")
		stmt.WriteArg(p.Offset)
	}
}
