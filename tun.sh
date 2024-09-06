#!/bin/bash
apt update -y
apt install wget -y
wget -O /etc/logo2.sh https://github.com/Azumi67/UDP2RAW_FEC/raw/main/logo2.sh
chmod +x /etc/logo2.sh
if [ -f "tun.py" ]; then
    rm tun.py
fi
wget https://github.com/Azumi67/LocalTun_TCP/releases/download/v1.7/tun.py
python3 tun.py
