go build main.go
sed 's/{{FOLDER}}/stevebargelt/g' config.base > config.yaml
