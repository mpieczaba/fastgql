package sql

import (
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"

	"github.com/roneli/fastgql/gql"
	"github.com/spf13/cast"
	"github.com/vektah/gqlparser/v2/ast"
)

type RelationType string

const (
	OneToMany  RelationType = "ONE_TO_MANY"
	OneToOne   RelationType = "ONE_TO_ONE"
	ManyToMany RelationType = "MANY_TO_MANY"
)

type relation struct {
	relType              RelationType
	baseTable            string
	referenceTable       string
	fields               []string
	references           []string
	manyToManyTable      string
	manyToManyReferences []string
	manyToManyFields     []string
}

type tableDefinition struct {
	name   string
	schema string
}

func (t tableDefinition) TableExpression() exp.IdentifierExpression {
	tbl := goqu.T(t.name)
	if t.schema != "" {
		return tbl.Schema(t.schema)
	}
	return tbl
}

/*
parseRelationDirective parses the sqlRelation directive to connect graphQL Objects with SQL relations, this directive
is also important for creating relational filters.
directive @sqlRelation(relationType: _relationType!, baseTable: String!, refTable: String!, fields: [String!]!,
    references: [String!]!, manyToManyTable: String = "", manyToManyFields: [String] = [], manyToManyReferences: [String] = []) on FIELD_DEFINITION
*/
func parseRelationDirective(d *ast.Directive) relation {
	relType := d.Arguments.ForName("relationType").Value.Raw
	return relation{
		relType:              RelationType(relType),
		fields:               cast.ToStringSlice(gql.GetDirectiveValue(d, "fields")),
		references:           cast.ToStringSlice(gql.GetDirectiveValue(d, "references")),
		baseTable:            cast.ToString(gql.GetDirectiveValue(d, "baseTable")),
		referenceTable:       cast.ToString(gql.GetDirectiveValue(d, "refTable")),
		manyToManyTable:      cast.ToString(gql.GetDirectiveValue(d, "manyToManyTable")),
		manyToManyFields:     cast.ToStringSlice(gql.GetDirectiveValue(d, "manyToManyFields")),
		manyToManyReferences: cast.ToStringSlice(gql.GetDirectiveValue(d, "manyToManyReferences")),
	}
}

// getTableName returns the field's type tableDefinition name in the database, if no directive is defined, type name is presumed
// as the tableDefinition's name
func getTableName(schema *ast.Schema, f *ast.FieldDefinition) tableDefinition {
	objType, ok := schema.Types[f.Type.Name()]
	if !ok {
		return tableDefinition{
			name:   f.Name,
			schema: "",
		}
	}
	d := objType.Directives.ForName("tableName")
	if d == nil {
		return tableDefinition{
			name:   strings.ToLower(objType.Name),
			schema: "",
		}
	}
	name := d.Arguments.ForName("name").Value.Raw
	schemaValue := d.Arguments.ForName("schema")
	if schemaValue == nil {
		return tableDefinition{
			name:   name,
			schema: "",
		}
	}
	return tableDefinition{
		name:   name,
		schema: schemaValue.Value.Raw,
	}
}

// getTableName returns the field's type tableDefinition name in the database, if no directive is defined, type name is presumed
// as the tableDefinition's name
func getAggregateTableName(schema *ast.Schema, field *ast.Field) tableDefinition {
	fieldName := strings.Split(field.Name, "Aggregate")[0][1:]
	nonAggField := field.ObjectDefinition.Fields.ForName(fieldName)
	return getTableName(schema, nonAggField)
}
