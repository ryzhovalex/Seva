package domains

import (
	"os"
	"seva/lib/rpc"
	"seva/lib/utils"

	"github.com/gin-gonic/gin"
)

func GetAll() ([]string, *utils.Error) {
	dir := "Var/Domains"

	files, be := os.ReadDir(dir)
	if be != nil {
		return nil, utils.BE(be)
	}

	r := []string{}
	for _, file := range files {
		if file.IsDir() {
			r = append(r, file.Name())
		}
	}
	return r, nil

}

func CheckDomainNotCreated(domain string) *utils.Error {
	e := CheckDomainCreated(domain)
	if e == nil {
		return utils.DE("Domain already registered: " + domain)
	}
	return nil
}

func CheckDomainCreated(domain string) *utils.Error {
	registered, e := IsDomainCreated(domain)
	if e != nil {
		return e
	}
	if !registered {
		return utils.DE("Domain not registered: " + domain)
	}
	return nil
}

func IsDomainCreated(domain string) (bool, *utils.Error) {
	domains, e := GetAll()
	if e != nil {
		return false, e
	}
	for _, d := range domains {
		if domain == d {
			return true, nil
		}
	}
	return false, nil
}

func GetDomainDir(domain string) string {
	return "Var/Domains/" + domain
}

func createDomain(domain string) *utils.Error {
	if domain == "" {
		return utils.DE("Domain name is empty.")
	}

	e := CheckDomainNotCreated(domain)
	if e != nil {
		return e
	}

	be := os.Mkdir("Var/Domains/"+domain, 0755)
	if be != nil {
		return utils.BE(be)
	}

	return nil
}

type CreateDomainData struct {
	Domain string
}

func RpcCreateDomain(c *gin.Context) {
	var data CreateDomainData
	be := c.Bind(&data)
	if be != nil {
		rpc.Error(c, utils.BE(be))
		return
	}

	e := createDomain(data.Domain)
	if e != nil {
		rpc.Error(c, e)
		return
	}
	rpc.Ok(c, nil)
}

func RpcGetDomains(c *gin.Context) {
	r, e := GetAll()
	if e != nil {
		rpc.Error(c, e)
	}
	rpc.Ok(c, r)
}
