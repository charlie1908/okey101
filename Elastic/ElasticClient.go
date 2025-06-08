package Elastic

//go get github.com/olivere/elastic/v7

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"net/http"
	"okey101/Core"
	"okey101/Model"
	shared "okey101/Shared"
	"time"
)

// InsertAuditLog inserts an AuditLog into the Elasticsearch audit index.
func InsertAuditLog(auditLog Model.AuditLog, indexName string) error {
	ctx := context.Background()

	esclient, err := GetESClient()
	if err != nil {
		fmt.Println("Error initializing Elasticsearch client: ", err)
		return err
	}

	// Serialize the AuditLog struct into JSON
	dataJSON, err := json.Marshal(auditLog)
	if err != nil {
		fmt.Println("Error marshaling AuditLog to JSON: ", err)
		return err
	}

	// Check if index exists
	exists, err := esclient.IndexExists(indexName).Do(ctx)
	if err != nil {
		return fmt.Errorf("error checking index existence: %w", err)
	}

	if !exists {
		// Generate new backing index name with timestamp
		newIndexName := fmt.Sprintf("%s-%d", indexName, makeTimestampMilli())

		// Retrieve mapping from ElasticMaps
		mapping, found := ElasticMaps[indexName]
		if !found {
			return fmt.Errorf("no mapping found for index: %s", indexName)
		}

		// Create the index with mapping
		_, err := esclient.CreateIndex(newIndexName).BodyString(mapping).Do(ctx)
		if err != nil {
			return fmt.Errorf("failed to create index %s: %w", newIndexName, err)
		}

		// Optional: Add alias to point to indexName
		_, err = esclient.Alias().Add(newIndexName, indexName).Do(ctx)
		if err != nil {
			return fmt.Errorf("failed to set alias for index %s: %w", newIndexName, err)
		}
	}

	// Insert the document into the index
	_, err = esclient.Index().
		Index(indexName).
		BodyJson(json.RawMessage(dataJSON)).
		Do(ctx)

	if err != nil {
		return fmt.Errorf("failed to insert document into index %s: %w", indexName, err)
	}

	fmt.Printf("[Elastic][Insert-%s] Insertion Successful\n", indexName)
	return nil
}

func GetESClient() (*elastic.Client, error) {
	//HTTPS yayinlarda TLS Certificat'i ihmal ediyor. Sadece Dev ortami icin tavsiye edilir.
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	elasticUrl, _ := Core.Decrypt(shared.Config.ELASTICURL, shared.Config.SECRETKEY)
	elasticPassword, _ := Core.Decrypt(shared.Config.ELASTICPASSWORD, shared.Config.SECRETKEY)
	elasticUser, _ := Core.Decrypt(shared.Config.ELASTICUSER, shared.Config.SECRETKEY)
	client, err := elastic.NewClient(elastic.SetURL(elasticUrl),
		elastic.SetHttpClient(httpClient),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false), elastic.SetBasicAuth(elasticUser, elasticPassword))

	fmt.Println("ES initialized...")

	return client, err
}

func unixMilli(t time.Time) int64 {
	return t.Round(time.Millisecond).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func makeTimestampMilli() int64 {
	return unixMilli(time.Now())
}
