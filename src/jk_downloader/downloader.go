package jk_downloader

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"jacky_go/src/jk_err"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

type TaskStauts int
const (
	TaskReady  TaskStauts =iota
	TaskRunning
	TaskPause
	TaskFinished
	TaskFailed
)

type DownloaderStatus int
const (
	DownloaderWaiting DownloaderStatus = iota
	DownloaderDownloading
	DownloaderFinished
)


type TaskField struct {
	Identify      int64
	sourceUrl     string
	destinatePath string
	append        bool
	Status        TaskStauts
	statusBlock   func(task *TaskField, status TaskStauts)
	progressBlock func(task *TaskField, progress float32, expectedSize int64)
}

func (t *TaskField)changeTaskStauts(status TaskStauts)  {
	t.Status = status
	if t.statusBlock != nil {
		t.statusBlock(t,status)
	}
}
func (t *TaskField)changeProgress(expectedSize int64, curSize int64)  {
	if t.progressBlock != nil {
		var progress float32
		if expectedSize > 0{
			progress = float32(curSize)/float32(expectedSize)
		}else {
			progress = 0
		}
		t.progressBlock(t,progress,expectedSize)
	}
}

type DownloaderInerface interface {
	DownloaderStautsListener(chan DownloaderStatus)
}


type DownloaderManager struct {
	tasks []TaskField
	taskChan chan *TaskField
	baseDownloadDir string
}

//func init() {
//	log.SetFlags(log.Ldate|log.Lshortfile|log.Lmicroseconds)
//	log.SetPrefix("【Downloader】")
//}

func BuildManager(baseDownloadDir string) *DownloaderManager {
	dM := &DownloaderManager{
		[]TaskField{},
		make(chan *TaskField),
		baseDownloadDir,
	}
	go startChannel(dM)
	return dM
}

func BuildTask(src string, dst string , append bool, statusBlock func(task *TaskField,status TaskStauts), progressBlock func(task *TaskField, progress float32, expectedSize int64)) (res *TaskField,err error)  {

	if len(src) <=0 || len(dst) <=0 {
		return nil,jk_err.JKErrInput(jk_err.Input)
	}

	res = &TaskField{getIdentify(),src,dst,append,TaskReady,statusBlock,progressBlock}
	res.changeTaskStauts(TaskReady)

	return res,nil
}

func AddTaskToManager(dM *DownloaderManager,task *TaskField) (err error) {
	dM.tasks = append(dM.tasks,*task)
	log.Println("add task " , task.Identify, " to channel")
	dM.taskChan <- task
	return err
}

func DeleteTaskInManager(dM *DownloaderManager,taskId int64) error {
	for pos, task := range dM.tasks {
		if task.Identify == taskId {
			if pos == 0 {
				dM.tasks = append(dM.tasks[:0],dM.tasks[1:]...)
			}else if pos == len(dM.tasks)-1 {
				dM.tasks = dM.tasks[:len(dM.tasks)-1]
			}else {
				dM.tasks = append(dM.tasks[:pos], dM.tasks[pos+1:]...)
			}
			return nil
		}
	}

	return errors.New("Task not found")
}

func startChannel(dM *DownloaderManager)  {
	for  {
		task, alive := <-dM.taskChan
		log.Println("task channel get a new task ",task.Identify)
		if alive == false {
			log.Println("task ",task.Identify," channel broken.. rebuild channel")
			continue
		}

		log.Println("start get url ",task.sourceUrl)
		go downloadWithTask(task)
	}
}

func downloadWithTask(task *TaskField)  {

	resp,err := http.Get(task.sourceUrl)
	if err!=nil{
		log.Println("task ",task.Identify," network error : ",err)
		return
	}

	task.changeTaskStauts(TaskRunning)

	defer resp.Body.Close()

	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	log.Println("Header",resp.Header)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		task.changeProgress(resp.ContentLength,int64(result.Len()))
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			log.Println("task ",task.Identify," write file error: ",err)
			task.changeTaskStauts(TaskFailed)
			return
		}
	}
	err = ioutil.WriteFile(task.destinatePath,result.Bytes(),os.ModePerm)
	if err!=nil {
		log.Println("failed download ",task.Identify," -- ",err)
		task.changeTaskStauts(TaskFailed)
	}else {
		log.Println("success -- ",task.Identify)
	}
	task.changeTaskStauts(TaskFinished)
}

//func filterHeader(){
	//"Content-Length"
	//"Content-Type"
	//Header map[Accept-Ranges:[bytes] Cache-Control:[max-age=315360000] Connection:[Keep-Alive] Content-Length:[7877] Content-Type:[image/png] Date:[Mon, 09 Dec 2019 06:24:38 GMT] Etag:["1ec5-502264e2ae4c0"] Expires:[Thu, 06 Dec 2029 06:24:38 GMT] Last-Modified:[Wed, 03 Sep 2014 10:00:27 GMT] P3p:[CP=" OTI DSP COR IVA OUR IND COM "] Server:[Apache] Set-Cookie:[BAIDUID=7E34657746BF6154A5CC0BB4D4805E97:FG=1; expires=Tue, 08-Dec-20 06:24:38 GMT; max-age=31536000; path=/; domain=.baidu.com; version=1]]
//}

//func checkDMStatus(dM DownloaderManager){
//	for _,task := range dM.tasks  {
//		if task.Status == TaskRunning {
//			dM.downloadStatusChan <- DownloaderDownloading
//			return
//		}
//	}
//
//	dM.downloadStatusChan <- DownloaderWaiting
//}

var identify int64 = 0
func getIdentify() int64 {
	atomic.AddInt64(&identify, 1)
	return identify
}

