package complicated_without_check

type User struct {
	Name          string  `validate:"min_len=3,max_len=64"`
	LastName      *string `validate:"nullable,min_len=1,max_len=5"`
	Age           uint    `validate:"min=18,max=105"`
	ChildrenCount *int    `validate:"not_null,min=0,max=15"`
	Float         float64 `validate:"min=-4.22,max=42.55"`
	Dog           Dog
	DogPointer    *Dog
	DogOptional   Dog      `validate:"method=ValidateOptional"`
	Urls          []string `validate:"min_items=1,item=[max_len=256]"`
	Dogs          []*Dog   `validate:"min_items=1,item=[nullable]"`
	Test          *[]int   `validate:"nullable,min_items=1,item=[min=4]"`
	Flag          bool
	Some          interface{}    `validate:"func=validateSome"`
	SomeArray     []interface{}  `validate:"min_items=1,item=[func=validateSome]"`
	Dict          map[string]int `validate:"min_items=2,key=[max_len=64],value=[min=-35,max=34]"`
	DictDogs      map[string]Dog `validate:"value=[method=ValidateOptional]"`
	Alias         DogsMapAlias
	AliasOnAlias  AliasOnDogsMapAlias
	MapOfMap      map[string]map[int]string `validate:"value=[min_items=1,value=[min_len=3]]"`
}

func validateSome(i interface{}) error {
	return nil
}

type Dog struct {
	Name string `validate:"min_len=1,max_len=64"`
}

func (Dog) ValidateOptional() error {
	return nil //len(dog.Name) is not validating here
}

type DogsMapAlias map[string]Dog

type AliasOnDogsMapAlias DogsMapAlias
