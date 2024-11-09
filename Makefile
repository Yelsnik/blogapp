mongodb:
	docker exec -it mongodb mongosh "mongodb+srv://kingsley:mahanta@cluster0.30vt0jd.mongodb.net/blog-app"
run:
	go run github.com/air-verse/air@latest
test:
	go test -v -cover ./...