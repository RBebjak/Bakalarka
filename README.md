# Praktická časť – Chat server v jazyku Go

## Požiadavky

Na spustenie aplikácie je potrebné mať nainštalované:

- Go vo verzii **1.21** alebo novšej
- TCP port **8080** nesmie byť obsadený inou aplikáciou
---

## Štruktúra projektu

```
goChat/
│
├── go.mod
├── client/
│   └── client.go
│
└── server/
    ├── server.go
    ├── client/
    │   └── connection.go
    ├── chatLog/
    │   └── chatLog.go
    └── message/
        └── message.go
```

---

## Spustenie servera

Prejdite do koreňového adresára projektu (adresár obsahujúci súbor `go.mod`) a spustite:

```bash
go run ./server/server.go
```

Po úspešnom spustení server počúva na porte **8080**.

---

## Spustenie klienta

Klienta je možné spustiť v ďalšom termináli.

Pripojenie na lokálny server:

```bash
go run ./client/client.go <username>
```

kde:

- `<username>` je používateľské meno.

Každý ďalší klient sa spúšťa v samostatnom termináli, napríklad:

```bash
go run ./client/client.go localhost:8080 Peter
```

---

## Používanie aplikácie

Po úspešnom pripojení je potrebné zadať názov miestnosti. Ak miestnosť neexistuje, bude automaticky vytvorená.

Podporované príkazy:

| Príkaz | Popis |
|--------|-------|
| `/fetch` | Načíta nové správy z aktuálnej miestnosti. |
| `/fetchAll` | Načíta správy zo všetkých navštívených miestností. |
| `/broadcast` | Odošle nasledujúcu správu do všetkých navštívených miestností. |
| `exit` | Opustí aktuálnu miestnosť. |
| `logout` | Odhlási používateľa a ukončí spojenie. |

Ak používateľ zadá ľubovoľný text, ktorý nie je príkazom, text sa odošle ako správa do aktuálnej miestnosti.

---

## Poznámky

- Server musí byť spustený pred pripojením klientov.
- Projekt využíva Go moduly (`go.mod`) a nevyžaduje inštaláciu externých knižníc.
- Celý projekt je implementovaný iba pomocou štandardnej knižnice jazyka Go.