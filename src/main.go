package main

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"io/ioutil"
	"jacky_go/src/go-socks5"
	"jacky_go/src/jk_source_page"
	"jacky_go/src/jk_utils"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//var baiduImg = "https://www.baidu.com/img/bd_logo1.png"
//var baiduapk = "http://p.gdown.baidu.com/01d3fbc183d1f6bf4fdbf830205af25147066776d515ad33141cf0cf4368df61d12c8d5e8a80bdd7773ecb8727f4050100e8a896fbed3192e5f43924b6962016fbf63e0a7615f9b262bd5625dfe35403aa54d6bef215b49be09bc83e1465545558cda5da5a384fe7455869d06e59607c"

//command paramater
var configFile = flag.String("config","","config file for port,path...etc")

var (
	uploadDir = ""
	uploadDirMaxSize = int64(1024 * 1024 * 1024) //1gb

	multiDownloadDir = []string{}

	logFile = "go_log.txt"
)


type ConfigFileStruct struct {
	Port 				int				`json:"port"`
	Download_config 	DownloadRoute 	`json:"download_config"`
	Upload_config 		UploadRoute		`json:"upload_config"`
	Routes 				[]EasyRoute		`json:"routes"`
	Socks5				Sock5Config		`json:"socks5"`
}

type Sock5Config struct {
	Protocol	string					`json:"protocol"`
	Addr		string					`json:"addr"`
	Auth		map[string]string		`json:"auth"`
}

type EasyRoute struct {
	Method string	 `json:"method"`
	Path string		 `json:"path"`
	Json_body string `json:"json_body"`
}

type FileRoute struct {
	Use_Default_Page bool   `json:"use_default_page"`
	Page_URL_Path    string `json:"page_url_path"`
	Page_File_Path   string `json:"page_file_path"`
}

type UploadRoute struct {
	Upload_URL_Path string `json:"upload_url_path"`
	Save_Dir_path   string `json:"save_dir_path"`
	Max_Size        int64  `json:"max_size"`
	FileRoute
}

type DownloadRoute struct {
	Download_Dir_Info_URL_Path string `json:"download_dir_info_url_path"`
	Download_Dir_path          string `json:"download_dir_path"`
	FileRoute
}

type DownloadFileStruct struct {
	Name string
	Size int64
	Mode  uint32
	ModTime int64
	IsDir bool
	DirPath string
	SubFiles []*DownloadFileStruct
}


var UserConfig = ConfigFileStruct{}

func main() {
	//InitSqliteDB("test.db",[]DBModel{new(DownloadShared),new(justATest)})
	//time.Sleep(2)


	//for i:=0 ;i < 10 ; i++{
	//	a := DownloadShared{0,int64(i),"123","321",123,"!23"}
	//	a.Save()
	//	fmt.Println("res - ",a)
	//}

	//res := &DownloadShared{}
	//fmt.Println("Out Model : ",res)
	//Get(res,"file_path = '123';")
	//fmt.Println("select result : ",res)

	//res := []*DownloadShared{}
	//fmt.Println("Out Model : ",res)
	//Get(&res,"file_path = '123' order by id desc;")
	//fmt.Println("select result : ",res)
	//for _,v := range res {
	//	fmt.Println("Final :",v);
	//}

	//jk_utils.Thumbnails()
	configConfig()
	go sockBuild()
	ginRoute()
	//ReadDirTree("/Users/jason/Desktop/goTest"
}

func sockBuild()  {

	if &UserConfig != nil && &UserConfig.Socks5 != nil && len(UserConfig.Socks5.Addr)>0 {
		// Create a SOCKS5 server
		var conf *socks5.Config
		if (len(UserConfig.Socks5.Auth) > 0) {
			conf = &socks5.Config{Credentials:socks5.StaticCredentials(UserConfig.Socks5.Auth),Logger:jk_utils.Trace}
		}else {
			conf = &socks5.Config{Logger:jk_utils.Trace}
		}

		//conf := &socks5.Config{}
		server, err := socks5.New(conf)
		if err != nil {
			jk_utils.Error.Println("Sock5 Error : ",err)
			return
		}

		// Create SOCKS5 proxy on localhost port 8000
		jk_utils.Info.Println("Sock5 Open in : ",UserConfig.Socks5.Addr + " With " + UserConfig.Socks5.Protocol)
		if err := server.ListenAndServe(UserConfig.Socks5.Protocol, UserConfig.Socks5.Addr); err != nil {
			jk_utils.Error.Println("Sock5 Error : ",err)
			return
		}
	}
}


func ginRoute()  {

	gin.SetMode(gin.ReleaseMode)

	route := gin.Default()

	//route.Use(CORSMiddleware())
	//pprof.Register(route,"debug/pprof")

	port := 8888
	configRoute(route)

	if &UserConfig != nil && UserConfig.Port > 0 {
		port = UserConfig.Port
	}
	portStr := ":"+strconv.Itoa(port)
	route.Run(portStr)
}

func configConfig()  {
	flag.Parse()
	logConfig()
	readConfig()
}

func logConfig()  {
	f,err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		jk_utils.Error.Println("Open Log File Failed : ",err);
		return;
	}
	w := io.MultiWriter(os.Stdout,f)
	gin.DefaultWriter = w
	jk_utils.SetWriter(w)

	jk_utils.Info.Println("Log File Path : ",filepath.Dir(os.Args[0]) + "/" + logFile)
}

func readConfig() {
	if len(*configFile) > 0 {
		jk_utils.Info.Println("Read Route Config File From Input : ",*configFile)
		UserConfig = *readConfigFile(*configFile)
	}else {
		filePossiblePath := "config.json"
		jk_utils.Info.Println("Read Route Config File From Surmise : ",filePossiblePath)
		UserConfig = *readConfigFile(filePossiblePath)
		if &UserConfig == nil {
			filePossiblePath = "config/config.json"
			jk_utils.Info.Println("Read Route Config File From Surmise : ",filePossiblePath)
			UserConfig = *readConfigFile(filePossiblePath)
		}
	}
}


func configRoute(route *gin.Engine){

	//play routes config
	if &UserConfig == nil {
		jk_utils.Error.Println("Read Route Config Failed")
		return;
	}

	jk_utils.Info.Println("Read Route Config Success")

	for _,usrRoute := range UserConfig.Routes{
		lowMethod := strings.ToLower(usrRoute.Method)
		if lowMethod == "get" {
			if len(usrRoute.Path) > 0 {
				route.GET(usrRoute.Path, func(c *gin.Context) {
					c.JSON(http.StatusOK,gin.H{
						"message":usrRoute.Json_body,
					})
				})
				continue
			}
		}

		if lowMethod == "post"{
			if len(usrRoute.Path) > 0 {
				route.POST(usrRoute.Path, func(c *gin.Context) {
					c.JSON(http.StatusOK,gin.H{
						"message":usrRoute.Json_body,
					})
				})
				continue
			}
		}

		jk_utils.Error.Println("Build Route From Config Error: method not support : ",lowMethod)
	}


	//upload route
	if &UserConfig.Upload_config != nil && len(UserConfig.Upload_config.Upload_URL_Path)>0 {
		jk_utils.Info.Println("Build Upload Route From Config >>> ")

		uploadPath := "/files/upload"
		uploadPath = UserConfig.Upload_config.Upload_URL_Path
		uploadDir,_ = filepath.Split(os.Args[0])

		route.MaxMultipartMemory = 32 << 20

		if len(UserConfig.Upload_config.Save_Dir_path) > 0 {
			uploadDir = UserConfig.Upload_config.Save_Dir_path
		}
		jk_utils.Info.Println("Upload Dir Path :",uploadDir)
		if UserConfig.Upload_config.Max_Size > 0 {
			uploadDirMaxSize = UserConfig.Upload_config.Max_Size * 1024 * 1024
		}
		jk_utils.Info.Println("Upload Dir Max Size Set To :",uploadDirMaxSize)

		route.POST(uploadPath, func(c *gin.Context) {

			dirSize, err := jk_utils.DirSize(uploadDir)

			var failedErr error
			var exist bool

			if err != nil {
				jk_utils.Info.Println("Upload Module : File Size Limited Error : ",err)
				failedErr = err
				goto FailedUpload
			}

			if c.ContentType() == "multipart/form-data" {
				exist,failedErr = UploadSaveAsForm(c,dirSize)

				if exist {
					if failedErr != nil {
						goto FailedUpload
					}else {
						return
					}
				}

				exist,failedErr = UploadSaveAsMultiForm(c,dirSize)

				if exist {
					if failedErr != nil {
						goto FailedUpload
					}else {
						return
					}
				}
			}

			if c.ContentType() == "application/octet-stream" {
				
			}

			FailedUpload:

				if failedErr == nil {
					failedErr = errors.New("counld not excute file upload format")
				}
				c.String(http.StatusExpectationFailed,failedErr.Error());

		})


		//upload web page path
		if len(UserConfig.Upload_config.Page_URL_Path) > 0 {

			if UserConfig.Upload_config.Use_Default_Page {

				jk_utils.Info.Println("Build Default Upload WebPage Route")

				route.GET(UserConfig.Upload_config.Page_URL_Path, func(c *gin.Context) {

					uploadUrlPath := struct {
						Path string
					}{UserConfig.Upload_config.Upload_URL_Path[1:]}

					tmpl, err := template.New("upload_page").Parse(jk_source_page.Upload_HTML)
					if err != nil {
						jk_utils.Error.Println("Upload Web Page Tmpl Err ", err);
						c.String(http.StatusNotFound,"this page is err")
					}else {
						c.Status(200)
						tmpl.Execute(c.Writer, uploadUrlPath)
					}
				})


			}else{
				jk_utils.Info.Println("Build Custom Upload WebPage Route From Config : ",UserConfig.Upload_config.Page_File_Path)

				route.LoadHTMLGlob(UserConfig.Upload_config.Page_File_Path)

				route.GET(UserConfig.Upload_config.Page_URL_Path, func(c *gin.Context) {

					var err error

					if len(UserConfig.Upload_config.Page_File_Path) > 0 {
						_,err = os.Stat(UserConfig.Upload_config.Page_File_Path)
					}else {
						err = errors.New("Page_File_Path not found")
					}

					if err != nil {
						c.String(http.StatusExpectationFailed,"path ",UserConfig.Upload_config.Page_URL_Path, "error: ",err)
					}else {
						c.HTML(http.StatusOK,filepath.Base(UserConfig.Upload_config.Page_File_Path),nil)
					}
				})
			}

		}

		jk_utils.Info.Println("Build Upload Route From Config <<< ")
	}

	//download route
	if &UserConfig.Download_config != nil && len(UserConfig.Download_config.Download_Dir_Info_URL_Path)>0 {

		jk_utils.Info.Println("Build Download Route From Config >>> ")

		downloadUrl := UserConfig.Download_config.Download_Dir_Info_URL_Path

		multiDownloadDir = strings.Split(UserConfig.Download_config.Download_Dir_path,",")

		//json path
		route.GET(downloadUrl, func(c *gin.Context) {

			downloadPath := strings.TrimSpace(c.Query("filePath"))

			if len(downloadPath) > 0 {
				if DownloadPathEnable(downloadPath) {

					fileStat,err := os.Stat(downloadPath)

					if err != nil {
						c.JSON(http.StatusBadRequest,"bad path")
						return
					}

					if fileStat.IsDir() {
						//need zip
						tmpPath := jk_utils.AppendPath(os.TempDir(),"/goserver/"+filepath.Base(downloadPath))
						//tmpPath := UserConfig.Download_config.Download_Dir_path + "/goserver/"+filepath.Base(downloadPath)
						_,tmpPathErr := os.Stat(tmpPath)
						if tmpPathErr != nil {
							jk_utils.Info.Println("Download Tmp Path Not Found, Create Path : ",tmpPath)
							os.MkdirAll(tmpPath,os.ModePerm)
						}

						tmpFilePath := tmpPath+"/"+strconv.FormatInt(fileStat.ModTime().Unix(),10)
						hasFile := false
						filepath.Walk(tmpPath, func(innerPath string, info os.FileInfo, err error) error {
							if err != nil {
								return err
							}
							if info.IsDir() {

							}else {
								if info.Name() == filepath.Base(tmpFilePath) {
									//os.RemoveAll(innerPath)
									hasFile = true
								}else {
									jk_utils.Info.Println("Download Remove Expired Tmp File :",innerPath)
									os.Remove(innerPath)
								}
							}

							return nil
						})

						var err error
						if hasFile {
							jk_utils.Info.Println("Download File Already Exist ",tmpFilePath," No Need Zip")
						}else {
							jk_utils.Info.Println("Download File ",downloadPath," --> Tmp Zip Pack :",tmpFilePath)
							err = jk_utils.Zip(downloadPath,tmpFilePath)
						}

						if err == nil {
							c.FileAttachment(tmpFilePath,filepath.Base(downloadPath)+".zip")
						}else{
							jk_utils.Error.Println("Zip ",downloadPath," Error : ",err)
							c.Status(http.StatusExpectationFailed)
						}
					}else{
						c.FileAttachment(downloadPath,filepath.Base(downloadPath))
					}
				}else{
					jk_utils.Error.Println("Client Request Unauthorize Path : ",downloadPath)
					//unauthorize path
					c.JSON(http.StatusBadRequest,"bad path")
				}
			}else {

				fileTree := []*DownloadFileStruct{}

				for _,singleDirPath := range multiDownloadDir {
					singleTree,err := ReadDirTree(singleDirPath)
					if err != nil {
						jk_utils.Error.Println("Fetch Download Dir ",singleDirPath," Error : ",err)
					}else {
						fileTree = append(fileTree, singleTree)
					}
				}

				if len(fileTree) <= 0 {
					c.JSON(http.StatusOK,gin.H{
						"message":"No Download Dir fetch",
					})
				}else {
					c.JSON(http.StatusOK,fileTree)
				}
			}
		})

		//download web page path
		if len(UserConfig.Download_config.Page_URL_Path) > 0 {

			if UserConfig.Download_config.Use_Default_Page {

				jk_utils.Info.Println("Build Default Download WebPage Route")

				route.GET(UserConfig.Download_config.Page_URL_Path, func(c *gin.Context) {

					downloadInfoPath := struct {
						JSON_Path string
					}{UserConfig.Download_config.Download_Dir_Info_URL_Path[1:]}

					tmpl, err := template.New("download_page").Parse(jk_source_page.DownloadPageHTML)
					if err != nil {
						jk_utils.Error.Println("Download Web Page Tmpl Err ", err);
						c.String(http.StatusNotFound,"something is wrong in this page")
					}else {
						c.Status(200)
						tmpl.Execute(c.Writer,downloadInfoPath)
					}
				})


			}else{

				jk_utils.Info.Println("Build Custom Download WebPage Route From Config ",UserConfig.Download_config.Page_File_Path)

				route.LoadHTMLGlob(UserConfig.Download_config.Page_File_Path)

				route.GET(UserConfig.Download_config.Page_URL_Path, func(c *gin.Context) {

					var err error

					if len(UserConfig.Download_config.Page_File_Path) > 0 {
						_,err = os.Stat(UserConfig.Download_config.Page_File_Path)
					}else {
						err = errors.New("Page_File_Path not found")
					}

					if err != nil {
						c.String(http.StatusExpectationFailed,"path ",UserConfig.Download_config.Page_URL_Path, "error: ",err)
					}else {
						c.HTML(http.StatusOK,filepath.Base(UserConfig.Download_config.Page_File_Path),nil)
					}
				})
			}

		}


		jk_utils.Info.Println("Build Download Route From Config <<< ")
	}


}
///curl -X POST http://localhost:8900/upload -F "file=@/Users/jason/Desktop/thanos_snap.png" -H "Content-Type: multipart/form-data"
///curl -X POST http://www.lifefordebug.com:8900/upload -F "file=@/Users/jason/Desktop/thanos_snap.png" -H "Content-Type: multipart/form-data"

func UploadSaveAsForm(c *gin.Context,dirSize int64) (exist bool,err error){
	//is form file?
	file,fileExistErr := c.FormFile("file")

	fileRelativePath ,_ := c.GetPostForm("relativePath");

	if fileExistErr != nil {
		return false,fileExistErr
	}

	if dirSize + file.Size > uploadDirMaxSize {
		return true, errors.New("file upload over dir limit")
	}

	if len(fileRelativePath) > 0{
		dirpath,fileName := filepath.Split(fileRelativePath);
		if len(dirpath) > 1 {
			dirpath = jk_utils.AppendPath(uploadDir,dirpath);
			_,dirExistErr := os.Stat(dirpath);
			if dirExistErr != nil {
				jk_utils.Info.Println("Upload mkdir : ",dirpath);
				os.MkdirAll(dirpath,os.ModePerm);
			}
		}
		jk_utils.Info.Println("Upload <",fileName,"> to relative path: ",dirpath);
	}else{
		fileRelativePath = file.Filename
		jk_utils.Info.Println("Upload <",fileRelativePath,"> to path: ",uploadDir);
	}

	uploadFilePath := jk_utils.AppendPath(uploadDir ,fileRelativePath)
	c.SaveUploadedFile(file,uploadFilePath)

	if err != nil {
		return true,err
	}
	jk_utils.Info.Printf("'%s' upload! \n",file.Filename)
	c.String(http.StatusOK,"'%s' upload! \n",file.Filename)
	return true,nil
}

func UploadSaveAsMultiForm(c *gin.Context,dirSize int64) (exist bool,err error) {
	//is Multipart form ?
	form, formFileExistErr := c.MultipartForm()

	if formFileExistErr != nil {
		return false, formFileExistErr
	}

	files := form.File["upload[]"]

	if dirSize + c.Request.ContentLength > uploadDirMaxSize {
		return true, errors.New("file upload over dir limit")
	}

	for _,file := range files {
		uploadFilePath := jk_utils.AppendPath(uploadDir ,file.Filename)
		log.Println(uploadFilePath)
		c.SaveUploadedFile(file,uploadFilePath)
	}
	jk_utils.Info.Printf("multi upload! \n")
	c.String(http.StatusOK,"multi upload! \n")
	return true,nil
}


func DownloadPathEnable(path string) bool {
	for _,enablePath := range multiDownloadDir {

		if strings.HasPrefix(path,enablePath) {
			return true
		}

	}
	return false
}



func readFileInfo(filepath string) (os.FileInfo,error)  {
	fileStat, err := os.Stat(filepath)

	if err != nil {
		return nil,err
	}
	//fmt.Println("File Name:", fileStat.Name())        // Base name of the file
	//fmt.Println("Size:", fileStat.Size())             // Length in bytes for regular files
	//fmt.Println("Permissions:", fileStat.Mode())      // File mode bits
	//fmt.Println("Last Modified:", fileStat.ModTime()) // Last modification time
	//fmt.Println("Is Directory: ", fileStat.IsDir())   // Abbreviation for Mode().IsDir()

	return fileStat,nil
}


func readConfigFile(filePath string) *ConfigFileStruct {

	res := ConfigFileStruct{}

	fileBlob,fileErr := ioutil.ReadFile(filePath)

	if fileErr != nil {
		jk_utils.Error.Println("Config Read Error : ",fileErr)
		return nil
	}

	json.Unmarshal(fileBlob,&res)

	return &res
}


func ReadDirTree(path string) (*DownloadFileStruct,error) {

	tmpMap := make(map[string]*DownloadFileStruct)

	err := filepath.Walk(path, func(curPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		//check parent path
		parentDir := filepath.Dir(curPath)
		parentStruct := tmpMap[parentDir]

		sonStruct := tmpMap[curPath]
		if sonStruct == nil {
			sonStruct = &DownloadFileStruct{info.Name(),info.Size(),uint32(info.Mode()),info.ModTime().Unix(),info.IsDir(),parentDir,nil}
			tmpMap[curPath] = sonStruct
		}

		if parentStruct != nil {
			if parentStruct.SubFiles == nil {
				parentStruct.SubFiles = []*DownloadFileStruct{}
			}
			parentStruct.SubFiles = append(parentStruct.SubFiles,sonStruct)

		}

		//fmt.Println("curpath -- ",curPath," ---> name -- ",info.Name()," ---> base -- ",parentDir)
		//fmt.Printf("curStruct : %p ",sonStruct,sonStruct,"\n")
		//fmt.Printf("parentStruct : %p ",parentStruct,parentStruct,"\n")
		//fmt.Println()

		return err
	})

	//jsonData, _ := json.Marshal(tmpMap[path])
	//fmt.Println("map:\n ",tmpMap)
	//fmt.Println()
	//fmt.Println("struct:\n",tmpMap[path])
	//fmt.Println()
	//fmt.Println("result:\n ",string(jsonData))
	//fmt.Println()

	return tmpMap[path],err
}


//test local cors middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

//func download(){
	//manager := jk_downloader.BuildManager("")
	//task1,_ := jk_downloader.BuildTask(baiduImg,"a.jpg",false, func(task *jk_downloader.TaskField, status jk_downloader.TaskStauts) {
	//	fmt.Println("task ", task.Identify, " outter get status",status)
	//}, func(task *jk_downloader.TaskField, progress float32, expectedSize int64) {
	//	fmt.Println("task ", task.Identify, " outter get progress :",progress)
	//})
	//
	//jk_downloader.AddTaskToManager(manager,task1)
	//
	//
	//task2,_ := jk_downloader.BuildTask(baiduapk,"baidu.apk",false, func(task *jk_downloader.TaskField, status jk_downloader.TaskStauts) {
	//	fmt.Println("task ", task.Identify, " outter get status",status)
	//}, func(task *jk_downloader.TaskField, progress float32, expectedSize int64) {
	//	fmt.Println("task ", task.Identify, " outter get progress :",progress)
	//})
	//
	//jk_downloader.AddTaskToManager(manager,task2)
	//
	//for  {
	//	time.Sleep(1*time.Second)
	//}
//}

//func buildDownload(doer *apkDownloader)  {
//
//
//	task,err := jk_downloader.BuildTask(baiduImg,"a.jpg",false,doer)
//
//	if err!=nil {
//		fmt.Println("Error")
//	}else {
//		fmt.Println("build success ",*task);
//	}
//
//	jk_downloader.AddDownloadTask(task)
//}


//func (a *apkDownloader)DownloaderStautsListener(downloaderChannel chan jk_downloader.DownloaderStatus) {
//	for  {
//		select {
//		case status := <-downloaderChannel:
//			switch status {
//			case jk_downloader.DownloaderWaiting:
//				fmt.Println("downloader waiting")
//			case jk_downloader.DownloaderDownloading:
//				fmt.Println("downloader downloading")
//			case jk_downloader.DownloaderFinished:
//				fmt.Println("downloader finished")
//			}
//		}
//	}
