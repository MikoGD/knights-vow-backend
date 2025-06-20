package files

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"knights-vow/internal/database"
	"knights-vow/internal/resources/users"
)

type File struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	CreatedDate   string `json:"createdDate"`
	OwnerID       int    `json:"ownerID"`
	OwnerUsername string `json:"ownerUsername"`
}

type FileChunkSaveMessage struct {
	message          string
	chunkNumber      int
	uploadPercentage int
}
