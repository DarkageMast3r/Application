cd service-discovery
go build -o app.exe
START /b app.exe
cd ../technology-selection
go build -o app.exe
start /B app.exe
cd ../ui
go build -o app.exe
start /B app.exe
