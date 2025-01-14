package resources

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/MaterializeInc/terraform-provider-materialize/pkg/materialize"
	"github.com/MaterializeInc/terraform-provider-materialize/pkg/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmoiron/sqlx"
	"golang.org/x/exp/slices"
)

type DefaultPrivilegeKey struct {
	objectType   string
	granteeId    string
	targetRoleId string
	databaseId   string
	schemaId     string
	privilege    string
}

func parseDefaultPrivilegeKey(id string) (DefaultPrivilegeKey, error) {
	ie := strings.Split(id, "|")

	if len(ie) != 7 {
		return DefaultPrivilegeKey{}, fmt.Errorf("%s: cannot be parsed correctly", id)
	}

	return DefaultPrivilegeKey{
		objectType:   ie[1],
		granteeId:    ie[2],
		targetRoleId: ie[3],
		databaseId:   ie[4],
		schemaId:     ie[5],
		privilege:    ie[6],
	}, nil
}

func grantDefaultPrivilegeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	i := d.Id()

	key, err := parseDefaultPrivilegeKey(i)
	if err != nil {
		log.Printf("[WARN] malformed privilege (%s), removing from state file", d.Id())
		d.SetId("")
		return nil
	}

	privileges, err := materialize.ScanDefaultPrivilege(meta.(*sqlx.DB), key.objectType, key.granteeId, key.targetRoleId, key.databaseId, key.schemaId)
	if err == sql.ErrNoRows {
		log.Printf("[WARN] grant (%s) not found, removing from state file", d.Id())
		d.SetId("")
		return nil
	} else if err != nil {
		return diag.FromErr(err)
	}

	// Check if default privilege has expected privilege
	privilegeMap, _ := materialize.MapDefaultGrantPrivileges(privileges)
	mapKey := strings.ToLower(key.objectType) + "|" + key.granteeId + "|" + key.targetRoleId + "|" + key.databaseId + "|" + key.schemaId

	if !slices.Contains(privilegeMap[mapKey], key.privilege) {
		log.Printf("[DEBUG] %s object does not contain privilege %s", i, key.privilege)
		// Remove id from state
		d.SetId("")
	}

	d.SetId(utils.TransformIdWithRegion(i))
	return nil
}
