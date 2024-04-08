

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

####Signup
Model/Domain
User
Token


#### Scrypt vs RSA: Scrypt is a key derivation function designed for securely deriving keys from passwords, providing resistance against brute-force attacks, whereas RSA is an asymmetric encryption algorithm used for key exchange and digital signatures.


#### Scrypt vs Hashed Shared Password: Scrypt provides stronger protection by incorporating computational and memory-intensive operations, making it more resistant to brute-force attacks compared to simple password hashing techniques used for shared password storage.

```
	password := []byte("mysecretpassword")
	salt := []byte("randomsalt123")
	N := 16384   // CPU/Memory cost (must be a power of 2 greater than 1)
	r := 8       // Block size
	p := 1       // Parallelization factor
	keyLen := 32 // Desired key length in bytes
```

### Authorization with JWTClient

Client  ------->   Account application

## Account application 
```
Two Roles
1 . Manages authentication and creation of authorization tokens

2 . Allows user to update a few person details

TOKEN SINGING

- idToken(vaild for 15 minutes) := RSA256 Signing 
- Signed with private ,verified with public

refresh Token (valid for 3 days) : HS256 Signing-
signeed and verified with secret string
```



Application in our domain or organization can verify idToken in Auth Header with public RSA Key. If valid , they can choose to grant access!

##### Generating an RSA Key
```
openssl genpkey -algorithm RSA -out rsa_private.pem -pkeyopt rsa_keygen_bits:2048

openssl rsa -in rsa_private.pem -pubout -out rsa_public.pem
```