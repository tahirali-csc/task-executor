package engine

type Metadata struct {
	UID       string `json:"uid,omitempty"`
	Namespace string
}

type Spec struct {
	Metadata Metadata
	Image    string
	Command  []string
	Args     []string
}
