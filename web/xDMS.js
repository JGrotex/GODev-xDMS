"use strict";

WhenReady(Scope => {

  Scope.objDir = {
    Dir: ''
  };
  var FileListElement = $('.filelist');

  Scope.GetFileList = dir => {

    var params = getQueryParams();
    var ref = params.ref;
    
    $('.ref-bar > input').value = ref;
    $('.ref-tools').innerHTML = '<a href="/upload?ref='+ref+'"><i class="material-icons">file_upload</i></a>';

    http.get('/list/?dir=' + dir, res => {
      if (res !== "Invalid Directory") {
        let FileList = JSON.parse(res);
        FileListElement.innerHTML = "";
        FileList.forEach(File => {
          let FileNameNoSpace = File.Name.replace(/ /g, "~!");
          if (File.IsDir !== true) {
            let FileSizeFormatted = formatBytes(File.Size);
            FileListElement.innerHTML += `
             <div>
             <span onclick="window.open('/file/?dir=${Scope.objDir.Dir}/${FileNameNoSpace}');"> ${File.Name} </span><span> ${FileSizeFormatted}</span>
             </div>
           `;
          } else {
            FileListElement.innerHTML += `
            <div class="dirdiv">
             <span onclick="WhenReady.Scope.objDir.Dir = '${Scope.objDir.Dir + "/" + FileNameNoSpace}'" class='dir'> ${File.Name}</span>
            </div>
          `;
          }
        });
      }
    });
  };

  Scope.GetFileList(Scope.objDir.Dir);

  Object.observe(Scope.objDir, changes => {
    changes.forEach(change => {
      if (change.name === "Dir" && Craft.isDefined(change.oldValue) || Craft.isntNull(change.oldValue)) $('.ref-bar > input').value = Scope.objDir.Dir;
    });
    Scope.GetFileList(Scope.objDir.Dir);
  }, ["update"]);

function formatBytes(bytes, decimals) {
    
  if (bytes == 0) return '0 Byte';
  let k = 1000;
  let dm = decimals + 1 || 3;
  let i = Math.floor(Math.log(bytes) / Math.log(k));
  let sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

  return (bytes / Math.pow(k, i)).toPrecision(dm) + ' ' + sizes[i];
}

function getQueryParams() { 
  var query = window.location.search.substring(1); 
  var qs = query.split("+").join(" "); 
  var params = {}, tokens, re = /[?&]?([^=]+)=([^&]*)/g; 
  while (tokens = re.exec(qs)) { 
    params[decodeURIComponent(tokens[1])]
    = decodeURIComponent(tokens[2]); 
  } 
  return params; 
} 

});
