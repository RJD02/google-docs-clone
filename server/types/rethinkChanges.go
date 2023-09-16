package types

import "github.com/RJD02/google-docs-clone/model"

var output string = `map[new_val:map[Content:Hello World, this is the first document CreatedAt:2023-09-15 19:00:06.219 +0530 +05:30
EditLink: Id:1 IsPublic:false Title:Something Interesting UpdatedAt:2023-09-15 19:00:06.219 +0530 +05:30 UserId:1 ViewLink: id:77901c28-f98c-4582-a452-f24fa6344391]
old_val:<nil>]`

type RethinkChange struct {
	new_val model.Document
	old_val *int
}
