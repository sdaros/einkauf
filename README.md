# Einkauf App

## Deployment Instructions

The einkauf app's `public/` directory can be served on a CDN or through a reverse proxy. Here is an example of the configuration for a caddy instance listening on port `8000` at the domain `einkauf.example.com` :

```
 encode zstd gzip
 root * /artifacts/sites/einkauf.example.com/public
 log
 handle /ws* {
   reverse_proxy localhost:9999
 }
```

and then make sure to have a systemd service which launches the web-socket server and you're good to go. 

```
[Unit]
Description=Let's go Einkaufen!
Type=simple

[Service]
Environment="PORT=9999"
ExecStart=/usr/bin/node /artifacts/sites/einkauf.example.com/node_modules/y-websocket/bin/server.js
```

All data is persisted on the server in the `db/` folder (in the example above, that would be in `/artifacts/sites/einkauf.example.com/db`)
