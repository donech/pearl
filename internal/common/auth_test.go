package common

import (
	"testing"
)

func init() {
	authKey = "donech0123456789"
}

func TestEncryptPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "No.1", args: args{password: "piao1234"}, want: "5AbwI06Aj84wiuDR1D08KA=="},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncryptPassword(tt.args.password); got != tt.want {
				t.Errorf("EncryptPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	type args struct {
		password        string
		encryptPassword string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "No.1", args: args{password: "piao1234", encryptPassword: "5AbwI06Aj84wiuDR1D08KA=="}, want: true},
		{name: "No.1", args: args{password: "piao1234", encryptPassword: "5AbwI06Aj84wiuDR1D08KA="}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidatePassword(tt.args.password, tt.args.encryptPassword); got != tt.want {
				t.Errorf("ValidatePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
