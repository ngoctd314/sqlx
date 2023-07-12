package sqlx

import (
	"reflect"
	"testing"
)

func TestOpen(t *testing.T) {
	type args struct {
		driverName     string
		dataSourceName string
	}
	tests := []struct {
		name    string
		args    args
		want    *DB
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Open(tt.args.driverName, tt.args.dataSourceName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Open() = %v, want %v", got, tt.want)
			}
		})
	}
}
