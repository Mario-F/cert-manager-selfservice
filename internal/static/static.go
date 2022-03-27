//go:generate go run generator.go

package static

type embedStatic struct {
	storage map[string][]byte
}

// Create new static for embed files
func newEmbedStatic() *embedStatic {
	return &embedStatic{storage: make(map[string][]byte)}
}

// Add a file to static
func (e *embedStatic) Add(file string, content []byte) {
	e.storage[file] = content
}

// Get file's content
func (e *embedStatic) Get(file string) []byte {
	if f, ok := e.storage[file]; ok {
		return f
	}
	return nil
}

// Embed static expose
var static = newEmbedStatic()

// Add a file content to static
func Add(file string, content []byte) {
	static.Add(file, content)
}

// Get a file from static
func Get(file string) []byte {
	return static.Get(file)
}
