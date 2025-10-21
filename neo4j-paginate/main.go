package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	// Initialize Neo4j driver
	ctx := context.Background()
	driver, err := neo4j.NewDriverWithContext(
		"neo4j://localhost:7687",
		neo4j.BasicAuth("neo4j", "your_password", ""),
	)
	if err != nil {
		log.Fatal("Error connecting to Neo4j:", err)
	}
	defer driver.Close(ctx)

	// Verify connectivity
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		log.Fatal("Error verifying connection:", err)
	}

	// Create session
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	// Insert 1000 sample entries
	_, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Create a batch of MERGE statements
		var queryBuilder strings.Builder
		queryBuilder.WriteString("CREATE ")

		for i := 0; i < 1000; i++ {
			if i > 0 {
				queryBuilder.WriteString(", ")
			}
			queryBuilder.WriteString(fmt.Sprintf("(p%d:Person {name: 'Person-%d', age: %d})",
				i, i, 20+(i%50))) // Ages will range from 20 to 69
		}

		result, err := tx.Run(ctx, queryBuilder.String(), nil)
		if err != nil {
			return nil, err
		}
		return result.Consume(ctx)
	})
	if err != nil {
		log.Fatal("Error creating sample data:", err)
	}

	fmt.Println("Successfully inserted 1000 records")

	start := time.Now()

	// Paginate records
	pageSize := 50 // Increased page size for larger dataset
	skip := 0

	for range 10 {
		result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			query := `
			MATCH (p:Person)
			RETURN p.name AS name, p.age AS age
			ORDER BY p.name
			SKIP $skip
			LIMIT $pageSize`

			params := map[string]any{
				"skip":     skip,
				"pageSize": pageSize,
			}

			records, err := tx.Run(ctx, query, params)
			if err != nil {
				return nil, err
			}

			var results []map[string]any
			for records.Next(ctx) {
				record := records.Record()
				results = append(results, record.AsMap())
			}

			return results, nil
		})
		if err != nil {
			log.Fatal("Error querying data:", err)
		}

		// Type assert and process results
		if records, ok := result.([]map[string]any); ok {
			if len(records) == 0 {
				break // No more records
			}

			fmt.Printf("\nPage %d (records %d-%d):\n",
				(skip/pageSize)+1,
				skip+1,
				skip+len(records))

			for _, record := range records {
				fmt.Printf("Name: %v, Age: %v\n", record["name"], record["age"])
			}
		}

		skip += pageSize
	}

	fmt.Println("Time taken to paginate through the entries:", time.Since(start))
}
