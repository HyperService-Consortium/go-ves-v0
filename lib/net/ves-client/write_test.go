package vesclient

import (
	"fmt"
	"testing"
)

func TestVesClient_readOpIntents(t *testing.T) {

	type args struct {
		filePath   string
		fileBuffer []byte
	}

	tests := []struct {
		name    string
		vc      *VesClient
		args    args
		wantErr bool
	}{
		{"test_easy", &VesClient{}, args{
			filePath:   "./test.json",
			fileBuffer: make([]byte, 5000),
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vc := tt.vc
			if n, err := vc.readOpIntents(tt.args.filePath, tt.args.fileBuffer); (err != nil) != tt.wantErr {
				t.Errorf("readOpIntents() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				fmt.Println(n)
			}
		})
	}
}
