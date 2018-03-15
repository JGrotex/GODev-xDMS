# GODev-xDMS, 
based on Ideas from a old Java Implementation created in 2010

## Available features
 - list files of a specific subfolder under a fix rootpath
 - download files, or open them in a new browser tab
 - upload new files to ref Folder

![xDMS image](screenshots/xDMS-golang.png?raw=true "xDMS Screenshot") 

## TODO
 - config file
 - test multible content types
 - add inline view for pdf / jpg / tiff

## Start/Build
```
start with > go run xDMS.go
build with > go build xDMS.go
```

## Commandline Options
command line options available to overwrite defaults

```
 -rootpath=<default c://temp>
 -port=<default 4180>
```

# Execute
``` 
 ./xDMS.exe (to use defaults)
 ./xDMS.exe -rootpath=c://dataroot -port=8080
```

![xDMS Execution image](screenshots/xDMS-startup.png?raw=true "xDMS Startup Screenshot") 