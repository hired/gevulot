# Gevulot — a PostgreSQL proxy for masking sensitive data

> Gvulot (Hebrew: גְּבוּלוֹת, lit. "Borders")

Gevulot is a TCP proxy that sits between your PostgreSQL database and client and proxies data back and forth.
It listens to all messages sent from database to the client and provides a mechanism for users to modify data in-transit and before received by the client. The main purpose of Gevulot is to mask personally identifiable information sent to clients.

## Usage

```
gevulot --listen 0.0.0.0:4241 --connect 0.0.0.0:5342
```