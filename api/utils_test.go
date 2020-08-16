package api

import (
	"reflect"
	"testing"

	mapset "github.com/deckarep/golang-set"
)

func Test_isValid(t *testing.T) {
	type args struct {
		data        map[string]interface{}
		validFields mapset.Set
	}
	tests := []struct {
		name            string
		args            args
		wantValid       bool
		wantInvalidList []string
	}{
		{
			name: "notes",
			args: args{
				data:        map[string]interface{}{"notes": "123"},
				validFields: mapset.NewSetFromSlice([]interface{}{"assignee", "status", "group", "priority_order", "notes"}),
			},
			wantValid: true,
		},
		{
			name: "invalid",
			args: args{
				data:        map[string]interface{}{"invalid": "123"},
				validFields: mapset.NewSetFromSlice([]interface{}{"assignee", "status", "group", "priority_order", "notes"}),
			},
			wantValid:       false,
			wantInvalidList: []string{"invalid"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := isValid(tt.args.data, tt.args.validFields)
			if got != tt.wantValid {
				t.Errorf("isValid() got = %v, want %v", got, tt.wantValid)
			}
			if !reflect.DeepEqual(got1, tt.wantInvalidList) {
				t.Errorf("isValid() got1 = %v, want %v", got1, tt.wantInvalidList)
			}
		})
	}
}
