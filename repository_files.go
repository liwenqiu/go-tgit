package tgit

import (
	"fmt"
	"net/http"
)

type RepositoryFilesService struct {
	client *Client
}

type File struct {
	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	Size     int    `json:"size"`
	Encoding string `json:"encoding"`
	Content  string `json:"content"`
	Ref      string `json:"ref"`
	BlobID   string `json:"blob_id"`
	CommitID string `json:"commit_id"`
}

func (r File) String() string {
	return Stringify(r)
}

type GetFileOptions struct {
	Ref      *string `url:"ref,omitempty""`
	FilePath *string `url:"file_path"`
}

func (s *RepositoryFilesService) GetFile(pid interface{}, opts *GetFileOptions) (*File, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/repository/files", project)

	req, err := s.client.NewRequest(http.MethodGet, u, opts)
	if err != nil {
		return nil, nil, err
	}

	f := new(File)
	resp, err := s.client.Do(req, f)
	if err != nil {
		return nil, resp, err
	}

	return f, resp, err
}

type FileInfo struct {
	FilePath   string `json:"file_path"`
	FileName   string `json:"file_name"`
	BranchName string `json:"branch_name"`
}

type CreateFileOptions struct {
	FilePath      *string `json:"file_path"`
	BranchName    *string `json:"branch_name"`
	Encoding      *string `json:"encoding,omitempty"`
	Content       *string `json:"content"`
	CommitMessage *string `json:"commit_message"`
}

func (s *RepositoryFilesService) CreateFile(pid interface{}, opts *CreateFileOptions) (*FileInfo, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/repository/files", project)

	req, err := s.client.NewRequest(http.MethodPost, u, opts)
	if err != nil {
		return nil, nil, err
	}

	f := new(FileInfo)
	resp, err := s.client.Do(req, f)
	if err != nil {
		return nil, resp, err
	}

	return f, resp, err
}

type UpdateFileOptions struct {
	FilePath      *string `json:"file_path"`
	BranchName    *string `json:"branch_name"`
	Encoding      *string `json:"encoding,omitempty"`
	Content       *string `json:"content"`
	CommitMessage *string `json:"commit_message"`
}

func (s *RepositoryFilesService) UpdateFile(pid interface{}, opts *UpdateFileOptions) (*FileInfo, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/repository/files", project)

	req, err := s.client.NewRequest(http.MethodPut, u, opts)
	if err != nil {
		return nil, nil, err
	}

	f := new(FileInfo)
	resp, err := s.client.Do(req, f)
	if err != nil {
		return nil, resp, err
	}

	return f, resp, err
}

type DeleteFileOptions struct {
	FilePath      *string `json:"file_path"`
	BranchName    *string `json:"branch_name"`
	CommitMessage *string `json:"commit_message"`
}

func (s *RepositoryFilesService) DeleteFile(pid interface{}, opts *DeleteFileOptions) (*FileInfo, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/repository/files", project)

	req, err := s.client.NewRequest(http.MethodDelete, u, opts)
	if err != nil {
		return nil, nil, err
	}

	f := new(FileInfo)
	resp, err := s.client.Do(req, f)
	if err != nil {
		return nil, resp, err
	}

	return f, resp, err
}
