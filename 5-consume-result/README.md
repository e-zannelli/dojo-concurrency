Here process receives a client to a service paginated with a next in the response, you need to

- call cCursor until the last page (Cursor.Next return an error)
- call the client c with all IDs from Cursor.Content.
