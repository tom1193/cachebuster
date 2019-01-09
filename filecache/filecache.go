package filecache

import (
	"fmt"
	"strings"
	"net/http"
	"github.com/tom1193/cachebuster/utils"
	"github.com/tom1193/cachebuster/proto"
)

const EnvError = "Invalid request, env must be 'dev' or 'prod'"

var DevFileCache = []string{}
var ProdFileCache = []string{}

func init() {
	//use AWS GET Bukcet (List Objects) API to populate file cache with current files from
	//https://s3.console.aws.amazon.com/s3/buckets/ph-mode-static/datavis-library/?region=us-east-1&tab=overview
}

func ReturnFileCacheEnv(env string) *[]string {
	if env == "dev" {
		return &DevFileCache
	} else if env == "prod" {
		return &ProdFileCache
	} else {
		return nil
	}
}

//init FileCache to have all files in the cloud
func UpdateFileCache(pr proto.PostRequest) (map[string]interface{}, int) {
	var fc = ReturnFileCacheEnv(pr.Env)
	if fc != nil {
		if pr.Filecache.Names != nil {
				*fc = pr.Filecache.Names
				fmt.Println(DevFileCache, ProdFileCache)
				return nil, http.StatusCreated
			} else {
				return utils.Message(false, "Invalid request, POST at least one file."), http.StatusBadRequest
			}
	} else {
		return utils.Message(false, EnvError), http.StatusBadRequest
	}
}

//receives file prefix and return full names of matching files
func RequestFileCache(filenames []string, env string) (map[string]interface{}, int) {
	var fc = ReturnFileCacheEnv(env)
	if fc != nil {
		if filenames != nil {
			var responseNames []string
			for i := 0; i < len(filenames); i++ {
				name := filenames[i]
				for j := 0; j<len(*fc); j++ {
					fullname := (*fc)[j] //https://flaviocopes.com/golang-does-not-support-indexing/
					prefix := fullname[:strings.IndexByte(fullname, '.')]
					if name == prefix {
						responseNames = append(responseNames, fullname)
					}
				}
			}
			res := utils.Message(true, "Returned matching files")
			res["filecache"] = proto.Filecache{responseNames}
			return res, http.StatusOK
		} else {
			return utils.Message(false, "Invalid request, GET at least one file."), http.StatusBadRequest
		}
	} else {
		return utils.Message(false, EnvError), http.StatusBadRequest
	}
}

func EchoFileCache(env string) (map[string]interface{}, int) {
	var fc = ReturnFileCacheEnv(env)
	if fc != nil {
		res := utils.Message(true, "Echoing file cache")
		res["filecache"] = proto.Filecache{*fc}
		return res, http.StatusOK
	} else {
		return utils.Message(false, EnvError), http.StatusBadRequest
	}

}
