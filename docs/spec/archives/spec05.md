**STEP 1**

Create User struct in @internal/data/login.go that mimics User table from docs/db/db.sql 
Create method in Servr GetUserByProvidedId which takes a string and searches for user with provided_user_id of this value  

**STEP 2**
Create method in Servr that creates user based on a value object of type data.User created at STEP 1


**STEP3**
Steps in process in post-login.go. Please create reasonable separate private methods, still targeting PostLoginEndpoint

1. get  login and password
2. send a post request against the url in ConfigData auth/url
3. get result
4. if error, return the error
5. if not error, you get a data.LoginResponse, get the rights
6. if ConfigData auth/right is among them, you are good, otherwise return error like 'you don't have the right to access this application'

in LoginResponse, userId is the provided id by the identity provider.

Use server method at step 1 to see if user exists, if not, create user with the method at **Step2**


Create entry in config.json and ConfigData etc called tokenValidity with value 1200

for the user in the database - either existent or new
- create random token using createRandomToken in @internal/endpoint/proc.go
- create entry in session table with that token, the user_id etc
- expired_ind is 'N'
- expired_dt will be time.Now() + tokenValidity from ConfigData
