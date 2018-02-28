all:
	go build ./...

fmt:
	go fmt ./...

check:
	go test ./...

%.html: %.adoc
	asciidoctor $<

doc: docs/itrex_go.html

clean:
	rm docs/*html
