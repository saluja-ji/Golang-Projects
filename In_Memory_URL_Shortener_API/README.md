# In-Memory URL Shortener API

A simple Go example that stores short URL codes in memory and redirects users to the original URL.

## How it works

- `POST /api/shorten` accepts JSON with `originalUrl` and returns a `shortCode`.
- `GET /{code}` looks up the short code and redirects to the stored original URL.

## Run the server

```bash
go run main.go
```

The server listens on port `8080`.

## Example request

```bash
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"originalUrl":"https://example.com"}'
```

Example response:

```json
{
  "shortCode": "abc123"
}
```

## Redirect

Open `http://localhost:8080/abc123` in your browser to be redirected to the original URL.

## Notes

- This implementation stores data only in memory, so all links are lost when the server stops.
- The current code is designed for learning and simple testing.
