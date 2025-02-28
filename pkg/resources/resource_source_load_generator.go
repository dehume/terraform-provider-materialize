package resources

import (
	"context"
	"fmt"
	"log"

	"github.com/MaterializeInc/terraform-provider-materialize/pkg/materialize"
	"github.com/MaterializeInc/terraform-provider-materialize/pkg/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmoiron/sqlx"
)

var tick_interval = &schema.Schema{
	Description: "The interval at which the next datum should be emitted. Defaults to one second.",
	Type:        schema.TypeString,
	Optional:    true,
	ForceNew:    true,
}

var scale_factor = &schema.Schema{
	Description: "The scale factor for the generator. Defaults to 0.01 (~ 10MB).",
	Type:        schema.TypeFloat,
	Optional:    true,
	Default:     0.01,
	ForceNew:    true,
}

var sourceLoadgenSchema = map[string]*schema.Schema{
	"name":               ObjectNameSchema("source", true, false),
	"schema_name":        SchemaNameSchema("source", false),
	"database_name":      DatabaseNameSchema("source", false),
	"qualified_sql_name": QualifiedNameSchema("source"),
	"comment":            CommentSchema(false),
	"cluster_name":       ObjectClusterNameSchema("source"),
	"size":               ObjectSizeSchema("source"),
	"load_generator_type": {
		Description:  fmt.Sprintf("The load generator types: %s.", loadGeneratorTypes),
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice(loadGeneratorTypes, true),
	},
	"counter_options": {
		Description: "Counter Options.",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"tick_interval": tick_interval,
				"scale_factor":  scale_factor,
				"max_cardinality": {
					Description: "Causes the generator to delete old values to keep the collection at most a given size. Defaults to unlimited.",
					Type:        schema.TypeInt,
					Optional:    true,
					ForceNew:    true,
				},
			},
		},
		Optional:      true,
		MinItems:      1,
		MaxItems:      1,
		ForceNew:      true,
		ConflictsWith: []string{"auction_options", "marketing_options", "tpch_options"},
	},
	"auction_options": {
		Description: "Auction Options.",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"tick_interval": tick_interval,
				"scale_factor":  scale_factor,
			},
		},
		Optional:      true,
		MinItems:      1,
		MaxItems:      1,
		ForceNew:      true,
		ConflictsWith: []string{"counter_options", "marketing_options", "tpch_options"},
	},
	"marketing_options": {
		Description: "Marketing Options.",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"tick_interval": tick_interval,
				"scale_factor":  scale_factor,
			},
		},
		Optional:      true,
		MinItems:      1,
		MaxItems:      1,
		ForceNew:      true,
		ConflictsWith: []string{"counter_options", "auction_options", "tpch_options"},
	},
	"tpch_options": {
		Description: "TPCH Options.",
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"tick_interval": tick_interval,
				"scale_factor":  scale_factor,
			},
		},
		Optional:      true,
		MinItems:      1,
		MaxItems:      1,
		ForceNew:      true,
		ConflictsWith: []string{"counter_options", "auction_options", "marketing_options"},
	},
	"expose_progress": IdentifierSchema("expose_progress", "The name of the progress subsource for the source. If this is not specified, the subsource will be named `<src_name>_progress`.", false),
	"subsource":       SubsourceSchema(),
	"ownership_role":  OwnershipRoleSchema(),
}

func SourceLoadgen() *schema.Resource {
	return &schema.Resource{
		Description: "A load generator source produces synthetic data for use in demos and performance tests.",

		CreateContext: sourceLoadgenCreate,
		ReadContext:   sourceRead,
		UpdateContext: sourceUpdate,
		DeleteContext: sourceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: sourceLoadgenSchema,
	}
}

func sourceLoadgenCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	sourceName := d.Get("name").(string)
	schemaName := d.Get("schema_name").(string)
	databaseName := d.Get("database_name").(string)

	o := materialize.MaterializeObject{ObjectType: "SOURCE", Name: sourceName, SchemaName: schemaName, DatabaseName: databaseName}
	b := materialize.NewSourceLoadgenBuilder(meta.(*sqlx.DB), o)

	if v, ok := d.GetOk("cluster_name"); ok {
		b.ClusterName(v.(string))
	}

	if v, ok := d.GetOk("size"); ok {
		b.Size(v.(string))
	}

	if v, ok := d.GetOk("expose_progress"); ok {
		e := materialize.GetIdentifierSchemaStruct(v)
		b.ExposeProgress(e)
	}

	if v, ok := d.GetOk("load_generator_type"); ok {
		b.LoadGeneratorType(v.(string))
	}

	if v, ok := d.GetOk("counter_options"); ok {
		o := materialize.GetCounterOptionsStruct(v)
		b.CounterOptions(o)
	}

	if v, ok := d.GetOk("auction_options"); ok {
		o := materialize.GetAuctionOptionsStruct(v)
		b.AuctionOptions(o)
	}

	if v, ok := d.GetOk("marketing_options"); ok {
		o := materialize.GetMarketingOptionsStruct(v)
		b.MarketingOptions(o)
	}

	if v, ok := d.GetOk("tpch_options"); ok {
		o := materialize.GetTPCHOptionsStruct(v)
		b.TPCHOptions(o)
	}

	// create resource
	if err := b.Create(); err != nil {
		return diag.FromErr(err)
	}

	// ownership
	if v, ok := d.GetOk("ownership_role"); ok {
		ownership := materialize.NewOwnershipBuilder(meta.(*sqlx.DB), o)

		if err := ownership.Alter(v.(string)); err != nil {
			log.Printf("[DEBUG] resource failed ownership, dropping object: %s", o.Name)
			b.Drop()
			return diag.FromErr(err)
		}
	}

	// object comment
	if v, ok := d.GetOk("comment"); ok {
		comment := materialize.NewCommentBuilder(meta.(*sqlx.DB), o)

		if err := comment.Object(v.(string)); err != nil {
			log.Printf("[DEBUG] resource failed comment, dropping object: %s", o.Name)
			b.Drop()
			return diag.FromErr(err)
		}
	}

	// set id
	i, err := materialize.SourceId(meta.(*sqlx.DB), o)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(utils.TransformIdWithRegion(i))

	return sourceRead(ctx, d, meta)
}
