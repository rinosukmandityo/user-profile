How to run
---
#### Set Environment Variable
User Profile is an application with following basic features:   

- Login
- Login with Google account
- Sign Up
- Sign Up with Google
- View user profile
- Update user profile
- Change Password
- Session management
- Token management

To run it we can use `make` command or manual run using `go run main.go`.   
Here is `make` command to run this application:   

- `make run` to run launch this application into port `:8000`, 
  if you want to use google sign up / login and send email feature 
  then before you run it make sure you change some environment variable
  in docker/docker-compose.yaml
-  `make test` to run the all test case and check the coverage  

This application support MySQL as a database, but we can implement a different database.  
By default, it will connect into our mySQL database with default host & port `localhost:3306` and database `users`  
To connect into a different database we need to set database information in environment variable.  

```cli
set url=root:Password.1@tcp(127.0.0.1:3306)/users
set timeout=10
set db=users
set driver=mysql
```
Setup google callback for login and sign up and also google client ID and google client secret from google account   

```cli
set GOOGLE_LOGIN_REDIRECT=http://localhost:8000/googlecallback  
set GOOGLE_SIGNUP_REDIRECT=http://localhost:8000/googlesignupcallback  
set GOOGLE_CLIENT_ID=[your_google_client_ID]  
set GOOGLE_CLIENT_SECRET=[your_google_client_secret]  
```
Setup SMTP info for sending email   

```cli
set SMTP_HOST=[your_SMTP_host]   
set SMTP_PORT=[your_SMTP_port]   
set SOURCE_EMAIL=[your_email]   
set EMAIL_PASSWORD=[your_email_password]   
```

After setting the database information we only need to run the main.go file  
`go run main.go`  

#### API List & Payloads
Here is few API List and its payload:  

1. [GET] **/user/{_user\_id_}**  
`/user`
2. [PUT] **/user/{_user\_id_}**  
`/user`
```json
{
	"Name":     "Name",  
	"Password": "Password.User",  
	"ID":       "userid01",  
	"Email":    "usermail01@gmail.com",  
	"Address":  "User Address 01",   
	"IsActive": false  
}
```
3. [POST] **/auth**  
```json
{  
	"Email": "usermail01@gmail.com",  
	"Password": "Password.User"
}
```

Project Structure
---
By implementing Hexagonal Architecture we also implement Dependency Inversion and Dependency Injection. Here is some explanations about project structure:

1. **api**  
contains handler for API
2. **models**  
contains data models
3. **repositories**  
contains **Port** interface for repository adapter
   - **mysql**  
contains MySQL **Adapter** that implement UserRepository interface. This package will store MySQL client and connect to MySQL server to handle database query or data manipulation
4. **serializer**  
contains **Port** interface for decode and encode serializer. It will be used in our API to decode and encode data.
   - **json**  
contains json **Adapter** that implement serializer interface to encode and decode data
5. **services**  
contains **Port** interface for our domain service and logic 
6. **logic**  
contains service **Adapter** that implement service interface to handle service logic like constructing repository parameter and calling repository interface to do data manipulation or query