package googlecloud

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/go-openapi/swag"
	"github.com/pottava/gpu-node-manager/src/app/util"
	"google.golang.org/api/iterator"
)

type Notebook struct {
	FirestoreID string    `json:"-"`
	Email       string    `json:"email"`
	Active      bool      `json:"active"`
	Menu        string    `json:"menu"`
	Runtime     string    `json:"runtime"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func GetNotebook(ctx context.Context, email, runtime string) (*Notebook, error) {
	notebooks, err := GetNotebooks(ctx, email)
	if err != nil {
		return nil, err
	}
	for _, note := range notebooks {
		if note.Runtime == runtime {
			return note, nil
		}
	}
	return nil, fmt.Errorf("no notebook found. Runtime: %s,", runtime)
}

func GetNotebooks(ctx context.Context, email string) ([]*Notebook, error) {
	projectID := firestore.DetectProjectID
	if !swag.IsZero(projectID) {
		projectID = util.ProjectID()
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
		note.FirestoreID = doc.Ref.ID
		records = append(records, note)
	}
	return records, nil
}

func SaveNotebook(ctx context.Context, name, email, menu string) error {
	projectID := firestore.DetectProjectID
	if !swag.IsZero(projectID) {
		projectID = util.ProjectID()
	}
	client, err := firestore.NewClient(ctx, projectID, clientOption())
	if err != nil {
		return err
	}
	defer client.Close()

	if _, _, err = client.Collection("notebooks").Add(ctx, map[string]interface{}{
		"email":      email,
		"active":     true,
		"menu":       menu,
		"runtime":    name,
		"created_at": firestore.ServerTimestamp,
		"updated_at": firestore.ServerTimestamp,
	}); err != nil {
		return err
	}
	return nil
}

func UpdateNotebook(ctx context.Context, ID string, updates map[string]interface{}) error {
	return updateRecord(ctx, "notebooks", ID, updates)
}

type VM struct {
	FirestoreID string    `json:"-"`
	Email       string    `json:"email"`
	Active      bool      `json:"active"`
	Menu        string    `json:"menu"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func GetVM(ctx context.Context, email, name string) (*VM, error) {
	vms, err := GetVMs(ctx, email)
	if err != nil {
		return nil, err
	}
	for _, vm := range vms {
		if vm.Name == name {
			return vm, nil
		}
	}
	return nil, fmt.Errorf("no vm found. Name: %s,", name)
}

func GetVMs(ctx context.Context, email string) ([]*VM, error) {
	projectID := firestore.DetectProjectID
	if !swag.IsZero(projectID) {
		projectID = util.ProjectID()
	}
	client, err := firestore.NewClient(ctx, projectID, clientOption())
	if err != nil {
		return nil, err
	}
	defer client.Close()

	records := []*VM{}
	iter := client.Collection("vms").Where("email", "==", email).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		bytes, _ := json.Marshal(doc.Data())

		vm := &VM{}
		if err = json.Unmarshal(bytes, vm); err != nil {
			return nil, err
		}
		vm.FirestoreID = doc.Ref.ID
		records = append(records, vm)
	}
	return records, nil
}

func SaveVM(ctx context.Context, name, email, menu string) error {
	projectID := firestore.DetectProjectID
	if !swag.IsZero(projectID) {
		projectID = util.ProjectID()
	}
	client, err := firestore.NewClient(ctx, projectID, clientOption())
	if err != nil {
		return err
	}
	defer client.Close()

	if _, _, err = client.Collection("vms").Add(ctx, map[string]interface{}{
		"email":      email,
		"active":     true,
		"menu":       menu,
		"name":       name,
		"created_at": firestore.ServerTimestamp,
		"updated_at": firestore.ServerTimestamp,
	}); err != nil {
		return err
	}
	return nil
}

func UpdateVM(ctx context.Context, ID string, updates map[string]interface{}) error {
	return updateRecord(ctx, "vms", ID, updates)
}

func updateRecord(ctx context.Context, table, ID string, updates map[string]interface{}) error {
	projectID := firestore.DetectProjectID
	if !swag.IsZero(projectID) {
		projectID = util.ProjectID()
	}
	client, err := firestore.NewClient(ctx, projectID, clientOption())
	if err != nil {
		return err
	}
	defer client.Close()

	ref := client.Collection(table).Doc(ID)
	if ref == nil {
		return errors.New("the record was not found")
	}
	values := []firestore.Update{{
		FieldPath: []string{"updated_at"},
		Value:     firestore.ServerTimestamp,
	}}
	for key, value := range updates {
		values = append(values, firestore.Update{
			FieldPath: []string{key},
			Value:     value,
		})
	}
	_, err = ref.Update(ctx, values)
	return err
}
