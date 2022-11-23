package controllers

import (
	"log"

	"cloud.google.com/go/firestore"
	"github.com/revel/revel"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type GoogleCloud struct {
	*revel.Controller
}

func (c GoogleCloud) Firestore() revel.Result {
	ctx := c.Request.Context()
	sa := option.WithCredentialsFile("key.json")
	client, err := firestore.NewClient(ctx, firestore.DetectProjectID, sa)
	if err != nil {
		log.Fatalf("Failtd to create client: %v", err)
	}
	defer client.Close()

	records := []map[string]interface{}{}
	iter := client.Collection("users").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		record := map[string]interface{}{}
		record[doc.Ref.ID] = doc.Data()
		records = append(records, record)
	}
	return c.RenderJSON(records)
}
