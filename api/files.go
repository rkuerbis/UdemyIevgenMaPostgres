package api

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"os"

	"github.com/disintegration/imaging"
	uuid "github.com/satori/go.uuid"
	"github.com/udemy_fileserver/models"
	"github.com/udemy_fileserver/server"
	"net/http"
)

var files = new(apiFiles)

//fmt.Printf("apiFiles struct: %T",files)
type apiFiles struct{}

func handleFiles(r *server.Router) {
	r.GET("/{name}", files.show)
	r.POST("", files.upload)
	r.GET("", files.list)
}

// FilesResponse first character capitalized for external struct
type FilesResponse struct {
	ID   int               `json:"id"`
	URLs map[string]string `json:"urls"`
}

// ImageResizeParams first character capitalized for external struct
type ImageResizeParams struct {
	Width   int
	Height  int
	Quality int
}

// ImageParams first character capitalized for external struct
var ImageParams = map[string]ImageResizeParams{
	"original": {0, 0, 0},
	"800x800":  {800, 800, 80},
	"400x400":  {400, 400, 80},
	"150x150":  {150, 150, 80},
	"50x50":    {50, 50, 80},
}

func getURLs(dir string) map[string]string {
	urls := map[string]string{}
	for k := range ImageParams {
		urls[k] = fmt.Sprintf("http://127.0.0.1:8080/img/%s/%s.jpg", dir, k)

	}
	return urls
}

func (apiFiles) response(f *models.File) *FilesResponse {
	return &FilesResponse{
		ID:   f.ID,
		URLs: getURLs(f.Name),
	}
}

func (af *apiFiles) multiResponse(fs []*models.File) []*FilesResponse {
	fr := []*FilesResponse{}
	for _, f := range fs {
		resp := af.response(f)
		if resp != nil {
			fr = append(fr, resp)
		}

	}
	return fr
}

// HTTPFile exported type
type HTTPFile interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}

func resizeAndStore(dir string, file HTTPFile, name string, params ImageResizeParams) error {
	if name == "original" {
		file.Seek(0, 0)
		buf := &bytes.Buffer{}
		_, err := buf.ReadFrom(file)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(fmt.Sprintf("assets/img/%s/%s.img", dir, name), buf.Bytes(), 0777)
		if err != nil {
			return err
		}
	} else {
		file.Seek(0, 0)
		img, _, err := image.Decode(file)
		if err != nil {
			return err
		}
		dImg := imaging.Fit(img, params.Width, params.Height, imaging.Lanczos)
		var b []byte
		buf := bytes.NewBuffer(b)
		jpeg.Encode(buf, dImg, &jpeg.Options{Quality: params.Quality})

		err = ioutil.WriteFile(fmt.Sprintf("assets/img/%s/%s.img", dir, name), buf.Bytes(), 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func (af *apiFiles) upload(c *server.Context) {
	httpFile, _, err := c.Request.FormFile("file")
	if err != nil {
		c.RenderError(http.StatusBadRequest, err)
		return
	}
	dirstring, _ := uuid.NewV1()
	dir := dirstring.String()
	dirtree := fmt.Sprintf("assets/img/%s", dir)
	err = os.Mkdir(dirtree, 0777)
	if err != nil {
		c.RenderError(http.StatusBadRequest, err)
		return
	}

	for k, p := range ImageParams {
		err = resizeAndStore(dir, httpFile, k, p)
		if err != nil {
			c.RenderError(http.StatusBadRequest, err)
			return
		}
	}

	file, err := models.Files.Create(dir)
	if err != nil {
		c.RenderError(http.StatusBadRequest, err)
		return
	}

	c.RenderJSON(http.StatusCreated, af.response(file))
}
func (af *apiFiles) show(c *server.Context) {
	file, err := models.Files.ByName(c.Param("name"))
	if err != nil {
		c.RenderError(http.StatusBadRequest, err)
		return
	}
	c.RenderJSON(http.StatusOK, af.response(file))

}
func (af *apiFiles) list(c *server.Context) {
	files, err := models.Files.List()

	if err != nil {
		c.RenderError(http.StatusBadRequest, err)
		return
	}
	c.RenderJSON(http.StatusOK, af.multiResponse(files))

}
