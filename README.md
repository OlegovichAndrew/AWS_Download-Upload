#AWS_Download-Upload
The program downloads the file from AWS by its bucket and key names, saves it,  
then pushes saved file back to the AWS bucket replacing the old one.

#Run

The application starts with flags  
```go run main.go -b [bucket] -f [key]```    
or with existing aws-dl-ul  
```./aws-dl-ul -b [bucket] -f [key]```