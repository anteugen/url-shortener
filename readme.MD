*** To attain an URL:
http://localhost:8080/r/(shortcode)

*** To create a new short link:
curl -X POST http://localhost:8080/shorten -H 'Content-Type: application/json' -d '{"url": "(url)"}'