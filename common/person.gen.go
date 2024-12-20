package common

type Person struct {
	ID   int    `json:"id" `
	Name string `json:"name" `
}

func (p Person) GetName() string {
	return p.Name
}

func GetOthersName(p Person,
) string {
	return p.Name
}
