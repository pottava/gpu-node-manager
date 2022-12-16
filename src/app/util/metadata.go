package util

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"
)

type GoogleCloudMetadata struct {
	Key string
	URI string
}

const (
	gceInstanceMetadataPrefix = "http://metadata.google.internal/computeMetadata/v1"
)

var (
	flavor = &http.Header{"Metadata-Flavor": []string{"Google"}}
	meta   = []*GoogleCloudMetadata{}
)

func init() {
	meta = append(meta, &GoogleCloudMetadata{
		Key: "project_id",
		URI: "/project/project-id",
	})
	meta = append(meta, &GoogleCloudMetadata{
		Key: "project_number",
		URI: "/project/numeric-project-id",
	})
	meta = append(meta, &GoogleCloudMetadata{
		Key: "region",
		URI: "/instance/region",
	})
	meta = append(meta, &GoogleCloudMetadata{
		Key: "instance_id",
		URI: "/instance/id",
	})
	meta = append(meta, &GoogleCloudMetadata{
		Key: "email",
		URI: "/instance/service-accounts/default/email",
	})
}

func InstanceMetadata(ctx context.Context) map[string]string {
	client := &http.Client{
		Transport: &http.Transport{MaxIdleConnsPerHost: len(meta)},
		Timeout:   time.Duration(250) * time.Millisecond,
	}
	wg := &sync.WaitGroup{}
	wg.Add(len(meta))

	result := map[string]string{}
	for _, m := range meta {
		go func(meta *GoogleCloudMetadata) {
			defer wg.Done()

			url := gceInstanceMetadataPrefix + meta.URI
			body, err := get(ctx, client, url, flavor)
			if err == nil {
				result[meta.Key] = string(body)
			}
		}(m)
	}
	wg.Wait()
	return result
}

const retryCount = 1

func get(ctx context.Context, client *http.Client, endpoint string, headers *http.Header) ([]byte, error) {
	var err error
	for i := 0; i < retryCount; i++ {
		resp, e := getOnce(ctx, client, http.MethodGet, endpoint, headers)
		if e == nil {
			return resp, nil
		}
		err = e
		time.Sleep(100 * time.Millisecond)
	}
	return nil, err
}

func getOnce(ctx context.Context, client *http.Client, method, endpoint string, headers *http.Header) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, endpoint, nil)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		req.Header = *headers
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
