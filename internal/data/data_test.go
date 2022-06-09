package data_test

import (
	"testing"
	"time"

	"github.com/drshapeless/shapeless-blog/internal/data"
)

var (
	timeFormat = "2006-01-02"
)

func TestData(t *testing.T) {
	db, err := data.OpenDB("shapeless-blog.db")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Resetting test database.")
	data.ResetDB(db)
	t.Log("Database reset.")

	m := data.NewModels(db)

	testpost := data.Post{
		Title:    "Test Post",
		URL:      "test-post",
		Content:  "Here is some test content.",
		CreateAt: time.Now().Add(time.Hour * -24).Format(timeFormat),
		UpdateAt: time.Now().Format(timeFormat),
	}

	t.Log("Inserting test posts...")
	err = m.Posts.Insert(&testpost)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Insert pass")

	t.Log("GetWithID 1...")
	p, err := m.Posts.GetWithID(1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(p)

	t.Log("GetWithTitle \"Test Post\"...")
	p, err = m.Posts.GetWithTitle("Test Post")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(p)

	t.Log("GetWithURL \"Test URL\"...")
	p, err = m.Posts.GetWithURL("test-post")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(p)

	t.Log("GetAll...")
	ps, err := m.Posts.GetAll(10, 1)
	if err != nil {
		t.Error(err)
		return
	}
	if len(ps) == 0 {
		t.Errorf("Get all fucked up\n")
		return
	}
	for _, pp := range ps {
		t.Log(pp)
	}

	t.Log("Updating...")
	p.Title = "Test Update"
	err = m.Posts.Update(p)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(p)

	t1 := data.Tag{
		PostID: 1,
		Tag:    "Test",
	}

	t2 := data.Tag{
		PostID: 1,
		Tag:    "Go",
	}

	t.Log("Inserting tags...")
	err = m.Tags.Insert(&t1)
	if err != nil {
		t.Error(err)
		return
	}
	err = m.Tags.Insert(&t2)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Inserted tags")

	t.Log("Getting all distinct tags...")
	tags, err := m.Tags.GetAllDistinctTags()
	if err != nil {
		t.Error(err)
		return
	}
	for _, tag := range tags {
		t.Logf("Tag: %s\n", tag)
	}

	t.Log("Getting post with tag")
	ps, err = m.Tags.GetPostsWithTag("Test", 10, 1)
	if err != nil {
		t.Error(err)
		return
	}
	for _, pp := range ps {
		t.Log(pp)
	}

	t.Log("Getting tags with post id...")
	tags, err = m.Tags.GetTagsWithPostID(1)
	if err != nil {
		t.Error(err)
		return
	}
	for _, tag := range tags {
		t.Logf("id: %s\n", tag)
	}

	t.Log("Deleting tag...")
	err = m.Tags.Delete("Test", 1)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("Deleting tag with post id...")
	err = m.Tags.DeleteAllWithPostID(1)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("Deleting post...")
	err = m.Posts.Delete(1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("Deleted")
}
