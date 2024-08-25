![R (2)](https://github.com/Azumi67/PrivateIP-Tunnel/assets/119934376/a064577c-9302-4f43-b3bf-3d4f84245a6f)
نام پروژه : لوکال تانل به صورت دایرکت یا ریورس - Tun interface و بر روی پورت TCP
---------------------------------------------------------------

**در حال نوشتن اسکریپت برای استفاده شخصی و گیم انلاین**(اسکریپت اماده است و بعدا اضافه میشود. در اسکریپت ترکیبات وایرگارد و geneve و vxlan هم خواهد بود)

**تغییراتی در authentication method انجام شد و از pub key & priv key استفاده خواهد شد**

**این پروژه برای استفاده شخصی خودم، گیم و یادگیری میباشد. لطفا استفاده نکنید**(در صورت تمایل به استفاده، اسکریپت آن آپلود خواهد شد)

**لاگ های تانل در مسیر etc/server.log/ یا etc/client.log/ ذخیره میشود**

**گزینه worker حذف شد**

**مورد Challenge n Response Authentication به همراه unique nonce و Sha 256 به همراه expiry time حذف شد**

**این مورد xtaci/smux دوباره با retry logic اضافه شد**

**این مورد tcpnodelay و logrus اضافه شد**



![check](https://github.com/Azumi67/PrivateIP-Tunnel/assets/119934376/13de8d36-dcfe-498b-9d99-440049c0cf14)
**امکانات**
- امکان لوکال تانل بین سرور و کلاینت به صورت دایرکت یا ریورس ( برای سرور های محدود)
- استفاده از ایپی پرایوت های ساخته شده در Tun interface برای تانل اصلی یا پورت فوروارد
- امکان اتصال بین سرور و کلاینت بر روی پورت TCP
- امکان اتصال سرور و کلاینت با پابلیک ایپی 4 یا native
- امکان انتخاب ایپی پرایوت انتخابی خودتان هم به صورت پرایوت ایپی 4 یا پرایوت ایپی 6
- امکان انتخاب subnet mask برای پرایوت ایپی های ساخته
- امکان وارد کردن mtu به صورت manual
- دارای ping interval و استفاده از Bin bash برای ریست سرویس ها
- دارای encryption های IPSEC (به زودی)
- دارای tcp nodelay و tcp keepalive
- دارای پرایوت و پابلیک key برای ارتباط بین سرور و کلاینت 
- دارای verbose برای نمایش لاگ (خطا)
- مناسب برای ترکیب با IPSEC > لینک : https://github.com/Azumi67/6TO4-GRE-IPIP-SIT
-----------------------
<div align="right">
  <details>
    <summary><strong>توضیحات</strong></summary>
  
------------------------------------ 
 <div align="right">
   
- هدف نوشتن این برنامه یادگیری و استفاده شخصی در گیم های خودم بوده است
- شما به صورت ریورس یا دایرکت، یک لوکال ایپی دریافت میکنید و سپس از آن پرایوت ایپی ها برای دایرکت تانل، پورت فوروارد یا ریورس استفاده مینمایید. 
- پس از انجام تانل‌ لوکال به صورت دایرکت یا ریورس، به طور مثال میتوانید از پورت فوروارد استفاده نمایید یا مثلا دایرکت تانل چیزل استفاده نمایید یا ریورس.
- در روش ریورس، سرور اصلی میتواند ایران باشد و کلاینت خارج و در روش دایرکت، سرور اصلی میتواند خارج باشد و کلاینت ایران. بدین صورت میتوان تانل لوکالی بر روی سرور های خارج محدود در ان سرور ایران(به صورت ریورس) هم ایجاد کرد.
- با ایپی 4 سرور و هم با ایپی 6 سرور و کلاینت میشود که وصل شد .
- پورت تنها برای ارتباط بین سرور و کلاینت میباشد و شما تنها باید از پرایوت ایپی ها برای تانل اصلی استفاده نمایید.
- اول دستورات سرور را اجرا کنید و سپس دستورات کلاینت . میتوانید هم به صورت دایرکت یا ریورس انجام دهید. یعنی سرور اصلی خارج و کلاینت ایران و یا سرور اصلی ایران و کلاینت خارج باشد
- این پروژه دارای اسکریپت و اموزش برای کسانی که مانند من مصرف شخصی دارند و گیم انلاین انجام میدهند، است.
  </details>
</div>

---------------------

  ![6348248](https://github.com/Azumi67/PrivateIP-Tunnel/assets/119934376/398f8b07-65be-472e-9821-631f7b70f783)
**روش اجرا**
-

- حتما در صورت استفاده از فایراول، پورت و پرایوت ایپی ها را در فایروال اضافه نمایید.
 <div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image">نصب </strong></summary>
  
<div align="left">
  
```
  apt update -y
  apt install wget -y
  apt install unzip -y
  ## amd64
  rm amd64.zip
  wget https://github.com/Azumi67/LocalTun_TCP/releases/download/v1.6/amd64.zip
  unzip amd64.zip -d /root/localTUN
  cd localTUN
  chmod +x tun-server 
  chmod +x tun-client  
 ```
 </details>
</div>
<div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image">دایرکت لوکال تانل پرایوت ایپی 4 - public ipv4 </strong></summary>

  - کامند های سرور (خارج)
  - برای verbose لاگ از کامند -v استفاده نمایید . 
 <div align="left">
   
```
./tun-server -server-port 800 -pub-key=/root/keys/public_key.pem -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -key azumi -mtu 1380 -tcpnodelay -keepalive 10 -smux
   
```
<div align="right">
  
- کامند های کلاینت (ایران)
- برای verbose لاگ از کامند -v استفاده نمایید . 
 <div align="left">
   
```
./tun-client -server-addr KHAREJ_IPV4 -server-port 800 -priv-key=/root/private_key.pem -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1280 -v -tcpnodelay -smux -keepalive 10
```
 <div align="right">
   
- نحوه ساختن سرویس
 <div align="left">
   
```
nano /etc/systemd/system/azumilocal.service
## put this config inside [ This is a sample]##

[Unit]
Description=Azumi local Service
After=network.target

[Service]
Type=simple
Restart=always    
LimitNOFILE=1048576
ExecStart=/root/localTUN/tun-server -server-port 800 -pub-key=/root/keys/public_key.pem -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -key azumi -mtu 1480 -v -tcpnodelay -keepalive 10 -smux
[Install]
WantedBy=multi-user.target
##### do not copy this ###
chmod u+x /etc/systemd/system/azumilocal.service
systemctl enable /etc/systemd/system/azumilocal.service
systemctl start azumilocal.service
 ```
<div align="right">
   
- نحوه ساختن سرویس ریست
 <div align="left">
   
```
nano /root/reset.sh
# copy this inside #
#!/bin/bash

while true; do
    ping -c 2 30.0.0.1 >/dev/null 2>&1 ##30.0.0.1 is your remote private ip address
    if [ $? -ne 0 ]; then
        systemctl restart azumilocal ## this is localtun service
        systemctl restart strong-azumi1 ## this is for ipsec
    fi
    sleep 5  #time for ping interval check
done

## do not copy this##

nano /etc/systemd/system/azumireset.service
## put this config inside [ This is a sample]##

[Unit]
Description=Azumi local Service reset
After=network.target

[Service]
Type=simple
Restart=always    
LimitNOFILE=1048576
ExecStart=/root/reset.sh
[Install]
WantedBy=multi-user.target

##### do not copy this ###
chmod u+x /etc/systemd/system/azumireset.service
systemctl enable /etc/systemd/system/azumireset.service
systemctl start azumireset.service
systemctl status azumireset.service
```
 </details>
</div>
<div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image">دایرکت لوکال تانل پرایوت ایپی 6 - public ipv4 </strong></summary>
```
  - کامند های سرور (خارج)
  - برای verbose لاگ از کامند -v استفاده نمایید 
 <div align="left">
   
```
./tun-server -server-port 800 -pub-key=/root/keys/public_key.pem -server-private 2001:db8::1 -client-private 2001:db8::2 -subnet 64 -device tun2 -key azumi -mtu 1380 -v -tcpnodelay -keepalive 10 -smux

```
<div align="right">
  
- کامند های کلاینت (ایران)
- برای verbose لاگ از کامند -v استفاده نمایید
 <div align="left">
   
```
./tun-client -server-addr KHAREJ_IPV4 -server-port 800 -priv-key=/root/private_key.pem -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1280 -tcpnodelay -smux -keepalive 10 -v
```
<div align="right">
  
- نحوه ساختن سرویس
 <div align="left">
   
```
nano /etc/systemd/system/azumilocal.service
## put this config inside [ This is a sample]##

[Unit]
Description=Azumi local Service
After=network.target

[Service]
Type=simple
Restart=always    
LimitNOFILE=1048576
ExecStart=/root/localTUN/tun-client -server-addr KHAREJ_IPV4 -server-port 800 -priv-key=/root/private_key.pem -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1280 -v -tcpnodelay -smux -keepalive 10
   

[Install]
WantedBy=multi-user.target
##### do not copy this ###

chmod u+x /etc/systemd/system/azumilocal.service
systemctl enable /etc/systemd/system/azumilocal.service
systemctl start azumilocal.service
 ```
<div align="right">
   
- نحوه ساختن سرویس ریست
 <div align="left">
   
```
nano /root/reset.sh
# copy this inside #
#!/bin/bash
while true; do
    ping -c 2 2001:db8::1 >/dev/null 2>&1 ##2001:db8::1 is your remote private ip address
    if [ $? -ne 0 ]; then
        systemctl restart azumilocal ## this is localtun service
        systemctl restart strong-azumi1 ## this is for ipsec
    fi
    sleep 5  #time for ping interval check
done
## do not copy this##

nano /etc/systemd/system/azumireset.service
## put this config inside [ This is a sample]##

[Unit]
Description=Azumi local Service reset
After=network.target

[Service]
Type=simple
Restart=always    
LimitNOFILE=1048576
ExecStart=/root/reset.sh
[Install]
WantedBy=multi-user.target

##### do not copy this ###
chmod u+x /etc/systemd/system/azumireset.service
systemctl enable /etc/systemd/system/azumireset.service
systemctl start azumireset.service
systemctl status azumireset.service
```
 </details>
</div>
<div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image">دایرکت لوکال تانل پرایوت ایپی 6 - public ipv4 </strong></summary>
```
 </details>
</div>
<div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image">دایرکت لوکال تانل پرایوت ایپی 4 - public ipv6 </strong></summary>

  - کامند های سرور (خارج)
  - برای verbose لاگ از کامند -v استفاده نمایید 
 <div align="left">
   
```
./tun-server_amd64 -server-port 800 -pub-key=/root/keys/public_key.pem -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -key azumi -mtu 1480 -v -tcpnodelay -smux -keepalive 10
```
<div align="right">
  
- کامند های کلاینت (ایران)
- برای verbose لاگ از کامند -v استفاده نمایید 
 <div align="left">
   
```
./tun-client -server-addr KHAREJ_IPV6 -server-port 800 -priv-key=/root/private_key.pem -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1280 -v -tcpnodelay -smux -keepalive 10
```
<div align="right">
  
- نحوه ساختن سرویس
 <div align="left">
   
```
nano /etc/systemd/system/azumilocal.service
## put this config inside [ This is a sample]##

[Unit]
Description=Azumi local Service
After=network.target

[Service]
Type=simple
Restart=always    
LimitNOFILE=1048576
ExecStart=/root/localTUN/tun-client -server-addr KHAREJ_IPV6 -server-port 800 -priv-key=/root/private_key.pem -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1380 -v -tcpnodelay -smux -keepalive 10
   

[Install]
WantedBy=multi-user.target
##### do not copy this ###
chmod u+x /etc/systemd/system/azumilocal.service
systemctl enable /etc/systemd/system/azumilocal.service
systemctl start azumilocal.service
 ```
<div align="right">
   
- نحوه ساختن سرویس ریست
 <div align="left">
   
```
nano /root/reset.sh
# copy this inside #
#!/bin/bash
while true; do
    ping -c 2 30.0.0.1 >/dev/null 2>&1 ##30.0.0.1 is your remote private ip address
    if [ $? -ne 0 ]; then
        systemctl restart azumilocal ## this is localtun service
        systemctl restart strong-azumi1 ## this is for ipsec
    fi
    sleep 5  #time for ping interval check
done

## do not copy this##

nano /etc/systemd/system/azumireset.service
## put this config inside [ This is a sample]##

[Unit]
Description=Azumi local Service reset
After=network.target

[Service]
Type=simple
Restart=always    
LimitNOFILE=1048576
ExecStart=/root/reset.sh
[Install]
WantedBy=multi-user.target

##### do not copy this ###
chmod u+x /etc/systemd/system/azumireset.service
systemctl enable /etc/systemd/system/azumireset.service
systemctl start azumireset.service
systemctl status azumireset.service
```
 </details>
</div>
<div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image">دایرکت لوکال تانل پرایوت ایپی 6 - public ipv6 </strong></summary>

  - کامند های سرور (خارج)
 <div align="left">
   
```
./tun-server -server-port 800 -pub-key=/root/keys/public_key.pem -server-private 2001:db8::1 -client-private 2001:db8::2 -subnet 64 -device tun2 -key azumi -mtu 1380 -v -tcpnodelay -smux -keepalive 10
```
<div align="right">
  
- کامند های کلاینت (ایران)
- برای verbose لاگ از کامند -v استفاده نمایید
 <div align="left">
   
```
./tun-client -server-addr [KHAREJ_IPV6] -server-port 800 -priv-key=/root/private_key.pem -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1380 -v -tcpnodelay -smux -keepalive 10
```
 <div align="right">
   
- نحوه ساختن سرویس
 <div align="left">
   
```
nano /etc/systemd/system/azumilocal.service
## put this config inside [ This is a sample]##

[Unit]
Description=Azumi local Service
After=network.target

[Service]
Type=simple
Restart=always    
LimitNOFILE=1048576
ExecStart=/root/localTUN/tun-client -server-addr [KHAREJ_IPV6] -server-port 800 -priv-key=/root/private_key.pem -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1280 -v -tcpnodelay -smux -keepalive 10
   

[Install]
WantedBy=multi-user.target
##### do not copy this ###
chmod u+x /etc/systemd/system/azumilocal.service
systemctl enable /etc/systemd/system/azumilocal.service
systemctl start azumilocal.service
 ```
<div align="right">
   
- نحوه ساختن سرویس ریست
 <div align="left">
   
```
nano /root/reset.sh
# copy this inside #
#!/bin/bash

while true; do
    ping -c 2 2001:db8::1 >/dev/null 2>&1 ##2001:db8::1 is your remote private ip address
    if [ $? -ne 0 ]; then
        systemctl restart azumilocal ## this is localtun service
        systemctl restart strong-azumi1 ## this is for ipsec
    fi
    sleep 5  #time for ping interval check
done
## do not copy this##

nano /etc/systemd/system/azumireset.service
## put this config inside [ This is a sample]##

[Unit]
Description=Azumi local Service reset
After=network.target

[Service]
Type=simple
Restart=always    
LimitNOFILE=1048576
ExecStart=/root/reset.sh
[Install]
WantedBy=multi-user.target

##### do not copy this ###
chmod u+x /etc/systemd/system/azumireset.service
systemctl enable /etc/systemd/system/azumireset.service
systemctl start azumireset.service
systemctl status azumireset.service
```
 </details>
</div>
<div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image">ریورس لوکال تانل پرایوت ایپی 4 - public ipv4 </strong></summary>

  - کامند های سرور ( ایران)
  - برای verbose لاگ از کامند -v استفاده نمایید
 <div align="left">
   
```
./tun-server -server-port 800 -pub-key=/root/keys/public_key.pem -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -key azumi -mtu 1380 -tcpnodelay -keepalive 10 -v -smux
```
<div align="right">
  
- کامند های کلاینت (خارج)
- برای verbose لاگ از کامند -v استفاده نمایید 
 <div align="left">
   
```
./tun-client -server-addr IRAN_IPV4 -server-port 800 -priv-key=/root/private_key.pem -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1280 -tcpnodelay -keepalive 10 -v -smux
```
<div align="right">
  
- نحوه ساختن سرویس
 <div align="left">
   
```
nano /etc/systemd/system/azumilocal.service
## put this config inside [ This is a sample]##

[Unit]
Description=Azumi local Service
After=network.target

[Service]
Type=simple
Restart=always    
LimitNOFILE=1048576
ExecStart=/root/localTUN/tun-server -server-port 800 -pub-key=/root/keys/public_key.pem -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -key azumi -mtu 1380 -tcpnodelay -keepalive 10 -smux -v
   

[Install]
WantedBy=multi-user.target
##### do not copy this ###
chmod u+x /etc/systemd/system/azumilocal.service
systemctl enable /etc/systemd/system/azumilocal.service
systemctl start azumilocal.service
 ```
<div align="right">
   
- نحوه ساختن سرویس ریست
 <div align="left">
   
```
nano /root/reset.sh
# copy this inside #
#!/bin/bash

while true; do
    ping -c 2 30.0.0.2 >/dev/null 2>&1 ##30.0.0.2 is your remote private ip address
    if [ $? -ne 0 ]; then
        systemctl restart azumilocal ## this is localtun service
        systemctl restart strong-azumi1 ## this is for ipsec
    fi
    sleep 5  #time for ping interval check
done
## do not copy this##

nano /etc/systemd/system/azumireset.service
## put this config inside [ This is a sample]##

[Unit]
Description=Azumi local Service reset
After=network.target

[Service]
Type=simple
Restart=always    
LimitNOFILE=1048576
ExecStart=/root/reset.sh
[Install]
WantedBy=multi-user.target

##### do not copy this ###
chmod u+x /etc/systemd/system/azumireset.service
systemctl enable /etc/systemd/system/azumireset.service
systemctl start azumireset.service
systemctl status azumireset.service
```
 </details>
</div>
<div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image">ریورس لوکال تانل پرایوت ایپی 6 - public ipv4 </strong></summary>

  - کامند های سرور (ایران)
  - برای verbose لاگ از کامند -v استفاده نمایید 
 <div align="left">
   
```
./tun-server -server-port 800 -pub-key=/root/keys/public_key.pem -server-private 2001:db8::1 -client-private 2001:db8::2 -subnet 64 -device tun2 -key azumi -mtu 1380 -v -tcpnodelay -smux -keepalive 10
```
<div align="right">
  
- کامند های کلاینت (خارج)
- برای verbose لاگ از کامند -v استفاده نمایید .
 <div align="left">
   
```
./tun-client -server-addr IRAN_IPV4 -server-port 800 -priv-key=/root/private_key.pem -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1280 -v -tcpnodelay -smux -keepalive 10
```
<div align="right">
  
- نحوه ساختن سرویس
 <div align="left">
   
```
nano /etc/systemd/system/azumilocal.service
## put this config inside [ This is a sample]##

[Unit]
Description=Azumi local Service
After=network.target

[Service]
Type=simple
Restart=always    
LimitNOFILE=1048576
ExecStart=/root/localTUN/tun-client -server-addr IRAN_IPV4 -server-port 800 -priv-key=/root/private_key.pem -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -mtu 1280 -v -tcpnodelay -smux -keepalive 10
   

[Install]
WantedBy=multi-user.target
##### do not copy this ###
chmod u+x /etc/systemd/system/azumilocal.service
systemctl enable /etc/systemd/system/azumilocal.service
systemctl start azumilocal.service
 ```
<div align="right">
   
- نحوه ساختن سرویس ریست
 <div align="left">
   
```
nano /root/reset.sh
# copy this inside #
#!/bin/bash

while true; do
    ping -c 2 2001:db8::1 >/dev/null 2>&1 ##2001:db8::1 is your remote private ip address
    if [ $? -ne 0 ]; then
        systemctl restart azumilocal ## this is localtun service
        systemctl restart strong-azumi1 ## this is for ipsec
    fi
    sleep 5  #time for ping interval check
done

## do not copy this##

nano /etc/systemd/system/azumireset.service
## put this config inside [ This is a sample]##

[Unit]
Description=Azumi local Service reset
After=network.target

[Service]
Type=simple
Restart=always    
LimitNOFILE=1048576
ExecStart=/root/reset.sh
[Install]
WantedBy=multi-user.target

##### do not copy this ###
chmod u+x /etc/systemd/system/azumireset.service
systemctl enable /etc/systemd/system/azumireset.service
systemctl start azumireset.service
systemctl status azumireset.service
```
 </details>
</div>
<div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image">ریورس لوکال تانل پرایوت ایپی 4 - public ipv6 </strong></summary>

  - کامند های سرور (ایران)
 <div align="left">
   
```
./tun-server -server-port 800 -pub-key=/root/keys/public_key.pem -server-private 30.0.0.1 -client-private 30.0.0.2 -pub-key=/root/keys/public_key.pem -subnet 24 -device tun2 -key azumi -mtu 1380 -v -tcpnodelay -smux -keepalive 10
```
<div align="right">
  
- کامند های کلاینت (خارج)
 <div align="left">
   
```
./tun-client -server-addr [IRAN_IPV6] -server-port 800 -priv-key=/root/private_key.pem -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1280 -v -tcpnodelay -smux -keepalive 10
```
<div align="right">
  
- نحوه ساختن سرویس
 <div align="left">
   
```
nano /etc/systemd/system/azumilocal.service
## put this config inside [ This is a sample]##

[Unit]
Description=Azumi local Service
After=network.target

[Service]
Type=simple
Restart=always    
LimitNOFILE=1048576
ExecStart=/root/localTUN/tun-client -server-addr [IRAN_IPV6] -server-port 800 -priv-key=/root/private_key.pem -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1280 -v -tcpnodelay -smux -keepalive 10
   

[Install]
WantedBy=multi-user.target
##### do not copy this ###
chmod u+x /etc/systemd/system/azumilocal.service
systemctl enable /etc/systemd/system/azumilocal.service
systemctl start azumilocal.service
 ```
<div align="right">
   
- نحوه ساختن سرویس ریست
 <div align="left">
   
```
nano /root/reset.sh
# copy this inside #
#!/bin/bash

while true; do
    ping -c 2 30.0.0.1 >/dev/null 2>&1 ##30.0.0.1 is your remote private ip address
    if [ $? -ne 0 ]; then
        systemctl restart azumilocal ## this is localtun service
        systemctl restart strong-azumi1 ## this is for ipsec
    fi
    sleep 5  #time for ping interval check
done

## do not copy this##

nano /etc/systemd/system/azumireset.service
## put this config inside [ This is a sample]##

[Unit]
Description=Azumi local Service reset
After=network.target

[Service]
Type=simple
Restart=always    
LimitNOFILE=1048576
ExecStart=/root/reset.sh
[Install]
WantedBy=multi-user.target

##### do not copy this ###
chmod u+x /etc/systemd/system/azumireset.service
systemctl enable /etc/systemd/system/azumireset.service
systemctl start azumireset.service
systemctl status azumireset.service
```
 </details>
</div>
<div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image">ریورس لوکال تانل پرایوت ایپی 6 - public ipv6 </strong></summary>

  - کامند های سرور (ایران)
  - برای verbose لاگ از کامند -v استفاده نمایید 
 <div align="left">
   
```
./tun-server -server-port 800 -pub-key=/root/keys/public_key.pem -server-private 2001:db8::1 -client-private 2001:db8::2 -subnet 64 -device tun2 -key azumi -mtu 1480 -v -tcpnodelay -keepalive 10 -smux
```
<div align="right">
  
- کامند های کلاینت (خارج)
 <div align="left">
   
```
./tun-client -server-addr [IRAN_IPV6] -server-port 800 -priv-key=/root/private_key.pem -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1280 -verbose -tcpnodelay -keepalive 10
```
<div align="right">
  
- نحوه ساختن سرویس
 <div align="left">
   
```
nano /etc/systemd/system/azumilocal.service
## put this config inside [ This is a sample]##

[Unit]
Description=Azumi local Service
After=network.target

[Service]
Type=simple
Restart=always    
LimitNOFILE=1048576
ExecStart=/root/localTUN/tun-client -server-addr [IRAN_IPV6] -server-port 800 -priv-key=/root/private_key.pem -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1280 -verbose -tcpnodelay -keepalive 10
   

[Install]
WantedBy=multi-user.target
##### do not copy this ###
chmod u+x /etc/systemd/system/azumilocal.service
systemctl enable /etc/systemd/system/azumilocal.service
systemctl start azumilocal.service
 ```
<div align="right">
   
- نحوه ساختن سرویس ریست
 <div align="left">
   
```
nano /root/reset.sh
# copy this inside #
#!/bin/bash

while true; do
    ping -c 2 2001:db8::1 >/dev/null 2>&1 ##2001:db8::1 is your remote private ip address
    if [ $? -ne 0 ]; then
        systemctl restart azumilocal ## this is localtun service
        systemctl restart strong-azumi1 ## this is for ipsec
    fi
    sleep 5  #time for ping interval check
done
## do not copy this##

nano /etc/systemd/system/azumireset.service
## put this config inside [ This is a sample]##

[Unit]
Description=Azumi local Service reset
After=network.target

[Service]
Type=simple
Restart=always    
LimitNOFILE=1048576
ExecStart=/root/reset.sh
[Install]
WantedBy=multi-user.target

##### do not copy this ###
chmod u+x /etc/systemd/system/azumireset.service
systemctl enable /etc/systemd/system/azumireset.service
systemctl start azumireset.service
systemctl status azumireset.service
```
 </details>
</div>



