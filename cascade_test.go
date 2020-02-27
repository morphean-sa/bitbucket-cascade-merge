package main

import (
	"math"
	"reflect"
	"testing"
)

func TestCascade_Next(t *testing.T) {
	type fields struct {
		Branches []string
		Current  int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "InBound", fields: fields{Branches: []string{"release/2", "release/3"}, Current: 0}, want: "release/3"},
		{name: "OutOfBound", fields: fields{Branches: []string{"release/2", "release/3"}, Current: 1}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cascade{
				Branches: tt.fields.Branches,
				Current:  tt.fields.Current,
			}
			if got := c.Next(); got != tt.want {
				t.Errorf("Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCascade_Append(t *testing.T) {
	type fields struct {
		BranchNames []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{name: "SortNumeric", fields: fields{BranchNames: []string{"release/3", "release/2"}}, want: []string{"release/2", "release/3"}},
		{name: "SortDevelop", fields: fields{BranchNames: []string{"develop", "release/3"}}, want: []string{"release/3", "develop"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cascade{
				Branches: make([]string, 0),
				Current:  0,
			}
			for _, n := range tt.fields.BranchNames {
				c.Append(n)
			}
			if !reflect.DeepEqual(c.Branches, tt.want) {
				t.Errorf("Next() = %v, want %v", c.Branches, tt.want)
			}
		})
	}
}

func Test_extractVersion(t *testing.T) {
	type args struct {
		b string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "valid int",
			args: args{b: "kind/10"},
			want: 10,
		}, {
			name: "invalid int",
			args: args{b: "kind/not-int"},
			want: math.MaxInt32,
		}, {
			name: "invalid int (float)",
			args: args{b: "kind/10.1"},
			want: math.MaxInt32,
		}, {
			name: "invalid format",
			args: args{b: "invalid format"},
			want: math.MaxInt32,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractVersion(tt.args.b); got != tt.want {
				t.Errorf("extractVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
