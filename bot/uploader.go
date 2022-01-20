package bot

import (
	"FilesToTarBot/database"
	"archive/tar"
	"context"
	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/uploader"
	"github.com/gotd/td/tg"
	"io"
	"log"
)

// uploadFiles uploads the user files to telegram
func uploadFiles(ctx context.Context, userID int64, entities tg.Entities, u message.AnswerableMessageUpdate) {
	// At first get the files from user
	files, totalSize := database.MainDatabase.GetFiles(userID)
	if len(files) == 0 {
		_, err := sender.Reply(entities, u).Text(ctx, "You have not submitted any files!")
		if err != nil {
			log.Printf("cannot send the error to user: %s\n", err)
		}
		return
	}
	// Setup a pipe to redirect the download to upload
	reader, writer := io.Pipe()
	defer func(reader *io.PipeReader) {
		_ = reader.Close()
	}(reader)
	// Write the files in tar
	go func() {
		tarFile := tar.NewWriter(writer)
		defer func() {
			_ = tarFile.Close()
			_ = writer.Close()
		}()
		for _, file := range files {
			err := tarFile.WriteHeader(&tar.Header{
				Name:   file.Name,
				Size:   file.Size,
				Mode:   0600,
				Format: tar.FormatGNU,
			})
			if err != nil {
				log.Printf("cannot write tar header: %s\n", err)
				return
			}
			_, err = downloader.NewDownloader().Download(api, &tg.InputDocumentFileLocation{
				ID:            file.ID,
				AccessHash:    file.AccessHash,
				FileReference: file.FileReference,
			}).Stream(ctx, tarFile)
			if err != nil {
				log.Printf("cannot write tar data: %s\n", err)
				return
			}
		}
	}()
	// Now setup an uploader
	uploadedFile, err := uploader.NewUploader(api).Upload(ctx, uploader.NewUpload("files.tar", reader, totalSize))
	if err != nil {
		log.Printf("cannot upload tar: %s\n", err)
	}
	// Send the uploaded file
	document := message.UploadedDocument(uploadedFile).
		MIME("application/x-tar").
		Filename("files.tar")
	_, err = sender.Reply(entities, u).Media(ctx, document)
	// Remove the files
	database.MainDatabase.Reset(userID)
}
