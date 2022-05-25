package image

type File struct {
	Width, Height int
	File          string
}

type Artifact struct {
	Files                         []File
	OriginalWidth, OriginalHeight int
}
