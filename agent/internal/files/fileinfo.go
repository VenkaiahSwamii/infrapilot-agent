// FileInfo represents a file or directory entry
package files

import "time"

type FileInfo struct {
    Name     string    `json:"name"`
    Path     string    `json:"path"`
    Size     int64     `json:"size"`
    IsDir    bool      `json:"is_dir"`
    Modified time.Time `json:"modified"`
}
