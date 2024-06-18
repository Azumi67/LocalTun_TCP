![R (2)](https://github.com/Azumi67/PrivateIP-Tunnel/assets/119934376/a064577c-9302-4f43-b3bf-3d4f84245a6f)
نام پروژه : لوکال تانل به صورت دایرکت یا ریورس - Tun interface و بر روی پورت TCP
---------------------------------------------------------------

![check](https://github.com/Azumi67/PrivateIP-Tunnel/assets/119934376/13de8d36-dcfe-498b-9d99-440049c0cf14)
**امکانات**
- امکان لوکال تانل بین سرور و کلاینت به صورت دایرکت یا ریورس ( برای سرور های محدود)
- استفاده از ایپی پرایوت های ساخته شده در Tun interface برای تانل اصلی یا پورت فوروارد
- امکان اتصال بین سرور و کلاینت بر روی پورت TCP
- امکان اتصال سرور و کلاینت با پابلیک ایپی 4 یا native
- امکان انتخاب ایپی پرایوت انتخابی خودتان هم به صورت پرایوت ایپی 4 یا پرایوت ایپی 6
- امکان انتخاب subnet mask برای پرایوت ایپی های ساخته
- امکان وارد کردن mtu به صورت manual
- دارای authentication key برای ارتباط بین سرور و کلاینت
- مناسب برای ترکیب با IPSEC > لینک : https://github.com/Azumi67/6TO4-GRE-IPIP-SIT
- اتصال بین چندین سرور ایران و خارج ( بعدا در اسکریپت)
- بعدا اسکریپت ساخته میشود
-----------------------
<div align="right">
  <details>
    <summary><strong>توضیحات</strong></summary>
  
------------------------------------ 
 <div align="right">
   
- از طریق tcp دو سرور به هم وصل میشوند و از طریق اینترفیس tun و پرایوت ایپی به هم دیگه متصل خواهند بود. encapsulation & decapsulation هم زمان انجام میشود.
- هدف نوشتن این برنامه برای این بوده است که از طریق پورت tcp و tun interface،‌ دو سرور به هم متصل شوند و از پرایوت آی‌پی های آنها برای تانل استقاده کرد و محدودیت بعضی از سرور ها به صورت ریورس برطرف شود.
- به عبارتی شما به صورت ریورس، یک لوکال ایپی دریافت میکنید و سپس از آن پرایوت ایپی ها برای دایرکت تانل، پورت فوروارد یا ریورس استفاده مینمایید. 
- پس از انجام تانل‌ لوکال به صورت دایرکت یا ریورس، به طور مثال میتوانید از این پروژه https://github.com/Azumi67/PortforwardSec برای پورت فوروارد استفاده نمایید یا مثلا دایرکت تانل چیزل استفاده نمایید.
- در روش ریورس، سرور اصلی میتواند ایران باشد و کلاینت خارج و در روش دایرکت، سرور اصلی میتواند خارج باشد و کلاینت ایران. بدین صورت میتوان تانل لوکالی بر روی سرور های خارج محدود در ان سرور ایران(به صورت ریورس) هم ایجاد کرد.
- با ایپی 4 سرور و هم با ایپی 6 سرور و کلاینت میشود که وصل شد .
- پورت تنها برای ارتباط بین سرور و کلاینت میباشد و شما تنها باید از پرایوت ایپی ها برای تانل اصلی استفاده نمایید.
- اول دستورات سرور را اجرا کنید و سپس دستورات کلاینت . میتوانید هم به صورت دایرکت یا ریورس انجام دهید. یعنی سرور اصلی خارج و کلاینت ایران و یا سرور اصلی ایران و کلاینت خارج باشد
- میتوانید تست کنید و بهم اطلاع بدید اما برای استفاده باید صبر کنید تا داخل اسکریپت برای ترکیب با ipsec آورده شود.
- کلید authentication اضافه شد که در کلاینت خوانده میشود و در سرور احراز هویت میشود
  </details>
</div>

---------------------

  ![6348248](https://github.com/Azumi67/PrivateIP-Tunnel/assets/119934376/398f8b07-65be-472e-9821-631f7b70f783)
**روش اجرا**
-

 <div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image">نصب </strong></summary>
  
<div align="left">
  
```
  apt update -y
  apt install wget -y
  apt install unzip -y
  wget https://github.com/Azumi67/LocalTun_TCP/releases/download/v1.1/tun-server.zip
  wget https://github.com/Azumi67/LocalTun_TCP/releases/download/v1.1/tun-client.zip
  unzip tun-server.zip -d /root/localTUN
  unzip tun-client.zip -d /root/localTUN
  cd localTUN
  chmod +x tun-server_amd64   << for amd64
  chmod +x tun-client_amd64   << for amd64
 ```
 </details>
</div>
<div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image">دایرکت لوکال تانل پرایوت ایپی 4 - public ipv4 </strong></summary>

  - کامند های سرور (خارج)
 <div align="left">
   
```
./tun-server_amd64 -server-port 800 -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -key azumi -mtu 1480
   
```
<div align="right">
  
- کامند های کلاینت (ایران)
 <div align="left">
   
```
./tun-client_amd64 -server-addr KHAREJ_IPV4 -server-port 800 -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1400
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
ExecStart=/root/localTUN/tun-server_amd64 -server-port 800 -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -key azumi -mtu 1480
   

[Install]
WantedBy=multi-user.target
##### do not copy this ###
chmod u+x /etc/systemd/system/azumilocal.service
systemctl enable /etc/systemd/system/azumilocal.service
systemctl start azumilocal.service
 ```
 </details>
</div>
<div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image">دایرکت لوکال تانل پرایوت ایپی 6 - public ipv4 </strong></summary>

  - کامند های سرور (خارج)
 <div align="left">
   
```
./tun-server_amd64 -server-port 800 -server-private 2001:db8::1 -client-private 2001:db8::2 -subnet 64 -device tun2 -key azumi -mtu 1480   
```
<div align="right">
  
- کامند های کلاینت (ایران)
 <div align="left">
   
```
./tun-client_amd64 -server-addr KHAREJ_IPV4 -server-port 800 -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1400
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
ExecStart=/root/localTUN/tun-client_amd64 -server-addr KHAREJ_IPV4 -server-port 800 -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1400
   

[Install]
WantedBy=multi-user.target
##### do not copy this ###
chmod u+x /etc/systemd/system/azumilocal.service
systemctl enable /etc/systemd/system/azumilocal.service
systemctl start azumilocal.service
 ```
 </details>
</div>
<div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image">دایرکت لوکال تانل پرایوت ایپی 4 - public ipv6 </strong></summary>

  - کامند های سرور (خارج)
 <div align="left">
   
```
./tun-server_amd64 -server-port 800 -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -key azumi -mtu 1480   
```
<div align="right">
  
- کامند های کلاینت (ایران)
 <div align="left">
   
```
./tun-client_amd64 -server-addr KHAREJ_IPV6 -server-port 800 -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1400
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
ExecStart=/root/localTUN/tun-client_amd64 -server-addr KHAREJ_IPV6 -server-port 800 -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1400
   

[Install]
WantedBy=multi-user.target
##### do not copy this ###
chmod u+x /etc/systemd/system/azumilocal.service
systemctl enable /etc/systemd/system/azumilocal.service
systemctl start azumilocal.service
 ```
 </details>
</div>
<div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image">دایرکت لوکال تانل پرایوت ایپی 6 - public ipv6 </strong></summary>

  - کامند های سرور (خارج)
 <div align="left">
   
```
./tun-server_amd64 -server-port 800 -server-private 2001:db8::1 -client-private 2001:db8::2 -subnet 64 -device tun2 -key azumi -mtu 1480
```
<div align="right">
  
- کامند های کلاینت (ایران)
 <div align="left">
   
```
./tun-client_amd64 -server-addr KHAREJ_IPV6 -server-port 800 -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1400
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
ExecStart=/root/localTUN/tun-client_amd64 -server-addr KHAREJ_IPV6 -server-port 800 -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1400
   

[Install]
WantedBy=multi-user.target
##### do not copy this ###
chmod u+x /etc/systemd/system/azumilocal.service
systemctl enable /etc/systemd/system/azumilocal.service
systemctl start azumilocal.service
 ```
 </details>
</div>
<div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image">ریورس لوکال تانل پرایوت ایپی 4 - public ipv4 </strong></summary>

  - کامند های سرور ( ایران)
 <div align="left">
   
```
./tun-server_amd64 -server-port 800 -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -key azumi -mtu 1480   
```
<div align="right">
  
- کامند های کلاینت (خارج)
 <div align="left">
   
```
./tun-client_amd64 -server-addr IRAN_IPV4 -server-port 800 -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1400
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
ExecStart=/root/localTUN/tun-server_amd64 -server-port 800 -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -key azumi -mtu 1480 
   

[Install]
WantedBy=multi-user.target
##### do not copy this ###
chmod u+x /etc/systemd/system/azumilocal.service
systemctl enable /etc/systemd/system/azumilocal.service
systemctl start azumilocal.service
 ```
 </details>
</div>
<div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image">ریورس لوکال تانل پرایوت ایپی 6 - public ipv4 </strong></summary>

  - کامند های سرور (ایران)
 <div align="left">
   
```
./tun-server_amd64 -server-port 800 -server-private 2001:db8::1 -client-private 2001:db8::2 -subnet 64 -device tun2 -key azumi -mtu 1480
```
<div align="right">
  
- کامند های کلاینت (خارج)
 <div align="left">
   
```
./tun-client_amd64 -server-addr IRAN_IPV4 -server-port 800 -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1400
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
ExecStart=/root/localTUN/tun-client_amd64 -server-addr IRAN_IPV4 -server-port 800 -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1400
   

[Install]
WantedBy=multi-user.target
##### do not copy this ###
chmod u+x /etc/systemd/system/azumilocal.service
systemctl enable /etc/systemd/system/azumilocal.service
systemctl start azumilocal.service
 ```
 </details>
</div>
<div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image">ریورس لوکال تانل پرایوت ایپی 4 - public ipv6 </strong></summary>

  - کامند های سرور (ایران)
 <div align="left">
   
```
./tun-server_amd64 -server-port 800 -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -key azumi -mtu 1480
```
<div align="right">
  
- کامند های کلاینت (خارج)
 <div align="left">
   
```
./tun-client_amd64 -server-addr IRAN_IPV6 -server-port 800 -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1400
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
ExecStart=/root/localTUN/tun-client_amd64 -server-addr IRAN_IPV6 -server-port 800 -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1400
   

[Install]
WantedBy=multi-user.target
##### do not copy this ###
chmod u+x /etc/systemd/system/azumilocal.service
systemctl enable /etc/systemd/system/azumilocal.service
systemctl start azumilocal.service
 ```
 </details>
</div>
<div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image">ریورس لوکال تانل پرایوت ایپی 6 - public ipv6 </strong></summary>

  - کامند های سرور (ایران)
 <div align="left">
   
```
./tun-server_amd64 -server-port 800 -server-private 2001:db8::1 -client-private 2001:db8::2 -subnet 64 -device tun2 -key azumi -mtu 1480
```
<div align="right">
  
- کامند های کلاینت (خارج)
 <div align="left">
   
```
./tun-client_amd64 -server-addr IRAN_IPV6 -server-port 800 -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1400
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
ExecStart=/root/localTUN/tun-client_amd64 -server-addr IRAN_IPV6 -server-port 800 -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1400
   

[Install]
WantedBy=multi-user.target
##### do not copy this ###
chmod u+x /etc/systemd/system/azumilocal.service
systemctl enable /etc/systemd/system/azumilocal.service
systemctl start azumilocal.service
 ```
 </details>
</div>



