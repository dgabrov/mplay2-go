@internal/data/business.go has a new type DescriptionUpdate

- create POST addPlaylist and POST updatePlaylist

- payload for both is DescriptionUpdate
- return a value object like that

```json
{
  "success": true
}
```

- check user authenticated - the token in the cookie 
- get userID with this occasion

addPlaylist => please insert new playlist in the table, you populate the userId, the rest you take from DescriptionUpdate
updatePlaylist=> same as addPlaylist but you update. If no record is there to update, don't err. 



