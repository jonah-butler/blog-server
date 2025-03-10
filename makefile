blogs:
	hurl --test blog.hurl

build:
	docker build -t blog_api .

run:
	docker run -p 8080:8080 --name blog_server --env-file .env blog_api

