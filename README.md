# GODev-xDMS, based on Ideas from a old Java Implementation created in 2010

##Available features
 - list files of a specific subfolder under a fix rootpath
 - download files, or open them in a new browser tab
 - upload new files to ref Folder

##TODO
 - config file
 - test multible content types
 - add inline view for pdf / jpg / tiff

start with > go run xDMS.go
build with > go build xDMS.go

command line options available to overwrite defaults
 -rootpath=<default c://temp>
 -port=<default 4180>

run
 ./xDMS.exe (to use defaults)
 ./xDMS.exe -rootpath=c://dataroot -port=8080