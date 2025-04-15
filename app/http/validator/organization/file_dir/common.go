package file_dir

type Dir struct {
	Dir string `form:"dir" json:"dir" binding:"omitempty,max=512"`
}
