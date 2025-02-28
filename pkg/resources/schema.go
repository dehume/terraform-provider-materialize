package resources

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	defaultSchema   = "public"
	defaultDatabase = "materialize"
)

func ObjectNameSchema(resource string, required, forceNew bool) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: fmt.Sprintf("The identifier for the %s.", resource),
		Required:    required,
		Optional:    !required,
		ForceNew:    forceNew,
	}
}

func SchemaNameSchema(resource string, required bool) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: fmt.Sprintf("The identifier for the %s schema. Defaults to `public`.", resource),
		Required:    required,
		Optional:    !required,
		ForceNew:    true,
		Default:     defaultSchema,
	}
}

func DatabaseNameSchema(resource string, required bool) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: fmt.Sprintf("The identifier for the %s database. Defaults to `MZ_DATABASE` environment variable if set or `materialize` if environment variable is not set.", resource),
		Required:    required,
		Optional:    !required,
		ForceNew:    true,
		DefaultFunc: schema.EnvDefaultFunc("MZ_DATABASE", defaultDatabase),
	}
}

func ClusterNameSchema() *schema.Schema {
	return &schema.Schema{
		Description: "The cluster whose resources you want to create an additional computation of.",
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
	}
}

func OwnershipRoleSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: "The owernship role of the object.",
		Optional:    true,
		Computed:    true,
	}
}

func QualifiedNameSchema(resource string) *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeString,
		Description: fmt.Sprintf("The fully qualified name of the %s.", resource),
		Computed:    true,
	}
}

func SizeSchema(resource string, required bool, forceNew bool) *schema.Schema {
	return &schema.Schema{
		Type:         schema.TypeString,
		Description:  fmt.Sprintf("The size of the %s.", resource),
		Required:     required,
		Optional:     !required,
		ForceNew:     forceNew,
		ValidateFunc: validation.StringInSlice(replicaSizes, true),
	}
}

func ValidateConnectionSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeBool,
		Description: "**Private Preview** If the connection should wait for validation.",
		Optional:    true,
		Default:     true,
	}
}

func IdentifierSchema(elem, description string, required bool) *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Description: fmt.Sprintf("The %s name.", elem),
					Type:        schema.TypeString,
					Required:    true,
				},
				"schema_name": {
					Description: fmt.Sprintf("The %s schema name. Defaults to `public`.", elem),
					Type:        schema.TypeString,
					Optional:    true,
					Default:     defaultSchema,
				},
				"database_name": {
					Description: fmt.Sprintf("The %s database name. Defaults to `MZ_DATABASE` environment variable if set or `materialize` if environment variable is not set.", elem),
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("MZ_DATABASE", defaultDatabase),
				},
			},
		},
		Required:    required,
		Optional:    !required,
		MinItems:    1,
		MaxItems:    1,
		ForceNew:    true,
		Description: description,
	}
}

func ValueSecretSchema(elem string, description string, required bool) *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"text": {
					Description:   fmt.Sprintf("The `%s` text value. Conflicts with `secret` within this block", elem),
					Type:          schema.TypeString,
					Optional:      true,
					Sensitive:     true,
					ConflictsWith: []string{fmt.Sprintf("%s.0.secret", elem)},
				},
				"secret": IdentifierSchema(
					elem,
					fmt.Sprintf("The `%s` secret value. Conflicts with `text` within this block.", elem),
					false,
				),
			},
		},
		Required:    required,
		Optional:    !required,
		MinItems:    1,
		MaxItems:    1,
		ForceNew:    true,
		Description: fmt.Sprintf("%s. Can be supplied as either free text using `text` or reference to a secret object using `secret`.", description),
	}
}

func FormatSpecSchema(elem string, description string, required bool) *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"avro": {
					Description: "Avro format.",
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"schema_registry_connection": IdentifierSchema("schema_registry_connection", "The name of a schema registry connection.", true),
							"key_strategy": {
								Description:  "How Materialize will define the Avro schema reader key strategy.",
								Type:         schema.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringInSlice(strategy, true),
							},
							"value_strategy": {
								Description:  "How Materialize will define the Avro schema reader value strategy.",
								Type:         schema.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringInSlice(strategy, true),
							},
						},
					},
				},
				"protobuf": {
					Description: "Protobuf format.",
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"schema_registry_connection": IdentifierSchema("schema_registry_connection", "The name of a schema registry connection.", true),
							"message": {
								Description: "The name of the Protobuf message to use for the source.",
								Type:        schema.TypeString,
								Required:    true,
								ForceNew:    true,
							},
						},
					},
				},
				"csv": {
					Description: "CSV format.",
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					MaxItems:    2,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"column": {
								Description: "The columns to use for the source.",
								Type:        schema.TypeInt,
								Optional:    true,
								ForceNew:    true,
							},
							"delimited_by": {
								Description: "The delimiter to use for the source.",
								Type:        schema.TypeString,
								Optional:    true,
								ForceNew:    true,
							},
							"header": {
								Description: "The number of columns and the name of each column using the header row.",
								Type:        schema.TypeList,
								Elem:        &schema.Schema{Type: schema.TypeString},
								Optional:    true,
								ForceNew:    true,
							},
						},
					},
				},
				"bytes": {
					Description: "BYTES format.",
					Type:        schema.TypeBool,
					Optional:    true,
					ForceNew:    true,
				},
				"text": {
					Description: "Text format.",
					Type:        schema.TypeBool,
					Optional:    true,
					ForceNew:    true,
				},
				"json": {
					Description: "JSON format.",
					Type:        schema.TypeBool,
					Optional:    true,
					ForceNew:    true,
				},
			},
		},
		Required:    required,
		Optional:    !required,
		MinItems:    1,
		MaxItems:    1,
		ForceNew:    true,
		Description: description,
	}
}

func SinkFormatSpecSchema(elem string, description string, required bool) *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"avro": {
					Description: "Avro format.",
					Type:        schema.TypeList,
					Optional:    true,
					ForceNew:    true,
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"schema_registry_connection": IdentifierSchema("schema_registry_connection", "The name of a schema registry connection.", true),
							"avro_key_fullname": {
								Description: "The full name of the Avro key schema.",
								Type:        schema.TypeString,
								Optional:    true,
								ForceNew:    true,
							},
							"avro_value_fullname": {
								Description: "The full name of the Avro value schema.",
								Type:        schema.TypeString,
								Optional:    true,
								ForceNew:    true,
							},
							"avro_doc_type": {
								Description: "**Private Preview** Add top level documentation comment to the generated Avro schemas.",
								Type:        schema.TypeList,
								MinItems:    1,
								MaxItems:    1,
								Optional:    true,
								ForceNew:    true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"object": IdentifierSchema("object", "The object to apply the Avro documentation.", true),
										"doc": {
											Description: "Documentation string.",
											Type:        schema.TypeString,
											Required:    true,
										},
										"key": {
											Description: "Applies to the key schema.",
											Type:        schema.TypeBool,
											Optional:    true,
										},
										"value": {
											Description: "Applies to the value schema.",
											Type:        schema.TypeBool,
											Optional:    true,
										},
									},
								},
							},
							"avro_doc_column": {
								Description: "**Private Preview** Add column level documentation comment to the generated Avro schemas.",
								Type:        schema.TypeList,
								MinItems:    1,
								Optional:    true,
								ForceNew:    true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"object": IdentifierSchema("object", "The object to apply the Avro documentation.", true),
										"column": {
											Description: "Name of the column in the Avro schema to apply to.",
											Type:        schema.TypeString,
											Required:    true,
										},
										"doc": {
											Description: "Documentation string.",
											Type:        schema.TypeString,
											Required:    true,
										},
										"key": {
											Description: "Applies to the key schema.",
											Type:        schema.TypeBool,
											Optional:    true,
										},
										"value": {
											Description: "Applies to the value schema.",
											Type:        schema.TypeBool,
											Optional:    true,
										},
									},
								},
							},
						},
					},
				},
				"json": {
					Description: "JSON format.",
					Type:        schema.TypeBool,
					Optional:    true,
					ForceNew:    true,
				},
			},
		},
		Required:    required,
		Optional:    !required,
		MinItems:    1,
		MaxItems:    1,
		ForceNew:    true,
		Description: description,
	}
}

func SubsourceSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Subsources of a source.",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Description: "The name of the subsource.",
					Type:        schema.TypeString,
					Required:    true,
				},
				"schema_name":   SchemaNameSchema("source", false),
				"database_name": DatabaseNameSchema("source", false),
			},
		},
		Computed: true,
	}
}

func ObjectClusterNameSchema(objectType string) *schema.Schema {
	return &schema.Schema{
		Description:   fmt.Sprintf("The cluster to maintain this %s. If not specified, the `size` option must be specified.", objectType),
		Type:          schema.TypeString,
		Optional:      true,
		Computed:      true,
		AtLeastOneOf:  []string{"cluster_name", "size"},
		ConflictsWith: []string{"size"},
		ForceNew:      true,
	}
}

func ObjectSizeSchema(objectType string) *schema.Schema {
	return &schema.Schema{
		Description:   fmt.Sprintf("The size of the %s. If not specified, the `cluster_name` option must be specified.", objectType),
		Type:          schema.TypeString,
		Optional:      true,
		Computed:      true,
		AtLeastOneOf:  []string{"cluster_name", "size"},
		ConflictsWith: []string{"cluster_name"},
		ValidateFunc:  validation.StringInSlice(sourceSizes, true),
	}
}

func IntrospectionIntervalSchema(forceNew bool, requiredWith []string) *schema.Schema {
	return &schema.Schema{
		Description:  "The interval at which to collect introspection data.",
		Type:         schema.TypeString,
		Optional:     true,
		ForceNew:     forceNew,
		Default:      "1s",
		RequiredWith: requiredWith,
	}
}

func IntrospectionDebuggingSchema(forceNew bool, requiredWith []string) *schema.Schema {
	return &schema.Schema{
		Description:  "Whether to introspect the gathering of the introspection data.",
		Type:         schema.TypeBool,
		Optional:     true,
		ForceNew:     forceNew,
		Default:      false,
		RequiredWith: requiredWith,
	}
}

func IdleArrangementMergeEffortSchema(forceNew bool, requiredWith []string) *schema.Schema {
	return &schema.Schema{
		Description:  "The amount of effort to exert compacting arrangements during idle periods. This is an unstable option! It may be changed or removed at any time.",
		Type:         schema.TypeInt,
		Optional:     true,
		ForceNew:     forceNew,
		RequiredWith: requiredWith,
	}
}

func GranteeNameSchema() *schema.Schema {
	return &schema.Schema{
		Description: "The role name that will gain the default privilege. Use the `PUBLIC` pseudo-role to grant privileges to all roles.",
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
	}
}

func GrantDefaultDatabaseNameSchema() *schema.Schema {
	return &schema.Schema{
		Description: "The default privilege will apply only to objects created in this database, if specified.",
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
	}
}

func GrantDefaultSchemaNameSchema() *schema.Schema {
	return &schema.Schema{
		Description: "The default privilege will apply only to objects created in this schema, if specified.",
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
	}
}

func RoleNameSchema() *schema.Schema {
	return &schema.Schema{
		Description: "The name of the role to grant privilege to.",
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
	}
}

func TargetRoleNameSchema() *schema.Schema {
	return &schema.Schema{
		Description: "The default privilege will apply to objects created by this role. If this is left blank, then the current role is assumed. Use the `PUBLIC` pseudo-role to target objects created by all roles.",
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
	}
}

func PrivilegeSchema(objectType string) *schema.Schema {
	return &schema.Schema{
		Description:  "The privilege to grant to the object.",
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validPrivileges(objectType),
	}
}

func DiskSchema(forceNew bool) *schema.Schema {
	return &schema.Schema{
		Description: "**Private Preview**. Whether or not the replica is a _disk-backed replica_.",
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    forceNew,
	}
}

func CommentSchema(forceNew bool) *schema.Schema {
	return &schema.Schema{
		Description: "**Private Preview** Comment on an object in the database.",
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    forceNew,
	}
}
