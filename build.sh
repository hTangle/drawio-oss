CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
version=$(cat version.md)
rm -rf output
mkdir -p output/editor
cp main output/editor/editor-web
mkdir -p output/editor/static/
cp -r static/css output/editor/static/
cp -r static/js output/editor/static/
cp -r templates output/editor/
docker build -t super_markdown_editor_web:${version} .
docker save super_markdown_editor_web:${version}>super_markdown_editor_web_${version}.tar.gz
rm -rf output
sha256sum super_markdown_editor_web_${version}.tar.gz
scp super_markdown_editor_web_${version}.tar.gz ahsup:~
ssh ahsup "docker load -i super_markdown_editor_web_${version}.tar.gz"
