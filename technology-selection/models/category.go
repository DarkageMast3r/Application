package models

type Category struct {
	Id       int        `json:"id"`
	Name     string     `json:"name"`
	Children []Category `json:"children"`
}

var categoryRoot Category

func Category_Init() {
	// Because Go does not permit declaring the full structure during variable declaration
	categoryRoot.Id = 1
	categoryRoot.Name = "Root"
	categoryRoot.Children = make([]Category, 2)
	categoryRoot.Children[0].Id = 2
	categoryRoot.Children[0].Name = "Category 1"
	categoryRoot.Children[1].Id = 3
	categoryRoot.Children[1].Name = "Category 2"
	categoryRoot.Children[1].Children = make([]Category, 2)
	categoryRoot.Children[1].Children[0].Id = 4
	categoryRoot.Children[1].Children[0].Name = "Category 2-1"
	categoryRoot.Children[1].Children[1].Id = 5
	categoryRoot.Children[1].Children[1].Name = "Category 2-2"
}

func (category Category) findChildById(id int) *Category {
	for _, child := range category.Children {
		if child.Id == id {
			return &child
		}
		found := child.findChildById(id)
		if found != nil {
			return found
		}
	}
	return nil
}

func Category_Get_All() *Category {
	return &categoryRoot
}

func Category_Get_By_Id(id int) *Category {
	if categoryRoot.Id == id {
		return &categoryRoot
	}
	return categoryRoot.findChildById(id)
}
