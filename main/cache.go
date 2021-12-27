package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"sort"
)

func (app *application) addToPostsCache(id int) {
	for _, value := range app.cache.PostIDs {
		if id == value {
			panic("Duplicated id.")
		}
	}
	app.cache.PostIDs = append(app.cache.PostIDs, id)
}

func (app *application) addToCategoriesCache(cats []string) {
	for _, n := range cats {
		if n == "" {
			continue
		}
		p := true
		for _, o := range app.cache.Categories {
			if o == n {
				p = false
				break
			}
		}
		if p {
			app.cache.Categories = append(app.cache.Categories, n)
			sort.Strings(app.cache.Categories)
		}
	}
}

func (app *application) loadPostsCache() error {
	temp, err := readFileAsIntArray(app.postsPath())
	if err != nil {
		// Should always panic if posts file is not found.
		panic(err)
	}

	app.cache.PostIDs = temp
	return nil
}

func (app *application) loadCategoriesCache() error {
	file, err := os.Open(app.categoriesPath())
	if err != nil {
		return err
	}

	defer file.Close()

	temp := []string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		if s != "" {
			temp = append(temp, s)
		}
	}

	app.cache.Categories = temp
	return nil
}

func (app *application) loadCache() error {
	err := app.loadPostsCache()
	if err != nil {
		return err
	}

	err = app.loadCategoriesCache()
	if err != nil {
		return err
	}

	err = app.loadTokenHash()
	if err != nil {
		return err
	}

	return nil
}

func (app *application) loadTokenHash() error {
	bytes, err := ioutil.ReadFile(app.tokenPath())
	if err != nil {
		return err
	}
	app.cache.TokenHash = bytes
	return nil
}
