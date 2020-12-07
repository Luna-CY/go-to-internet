package acme

const template = "server {\n" +
    "    listen 80; # ipv4端口\n" +
    "    listen [::]80; # ipv6端口\n" +
    "    server_name {host};\n" +
    "    location / {\n" +
    "        root /var/www/html;\n" +
    "    }\n" +
    "}"
