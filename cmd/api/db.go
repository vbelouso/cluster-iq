package main

import (
	"context"
	"fmt"
	"os"

	"github.com/RHEcosystemAppEng/cluster-iq/internal/inventory"
	"github.com/jackc/pgx/v5"
)

const (
	// SelectInstancesQuery returns every instance in the inventory ordered by ID
	SelectInstancesQuery = "SELECT * FROM instances ORDER BY id"

	// SelectInstancesByIDQuery returns an instance by its ID
	SelectInstancesByIDQuery = "SELECT * FROM instances WHERE id = $1 ORDER BY id"

	// SelectClustersQuery returns every cluster in the inventory ordered by Name
	SelectClustersQuery = "SELECT * FROM clusters ORDER BY name"

	// SelectClustersByNameQuery returns an cluster by its Name
	SelectClustersByNameQuery = "SELECT * FROM clusters WHERE name = $1 ORDER BY name"

	// SelectInstancesOnClusterQuery returns every instance belonging to a acluster
	SelectInstancesOnClusterQuery = "SELECT * FROM instances WHERE cluster_name = $1 ORDER BY id"

	// SelectAccountsQuery returns every instance in the inventory ordered by Name
	SelectAccountsQuery = "SELECT * FROM accounts ORDER BY name"

	// SelectAccountsByNameQuery returns an instance by its Name
	SelectAccountsByNameQuery = "SELECT * FROM accounts WHERE name = $1 ORDER BY name"

	// SelectClustersOnAccountQuery returns an cluster by its Name
	SelectClustersOnAccountQuery = "SELECT * FROM clusters WHERE account_name = $1 ORDER BY name"

	// InsertInstanceQuery inserts into a new instance in its table
	InsertInstanceQuery = "INSERT INTO instances (id, name, provider, instance_type, region, state, cluster_name) VALUES (:id, :name, :provider, :instance_type, :region, :state, :cluster_name)"

	// DeleteInstanceQuery inserts into a new instance in its table
	DeleteInstanceQuery = "DELETE FROM instances WHERE id=$1"
)

// getAccounts returns every account in Stock
func getInstances() ([]inventory.Instance, error) {
	var instances []inventory.Instance
	rows, err := dbpool.Query(context.Background(), SelectInstancesQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	instances, err = pgx.CollectRows(rows, pgx.RowToStructByName[inventory.Instance])
	if err != nil {
		return nil, err
	}

	return instances, nil
}

func getInstanceByID(instanceID string) ([]inventory.Instance, error) {
	rows, err := dbpool.Query(context.Background(), SelectInstancesByIDQuery, instanceID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	var instance inventory.Instance
	instance, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[inventory.Instance])
	if err != nil {
		return nil, err
	}

	return []inventory.Instance{instance}, nil
}
