package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	files_pb "github.com/reversersed/LitGO-proto/gen/go/files"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO write tests
func (f *fileServer) GetBookCover(c context.Context, r *files_pb.FileRequest) (*files_pb.FileResponse, error) {
	if err := f.validator.StructValidation(r); err != nil {
		return nil, err
	}

	filePath := fmt.Sprintf("/files/%s/%s", f.fileCfg.BookCoversFolder, r.GetFileName())
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.FileMode(0777))

	if err != nil {
		f.logger.Errorf("book cover file not found: %v (%s)", err, r.GetFileName())
		return nil, status.Error(codes.NotFound, "file not found")
	}
	defer file.Close()

	fileBytes := []byte{}
	buf := make([]byte, 1024*32)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			f.logger.Errorf("failed to read file: %v", err)
			return nil, status.Error(codes.Internal, "failed to read file")
		}

		fileBytes = append(fileBytes, buf[:n]...)
	}
	names := strings.Split(r.GetFileName(), ".")
	f.logger.Infof("file sended to client: %v from %d bytes (%d kB)", file.Name(), len(fileBytes), len(fileBytes)/1024)
	return &files_pb.FileResponse{File: fileBytes, Mimetype: fmt.Sprintf("image/%s", names[len(names)-1])}, nil
}
func (f *fileServer) GetBookFile(c context.Context, r *files_pb.FileRequest) (*files_pb.FileResponse, error) {
	if err := f.validator.StructValidation(r); err != nil {
		return nil, err
	}

	filePath := fmt.Sprintf("/files/%s/%s", f.fileCfg.BooksFolder, r.GetFileName())
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.FileMode(0777))

	if err != nil {
		f.logger.Errorf("book file not found: %v (%s)", err, r.GetFileName())
		return nil, status.Error(codes.NotFound, "file not found")
	}
	defer file.Close()

	fileBytes := []byte{}
	buf := make([]byte, 1024*32)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			f.logger.Errorf("failed to read file: %v", err)
			return nil, status.Error(codes.Internal, "failed to read file")
		}

		fileBytes = append(fileBytes, buf[:n]...)
	}

	f.logger.Infof("file sended to client: %v from %d bytes (%d kB)", file.Name(), len(fileBytes), len(fileBytes)/1024)
	return &files_pb.FileResponse{File: fileBytes, Mimetype: "application/epub+zip"}, nil
}
