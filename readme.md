# Installatie
1) Installeer databasen
Installeer Mysql op de relevante server. Zoek naar de Deploy map van de relevante microservice, indien aanwezig, en vor dan CreateData.sql uit.
Stel hierna de connection_string in config.json in zodat die verbinding maakt met de relevante database.

2) Relevante bestanden kopiëren
Laat het project bouwen via github actions, en kies dan of je de windows of linux release wilt gebruiken.
Kopiëer dan de folders van de services die je wilt hebben over naar de server, en voer de uitvoerbare bestanden uit.
In linux moeten alle app_service_XXX eerst nog uitvoerbaarheid worden door chmod +x app_service_XXXX.
Alle uitvoerbare bestanden moeten worden geplaatst in hun respectieve microservice folder.
Ook moet een server certificaat en key aanwezig zijn. Deze kunnen op windows worden gegenereerd via genkeys.bat, maar het is beter om via een acme client de server.crt en server.key bestanden te verkrijgen.
De uiteindelijke folder structuur ziet er dan zo uit:
applicatie
- server.crt
- server.key
- service-XXX
  - app_service_XXX
  - config.json
- service-YYY
  - app_service_YYY
  - config.json


3) Server configureren
Kies het domein waarop je server gaat draaien, en stel dan de volgende twee instellingen in in alle config.json bestanden. Deze moeten precies hetzelfde zijn in alle bestanden.
    "service_discovery_root": "localhost",
    "service_discovery_port": 443,
    "rabbitmq": "amqp://guest:guest@localhost:5672/",
Ten slotte hebben alle servers 
    "allow_insecure": false

4) Server starten
Alle uitvoerbare bestanden kunnen nu op de achtergrond worden uitgevoerd,

# Werk bijdrage
- Implementatie: Arjan
- service-discovery: Tijmen
- service-financiering: David
- service-signalering: Pelle
- technology-selection: Tijmen
- ui: Pelle
- ZorgTechCatalogus: Arjan