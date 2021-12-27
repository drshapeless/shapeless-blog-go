package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"

	"github.com/drshapeless/shapeless-blog/data"
)

func (app *application) Test() {
}

func (app *application) createPost(blog data.Blog) (int, error) {
	id, err := app.generateNewID()
	if err != nil {
		return 0, err
	}

	blog.Metadata.ID = id
	fmt.Printf("Creating post %d.\n", id)
	app.writeMetadata(blog.Metadata)
	app.writeContent(id, blog.Content)
	app.addToPosts(id)
	app.addNewCategories(blog.Metadata.Category, id)

	return id, nil
}

func (app *application) deletePost(id int) error {
	if app.isPostExist(id) {
		fmt.Printf("Deleting post %d.\n", id)
		cats := app.getMetadata(id).Category
		app.removeFromCategories(cats, id)
		app.removeFromPosts(id)
		app.deleteMetadata(id)
		app.deleteContent(id)
	} else {
		return fmt.Errorf("post with id:%d does not exist", id)
	}

	return nil
}

func (app *application) updatePost(blog data.Blog) error {
	id := blog.Metadata.ID
	if app.isPostExist(id) {
		fmt.Printf("Editing post %d.\n", id)
		old := app.getMetadata(id).Category
		app.writeMetadata(blog.Metadata)
		app.writeContent(id, blog.Content)
		if !isStringSlicesEqual(old, blog.Metadata.Category) {
			removed, added := removedAndAddedValue(old, blog.Metadata.Category)
			app.removeFromCategories(removed, id)
			app.addNewCategories(added, id)
		}
	} else {
		return fmt.Errorf("post with id: %d does not exist", blog.Metadata.ID)
	}

	return nil
}

func (app *application) generateNewID() (int, error) {
	current, err := app.readIDCounter()
	if err != nil {
		return current, err
	}
	new := current + 1
	err = app.writeIDCounter(new)
	if err != nil {
		return 0, err
	}
	return new, nil
}

func (app *application) readIDCounter() (int, error) {
	file, err := os.Open(app.IDCounterPath())
	if err != nil {
		return 0, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	s := scanner.Text()

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func (app *application) writeIDCounter(id int) error {
	file, err := os.Create(app.IDCounterPath())
	if err != nil {
		return err
	}

	defer file.Close()

	s := fmt.Sprintf("%d\n\nThis file should never be modified by hand.\n", id)
	_, err = file.WriteString(s)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) writeMetadata(md data.PostMetadata) error {
	file, err := os.Create(app.metadataPath(md.ID))
	if err != nil {
		return err
	}
	defer file.Close()

	s, err := json.MarshalIndent(md, "", "  ")
	if err != nil {
		return err
	}
	_, err = file.Write(s)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) writeContent(id int, cnt string) error {
	file, err := os.Create(app.contentPath(id))
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(cnt)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) addToPosts(id int) error {
	app.addToPostsCache(id)
	return app.WritePosts()
}

func (app *application) WritePosts() error {
	file, err := os.Create(app.postsPath())
	if err != nil {
		return err
	}
	defer file.Close()

	for _, value := range app.cache.PostIDs {
		_, err = file.WriteString(fmt.Sprintf("%d\n", value))
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *application) addNewCategories(cats []string, id int) error {
	oldcat := app.cache.Categories
	app.addToCategoriesCache(cats)
	if len(oldcat) != len(app.cache.Categories) {
		err := app.writeCategories()
		if err != nil {
			return err
		}
	}

	for _, value := range cats {
		arr, err := readFileAsIntArray(app.categoryPath(value))
		if err != nil {
			arr = []int{}
		}
		if isIDUnique(id, arr) {
			err = app.appendToCategory(id, value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (app *application) writeCategories() error {
	file, err := os.Create(app.categoriesPath())
	if err != nil {
		return err
	}

	defer file.Close()

	for _, value := range app.cache.Categories {
		if value == "" {
			continue
		}
		_, err = file.WriteString(fmt.Sprintf("%s\n", value))
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *application) appendToCategory(id int, cat string) error {
	file, err := os.OpenFile(app.categoryPath(cat), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%d\n", id))
	if err != nil {
		return err
	}

	return nil
}

func (app *application) isPostExist(id int) bool {
	return !isIDUnique(id, app.cache.PostIDs)
}

func (app *application) getMetadata(id int) data.PostMetadata {
	file, err := os.Open(app.metadataPath(id))
	if err != nil {
		// When calling this function, we assume that there is a file.
		panic(err)
	}

	defer file.Close()

	var md data.PostMetadata
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(byteValue, &md)
	if err != nil {
		// This must be a json file.
		panic(err)
	}

	return md
}

func (app *application) removeFromCategory(cat string, id int) error {
	arr, err := readFileAsIntArray(app.categoryPath(cat))
	if err != nil {
		return err
	}

	arr = removeInt(arr, id)
	sort.Ints(arr)

	file, err := os.Create(app.categoryPath(cat))
	if err != nil {
		return err
	}
	defer file.Close()

	for _, v := range arr {
		_, err = file.WriteString(fmt.Sprintf("%d\n", v))
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *application) removeFromCategories(cats []string, id int) error {
	for _, v := range cats {
		err := app.removeFromCategory(v, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *application) deleteMetadata(id int) error {
	return os.Remove(app.metadataPath(id))
}

func (app *application) deleteContent(id int) error {
	return os.Remove(app.contentPath(id))
}

func (app *application) removeFromPosts(id int) error {
	app.cache.PostIDs = removeInt(app.cache.PostIDs, id)
	return app.WritePosts()
}

func (app *application) hardReset() error {
	// There is no going back when using this.
	// Only use in testing.
	dir, err := ioutil.ReadDir(app.postsDir())
	if err != nil {
		return err
	}

	for _, d := range dir {
		fmt.Printf("Removing %s\n", d.Name())
		os.RemoveAll(app.postsDir() + "/" + d.Name())
	}

	fmt.Printf("Resetting idcounter.\n")
	app.writeIDCounter(0)

	fmt.Printf("Creating necessary files.\n")
	os.Create(app.postsPath())
	os.Create(app.categoriesPath())
	os.Mkdir(app.metadataDir(), os.ModePerm)
	os.Mkdir(app.contentDir(), os.ModePerm)
	os.Mkdir(app.categoryDir(), os.ModePerm)

	return nil
}

func (app *application) getBlog(id int) (data.Blog, error) {
	if !app.isPostExist(id) {
		return data.Blog{}, fmt.Errorf("post with id: %d does not exist", id)
	}
	blog := data.Blog{
		Content:  app.getContent(id),
		Metadata: app.getMetadata(id),
	}
	return blog, nil
}

func (app *application) getContent(id int) string {
	bv, err := os.ReadFile(app.contentPath(id))
	if err != nil {
		// We assume that the file exists.
		panic(err)
	}
	s := string(bv)
	return s
}
