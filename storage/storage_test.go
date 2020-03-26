package storage

import (
	"fmt"
	"testing"
)

func TestUpload(t *testing.T) {
	storage := Storage{"mycats-ba2ef.appspot.com", "test", "/Users/stevebargelt/Downloads/mycats-ba2ef-2f24ef007822.json"}

	url, err := storage.Upload("/Users/stevebargelt/go/src/github.com/stevebargelt/cameraPoller/integrationtests/testdata/bear-in-day-00.jpg")
	if err != nil {
		t.Errorf("got %v want no error!", err)
	}
	if url == "" {
		t.Errorf("got empty string want URL!")
	}

	fmt.Printf("%s\n", url)
}
