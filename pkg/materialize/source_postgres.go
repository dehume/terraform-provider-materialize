package materialize

import (
	"fmt"
	"strings"
)

type TablePostgres struct {
	Name  string
	Alias string
}

type SourcePostgresBuilder struct {
	Source
	clusterName        string
	size               string
	postgresConnection IdentifierSchemaStruct
	publication        string
	textColumns        []string
	table              []TablePostgres
}

func NewSourcePostgresBuilder(sourceName, schemaName, databaseName string) *SourcePostgresBuilder {
	return &SourcePostgresBuilder{
		Source: Source{sourceName, schemaName, databaseName},
	}
}

func (b *SourcePostgresBuilder) ClusterName(c string) *SourcePostgresBuilder {
	b.clusterName = c
	return b
}

func (b *SourcePostgresBuilder) Size(s string) *SourcePostgresBuilder {
	b.size = s
	return b
}

func (b *SourcePostgresBuilder) PostgresConnection(p IdentifierSchemaStruct) *SourcePostgresBuilder {
	b.postgresConnection = p
	return b
}

func (b *SourcePostgresBuilder) Publication(p string) *SourcePostgresBuilder {
	b.publication = p
	return b
}

func (b *SourcePostgresBuilder) TextColumns(t []string) *SourcePostgresBuilder {
	b.textColumns = t
	return b
}

func (b *SourcePostgresBuilder) Table(t []TablePostgres) *SourcePostgresBuilder {
	b.table = t
	return b
}

func (b *SourcePostgresBuilder) Create() string {
	q := strings.Builder{}
	q.WriteString(fmt.Sprintf(`CREATE SOURCE %s`, b.QualifiedName()))

	if b.clusterName != "" {
		q.WriteString(fmt.Sprintf(` IN CLUSTER %s`, QuoteIdentifier(b.clusterName)))
	}

	q.WriteString(fmt.Sprintf(` FROM POSTGRES CONNECTION %s`, b.postgresConnection.QualifiedName()))

	// Publication
	p := fmt.Sprintf(`PUBLICATION %s`, QuoteString(b.publication))

	if len(b.textColumns) > 0 {
		s := strings.Join(b.textColumns, ", ")
		p = p + fmt.Sprintf(`, TEXT COLUMNS (%s)`, s)
	}

	q.WriteString(fmt.Sprintf(` (%s)`, p))

	if len(b.table) > 0 {
		q.WriteString(` FOR TABLES (`)
		for i, t := range b.table {
			if t.Alias == "" {
				t.Alias = t.Name
			}
			q.WriteString(fmt.Sprintf(`%s AS %s`, t.Name, t.Alias))
			if i < len(b.table)-1 {
				q.WriteString(`, `)
			}
		}
		q.WriteString(`)`)
	} else {
		q.WriteString(` FOR ALL TABLES`)
	}

	if b.size != "" {
		q.WriteString(fmt.Sprintf(` WITH (SIZE = %s)`, QuoteString(b.size)))
	}

	q.WriteString(`;`)
	return q.String()
}

func (b *SourcePostgresBuilder) Rename(newName string) string {
	n := QualifiedName(b.DatabaseName, b.SchemaName, newName)
	return fmt.Sprintf(`ALTER SOURCE %s RENAME TO %s;`, b.QualifiedName(), n)
}

func (b *SourcePostgresBuilder) UpdateSize(newSize string) string {
	return fmt.Sprintf(`ALTER SOURCE %s SET (SIZE = %s);`, b.QualifiedName(), QuoteString(newSize))
}

func (b *SourcePostgresBuilder) Drop() string {
	return fmt.Sprintf(`DROP SOURCE %s;`, b.QualifiedName())
}

func (b *SourcePostgresBuilder) ReadId() string {
	return ReadSourceId(b.SourceName, b.SchemaName, b.DatabaseName)
}