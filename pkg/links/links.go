package links

// Links contains links related to the resource
type Links struct {
	Self string `json:"self"`
}

// Resource is a resource with links
type Resource struct {
	Links Links `json:"links"`
}
