- create endpoint GET /searchMedia
- query searchMedia=value
- returns array of Media

The Media type already exists
Check first the token and retrieve the userId
If any error, return
If not, make sure you filter by the userID retrieved


searchMedia parameter: this is a valid value

"value here correc?ted*"
- transform ? in _ and * in % for search
- split in words
- each word add % at beginning and end if not already there

for the example above you will have

where (description like ? and description like ? and description like ?)

where the bound parameters will be
%value%
%here%
%correc_ted%"

Create any Servr methods you need for this, and the db as you know is here @docs/db/db.sql
