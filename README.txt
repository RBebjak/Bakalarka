Elektronická príloha bakalárskej práce obsahuje:

text/
    Adresár obsahujúci dokumentáciu bakalárskej práce.

    bebjak.pdf
        Finálna verzia bakalárskej práce vo formáte PDF.

    Bebjak.zip
        Zdrojové súbory bakalárskej práce vrátane dokumentu vo formáte
        LaTeX, bibliografických záznamov a vložených obrázkov.

chatServer.zip
    Zdrojové kódy implementovanej klient-server aplikácie v programovacom
    jazyku Go.


PRAKTICKÁ ČASŤ – CHAT SERVER V JAZYKU GO

1. Požiadavky

Na spustenie aplikácie je potrebné mať:

- Go vo verzii 1.21 alebo novšej,
- voľný TCP port 8080.

Verziu jazyka Go je možné overiť príkazom:

go version


2. Štruktúra projektu

goChat/
|
|-- go.mod
|-- client/
    |-- client.go
|
|-- server/
    |-- server.go
    |-- connection/
    |   |-- connection.go
    |-- chatLog/
    |   |-- chatLog.go
    |-- message/
    |   |-- message.go
    |-- broadcast/
    |   |-- broadcast.go
    |-- fetchAll/
    |   |--fetchAll.go
    


3. Spustenie servera

Prejdite do koreňového adresára (goChat/) projektu (adresár obsahujúci súbor go.mod) a spustite príkaz:

go run ./server/server.go

Po úspešnom spustení server počúva na TCP porte 8080.


4. Spustenie klienta

Klienta spustite v ďalšom termináli z koreňového adresára projektu:

go run ./client/client.go <username>

kde <username> predstavuje meno používateľa.

Príklad:

go run ./client/client.go Adam


5. Používanie aplikácie

Po úspešnom pripojení je potrebné zadať názov komunikačnej miestnosti.
Ak miestnosť neexistuje, server ju automaticky vytvorí.

Podporované príkazy:

/fetch
    Načíta nové správy z aktuálnej miestnosti.

/fetchAll
    Načíta správy zo všetkých navštívených miestností.

/broadcast
    Odošle nasledujúcu správu do všetkých navštívených miestností.

exit
    Opustí aktuálnu miestnosť.

logout
    Odhlási používateľa a ukončí spojenie.

Ak používateľ zadá text, ktorý nie je príkazom, text bude odoslaný ako správa do aktuálnej miestnosti.


6. Poznámky

- Server musí byť spustený pred pripojením klientov.
- Projekt využíva Go moduly (go.mod).
- Projekt nevyžaduje inštaláciu externých knižníc.
- Celá aplikácia je implementovaná výhradne pomocou štandardnej knižnice programovacieho jazyka Go.