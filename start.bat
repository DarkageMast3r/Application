cd service-discovery
go build -o app.exe
START /b app.exe
cd ../test-service
go build -o app.exe
start /B app.exe
cd ../test-service-2
go build -o app.exe
start /B app.exe
cd ../ui
go build -o app.exe
start /B app.exe