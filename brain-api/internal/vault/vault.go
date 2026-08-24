// The vault on disk is the authoritative store, this package reads and writes md
package vault

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	envVaultPath = "VAULT_PATH"
	tmpPattern   = ".*.md.tmp"
)

var skipDirNames = map[string]struct{}{
	".git":      {},
	".obsidian": {},
	".trash":    {},
}

// Note is one markdown file under the vault root
type Note struct {
	// Path is relative to the vault root, slash-separated
	Path string
	// frontmatter is the YAML block, if any ever nil after Read
	Frontmatter map[string]any
	Body        string
	ModTime     time.Time
}

// Store reads and writes notes as files under Root
type Store struct {
	Root string
}

// Open returns a Store for an existing dir
func Open(root string) (*Store, error) {
	path , err := os.Stat(root) 
	
	if (err != nil) {
		return nil, fmt.Errorf("error occurred, %w", err)
	}
	if !(path.IsDir()) {
		return nil, fmt.Errorf("not directory at %s", root)
	}
	return &Store{Root: root}, nil
}
// OpenFromEnv opens VAULT_PATH
func OpenFromEnv() (*Store, error) {
	env_path, ok := os.LookupEnv(envVaultPath)
    
	if !ok|| value == "" {
		return nil, fmt.Errorf("not path found in .env occurred, %s", env_path)
	}
	return Open(env_path)
}

// Read loads one note by relative path
func (s *Store) Read(relPath string) (*Note, error) {

}
// Write creates or replaces a note. The file is written to a sibling temp name
// then renamed so livesync-bridge never sees a partial note
func (s *Store) Write(n Note) error {

}

// Walk visits every .md file under the vault, skipping .git, .obsidian, and
// .trash. fn receives a fully read Note. Returning a non-nil error stops the walk
func (s *Store) Walk(fn func(Note) error) error {
	
}

func (s *Store) resolve(relPath string) (string, error) {

}

func toSlash(relPath string) (string, error) {
	if relPath == "" {
		return "", errors.New("error, the path is not real")
	}
	cleanPath := filepath.ToSlash(relPath)
	cleanPath := path.Clean(cleanPath)

	return cleanPath, nil

}
