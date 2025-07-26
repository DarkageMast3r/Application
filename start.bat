cd service-discovery
go build -o app.exe
START /b app.exe
cd ..

cd service-financiering
go build -o app.exe
START /b app.exe
cd ..

cd service-signalering
go build -o app.exe
START /b app.exe
cd ..

cd technology-selection
go build -o app.exe
start /B app.exe
cd ..

cd ui
go build -o app.exe
start /B app.exe
cd ..

cd ZorgTechCatalogus/cmd/server
go build -o app.exe
START /b app.exe
cd ../../..

cd implementatie/cmd/server
go build -o app.exe
START /b app.exe
cd ../../..