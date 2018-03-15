"use strict";

WhenReady(Scope => {

  Scope.objDir = {
    Dir: ''
  };
  
  var params = getQueryParams();
  var ref = params.ref;

  $('.ref-bar > input').value = ref;
  $('.ref-tools').innerHTML = '<a href="/?ref='+ref+'"><i class="material-icons">backspace</i></a>';
  $('.uploadform').innerHTML = '<form enctype="multipart/form-data" action="/fupload?ref='+ ref +
    '" method="post"> &nbsp; <input type="file" name="uploadfile" /> <input type="hidden" name="token" value="{{.}}"/> <input type="submit" value=" >>> " /></form>'

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
