package provider

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/MaterializeInc/terraform-provider-materialize/pkg/materialize"
	"github.com/MaterializeInc/terraform-provider-materialize/pkg/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jmoiron/sqlx"
)

func TestAccClusterReplica_basic(t *testing.T) {
	clusterName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	replicaName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterReplicaResource(clusterName, replicaName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterReplicaExists("materialize_cluster_replica.test"),
					resource.TestMatchResourceAttr("materialize_cluster_replica.test", "id", terraformObjectIdRegex),
					resource.TestCheckResourceAttr("materialize_cluster_replica.test", "cluster_name", clusterName),
					resource.TestCheckResourceAttr("materialize_cluster_replica.test", "name", replicaName),
					resource.TestCheckResourceAttr("materialize_cluster_replica.test", "size", "3xsmall"),
					resource.TestCheckResourceAttr("materialize_cluster_replica.test", "introspection_interval", "1s"),
					resource.TestCheckResourceAttr("materialize_cluster_replica.test", "introspection_debugging", "false"),
					resource.TestCheckResourceAttr("materialize_cluster_replica.test", "disk", "true"),
					resource.TestCheckNoResourceAttr("materialize_cluster_replica.test", "idle_arrangement_merge_effort"),
					resource.TestCheckResourceAttr("materialize_cluster_replica.test", "comment", ""),
				),
			},
			{
				ResourceName:            "materialize_cluster_replica.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"introspection_debugging", "introspection_interval"},
			},
		},
	})
}

func TestAccClusterReplica_update(t *testing.T) {
	clusterName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	replicaName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	comment := "cluster replica comment"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterReplicaResource(clusterName, replicaName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterReplicaExists("materialize_cluster_replica.test"),
					resource.TestCheckResourceAttr("materialize_cluster_replica.test", "comment", ""),
				),
			},
			{
				Config: testAccClusterReplicaWithComment(clusterName, replicaName, comment),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterReplicaExists("materialize_cluster_replica.test"),
					resource.TestCheckResourceAttr("materialize_cluster_replica.test", "comment", comment),
				),
			},
		},
	})
}

func TestAccClusterReplica_disappears(t *testing.T) {
	clusterName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	replicaName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAllClusterReplicaDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterReplicaResource(clusterName, replicaName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterReplicaExists("materialize_cluster_replica.test"),
					testAccCheckObjectDisappears(
						materialize.MaterializeObject{
							ObjectType:  "CLUSTER REPLICA",
							Name:        replicaName,
							ClusterName: clusterName,
						},
					),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccClusterReplicaResource(clusterName, clusterReplica string) string {
	return fmt.Sprintf(`
	resource "materialize_cluster" "test" {
		name = "%[1]s"
	}

	resource "materialize_cluster_replica" "test" {
		cluster_name = materialize_cluster.test.name
		name = "%[2]s"
		size = "3xsmall"
		disk = true
	}
	`, clusterName, clusterReplica)
}

func testAccClusterReplicaWithComment(clusterName, clusterReplica, comment string) string {
	return fmt.Sprintf(`
	resource "materialize_cluster" "test" {
		name = "%[1]s"
	}

	resource "materialize_cluster_replica" "test" {
		cluster_name = materialize_cluster.test.name
		name = "%[2]s"
		size = "3xsmall"
		disk = true
		comment = "%[3]s"
	}
	`, clusterName, clusterReplica, comment)
}

func testAccCheckClusterReplicaExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		db := testAccProvider.Meta().(*sqlx.DB)
		r, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("cluster replica not found: %s", name)
		}
		_, err := materialize.ScanClusterReplica(db, utils.ExtractId(r.Primary.ID))
		return err
	}
}

func testAccCheckAllClusterReplicaDestroyed(s *terraform.State) error {
	db := testAccProvider.Meta().(*sqlx.DB)

	for _, r := range s.RootModule().Resources {
		if r.Type != "materialize_cluster_replica" {
			continue
		}

		_, err := materialize.ScanClusterReplica(db, utils.ExtractId(r.Primary.ID))
		if err == nil {
			return fmt.Errorf("Cluster replica %v still exists", utils.ExtractId(r.Primary.ID))
		} else if err != sql.ErrNoRows {
			return err
		}
	}

	return nil
}
