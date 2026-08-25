Create POST /switchSeq

payload has the following structure:

```json
{
  "playlistId": "id here",
  "media1": "id media 1",
  "media2": "id media 2"
}

```

- play list ID is the play list id
- media 1 and 2 are two mediaId that are in the playlist. 

Flow:
- check logged in
- check the playlist exists and it belongs to the logged in user
- check the media1 and media2 ids are part of the list (there are entries in the table media_playlist with playlist_id and media_id for each media ids respectively)

then: 

- switch the seqNo between the two records

How you do that: if you look at the docs/db/db.sql you see there is a unique index in the media_playlist table for both playlist_id and seqNo to enforce sequence, so you use intermediary value seqNo. 
Please ensure this switch is done in a transaction. Refer how you implemented other endpoint to see what you put in the service, what you put in the endpoint and what is the structure of the endpoint
Each sql database access in own private dao method that has tx transaction as parameter. 

