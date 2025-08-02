package database

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"go-fiber-template/models"

	"gorm.io/gorm"
)

// ModelInfo represents information about a database model
type ModelInfo struct {
	TableName string
	Model     interface{}
	Fields    []FieldInfo
}

// FieldInfo represents information about a model field
type FieldInfo struct {
	Name          string
	Type          string
	Size          int
	NotNull       bool
	Default       interface{}
	Unique        bool
	Index         bool
	PrimaryKey    bool
	AutoIncrement bool
	ForeignKey    string
	GormTag       string
	JsonTag       string
	// Foreign key constraint details
	ReferencedTable  string
	ReferencedColumn string
	OnUpdate         string
	OnDelete         string
}

// MigrationOperation represents a database migration operation
type MigrationOperation struct {
	Type        string // "add_column", "drop_column", "modify_column", "add_index", "drop_index", "add_constraint", "drop_constraint"
	TableName   string
	ColumnName  string
	OldField    *FieldInfo
	NewField    *FieldInfo
	SQL         string
	Description string
	// Foreign key constraint fields
	ConstraintName   string
	ReferencedTable  string
	ReferencedColumn string
	OnUpdate         string
	OnDelete         string
}

// ColumnInfo represents information about an existing database column
type ColumnInfo struct {
	Name         string
	Type         string
	IsNullable   bool
	Default      interface{}
	IsPrimaryKey bool
}

// DynamicMigrator handles dynamic database migrations
type DynamicMigrator struct {
	db     *gorm.DB
	models []ModelInfo
}

// NewDynamicMigrator creates a new dynamic migrator instance
func NewDynamicMigrator(db *gorm.DB) *DynamicMigrator {
	return &DynamicMigrator{
		db:     db,
		models: getRegisteredModels(),
	}
}

// getRegisteredModels returns all registered models for migration
func getRegisteredModels() []ModelInfo {
	modelsToRegister := []interface{}{
		// Core models
		&models.User{},
		&models.Log{},
		// Add other models here as they are created
	}

	var modelInfos []ModelInfo
	for _, model := range modelsToRegister {
		modelInfo := extractModelInfo(model)
		modelInfos = append(modelInfos, modelInfo)
	}

	return modelInfos
}

// GetRegisteredModels is a public wrapper for getRegisteredModels
func GetRegisteredModels() []ModelInfo {
	return getRegisteredModels()
}

// extractModelInfo extracts field information from a model
func extractModelInfo(model interface{}) ModelInfo {
	modelType := reflect.TypeOf(model).Elem()

	// Get table name from GORM
	stmt := &gorm.Statement{DB: DB}
	stmt.Parse(model)
	tableName := stmt.Schema.Table

	var fields []FieldInfo

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Skip embedded structs that are not database fields
		if field.Anonymous {
			continue
		}

		fieldInfo := extractFieldInfo(field)
		if fieldInfo.Name != "" {
			fields = append(fields, fieldInfo)
		}
	}

	// Extract foreign key relationships from struct fields
	extractForeignKeyRelationships(modelType, &fields, tableName)

	return ModelInfo{
		TableName: tableName,
		Model:     model,
		Fields:    fields,
	}
}

// toSnakeCase converts CamelCase to snake_case
func toSnakeCase(str string) string {
	// Handle common special cases for full words
	switch str {
	case "ID":
		return "id"
	case "URL":
		return "url"
	case "API":
		return "api"
	case "HTTP":
		return "http"
	case "JSON":
		return "json"
	case "XML":
		return "xml"
	case "SQL":
		return "sql"
	case "UUID":
		return "uuid"
	}

	// Handle cases where common acronyms are at the end of a word
	acronyms := []string{"ID", "URL", "API", "HTTP", "JSON", "XML", "SQL", "UUID"}
	for _, acronym := range acronyms {
		if strings.HasSuffix(str, acronym) && len(str) > len(acronym) {
			prefix := str[:len(str)-len(acronym)]
			return toSnakeCase(prefix) + "_" + strings.ToLower(acronym)
		}
	}

	// Handle cases where common acronyms are at the beginning of a word
	for _, acronym := range acronyms {
		if strings.HasPrefix(str, acronym) && len(str) > len(acronym) {
			suffix := str[len(acronym):]
			// Check if the next character is uppercase (indicating start of new word)
			if len(suffix) > 0 && suffix[0] >= 'A' && suffix[0] <= 'Z' {
				return strings.ToLower(acronym) + "_" + toSnakeCase(suffix)
			}
		}
	}

	var result strings.Builder
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		if r >= 'A' && r <= 'Z' {
			result.WriteRune(r - 'A' + 'a')
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// extractForeignKeyRelationships extracts foreign key relationships from model struct
func extractForeignKeyRelationships(modelType reflect.Type, fields *[]FieldInfo, _ string) {
	// Map to store table name mappings for different models
	tableNameMap := map[string]string{
		"User": "users",
		"Log":  "logs",
		// Add other models here as they are created
	}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		gormTag := field.Tag.Get("gorm")

		// Skip if no gorm tag or marked as ignored
		if gormTag == "" || gormTag == "-" {
			continue
		}

		// Skip association fields - we only want to process the foreign key fields
		if isGormAssociationField(field, gormTag) {
			// Parse this association field to find foreign key relationships
			var foreignKeyColumn, referencedTable, referencedColumn, onUpdate, onDelete string

			// Extract foreign key column name and constraint details
			tags := strings.Split(gormTag, ";")
			for _, tag := range tags {
				tag = strings.TrimSpace(tag)
				if strings.HasPrefix(tag, "foreignKey:") {
					foreignKeyColumn = strings.TrimPrefix(tag, "foreignKey:")
				} else if strings.HasPrefix(tag, "constraint:") {
					constraintStr := strings.TrimPrefix(tag, "constraint:")
					parts := strings.Split(constraintStr, ",")
					for _, part := range parts {
						part = strings.TrimSpace(part)
						if strings.HasPrefix(part, "OnUpdate:") {
							onUpdate = strings.TrimPrefix(part, "OnUpdate:")
						} else if strings.HasPrefix(part, "OnDelete:") {
							onDelete = strings.TrimPrefix(part, "OnDelete:")
						}
					}
				}
			}

			if foreignKeyColumn != "" {
				// Determine referenced table from field type
				fieldTypeName := field.Type.String()
				if strings.Contains(fieldTypeName, "*") {
					fieldTypeName = strings.TrimPrefix(fieldTypeName, "*")
				}
				if strings.Contains(fieldTypeName, "[]") {
					fieldTypeName = strings.TrimPrefix(fieldTypeName, "[]")
				}

				// Extract type name from package.Type format
				if dotIndex := strings.LastIndex(fieldTypeName, "."); dotIndex != -1 {
					fieldTypeName = fieldTypeName[dotIndex+1:]
				}

				// Get referenced table name
				if mappedTable, exists := tableNameMap[fieldTypeName]; exists {
					referencedTable = mappedTable
				} else {
					// Default conversion: convert to snake_case and pluralize
					referencedTable = toSnakeCase(fieldTypeName) + "s"
				}

				// Default referenced column is usually 'id'
				referencedColumn = "id"

				// Find the foreign key field in our fields slice and update it
				foreignKeyFieldName := toSnakeCase(foreignKeyColumn)
				for j, f := range *fields {
					if f.Name == foreignKeyFieldName {
						(*fields)[j].ReferencedTable = referencedTable
						(*fields)[j].ReferencedColumn = referencedColumn
						(*fields)[j].OnUpdate = onUpdate
						(*fields)[j].OnDelete = onDelete
						break
					}
				}
			}
		}
	}
}

// isGormAssociationField checks if a field is a GORM association (relationship) field
func isGormAssociationField(field reflect.StructField, gormTag string) bool {
	// Check if the field type is a struct or slice of structs (associations)
	fieldType := field.Type

	// Handle pointer types
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	// Handle slice types (for has-many relationships)
	if fieldType.Kind() == reflect.Slice {
		fieldType = fieldType.Elem()
		// Handle slice of pointers
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
	}

	// If the field type is a struct (not a basic type), it's likely an association
	if fieldType.Kind() == reflect.Struct {
		// Exclude time.Time which is a struct but should be stored as a database column
		if fieldType == reflect.TypeOf(time.Time{}) {
			return false
		}

		// Check if it has foreignKey tag (indicates this is an association field)
		if strings.Contains(gormTag, "foreignKey:") {
			return true
		}

		// Check for other association indicators
		if strings.Contains(gormTag, "references:") ||
			strings.Contains(gormTag, "many2many:") ||
			strings.Contains(gormTag, "polymorphic:") ||
			strings.Contains(gormTag, "joinForeignKey:") {
			return true
		}

		// If it's a struct type from our models packages, it's likely an association
		typeName := fieldType.String()
		if strings.Contains(typeName, "models.User") ||
			strings.Contains(typeName, "models.Log") {
			return true
		}

		return true // Default: if it's a struct, treat as association
	}

	return false
}

// extractFieldInfo extracts field information from a struct field
func extractFieldInfo(field reflect.StructField) FieldInfo {
	gormTag := field.Tag.Get("gorm")
	jsonTag := field.Tag.Get("json")

	// Skip fields marked with gorm:"-"
	if gormTag == "-" {
		return FieldInfo{}
	}

	// Skip GORM association fields (relationships)
	if isGormAssociationField(field, gormTag) {
		return FieldInfo{}
	}

	fieldInfo := FieldInfo{
		Name:    getFieldName(field.Name, gormTag),
		Type:    getFieldType(field.Type, gormTag),
		GormTag: gormTag,
		JsonTag: jsonTag,
	}

	// Parse GORM tags
	fieldInfo.parseGormTags(gormTag)

	return fieldInfo
}

// getFieldName extracts the database field name
func getFieldName(fieldName, gormTag string) string {
	// Check if column name is specified in gorm tag
	tags := strings.Split(gormTag, ";")
	for _, tag := range tags {
		if strings.HasPrefix(tag, "column:") {
			return strings.TrimPrefix(tag, "column:")
		}
	}

	// Convert to snake_case
	return toSnakeCase(fieldName)
}

// getFieldType determines the database field type
func getFieldType(fieldType reflect.Type, gormTag string) string {
	// Handle pointers
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	// Check for explicit type in gorm tag
	tags := strings.Split(gormTag, ";")
	for _, tag := range tags {
		if strings.HasPrefix(tag, "type:") {
			return strings.TrimPrefix(tag, "type:")
		}
	}

	// Check if this is an auto-increment field
	isAutoIncrement := false
	for _, tag := range tags {
		if strings.TrimSpace(tag) == "autoIncrement" {
			isAutoIncrement = true
			break
		}
	}

	// Map Go types to PostgreSQL types
	switch fieldType.Kind() {
	case reflect.String:
		// Check for size specification
		for _, tag := range tags {
			if strings.HasPrefix(tag, "size:") {
				size := strings.TrimPrefix(tag, "size:")
				return fmt.Sprintf("varchar(%s)", size)
			}
		}
		return "text"
	case reflect.Int, reflect.Int32:
		if isAutoIncrement {
			return "serial"
		}
		return "integer"
	case reflect.Int64:
		if isAutoIncrement {
			return "bigserial"
		}
		return "bigint"
	case reflect.Uint:
		// uint in Go should map to bigint in PostgreSQL to handle large values
		if isAutoIncrement {
			return "bigserial"
		}
		return "bigint"
	case reflect.Uint32:
		if isAutoIncrement {
			return "bigserial"
		}
		return "bigint"
	case reflect.Uint64:
		if isAutoIncrement {
			return "bigserial"
		}
		return "bigint"
	case reflect.Float32:
		return "real"
	case reflect.Float64:
		return "double precision"
	case reflect.Bool:
		return "boolean"
	default:
		if fieldType == reflect.TypeOf(time.Time{}) {
			return "timestamp"
		}
		return "text"
	}
}

// parseGormTags parses GORM tags and sets field properties
func (fi *FieldInfo) parseGormTags(gormTag string) {
	tags := strings.Split(gormTag, ";")

	for _, tag := range tags {
		tag = strings.TrimSpace(tag)

		switch {
		case tag == "primaryKey":
			fi.PrimaryKey = true
			fi.NotNull = true // Primary keys are always NOT NULL
		case tag == "autoIncrement":
			fi.AutoIncrement = true
		case tag == "not null":
			fi.NotNull = true
		case tag == "unique":
			fi.Unique = true
		case tag == "index":
			fi.Index = true
		case strings.HasPrefix(tag, "size:"):
			fmt.Sscanf(tag, "size:%d", &fi.Size)
		case strings.HasPrefix(tag, "default:"):
			fi.Default = strings.TrimPrefix(tag, "default:")
		case strings.HasPrefix(tag, "foreignKey:"):
			fi.ForeignKey = strings.TrimPrefix(tag, "foreignKey:")
		case strings.HasPrefix(tag, "constraint:"):
			// Parse constraint details: constraint:OnUpdate:CASCADE,OnDelete:SET NULL
			constraintStr := strings.TrimPrefix(tag, "constraint:")
			fi.parseConstraintDetails(constraintStr)
		}
	}
}

// parseConstraintDetails parses foreign key constraint details
func (fi *FieldInfo) parseConstraintDetails(constraintStr string) {
	// Split by comma: OnUpdate:CASCADE,OnDelete:SET NULL
	parts := strings.Split(constraintStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "OnUpdate:") {
			fi.OnUpdate = strings.TrimPrefix(part, "OnUpdate:")
		} else if strings.HasPrefix(part, "OnDelete:") {
			fi.OnDelete = strings.TrimPrefix(part, "OnDelete:")
		}
	}
}

// DetectChanges detects schema changes by comparing current models with database schema
func (dm *DynamicMigrator) DetectChanges() ([]MigrationOperation, error) {
	var operations []MigrationOperation

	for _, modelInfo := range dm.models {
		// Check if table exists
		tableExists := dm.db.Migrator().HasTable(modelInfo.TableName)

		if !tableExists {
			// Table doesn't exist, create it
			operation := MigrationOperation{
				Type:        "create_table",
				TableName:   modelInfo.TableName,
				Description: fmt.Sprintf("Create table %s", modelInfo.TableName),
			}
			operations = append(operations, operation)
			continue
		}

		// Table exists, check for column differences
		for _, field := range modelInfo.Fields {
			columnExists := dm.db.Migrator().HasColumn(modelInfo.TableName, field.Name)

			if !columnExists {
				// Column doesn't exist, add it
				operation := MigrationOperation{
					Type:        "add_column",
					TableName:   modelInfo.TableName,
					ColumnName:  field.Name,
					NewField:    &field,
					Description: fmt.Sprintf("Add column %s to table %s", field.Name, modelInfo.TableName),
				}
				operations = append(operations, operation)
			}
			// TODO: Add logic to detect column modifications
		}
	}

	return operations, nil
}

// ExecuteMigrations executes the detected migration operations
func (dm *DynamicMigrator) ExecuteMigrations(operations []MigrationOperation) error {
	for _, operation := range operations {
		switch operation.Type {
		case "create_table":
			// Find the model for this table
			var model interface{}
			for _, modelInfo := range dm.models {
				if modelInfo.TableName == operation.TableName {
					model = modelInfo.Model
					break
				}
			}

			if model != nil {
				if err := dm.db.AutoMigrate(model); err != nil {
					return fmt.Errorf("failed to create table %s: %w", operation.TableName, err)
				}
			}

		case "add_column":
			// Use GORM's migrator to add column
			if operation.NewField != nil {
				// Create a temporary struct to represent the new column
				// For now, we'll use AutoMigrate to handle this
				var model interface{}
				for _, modelInfo := range dm.models {
					if modelInfo.TableName == operation.TableName {
						model = modelInfo.Model
						break
					}
				}

				if model != nil {
					if err := dm.db.AutoMigrate(model); err != nil {
						return fmt.Errorf("failed to add column %s to table %s: %w",
							operation.ColumnName, operation.TableName, err)
					}
				}
			}

		case "modify_column":
			// Handle column modifications
			// TODO: Implement column modification logic

		case "add_index":
			// Handle index creation
			// TODO: Implement index creation logic

		default:
			return fmt.Errorf("unsupported migration operation: %s", operation.Type)
		}
	}

	return nil
}
