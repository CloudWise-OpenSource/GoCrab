# httplib
httplib is an libs help you to curl remote url.

# How to use?

## GET
you can use Get to crawl data.

	import "github.com/CloudWise-OpenSource/GoCrab/Libs/httplib"
	
	str, err := httplib.Get("http://GoCrab.me/").String()
	if err != nil {
        	// error
	}
	fmt.Println(str)
	
## POST
POST data to remote url

	req := httplib.Post("http://GoCrab.me/")
	req.Param("username","Neeke")
	req.Param("password","123456")
	str, err := req.String()
	if err != nil {
        	// error
	}
	fmt.Println(str)

## Set timeout

The default timeout is `60` seconds, function prototype:

	SetTimeout(connectTimeout, readWriteTimeout time.Duration)

Exmaple:

	// GET
	httplib.Get("http://GoCrab.me/").SetTimeout(100 * time.Second, 30 * time.Second)
	
	// POST
	httplib.Post("http://GoCrab.me/").SetTimeout(100 * time.Second, 30 * time.Second)


## Debug

If you want to debug the request info, set the debug on

	httplib.Get("http://GoCrab.me/").Debug(true)
	
## Set HTTP Basic Auth

	str, err := Get("http://GoCrab.me/").SetBasicAuth("user", "passwd").String()
	if err != nil {
        	// error
	}
	fmt.Println(str)
	
## Set HTTPS

If request url is https, You can set the client support TSL:

	httplib.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	
More info about the `tls.Config` please visit http://golang.org/pkg/crypto/tls/#Config	

## Set HTTP Version

some servers need to specify the protocol version of HTTP

	httplib.Get("http://GoCrab.me/").SetProtocolVersion("HTTP/1.1")
	
## Set Cookie

some http request need setcookie. So set it like this:

	cookie := &http.Cookie{}
	cookie.Name = "username"
	cookie.Value  = "Neeke"
	httplib.Get("http://GoCrab.me/").SetCookie(cookie)

## Upload file

httplib support mutil file upload, use `req.PostFile()`

	req := httplib.Post("http://GoCrab.me/")
	req.Param("username","Neeke")
	req.PostFile("uploadfile1", "httplib.pdf")
	str, err := req.String()
	if err != nil {
        	// error
	}
	fmt.Println(str)


See godoc for further documentation and examples.

* [godoc.org/github.com/CloudWise-OpenSource/GoCrab/Core/httplib](https://godoc.org/github.com/CloudWise-OpenSource/GoCrab/Core/httplib)
