package handler

import (
	"encoding/json"

	"code.cloudfoundry.org/lager"
	"github.com/gin-gonic/gin"

	"io/ioutil"
	"net/http"
	"os"
)

const JSON = "application/json"

var log lager.Logger
var ssoHost, dfHost string

func init() {
	//dfHost = getENV("DFHOST")
	log = lager.NewLogger("Token_Proxy")
	log.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
}

type Users struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Cluster(c *gin.Context) {

	//token := getToken(c)
	rBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error("CreateBC Read Request.Body error", err)
	}
	var user Users
	_ = json.Unmarshal(rBody, &user)
	//fmt.Printf("%#v\n",user)

	//users :=c.Request.Header["Authorization"][0]

	errorRsp := ErrorResponse{}

	result := ""

	dfRsp, err := trRequest("GET", "https://"+dfHost+"/oauth/authorize?client_id=openshift-challenging-client&response_type=token", user)

	if err != nil {
		log.Error("GetToken request fail error", err)
		errorRsp.Description = "GetToken request fail error"
		errorRsp.Error = err.Error()
		c.JSON(500, errorRsp)
		return
	}

	if value, ok := dfRsp.Header["Location"]; ok {
		result = getDFtoken(value[0])
	} else {
		log.Debug("GetToken respose header fail error")
		errorRsp.Description = "GetToken respose header fail error"
		errorRsp.Error = err.Error()
		c.JSON(500, errorRsp)
		return
	}

	if len(result) == 0 {
		log.Debug("GetToken respose header[location]  fail error")
		errorRsp.Description = "GetToken respose header[location]  fail error"
		errorRsp.Error = err.Error()
		c.JSON(500, errorRsp)
		return
	}

	rsp := RstResponse{}
	rsp.User = "admin"
	rsp.Token = result

	c.JSON(http.StatusOK, rsp)
	return

}
