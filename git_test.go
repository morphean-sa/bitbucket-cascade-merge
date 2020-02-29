package main

import (
	"github.com/libgit2/git2go"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

var path = filepath.Join(os.TempDir(), "cascade-"+time.Nanosecond.String()+".git")
var bare *git.Repository

func TestCascadeMerge(t *testing.T) {
	var err error

	os.RemoveAll(path)

	// we need a bare repository in our tests because libgit2 does not support local push to non bare repository yet.
	bare, err = git.InitRepository(path, true)
	CheckFatal(err, t)
	defer os.RemoveAll(path)
	defer bare.Free()

	err = WorkOnBareRepository(bare,
		&InitializeWithReadmeTask{
			t: t,
		},
		&CreateDummyFileOnBranchTask{
			BranchName: "release/48",
			Filename:   "foo",
			t:          t,
		},
		&CreateDummyFileOnBranchTask{
			BranchName: "release/49",
			Filename:   "bar",
			t:          t,
		},
		&CreateDummyFileOnBranchTask{
			BranchName: "develop",
			Filename:   "baz",
			t:          t,
		},
	)
	CheckFatal(err, t)

	t.Run("NoConflict", CascadeNoConflict)
	t.Run("Conflict", CascadeConflict)
	t.Run("AutoResolveNotWorking", CascadeAutoResolveNotWorking)
}

func CascadeNoConflict(t *testing.T) {
	err := WorkOnBareRepository(bare, &CreateDummyFileOnBranchTask{
		BranchName: "release/48",
		Filename:   "patch-1",
		t:          t,
	})
	CheckFatal(err, t)

	work, err := git.Clone(bare.Path(), filepath.Join(filepath.Dir(path), "cascade"), &git.CloneOptions{})
	CheckFatal(err, t)
	defer os.RemoveAll(work.Workdir())
	defer work.Free()

	client := &Client{
		Repository: work,
		Author: &Author{
			Name:  "Jon Snow",
			Email: "jon.snow@winterfell.net",
		},
	}

	// assume someone else push a new commit to the same branch
	err = WorkOnBareRepository(bare, &CreateDummyFileOnBranchTask{
		BranchName: "release/48",
		Filename:   "patch-2",
		t:          t,
	})
	CheckFatal(err, t)

	err = client.CascadeMerge("release/48", nil)

	stat, err := os.Stat(filepath.Join(work.Workdir(), "patch-1"))
	CheckFatal(err, t)
	if !stat.Mode().IsRegular() {
		t.Fail()
	}

	stat, err = os.Stat(filepath.Join(work.Workdir(), "patch-2"))
	CheckFatal(err, t)
	if !stat.Mode().IsRegular() {
		t.Fail()
	}

	CheckFatal(err, t)
}

func CascadeConflict(t *testing.T) {
	err := WorkOnBareRepository(bare,
		&ChangeFileOnBranchTask{
			BranchName: "release/48",
			Filename:   "foo",
			Content:    "foo-edit-48",
			t:          t,
		},
		&ChangeFileOnBranchTask{
			BranchName: "develop",
			Filename:   "foo",
			Content:    "foo-edit-develop",
			t:          t,
		},
	)
	CheckFatal(err, t)

	client, err := NewClient(&ClientOptions{
		Path: filepath.Join(filepath.Dir(path), "cascade"),
		URL:  bare.Path(),
		Author: &Author{
			Name:  "Jon Snow",
			Email: "jon.snow@winterfell.net",
		},
	})
	CheckFatal(err, t)
	defer os.RemoveAll(filepath.Join(filepath.Dir(path), "cascade"))
	defer client.Close()

	err = client.CascadeMerge("release/48", nil)
	if err == nil {
		t.Fail()
	}

	err = client.Checkout("release/49")
	CheckFatal(err, t)
	bytes49, err := client.ReadFile("foo")
	CheckFatal(err, t)

	if !reflect.DeepEqual(bytes49, []byte("foo-edit-48")) {
		t.Fail()
	}

	err = client.Checkout("develop")
	CheckFatal(err, t)
	bytesDevelop, err := client.ReadFile("foo")
	CheckFatal(err, t)

	if !reflect.DeepEqual(bytesDevelop, []byte("foo-edit-develop")) {
		t.Fail()
	}
}

func CascadeAutoResolveNotWorking(t *testing.T) {
	err := WorkOnBareRepository(bare,
		&ChangeFileOnBranchTask{
			BranchName: "release/48",
			Filename:   "foo",
			Content:    "foo-same-edit",
			t:          t,
		},
		&ChangeFileOnBranchTask{
			BranchName: "release/49",
			Filename:   "foo",
			Content:    "foo-same-edit",
			t:          t,
		},
	)
	CheckFatal(err, t)

	client, err := NewClient(&ClientOptions{
		Path: filepath.Join(filepath.Dir(path), "cascade"),
		URL:  bare.Path(),
		Author: &Author{
			Name:  "Jon Snow",
			Email: "jon.snow@winterfell.net",
		},
	})
	CheckFatal(err, t)
	defer os.RemoveAll(filepath.Join(filepath.Dir(path), "cascade"))
	defer client.Close()

	err = client.CascadeMerge("release/48", nil)
	if err == nil {
		t.Fail()
	}

	err = client.Checkout("release/49")
	CheckFatal(err, t)
	bytes49, err := client.ReadFile("foo")
	CheckFatal(err, t)

	if !reflect.DeepEqual(bytes49, []byte("foo-same-edit")) {
		t.Fail()
	}

	err = client.Checkout("develop")
	CheckFatal(err, t)
	bytesDevelop, err := client.ReadFile("foo")
	CheckFatal(err, t)

	if !reflect.DeepEqual(bytesDevelop, []byte("foo-edit-develop")) {
		t.Fail()
	}
}

func CheckFatal(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

func (c *Client) CommitDummyFile(filename string, t *testing.T) {
	c.NewFile(filename, filename+"\n", t)
	_, err := c.Commit("add "+filename, filename)
	CheckFatal(err, t)
}

func (c *Client) NewFile(filename string, content string, t *testing.T) {
	err := c.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		CheckFatal(err, t)
	}
}

func (c *Client) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(c.Repository.Workdir(), filename))
}

func (c *Client) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filepath.Join(c.Repository.Workdir(), filename), data, perm)
}

type Task interface {
	Do(client *Client)
}

type InitializeWithReadmeTask struct {
	t *testing.T
}

func (t *InitializeWithReadmeTask) Do(client *Client) {
	filename := "README.md"
	client.NewFile(filename, "# Cascade Merge\n", t.t)

	_, err := client.Commit("initial commit", filename)
	CheckFatal(err, t.t)
	err = client.Push("master")
	CheckFatal(err, t.t)
}

type CreateDummyFileOnBranchTask struct {
	BranchName string
	Filename   string
	t          *testing.T
}

func (t *CreateDummyFileOnBranchTask) Do(client *Client) {
	err := client.Checkout(t.BranchName)
	CheckFatal(err, t.t)
	client.CommitDummyFile(t.Filename, t.t)
	err = client.Push(t.BranchName)
	CheckFatal(err, t.t)
}

type ChangeFileOnBranchTask struct {
	BranchName string
	Filename   string
	Content    string
	t          *testing.T
}

func (t *ChangeFileOnBranchTask) Do(client *Client) {
	err := client.Checkout(t.BranchName)
	CheckFatal(err, t.t)

	err = ioutil.WriteFile(t.Filename, []byte(t.Content), 0644)
	CheckFatal(err, t.t)

	_, err = client.Commit("edit "+t.Filename, t.Filename)
	CheckFatal(err, t.t)

	err = client.Push(t.BranchName)
	CheckFatal(err, t.t)
}

func WorkOnBareRepository(bare *git.Repository, tasks ...Task) error {

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	tmp, err := ioutil.TempDir(os.TempDir(), "cascade-test-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	client, err := NewClient(&ClientOptions{
		Path: tmp,
		URL:  bare.Path(),
		Author: &Author{
			Name:  "Jon Snow",
			Email: "jon.snow@winterfell.net",
		},
	})
	if err != nil {
		return err
	}
	defer client.Close()

	err = os.Chdir(tmp)
	if err != nil {
		return err
	}
	defer os.Chdir(cwd)

	for _, t := range tasks {
		t.Do(client)
	}

	return nil
}

func TestClientOptions_Validate(t *testing.T) {
	type fields struct {
		Author *Author
		Path   string
		URL    string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Valid",
			fields: fields{
				Author: &Author{Name: "Jon Snow", Email: "jon@snow.nl"},
				Path:   "907ab3da-653e-460e-bb5b-11b0b0b95e3f",
				URL:    "https://git.com/winterfell.git",
			},
			want: true,
		}, {
			name: "Invalid",
			fields: fields{
				Author: &Author{Name: "Jon Snow", Email: "jon@snow.nl"},
				Path:   "907ab3da-653e-460e-bb5b-11b0b0b95e3f",
				URL:    "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &ClientOptions{
				Author: tt.fields.Author,
				Path:   tt.fields.Path,
				URL:    tt.fields.URL,
			}
			if got := o.Validate(); got != tt.want {
				t.Errorf("Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
