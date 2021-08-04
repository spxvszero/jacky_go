# jacky_go
a go project for easy setup my server, not open at all.

## Quick

```
#linux
$ curl -O -L https://raw.githubusercontent.com/spxvszero/jacky_go/master/go_server.sh && chmod +x go_server.sh
$ sh go_server.sh

```


## Usage

* build `go_build` file
* run `./go_build --config=xxx.json` or run `./go_build` if json file is in same path.

```
//json config format, use it with no descriptions in this file.
  {
	"port":10000,
	"Download_config":{
		"download_dir_info_url_path":"/download/json",
		"download_dir_path":"/Users/username/Desktop/goTest,/Users/username/Desktop/test",
		"use_default_page": true,
		"page_url_path": "/download",
		"page_file_path": ""
	},
	"Upload_config":{
		"upload_url_path":"/upload/interface",
		"save_dir_path":"/Users/username/Desktop/goTest",
		"max_size":1000,	//size MB
		"use_default_page": true,
		"page_url_path": "/upload",
		"page_file_path": ""
	},
	//this routes is simple http request
	"routes":[
		{
			"method":"get",
			"path":"/json",
			"json_body":"{\"oh\":\"no!!\"}"
		}
	],
	//quick build sock5 server , auth is not required
	"socks5":{
		"protocol": "tcp",
		"addr": "127.0.0.1:12121",
		"auth": {
			"usr":"pwd"
		}
	}
}
```

## Example

You can see this in my temp link : No More.

WebPage is Embed in go project, source in : [WebSouce](https://github.com/spxvszero/jacky_go_web_source).

