<# Simple script to ease development on Windows #>
Remove-Item '.\frontend_build\' -Recurse -ErrorAction SilentlyContinue
Set-Location frontend
yarn build
Set-Location ..
Move-Item frontend\build frontend_build
go build -o main.exe .\cmd\server\main.go
.\main.exe
