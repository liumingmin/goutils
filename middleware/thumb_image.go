package middleware

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"goutils/utils"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

type HttpFs interface {
	http.FileSystem
	Exists(prefix string, path string) (bool, string)
}

type localHttpFs struct {
	http.FileSystem
	root    string
	indexes bool
}

func GinHttpFs(root string, indexes bool) *localHttpFs {
	return &localHttpFs{
		FileSystem: gin.Dir(root, indexes),
		root:       root,
		indexes:    indexes,
	}
}

func (l *localHttpFs) Exists(prefix string, filepath string) (bool, string) {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		name := path.Join(l.root, p)
		stats, err := os.Stat(name)
		if err != nil {
			return false, name
		}
		if !l.indexes && stats.IsDir() {
			return false, name
		}
		return true, name
	}
	return false, ""
}

func ThumbImageServeRoot(urlPrefix, root string) gin.HandlerFunc {
	return ThumbImageServe(urlPrefix, GinHttpFs(root, false))
}

func ThumbImageServe(urlPrefix string, fs HttpFs) gin.HandlerFunc {
	fileserver := http.FileServer(fs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *gin.Context) {

		if isok, realpath := fs.Exists(urlPrefix, c.Request.URL.Path); isok {
			wstr := c.DefaultQuery("w", "")
			hstr := c.DefaultQuery("h", "")

			if wstr != "" && hstr != "" {
				w, _ := strconv.Atoi(wstr)
				h, _ := strconv.Atoi(hstr)

				thumbnailPath := realpath
				thumbnailUrl := c.Request.URL.Path
				idx := strings.LastIndex(realpath, ".")
				if idx >= 0 {
					thumbnailPath = fmt.Sprintf("%s-%sx%s%s", realpath[0:idx], wstr, hstr, realpath[idx:])
				}

				origPath := c.Request.URL.Path
				origIdx := strings.LastIndex(origPath, ".")
				if origIdx >= 0 {
					thumbnailUrl = fmt.Sprintf("%s-%sx%s%s", origPath[0:origIdx], wstr, hstr, origPath[origIdx:])
				}

				if !utils.FileExist(thumbnailPath) {
					img, err2 := utils.ReadImage(realpath)
					if err2 == nil {
						img2 := resize.Thumbnail(uint(w), uint(h), img, resize.NearestNeighbor)
						utils.WriteImage(img2, thumbnailPath)
					}
				}

				if utils.FileExist(thumbnailPath) {
					c.Request.URL.Path = thumbnailUrl
				}
			}
			//fmt.Println(c.Request.URL.Path)
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}
