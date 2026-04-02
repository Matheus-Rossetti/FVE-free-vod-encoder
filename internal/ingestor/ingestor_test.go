package ingestor

import (
	"context"
	"testing"
)

// func TestPrepareFileForDownload(t *testing.T) {
// 	testCases := []struct {
// 		name      string
// 		fileState string
// 		wantErr   bool
// 	}{
// 		{
// 			name:      "opened file",
// 			fileState: "open",
// 			wantErr:   false,
// 		},
// 		{
// 			name:      "closed file",
// 			fileState: "closed",
// 			wantErr:   true,
// 		},
// 	}

// 	for _, tt := range testCases {
// 		t.Run(tt.name, func(t *testing.T) {
// 			file, err := os.CreateTemp("", "temp-file-*")
// 			if err != nil {
// 				t.Errorf("failed to create a temp file to run tests with")
// 			}
// 			defer file.Close()
// 			defer os.Remove(file.Name())

// 			if tt.fileState == "closed" {
// 				file.Close()
// 			}

// 			ingestor := ingestor{}

// 			err = ingestor.prepareFileForDownload(file)

// 			hasErr := (err != nil)
// 			if tt.wantErr != hasErr {
// 				t.Errorf("expected error to be %v, got %v", tt.wantErr, hasErr)
// 				return // no need to continue test
// 			}

// 			meta, _ := file.Stat()
// 			if meta.Size() != 0 {
// 				t.Errorf("expected file size to be 0, got %v", meta.Size())
// 			}

// 			offSet, _ := file.Seek(0, 1)
// 			if offSet != 0 {
// 				t.Errorf("expected offset of file to be at 0, got %v", offSet)
// 			}
// 		})
// 	}
// }

func TestCheckFileType(t *testing.T) {

	testCases := []struct {
		name     string
		fileType string
		wantErr  bool
	}{
		{
			name:     "video file type",
			fileType: "video",
			wantErr:  false,
		},
		{
			name:     "not video file type",
			fileType: "not video",
			wantErr:  true,
		},
	}

	ingestor := ingestor{
		ctx: context.Background(),
	}

	buffer := make([]byte, 512)

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			if tt.fileType == "video" {
				buffer = videoMagicBytes()
			} else {
				buffer = textMagicBytes()
			}

			err := ingestor.checkForVideo(buffer)

			hasErr := (err != nil)
			if tt.wantErr != hasErr {
				t.Errorf("Expected error to be %v, got %v",
					tt.wantErr, hasErr)
			}
		})
	}
}

// Files have something called "magic bytes", which are the first few bytes of a file
// they are used to identify the file type.
// Here are some examples: https://gist.github.com/iahu/396eaf109ed0969382abdbc9c3f0f029
func videoMagicBytes() []byte {
	return []byte{0x00, 0x00, 0x00, 0x18, 'f', 't', 'y', 'p', 'm', 'p', '4', '2'}
}
func textMagicBytes() []byte {
	return []byte("This is not a videofile")
}
