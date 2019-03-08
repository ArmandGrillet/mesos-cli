package mesos

// File holds information about files in a sandbox.
type File struct {
	GID   string  `json:"gid"`
	Mode  string  `json:"mode"`
	MTime float64 `json:"mtime"`
	NLink int     `json:"nlink"`
	Path  string  `json:"path"`
	Size  int     `json:"size"`
	UID   string  `json:"uid"`
}
