package common

import (
	"testing"
)

func TestSplitName(t *testing.T) {
	type args struct {
		fullName string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			"default name",
			args{
				"John Doe",
			},
			"John",
			"Doe",
		},
		{
			"single name",
			args{
				"Mark",
			},
			"Mark",
			"",
		},
		{
			"long name",
			args{
				"Mark Anthony Something",
			},
			"Mark",
			"Anthony Something",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := SplitName(tt.args.fullName)
			if got != tt.want {
				t.Errorf("SplitName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SplitName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestExtractHostFromUrl(t *testing.T) {
	type args struct {
		rawUrl string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"magaluObjs",
			args{
				"https://br-se1.magaluobjects.com",
			},
			"br-se1.magaluobjects.com",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractHostFromUrl(tt.args.rawUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractHostFromUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractHostFromUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
