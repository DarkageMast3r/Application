cd implementatie/cmd/server
go build -o app.exe
START /b app.exe
cd ../../service-discovery
go build -o app.exe
START /b app.exe
cd service-financiering
go build -o app.exe
START /b app.exe
cd service-signalering
go build -o app.exe
START /b app.exe
cd ../technology-selection
go build -o app.exe
start /B app.exe
cd ../ui
go build -o app.exe
start /B app.exe
cd ZorgTechCatalogus_
go build -o app.exe
START /b app.exe