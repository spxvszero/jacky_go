package jk_source_page

const Upload_HTML = `

<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css" integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous">
    <script src="https://code.jquery.com/jquery-3.4.1.slim.min.js" integrity="sha384-J6qa4849blE2+poT4WnyKhv5vZF5SrPo0iEjwBvKU7imGFAV0wwj1yYfoRSJoZ+n" crossorigin="anonymous"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.16.0/umd/popper.min.js"></script>
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.4.1/js/bootstrap.min.js"></script>
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <style>
        .hovering {
            background-color: lightgray;
            opacity: 0.5;
        }
        .interactive_btn_class{
            z-index: 5;
        }

        .revert_more_btn_anim{
            -webkit-animation-direction: revert;
            -moz-animation-direction: revert;
            -o-animation-direction: revert;
            animation-direction: revert;
        }
        .success_btn{
            color: green;
        }
        .failed_btn{
            color: darkgray;
        }

        .more_btn_anim{
            animation: more_btn_spread 350ms;
            -moz-animation: more_btn_spread 350ms;	/* Firefox */
            -webkit-animation: more_btn_spread 350ms;	/* Safari 和 Chrome */
            -o-animation: more_btn_spread 350ms;	/* Opera */

            animation-fill-mode:forwards;
            -webkit-animation-fill-mode: forwards;
            -moz-animation-fill-mode: forwards;
            -o-animation-fill-mode: forwards;
        }
        @keyframes more_btn_spread
        {
            0%   {transform: rotateZ(0deg);}
            100% {transform: rotateZ(90deg);}
        }

        @-moz-keyframes more_btn_spread /* Firefox */
        {
            0%   {transform: rotateZ(0deg);}
            100% {transform: rotateZ(90deg);}
        }

        @-webkit-keyframes more_btn_spread /* Safari 和 Chrome */
        {
            0%   {transform: rotateZ(0deg);}
            100% {transform: rotateZ(90deg);}
        }

        @-o-keyframes more_btn_spread /* Opera */
        {
            0%   {transform: rotateZ(0deg);}
            100% {transform: rotateZ(90deg);}
        }
    </style>
</head>
<body>

<div class="container-fluid">
    <div id="list-group" class="row list-group list-group-flush" style="min-height: 500px;position: relative">
        <div id="drop-zone" class="container border rounded align-self-center" style="height: 100%;position: absolute;">
        </div>
    </div>

    <div id="control-space" style="height: 10px"></div>
    <div class="container">
        <div class="row justify-content-between">
            <div class="col-sm-4">
                <input id="open-file" type="file" name="file" multiple="multiple" onchange="addFiles(this.files)" style="display: none"/>
                <button class="btn btn-primary" onclick="openFile()">新增文件</button>
                <button class="btn btn-primary" onclick="clearAll()">移除所有</button>
            </div>
            <div class="col-sm-4">

            </div>
            <div class="col-sm-2 align-items-center">
                <div>
                    <input id="auto-upload-checkbox" type="checkbox" aria-label="自动开始上传"/>
                    <span class="align-middle">自动开始上传</span>
                </div>
            </div>
            <div class="col-sm-2">
                <button class="btn btn-primary" onclick="clearUpTask()">清理</button>
                <button class="btn btn-primary" onclick="startUploadAll()">开始上传</button>
            </div>
        </div>
    </div>


</div>

</body>

<script type="text/javascript">
var dropZone=document.getElementById("drop-zone");if("draggable"in dropZone&&"ondragenter"in dropZone&&"ondragleave"in dropZone&&"ondragover"in dropZone&&window.File&&window.FileList&&window.FileReader){function handleFileDragEnter(e){jkconsole("file in ",e),e.stopPropagation(),e.preventDefault(),this.classList.add("hovering"),$(".interactive_btn_class").css("z-index","0")}function handleFileDragLeave(e){jkconsole("file out ",e),e.stopPropagation(),e.preventDefault(),this.classList.remove("hovering"),$(".interactive_btn_class").css("z-index","5")}function handleFileDragOver(e){e.stopPropagation(),e.preventDefault(),e.dataTransfer.dropEffect="copy"}function handleFileDrop(e){e.stopPropagation(),e.preventDefault(),this.classList.remove("hovering"),$(".interactive_btn_class").css("z-index","5"),jkconsole("drag finish -- ",e,e.target.files,e.dataTransfer.files,e.dataTransfer.items);var t=!0;try{e.dataTransfer.items[0].webkitGetAsEntry()}catch(e){jkconsole("不支持读取文件夹"),t=!1}t?(addFilesFromDragEvt(e.dataTransfer.items),jkconsole("dididididid --- ",e.dataTransfer.items," files --- ",e.dataTransfer.files)):addFiles(e.dataTransfer.files)}dropZone.addEventListener("dragenter",handleFileDragEnter,!1),dropZone.addEventListener("dragleave",handleFileDragLeave,!1),dropZone.addEventListener("dragover",handleFileDragOver,!1),dropZone.addEventListener("drop",handleFileDrop,!1)}async function getAllFileEntries(e){let t=[],a=[];for(let t=0;t<e.length;t++)a.push(e[t].webkitGetAsEntry());for(;a.length>0;){let e=a.shift();e.isFile?t.push(e):e.isDirectory&&a.push(...await readAllDirectoryEntries(e.createReader()))}return t}async function readAllDirectoryEntries(e){let t=[],a=await readEntriesPromise(e);for(;a.length>0;)t.push(...a),a=await readEntriesPromise(e);return t}async function readEntriesPromise(e){try{return await new Promise((t,a)=>{e.readEntries(t,a)})}catch(e){jkconsole(e)}}var tasksObj=new Object,groupObj=new Object,progressName="upload-progress-",percentageName="upload-percentage-",timeName="upload-time-",listItemName="list-item-",cancelBtnName="cancel-btn-";let taskId=0;function startUploadAll(){Object.keys(tasksObj).forEach(function(e){jkconsole("start upload ",e);var t=tasksObj[e];(t&&null==t.xhr||200!=t.xhr.status&&4==t.xhr.readyState)&&sendUploadRequest(t)})}function openFile(){$("#open-file").click()}function addFiles(e){for(jkconsole("files - ",e),i=0;i<e.length;i++)addTask(e[i])}function addFilesFromDragEvt(e){getAllFileEntries(e).then(e=>{var t=e.length;e.map((e,a)=>{var n=e,s=n.fullPath;n.isFile;n.file(e=>{buildTaskWithGroup(e,s),0==--t&&jkconsole("finish map -- ",groupObj)})})}).finally(()=>{jkconsole("groupObject ",groupObj)})}function literateGroup(e,t,a){var n=t.split("/");if(n.pop(),n.length<=1)return!1;for(var s=groupObj,r=!1,o=0;o<n.length;o++){var i=n[o];if(!(i.length<=0)){var l=s[i];if(a)l||((l=new Object).allTasks=new Array,s[i]=l),l.allTasks.push(e);else for(var d=0;d<l.allTasks.length;d++)if(l.allTasks[d]==e){if(l.allTasks.splice(d,1),r=!0,l.allTasks.length<=0)break;d--}s=l}}return a?n[1]:r}function buildTaskWithGroup(e,t){taskId++;var a=e.name;tasksObj[taskId]={taskId:taskId,groupPath:t,fileName:a,xhr:null,fileObject:e};var n=literateGroup(taskId,t,!0);return n?addGroupListItem(n):addListItem(a,taskId),autoUpload(tasksObj[taskId]),taskId}function buildTask(e){taskId++;var t=e.name;return tasksObj[taskId]={taskId:taskId,fileName:t,xhr:null,fileObject:e},addListItem(t,taskId),autoUpload(tasksObj[taskId]),taskId}function addTask(e){jkconsole("add task - ",e);var t=new FileReader;t.onload=function(t){buildTask(e)},t.onerror=function(e){jkconsole("read file error : ",e)},t.readAsText(e)}function removeTask(e){var t=tasksObj[e];if(t){if(removeListItem(e),t.groupPath){var a=getGroupName(t.groupPath,0);if(literateGroup(e,t.groupPath,!1))if(!groupObj[a].allTasks||groupObj[a].allTasks.length<=0)removeGroupListItem(a),delete groupObj[a];else{let e="title-"+a;document.getElementById(e).innerText=a+"("+groupObj[a].allTasks.length+")",updataGroupProgress(a)}}delete tasksObj[e]}}function addGroupListItem(e){var t=document.getElementById("list-group");if(document.getElementById("group-"+e)){let t="title-"+e;document.getElementById(t).innerText=e+"("+groupObj[e].allTasks.length+")"}else{var a='<div id="group-'+e+'" class="container list-group-item">\n            <div class="row justify-content-between align-items-center">\n                <div class="interactive_btn_class col-10">\n<button type=\'button\' class=\'btn\' onclick="spreadGroup(\''+e+'\')"><svg id="icon-'+e+'" width="15px" height="15px" viewBox="0 0 15 15">\n    <g stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">\n        <polygon id="Path" fill="#000000" points="0 0 0 15 15 7.5"></polygon>\n    </g>\n</svg>                    <span id="title-'+e+"\" class='align-middle'>"+e+'</span>\n </button><div id="inner-'+e+'"></div>                     <div id="'+progressName+"group-"+e+'" class="progress">\n                        <div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" aria-valuenow="75" aria-valuemin="0" aria-valuemax="100" style="width: 0%"></div>\n                    </div>\n                    <span  id="'+percentageName+"percent-"+e+'">0%</span><span></span>\n                </div>\n                <div class="interactive_btn_class col-2 text-center">\n                    <button  class="btn btn-primary" type="button" onclick="cancelGroupUploadFile(\''+e+"')\">cancel</button>\n                </div>\n            </div>\n        </div>";t.insertAdjacentHTML("afterbegin",a)}}function removeGroupListItem(e){document.getElementById("group-"+e).remove()}function addListItem(e,t){var a=document.getElementById("list-group"),n='<div id="'+listItemName+t+'" class="container list-group-item">\n            <div class="row justify-content-between align-items-center">\n                <div class="col-10">\n                    <span>'+e+'</span>\n                    <div id="'+progressName+t+'" class="progress">\n                        <div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" aria-valuenow="75" aria-valuemin="0" aria-valuemax="100" style="width: 0%"></div>\n                    </div>\n                    <span id="'+percentageName+t+'">0%</span><span id="'+timeName+t+'"></span>\n                </div>\n                <div class="interactive_btn_class col-2 text-center">\n                    <button id="'+cancelBtnName+t+'"  class="btn btn-primary" type="button" onclick="cancelUploadFile('+t+')">cancel</button>\n                </div>\n            </div>\n        </div>';a.insertAdjacentHTML("afterbegin",n)}function removeListItem(e){$("#"+listItemName+e).remove()}function spreadGroup(e){var t=document.getElementById("icon-"+e);t.classList.contains("more_btn_anim")?(t.classList.remove("more_btn_anim"),clearGroupInnerListItem(e)):(t.classList.add("more_btn_anim"),addGroupInnerListItem(e))}function addGroupInnerListItem(e){var t=document.getElementById("inner-"+e);if(t)for(var a=groupObj[e].allTasks,n=0;n<a.length;n++){var s=a[n],r=tasksObj[s],o="0%",i="cancel",l="";r.xhr&&200==r.xhr.status&&(o="100%",i="success",l="success_btn");var d='<div id="'+listItemName+s+'" class="container" style=\'width: 80%\'>\n            <div class="row justify-content-between align-items-center">\n                <div class="col-10">\n                    <span>'+r.fileName+"</span>\n                    <span style='font-style: italic;color: lightgray;' id=\""+percentageName+s+'">'+o+"</span>\n                    <span style='font-size: 10px;color: lightgray;'>"+r.groupPath+'</span>\n                </div>\n                <div class="col-2 interactive_btn_class text-center">\n                    <button id="'+cancelBtnName+s+'"  class="btn btn-link '+l+'" type="button" onclick="cancelUploadFile('+s+')">'+i+"</button>\n                </div>\n            </div>\n        </div>";t.insertAdjacentHTML("afterbegin",d)}}function updataGroupProgress(e){var t=progressName+"group-"+e,a=percentageName+"percent-"+e,n=document.getElementById(t),s=document.getElementById(a),r=groupObj[e].allTasks;if(r&&!(r.length<=0)){for(var o=0,i=0;i<r.length;i++){var l=r[i];tasksObj[l].xhr&&200==tasksObj[l].xhr.status&&o++}var d=o/r.length*100+"%";n&&$(n).children().css("width",d),s&&(s.innerText=d)}}function clearGroupInnerListItem(e){var t=document.getElementById("inner-"+e);t&&(t.innerText="")}function getGroupName(e,t){var a=e.split("/");return a.pop(),t+1>=a.length?null:a[t+1]}function autoUpload(e){$("#auto-upload-checkbox")[0].checked&&sendUploadRequest(e)}function clearAll(){Object.keys(tasksObj).forEach(function(e){jkconsole("start upload ",e);var t=tasksObj[e];t&&(t.xhr?(4==t.xhr.readyState||t.xhr.abort(),removeTask(e)):removeTask(e))})}function clearUpTask(){Object.keys(tasksObj).forEach(function(e){jkconsole("clear up task",e);var t=tasksObj[e];t&&t.xhr&&4==t.xhr.readyState&&removeTask(e)})}function retryTask(e){delete e.xhr;var t=document.getElementById(cancelBtnName+e.taskId);$(t)[0].innerText="cancel",$(t)[0].classList.contains("btn-primary")?$(t)[0].classList.remove("btn-secondary"):$(t)[0].classList.remove("failed_btn");var a=document.getElementById(progressName+e.taskId);a&&$(a).children().removeClass("bg-secondary"),sendUploadRequest(e)}function sendUploadRequest(e){var t;let a=new FormData;a.append("file",e.fileObject),e.groupPath&&a.append("relativePath",e.groupPath),(t=new XMLHttpRequest).open("post",{{.Path}},!0),t.addEventListener("load",function(t){uploadComplete(t,e.taskId)},!1),t.addEventListener("error",function(t){uploadFailed(t,e.taskId)}),t.upload.addEventListener("progress",function(t){progressFunction(t,e.taskId)}),t.upload.onloadstart=function(){e.beginTime=(new Date).getTime(),e.ot=(new Date).getTime(),e.oloaded=0,e.maxspeed=0,e.maxspeedStr="b/s"},t.upload.onreadystatechange=function(){4==t.readyState&&200==t.status&&(jkconsole("upload complete"),jkconsole("response: "+t.responseText))},e.xhr=t,t.send(a)}function uploadComplete(e,t){jkconsole("上传---结果 ",e,t,tasksObj[t]);var a=cancelBtnName+t,n=document.getElementById(a);if(200==e.target.status){n&&($(n)[0].innerText="success",$(n)[0].classList.contains("btn-primary")?$(n)[0].classList.add("btn-success"):$(n)[0].classList.add("success_btn"));var s=document.getElementById(progressName+t);s&&$(s).children().addClass("bg-success");var r=tasksObj[t],o=document.getElementById(timeName+t);if(o){var i=((new Date).getTime()-r.beginTime)/1e3,l=i<=0?r.fileObject.size:r.fileObject.size/i,d="b/s";l/1024>1&&(l/=1024,d="k/s"),l/1024>1&&(l/=1024,d="M/s"),l=l.toFixed(1),o.innerHTML="，平均速度："+l+d+"，最大速度："+r.maxspeedStr+"，所花时间："+i+"s"}var c=getGroupName(r.groupPath,0);c&&updataGroupProgress(c)}else uploadFailed(e,t)}function uploadFailed(e,t){var a=document.getElementById(cancelBtnName+t);$(a)[0].innerText="retry",$(a)[0].classList.contains("btn-primary")?$(a)[0].classList.add("btn-secondary"):$(a)[0].classList.add("failed_btn");var n=document.getElementById(progressName+t);n&&$(n).children().addClass("bg-secondary")}function cancelGroupUploadFile(e){for(var t=groupObj[e].allTasks.slice(),a=0;a<t.length;a++)cancelUploadFile(t[a])}function cancelUploadFile(e){jkconsole("cancel click -- ",e);var t=tasksObj[e];t&&(t.xhr?200==t.xhr.status?removeTask(e):4==t.xhr.readyState?retryTask(t):(t.xhr.abort(),jkconsole("cancel ",e," success"),removeTask(e)):removeTask(e))}function progressFunction(e,t){var a=document.getElementById(progressName+t),n=document.getElementById(percentageName+t),s=tasksObj[t];if(e.lengthComputable){var r=Math.round(e.loaded/e.total*100)+"%";a&&$(a).children().css("width",r),n&&(n.innerHTML=r)}var o=document.getElementById(timeName+t),i=((new Date).getTime()-s.ot)/1e3;s.ot=(new Date).getTime();var l=e.loaded-s.oloaded;s.oloaded=e.loaded;var d=i<=0?l:l/i,c=d,p="b/s";d/1024>1&&(d/=1024,p="k/s"),d/1024>1&&(d/=1024,p="M/s"),d=d.toFixed(1),c>s.maxspeed&&(s.maxspeed=c,s.maxspeedStr=d+p);var u=((e.total-e.loaded)/c).toFixed(1);o&&(o.innerHTML="，速度："+d+p+"，剩余时间："+u+"s",0==c&&(o.innerHTML="，速度：-"+p+"，剩余时间：-s"))}function jkconsole(...e){}
</script>
</html>

`