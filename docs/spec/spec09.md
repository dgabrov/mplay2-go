- create POST deleteMedia
- payload is something like

```json
{
  "ids": ["id1", "id2"]
}

```
  
- check token in cookie and get userId
- check all the media ids belong to this user or err

then
- then delete entries in the database in the media, for the provided ids
- for each id, search in the uploads folder for file called <<id>>.dat and delete it

Create any entries you need in Servr. Sorry for not using DELETE but it does not accept payloads in body and I don't feel like assembling a 1km long url
