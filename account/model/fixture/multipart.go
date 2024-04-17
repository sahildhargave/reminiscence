package fixture

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"mime/multipart"
	"net/textproto"
	"mime"
)

type MultipartImage struct{
	imagePath  string
	ImageFile  *os.File
	MultipartBody *bytes.Buffer
	ContentType  string
}

func NewMultipartImage(fileName string, contentType string) *MultipartImage{
_,b, _,_ := runtime.Caller(0)
dir:= filepath.Dir(b)

imagePath := filepath.Join(dir, fileName)

f:= createImage(imagePath)
defer f.Close()

body := &bytes.Buffer{}
writer:= multipart.NewWriter(body)

h:= make(textproto.MIMEHeader)
h.Set(
	"Content-Disposition",
	fmt.Sprintf(`form-data;name="%s"; filename="%s"`,"imageFile", fileName),

)
h.Set("Context-Type", contentType)
part,_:= writer.CreatePart(h)

io.Copy(part, f)
writer.Close()

return &MultipartImage{
	imagePath:  imagePath,
	ImageFile:  f,
	MultipartBody: body,
	ContentType: writer.FormDataContentType(),
}

}


func (m *MultipartImage) GetFormFile() *multipart.FileHeader{
	_,params,_ := mime.ParseMediaType(m.ContentType)
	mr:= multipart.NewReader(m.MultipartBody, params["boundry"])
	form,_:= mr.ReadForm(1024)

	files := form.File["imageFile"]

	return files[0]
}


func (m *MultipartImage)Close(){
	m.ImageFile.Close()
	os.Remove(m.imagePath)
}

func createImage(imagePath string) *os.File{
	rect := image.Rect(0,0,1,1)
	img:= image.NewRGBA(rect)

	f,_ := os.Create(imagePath)
	png.Encode(f, img)

	return f
}