package schema

// TODO: in go 1.16 use //go:embed

const FastgqlSchema = `# Used if Object/Interface type name is different then the actual table name
directive @tableName(name: String!, schema: String) on OBJECT | INTERFACE

# Generate filter input on an object
directive @generateFilterInput(name: String!, description: String) on OBJECT | INTERFACE

# Generate arguments for a given field or all object fields
directive @generate(filter: Boolean = True, pagination: Boolean = True, ordering: Boolean = True, aggregate: Boolean = True, recursive: Boolean = True) on OBJECT

directive @generateMutations(create: Boolean = True) on OBJECT

enum _relationType {
    ONE_TO_ONE
    ONE_TO_MANY
    MANY_TO_MANY
}

directive @sqlRelation(relationType: _relationType!, baseTable: String!, refTable: String!, fields: [String!]!,
    references: [String!]!, manyToManyTable: String = "", manyToManyFields: [String] = [], manyToManyReferences: [String] = []) on FIELD_DEFINITION

directive @skipGenerate(resolver: Boolean = True) on FIELD_DEFINITION

enum _OrderingTypes {
    ASC
    DESC
    ASC_NULL_FIRST
    DESC_NULL_FIRST
}

type _AggregateResult {
    count: Int!
}

input StringComparator {
    eq: String
    neq: String
    contains: [String]
    not_contains: [String]
    like: String
    ilike: String
    suffix: String
    prefix: String
}

input StringListComparator {
    eq: [String]
    neq: [String]
    contains: [String]
    containedBy: [String]
    overlap: [String]
}

input IntComparator {
    eq: Int
    neq: Int
    gt: Int
    gte: Int
    lt: Int
    lte: Int
}

input IntListComparator {
    eq: [Int]
    neq: [Int]
    contains: [Int]
    contained: [Int]
    overlap: [Int]
}

input BooleanComparator {
    eq: Boolean
    neq: Boolean
}

input BooleanListComparator {
    eq: [Boolean]
    neq: [Boolean]
    contains: [Boolean]
    contained: [Boolean]
    overlap: [Boolean]
}`
