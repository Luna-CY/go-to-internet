package acme

const template = "server {\n" +
    "    listen 80;\n" +
    "    server_name {host};\n" +
    "    location / {\n" +
    "        root /var/www/html;\n" +
    "    }\n" +
    "}"
