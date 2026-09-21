DOCS_DIR := docs

PHONY: build-docs serve-docs docs-theme
build-docs:
	hugo --source $(DOCS_DIR) --gc --minify

serve-docs:
	hugo server --source $(DOCS_DIR) --disableFastRender

docs-theme:
	cd $(DOCS_DIR) && hugo mod get -u github.com/imfing/hextra && hugo mod tidy
