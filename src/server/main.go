// Listen to 10389 port for LDAP Request
// and route bind request to the handleBind func
package main

import (
	"net/http"
	"log"
	"os"
	"os/signal"
	"syscall"
	"io/ioutil"
	"net/url"
	"encoding/json"
	"strconv"
	"strings"
	"github.com/vjeantet/goldap/message"
	"flag"

	ldap "github.com/vjeantet/ldapserver"
)

var (
	host string
	port string
	auth_url string
	auth_token string
)
func main() {
	key_host := flag.String("host", "127.0.0.1", "listen host")
	key_port := flag.String("port", "10389", "listen prot")
	key_auth_url := flag.String("auth_url", "https://127.0.0.1/", "request url for remote auth center ")
	key_auth_token := flag.String("auth_token", "token", "request token for remote auth center")
	flag.Parse()

	host = *key_host
	port = *key_port
	auth_url = *key_auth_url
	auth_token = *key_auth_token


	//ldap logger
	ldap.Logger = log.New(os.Stdout, "[server] ", log.LstdFlags)

	//Create a new LDAP Server
	server := ldap.NewServer()

	routes := ldap.NewRouteMux()
	routes.Bind(handleBind).Label("Bind")
	routes.Search(handleSearch).Label("Search - Generic")

	server.Handle(routes)

	// listen on 10389
	go server.ListenAndServe(*host+":"+*port)

	// When CTRL+C, SIGINT and SIGTERM signal occurs
	// Then stop server gracefully
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	close(ch)

	server.Stop()
}

func handleBind(w ldap.ResponseWriter, m *ldap.Message) {
	request := m.GetBindRequest()
	dn := string(request.Name())
	password := string(request.AuthenticationSimple())

	log.Printf("Bind request data DN=%s, Pass=%s", dn, password)

	email := getEmailFromBaseDN(dn)

	if "" == email {
		log.Printf("no email in bind DN=%s, Pass=%s", dn, password)

		result := ldap.NewBindResponse(ldap.LDAPResultInvalidCredentials)
		result.SetDiagnosticMessage("invalid credentials")

		w.Write(result)
		return
	}

	authResult := authRequest(email, password)

	if authResult {
		result := ldap.NewBindResponse(ldap.LDAPResultSuccess)

		w.Write(result)
	} else  {
		log.Printf("Bind failed User=%s, Pass=%s", dn, password)

		result := ldap.NewBindResponse(ldap.LDAPResultInvalidCredentials)
		result.SetDiagnosticMessage("invalid credentials")

		w.Write(result)
	}
}

func handleSearch(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetSearchRequest()

	log.Printf("Request BaseDn=%s", r.BaseObject())
	log.Printf("Request Filter=%s", r.Filter())
	log.Printf("Request FilterString=%s", r.FilterString())
	log.Printf("Request Attributes=%s", r.Attributes())
	log.Printf("Request TimeLimit=%d", r.TimeLimit().Int())

	select {
	case <-m.Done:
		log.Print("Leaving handleSearch...")
		return
	default:
	}

	//通过解析filterString获取email
	filterStr := r.FilterString()
	filterSplit := strings.Split(filterStr, "=")
	filterSplit = strings.Split(filterSplit[1], ")")
	email := filterSplit[0]

	//配置的管理员账号的search请求
	if "*" == email {
		dn := string(r.BaseObject())
		email = getEmailFromBaseDN(dn)

		if "" == email {
			log.Printf("no email in search BaseDN=%s", dn)

			result := ldap.NewSearchResultDoneResponse(ldap.LDAPResultNoSuchObject)
			w.Write(result)

			return
		}
	}

	userInfo := queryUserInfoRequest(email)

	if nil == userInfo {
		res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultNoSuchObject)
		w.Write(res)

		return
	}

	uid := message.AttributeValue(userInfo["uid"].(string))
	mail := message.AttributeValue(email)

	e := ldap.NewSearchResultEntry("CN=" + email + "," + string(r.BaseObject()))
	e.AddAttribute("uid", uid)
	e.AddAttribute("email", mail)

	e.AddAttribute("sAMAccountName", mail)
	e.AddAttribute("mail", mail)
	e.AddAttribute("cn", mail)
	w.Write(e)

	res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
	w.Write(res)

}

func authRequest(email string, password string) bool {
	authUrl := auth_url+"/auth"
	params := url.Values{
		"email": {email},
		"password": {password},
		"is_otp" : {"1"},
	}
	resp, err := http.PostForm(authUrl, params)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer resp.Body.Close()

	respJson, _ := ioutil.ReadAll(resp.Body)

	var authResult map[string]interface{}
	_ = json.Unmarshal(respJson, &authResult)

	var errCode = strconv.FormatFloat(authResult["retcode"].(float64), 'f', -1, 64)

	return ("2000000" == errCode)
}

func queryUserInfoRequest(email string) map[string]interface{} {
	authUrl := auth_url+"/user?token="+auth_token
	params := url.Values{
		"email": {email},
	}
	resp, err := http.PostForm(authUrl, params)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer resp.Body.Close()

	respStr, _ := ioutil.ReadAll(resp.Body)

	var respData map[string]interface{}
	_ = json.Unmarshal(respStr, &respData)

	var errCode = strconv.FormatFloat(respData["retcode"].(float64), 'f', -1, 64)
	if "2000000" == errCode {
		return respData["user_info"].(map[string]interface{})
	} else {
		return nil
	}
}

// 从DN中获取email信息
// 默认CN项目为email
func getEmailFromBaseDN(baseDN string) string  {
	baseDNSplit := strings.Split(string(baseDN), ",")

	email := ""
	for _, item := range baseDNSplit {
		itemSplit := strings.Split(item, "=")

		if "CN" == itemSplit[0] {
			email = itemSplit[1]
			break
		}
	}

	return email
}