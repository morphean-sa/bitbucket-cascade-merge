package main

import (
	"github.com/libgit2/git2go"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var path = "/Users/sc/Desktop/repo.git"
var bare *git.Repository

func TestCascadeMerge(t *testing.T) {
	var err error

	os.RemoveAll(path)
	os.RemoveAll(filepath.Join(filepath.Dir(path), "repo"))

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
		Name:       "Jon Snow",
		Email:      "jon.snow@winterfell.net",
	}

	err = client.CascadeMerge("release/48", nil)

	stat, err := os.Stat(filepath.Join(work.Workdir(), "patch-1"))
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

	work, err := git.Clone(bare.Path(), filepath.Join(filepath.Dir(path), "cascade"), &git.CloneOptions{})
	CheckFatal(err, t)
	defer os.RemoveAll(work.Workdir())
	defer work.Free()

	os.Chdir(work.Workdir())

	client := &Client{
		Repository: work,
		Name:       "Jon Snow",
		Email:      "jon.snow@winterfell.net",
	}

	err = client.CascadeMerge("release/48", nil)
	if err == nil {
		t.Fail()
	}

	err = client.Checkout("release/49")
	CheckFatal(err, t)
	bytes49, err := ioutil.ReadFile("foo")
	CheckFatal(err, t)

	if !reflect.DeepEqual(bytes49, []byte("foo-edit-48")) {
		t.Fail()
	}

	err = client.Checkout("develop")
	CheckFatal(err, t)
	bytesDevelop, err := ioutil.ReadFile("foo")
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

	work, err := git.Clone(bare.Path(), filepath.Join(filepath.Dir(path), "cascade"), &git.CloneOptions{})
	CheckFatal(err, t)
	defer os.RemoveAll(work.Workdir())
	defer work.Free()

	os.Chdir(work.Workdir())

	client := &Client{
		Repository: work,
		Name:       "Jon Snow",
		Email:      "jon.snow@winterfell.net",
	}

	err = client.CascadeMerge("release/48", nil)
	if err == nil {
		t.Fail()
	}

	err = client.Checkout("release/49")
	CheckFatal(err, t)
	bytes49, err := ioutil.ReadFile("foo")
	CheckFatal(err, t)

	if !reflect.DeepEqual(bytes49, []byte("foo-same-edit")) {
		t.Fail()
	}

	err = client.Checkout("develop")
	CheckFatal(err, t)
	bytesDevelop, err := ioutil.ReadFile("foo")
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

func CommitDummyFile(filename string, client *Client, t *testing.T) {
	NewFile(filename, filename+"\n", t)
	_, err := client.Commit("add "+filename, filename)
	CheckFatal(err, t)
}

func NewFile(filename string, content string, t *testing.T) {
	err := ioutil.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		CheckFatal(err, t)
	}
}

type Task interface {
	Do(client *Client)
}

type InitializeWithReadmeTask struct {
	t *testing.T
}

func (t *InitializeWithReadmeTask) Do(client *Client) {
	filename := "README.md"
	NewFile(filename, "# Cascade Merge\n", t.t)

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
	CommitDummyFile(t.Filename, client, t.t)
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
	defer os.RemoveAll(tmp)

	work, err := git.Clone(bare.Path(), tmp, &git.CloneOptions{})
	if err != nil {
		return err
	}
	defer work.Free()

	err = os.Chdir(tmp)
	if err != nil {
		return err
	}
	defer os.Chdir(cwd)

	client := &Client{
		Repository: work,
		Name:       "Jon Snow",
		Email:      "jon.snow@winterfell.net",
	}

	for _, t := range tasks {
		t.Do(client)
	}

	return nil
}
