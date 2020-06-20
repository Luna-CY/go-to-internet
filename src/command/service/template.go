package service

const serverTemplate = "[Unit]\n" +
    "Description=Go To Internet\n" +
    "After=network-online.target\n" +
    "\n" +
    "[Service]\n" +
    "Type=simple\n" +
    "ExecStart=EXEC_CMD -H YOUR_HOST\n" +
    "\n" +
    "[Install]\n" +
    "WantedBy=multi-user.target\n"

const clientTemplate = "[Unit]\n" +
    "Description=Go To Internet\n" +
    "After=network-online.target\n" +
    "\n" +
    "[Service]\n" +
    "Type=simple\n" +
    "ExecStart=EXEC_CMD -sh YOUR_HOST -u USERNAME -p PASSWORD\n" +
    "\n" +
    "[Install]\n" +
    "WantedBy=multi-user.target\n"
