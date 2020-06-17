package service

const template = "[Unit]\n" +
    "Description=Go To Internet\n" +
    "After=network-online.target\n" +
    "\n" +
    "[Service]\n" +
    "Type=simple\n" +
    "ExecStart=EXEC_CMD -H YOUR_HOST\n" +
    "\n" +
    "[Install]\n" +
    "WantedBy=multi-user.target\n"
