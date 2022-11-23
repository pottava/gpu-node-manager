package googlecloud

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/firestore"
	"github.com/go-openapi/swag"
	"google.golang.org/api/iterator"
)

type Notebook struct {
	Email   string `json:"email"`
	Menu    string `json:"menu"`
	Runtime string `json:"runtime"`
}

func GetNotebooks(ctx context.Context, email string) ([]*Notebook, error) {
	projectID := firestore.DetectProjectID
	if !swag.IsZero(ProjectID) {
		projectID = ProjectID
	}
	client, err := firestore.NewClient(ctx, projectID, clientOption())
	if err != nil {
		return nil, err
	}
	defer client.Close()

	records := []*Notebook{}
	iter := client.Collection("notebooks").Where("email", "==", email).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		bytes, _ := json.Marshal(doc.Data())

		note := &Notebook{}
		if err = json.Unmarshal(bytes, note); err != nil {
			return nil, err
		}
		records = append(records, note)
	}
	return records, nil
}

func SaveNotebook(ctx context.Context, name, email, menu string) error {
	projectID := firestore.DetectProjectID
	if !swag.IsZero(ProjectID) {
		projectID = ProjectID
	}
	client, err := firestore.NewClient(ctx, projectID, clientOption())
	if err != nil {
		return err
	}
	defer client.Close()

	if _, _, err = client.Collection("notebooks").Add(ctx, map[string]interface{}{
		"email":   email,
		"menu":    menu,
		"runtime": name,
	}); err != nil {
		return err
	}
	return nil
}
