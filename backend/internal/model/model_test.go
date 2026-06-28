package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&User{}, &Category{}, &Tag{}, &Article{}, &Comment{})
	require.NoError(t, err)
	SetDB(db)
	return db
}

func TestUserModel(t *testing.T) {
	db := setupTestDB(t)
	user := User{Username: "admin", PasswordHash: "hash", Nickname: "Admin"}
	err := db.Create(&user).Error
	require.NoError(t, err)
	assert.NotZero(t, user.ID)

	var found User
	err = db.Where("username = ?", "admin").First(&found).Error
	require.NoError(t, err)
	assert.Equal(t, "admin", found.Username)
}

func TestCategoryModel(t *testing.T) {
	db := setupTestDB(t)
	cat := Category{Name: "Go", Slug: "go", Description: "Go语言"}
	err := db.Create(&cat).Error
	require.NoError(t, err)

	var found Category
	err = db.First(&found, cat.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "Go", found.Name)
}

func TestTagModel(t *testing.T) {
	db := setupTestDB(t)
	tag := Tag{Name: "Gin", Slug: "gin"}
	err := db.Create(&tag).Error
	require.NoError(t, err)

	var found Tag
	err = db.First(&found, tag.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "Gin", found.Name)
}

func TestArticleModelWithRelations(t *testing.T) {
	db := setupTestDB(t)

	cat := Category{Name: "Tech", Slug: "tech"}
	require.NoError(t, db.Create(&cat).Error)

	tag1 := Tag{Name: "Go", Slug: "go"}
	tag2 := Tag{Name: "Gin", Slug: "gin"}
	require.NoError(t, db.Create(&tag1).Error)
	require.NoError(t, db.Create(&tag2).Error)

	article := Article{
		Title:       "Hello World",
		Slug:        "hello-world",
		ContentMD:   "# Hello",
		ContentHTML: "<h1>Hello</h1>",
		Status:      1,
		CategoryID:  &cat.ID,
		Tags:        []Tag{tag1, tag2},
	}
	err := db.Create(&article).Error
	require.NoError(t, err)
	assert.NotZero(t, article.ID)

	var found Article
	err = db.Preload("Category").Preload("Tags").First(&found, article.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "Hello World", found.Title)
	assert.NotNil(t, found.Category)
	assert.Equal(t, "Tech", found.Category.Name)
	assert.Len(t, found.Tags, 2)
}

func TestCommentModel(t *testing.T) {
	db := setupTestDB(t)

	article := Article{
		Title:       "Test",
		Slug:        "test",
		ContentMD:   "test",
		ContentHTML: "<p>test</p>",
		Status:      1,
	}
	require.NoError(t, db.Create(&article).Error)

	comment := Comment{
		ArticleID: article.ID,
		Author:    "Visitor",
		Content:   "Great post!",
		Status:    1,
	}
	err := db.Create(&comment).Error
	require.NoError(t, err)
	assert.NotZero(t, comment.ID)

	reply := Comment{
		ArticleID: article.ID,
		ParentID:  &comment.ID,
		Author:    "Author",
		Content:   "Thanks!",
		Status:    1,
	}
	require.NoError(t, db.Create(&reply).Error)

	var found Comment
	err = db.Preload("Replies").First(&found, comment.ID).Error
	require.NoError(t, err)
	assert.Len(t, found.Replies, 1)
	assert.Equal(t, "Thanks!", found.Replies[0].Content)
}
