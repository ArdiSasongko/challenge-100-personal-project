package model

type (
	UploadModel struct {
		ID       int64  `db:"id" json:"id"`
		FileName string `db:"file_name" json:"file_name"`
		URL      string `db:"url" json:"url"`
	}
)
