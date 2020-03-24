package db

import (
	"reflect"
	"testing"
)

func TestRegister(t *testing.T) {
	type args struct {
		name    string
		factory dataStoreFactory
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "register fake store",
			args: args{
				name:    "some-fake-store",
				factory: NewFakeStore},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Register(tt.args.name, tt.args.factory)
			_, isRegistered := datastoreFactories["some-fake-store"]
			if !isRegistered {
				t.Errorf("Register(), store %s not registered", tt.name)
			}
		})
	}
}

func TestCreateDatastore(t *testing.T) {
	storeName := "some-fake-store"
	Register(storeName, NewFakeStore)
	store, err := CreateDatastore(storeName)
	if err != nil {
		t.Fatalf("unable to create datastore because =%v", err)
	}
	type args struct {
		datastoreType string
	}
	tests := []struct {
		name    string
		args    args
		want    DataStore
		wantErr bool
	}{
		{
			name: "create registered fake datastore",
			args: args{
				datastoreType: storeName,
			},
			want: store,
		},
		{
			name: "create unregister fake datastore",
			args: args{
				datastoreType: "unreg-fake-store",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateDatastore(tt.args.datastoreType)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDatastore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateDatastore() = %v, want %v", got, tt.want)
			}
		})
	}
}
