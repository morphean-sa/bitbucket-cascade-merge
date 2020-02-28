package bitbucket

import (
	"reflect"
	"testing"
)

func TestRepository_URL(t *testing.T) {
	type fields struct {
		Links Links
	}
	type args struct {
		protocols string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "FindHttps",
			fields: fields{
				Links: Links{
					Clone: []*Link{{
						Name: "https",
						Href: "https://git.com/winterfell.git",
					}},
				},
			},
			args:    args{protocols: "https"},
			want:    "https://git.com/winterfell.git",
			wantErr: false,
		}, {
			name: "FindDefault",
			fields: fields{
				Links: Links{
					Clone: []*Link{{
						Name: "file",
						Href: "/tmp/git/winterfell.git",
					}},
				},
			},
			args:    args{},
			want:    "/tmp/git/winterfell.git",
			wantErr: false,
		}, {
			name: "Missing",
			fields: fields{
				Links: Links{
					Clone: []*Link{{
						Name: "file",
						Href: "/tmp/git/winterfell.git",
					}},
				},
			},
			args:    args{protocols: "https"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Links: tt.fields.Links,
			}
			got, err := r.URL(tt.args.protocols)
			if (err != nil) != tt.wantErr {
				t.Errorf("URL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("URL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
