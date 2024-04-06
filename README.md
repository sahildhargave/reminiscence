

 Application Architecture
 Go Application With Gin FrameWork
 Model/Domain
 User
 Token
 Interfaces
 Errors
   |
   |
   |
   |= Handler            Parse/validate incoming requests, call services
   |= Service/Usecase    User ,Token
   |= Repository/Data    User , Image, Token, Events
   |= Data Sources       Redis ,Postgres, Cloud Storage PubSub

PS E:\advance\april\memories\account> curl.exe -X POST "http://malcorp.test/api/account/signin"
{"hello":"it's signin"}
PS E:\advance\april\memories\account> curl.exe -X POST "http://malcorp.test/api/account/signUP"
404 page not found
PS E:\advance\april\memories\account> curl.exe -X POST "http://malcorp.test/api/account/signup"
{"hello":"it's signup"}
PS E:\advance\april\memories\account> curl.exe -X POST "http://malcorp.test/api/account/image" 
{"hello":"it's image"}
PS E:\advance\april\memories\account> curl.exe -X POST "http://malcorp.test/api/account/details"
404 page not found
PS E:\advance\april\memories\account> curl.exe -X POST "http://malcorp.test/api/account/detail" 
404 page not found
PS E:\advance\april\memories\account> curl.exe -X PUT "http://malcorp.test/api/account/details" 
{"hello":"it's Details"}
PS E:\advance\april\memories\account> curl.exe -X DELETE "http://malcorp.test/api/account/image"
{"hello":"it's Delete Image"}
PS E:\advance\april\memories\account> 