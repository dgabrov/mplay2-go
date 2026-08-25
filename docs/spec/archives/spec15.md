**step 1**

create POST endpoint /updateMedia. Do whatever changes in the persistence (Servr) to accommodate your actions

- check user is logged in

it is a multipart request, and has two parts, 

part 1 contains a json with the following structure:

```json
{
  "adding": true,
  "id": "idwill behere",
  "description": "stuff"
}
```

part two contains a file upload. The file might be there or not. 

- adding is true:
  - description should be filled out
  - if not filled out, then you take the name of the file without extension and fill the description with it
  - file attachment must be present or else throw error "cannot add media without the actual content"

- adding is false:
  - check media_id exists (id goes to media_id)
  - check it belongs to the logged in user
  - if file does not exist in the request, only update the description
  - if file exists in the request, replace existing file with it besides updating the description

The files in all cases have the name <<id>>.dat and they reside in the folder that is defined with UploadsFolder in ConfigData

Field 

**step 2**
- when file is uploaded, try to probe it to see what type of media it is, should be either audio or video. Based on that, fill / update the field content_type in table media

**step 3**
- for all files, store the size in bytes in the size field
- in case the media file uploaded is a video, try to get the width and height and put the values in the width and height fields
- if it was a video before and now I update with an audio media, no issue, update width and height in media table with 0

**step 4**
- I already placed the file @/home/daniel/IdeaProjects/mplay2-go/docs/http/sample.mp4 In the same folder create an http that would push a request like the one above

