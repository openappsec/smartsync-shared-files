package rest

import (
	"encoding/xml"
	"io"
	"net/http"
	"strings"
	"time"

	"openappsec.io/smartsync-shared-files/internal/app/utils"
	"openappsec.io/errors"
	"openappsec.io/httputils/responses"
	"openappsec.io/log"
)

const (
	internalErrorBodyKey = "internal-error"
)

// PutFile stores the body in file with given path in uri
func (a *Adapter) PutFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	path := strings.TrimPrefix(r.URL.Path, "/api/")
	log.WithContextAndEventID(ctx, "67305fca-e3cb-4c3c-8537-fc633cc4742d").Infof("put file: %v", path)
	defer r.Body.Close()
	content, err := io.ReadAll(r.Body)
	if err != nil {
		log.WithContextAndEventID(
			ctx, "5efd2539-ebce-441f-a372-d801f354ca1b",
		).Errorf("failed to read request body. err: %v", err)
		errString := utils.CreateErrorBody(ctx, internalErrorBodyKey)
		responses.HTTPReturn(ctx, w, http.StatusInternalServerError, []byte(errString), true)
		return
	}
	err = a.svc.PutFile(ctx, path, content)
	if err != nil {
		log.WithContextAndEventID(ctx, "9de9ba8b-7e94-4ddb-befb-7cb02bdb5bf4").Errorf(
			"failed to put file. err: %v", err,
		)
		errString := utils.CreateErrorBody(ctx, internalErrorBodyKey)
		responses.HTTPReturn(ctx, w, http.StatusInternalServerError, []byte(errString), true)
		return
	}
	log.WithContextAndEventID(ctx, "f5ab58b3-0722-4525-a661-e819af8eb12f").Infof("put file %v success", path)
	responses.HTTPReturn(ctx, w, http.StatusOK, nil, true)
}

// GetFile returns the file content
func (a *Adapter) GetFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	path := strings.TrimPrefix(r.URL.Path, "/api/")
	log.WithContextAndEventID(ctx, "e2e5e899-bd6e-41d1-9ae6-8ee2c96e3a14").Infof("get file: %v", path)
	fileContent, err := a.svc.GetFile(ctx, path)
	if err != nil {
		if errors.IsClass(err, errors.ClassNotFound) {
			log.WithContextAndEventID(ctx, "12f72909-a816-444b-80a1-f48bfb286be7").Infof("file %v not found", path)
			responses.HTTPReturn(ctx, w, http.StatusNotFound, []byte{}, true)
			return
		}
		log.WithContextAndEventID(
			ctx, "4c7190fb-61fa-434f-80d1-cf00eb4a3595",
		).Errorf("unexpected error on get file: %v, err: %v", path, err)
		errString := utils.CreateErrorBody(ctx, internalErrorBodyKey)
		responses.HTTPReturn(ctx, w, http.StatusInternalServerError, []byte(errString), true)
		return
	}
	log.WithContextAndEventID(ctx, "56a4d207-3993-4282-9f83-d946c36a4afb").Infof(
		"got file ok, file length: %v", len(fileContent),
	)
	responses.HTTPReturn(ctx, w, http.StatusOK, fileContent, false)
}

type contents struct {
	Key          string
	LastModified string
}

type filesList struct {
	XMLName  xml.Name `xml:"ListBucketResult"`
	KeyCount int
	Contents []contents
}

// GetFilesList lists the files with given prefix
func (a *Adapter) GetFilesList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	path := r.URL.Query().Get("prefix")
	log.WithContextAndEventID(ctx, "3120a134-d0a9-4ce6-9317-b25631376d54").Infof(
		"listing files with prefix: %v", path,
	)
	files, err := a.svc.GetFilesList(ctx, path)
	if err != nil {
		log.WithContextAndEventID(ctx, "0954d824-c7e1-40b2-9005-b32015ffe7a8").Errorf(
			"failed to list files. err: %v", err,
		)
		errString := utils.CreateErrorBody(ctx, internalErrorBodyKey)
		responses.HTTPReturn(ctx, w, http.StatusInternalServerError, []byte(errString), true)
		return
	}
	filesListRes := filesList{
		KeyCount: len(files),
		Contents: make([]contents, len(files)),
	}
	for i, file := range files {
		filesListRes.Contents[i] = contents{
			Key:          file.Path,
			LastModified: file.LastModified.Format(time.RFC3339),
		}
	}
	response, err := xml.Marshal(filesListRes)
	if err != nil {
		log.WithContextAndEventID(
			ctx, "de79b434-b6b1-4626-85d9-11935b900756",
		).Errorf("failed to marshal list files %+v. err: %v", filesListRes, err)
		errString := utils.CreateErrorBody(ctx, internalErrorBodyKey)
		responses.HTTPReturn(ctx, w, http.StatusInternalServerError, []byte(errString), true)
		return
	}
	responses.HTTPReturn(ctx, w, http.StatusOK, response, true)
}
