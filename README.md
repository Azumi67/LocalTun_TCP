![R (2)](https://github.com/Azumi67/PrivateIP-Tunnel/assets/119934376/a064577c-9302-4f43-b3bf-3d4f84245a6f)
نام پروژه : لوکال تانل به صورت دایرکت یا ریورس - Tun interface و بر روی پورت TCP
---------------------------------------------------------------

**اسکریپت در پروسه نوشته شدن می‌باشد**

**بدون ipsec استفاده نشود. این پروژه برای یادگیری خودم میباشد و نیاز به کامل شدن در زمان دارد. لطفا استفاده نکنید**

**فعلا این پروژه در حالت پیش نویس است و باید صبر کنید تا کامل شود**

**این مورد xtaci/smux به پروژه اضافه شد**

**این مورد tcpnodelay و logrus اضافه شد**

**گزینه heartbeat و heartbeat interval و ping interval اضافه شد**

**گزینه سرویس نیم برای ریستارت کلاینت در صورتی که سرور و کلاینت به هم دیگر پینگ نداشتند، اضافه شد به طور مثال service-name azumilocal-**

**بعدا انتخاب Ipsec encryption یا chacha را خواهید داشت و اضافه خواهد شد. worker اضافه خواهد شد.**

**در حال حاضر encapsulation استریم های دیتا tcp بر روی اینترفیس یا device tun انجام میشود. این عمل شامل اضافه کردن packet length میباشد که کلاینت یا دریافت کننده دقیقا همان مقدار از دیتا را از استریم read و دریافت میکند.(ممکن است در این قسمت همچنان مشکلاتی باشد)**




![check](https://github.com/Azumi67/PrivateIP-Tunnel/assets/119934376/13de8d36-dcfe-498b-9d99-440049c0cf14)
**امکانات**
- امکان لوکال تانل بین سرور و کلاینت به صورت دایرکت یا ریورس ( برای سرور های محدود)
- استفاده از ایپی پرایوت های ساخته شده در Tun interface برای تانل اصلی یا پورت فوروارد
- امکان اتصال بین سرور و کلاینت بر روی پورت TCP
- امکان اتصال سرور و کلاینت با پابلیک ایپی 4 یا native
- امکان انتخاب ایپی پرایوت انتخابی خودتان هم به صورت پرایوت ایپی 4 یا پرایوت ایپی 6
- امکان انتخاب subnet mask برای پرایوت ایپی های ساخته
- امکان وارد کردن mtu به صورت manual
- دارای smux
- دارای heartbeat & heartbeat interval
- دارای ping interval و انتخاب نام سرویس برای ریست کردن کلاینت در صورت disconnection
- دارای encryption های IPSEC یا chacha20 (به زودی)
- دارای worker (به زودی)
- دارای authentication key برای ارتباط بین سرور و کلاینت
- دارای verbose برای نمایش لاگ (خطا)
- مناسب برای ترکیب با IPSEC > لینک : https://github.com/Azumi67/6TO4-GRE-IPIP-SIT
- اتصال بین چندین سرور ایران و خارج ( بعدا در اسکریپت)
- بعدا اسکریپت ساخته میشود ( در اسکریپت با IPsec هم میتوان ترکیب کرد)- باید با IPSEC استفاده شود
-----------------------
<div align="right">
  <details>
    <summary><strong>توضیحات</strong></summary>
  
------------------------------------ 
 <div align="right">
   
- از طریق tcp دو سرور به هم وصل میشوند و از طریق اینترفیس tun و پرایوت ایپی به هم دیگه متصل خواهند بود. encapsulation & decapsulation هم زمان انجام میشود.
- هدف نوشتن این برنامه برای این بوده است که از طریق پورت tcp و tun interface،‌ دو سرور به هم متصل شوند و از پرایوت آی‌پی های آنها برای تانل استقاده کرد و محدودیت بعضی از سرور ها به صورت ریورس برطرف شود.
- به عبارتی شما به صورت ریورس، یک لوکال ایپی دریافت میکنید و سپس از آن پرایوت ایپی ها برای دایرکت تانل، پورت فوروارد یا ریورس استفاده مینمایید. 
- پس از انجام تانل‌ لوکال به صورت دایرکت یا ریورس، به طور مثال میتوانید از پورت فوروارد استفاده نمایید یا مثلا دایرکت تانل چیزل استفاده نمایید یا ریورس.
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
  wget https://github.com/Azumi67/LocalTun_TCP/releases/download/v1.3/amd64.zip
  unzip amd64.zip -d /root/localTUN
  cd localTUN
  chmod +x tun-server_amd64   
  chmod +x tun-client_amd64   
 ```
 </details>
</div>
<div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image">دایرکت لوکال تانل پرایوت ایپی 4 - public ipv4 </strong></summary>

  - کامند های سرور (خارج)
 <div align="left">
   
```
./tun-server_amd64 -server-port 800 -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -key azumi -mtu 1480 -verbose true -smux true -heartbeat true -tcp-nodelay true -service-name azumilocal
   
```
<div align="right">
  
- کامند های کلاینت (ایران)
 <div align="left">
   
```
./tun-client_amd64 -server-addr KHAREJ_IPV4 -server-port 800 -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1400 -verbose true -smux true -heartbeat true -tcp-nodelay true -service-name azumilocal
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
ExecStart=/root/localTUN/tun-server_amd64 -server-port 800 -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -key azumi -mtu 1480 -verbose true -smux true -tcp-nodelay true -heartbeat true -service-name azumilocal

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
./tun-server_amd64 -server-port 800 -server-private 2001:db8::1 -client-private 2001:db8::2 -subnet 64 -device tun2 -key azumi -mtu 1480 -verbose true -smux true -tcp-nodelay true -heartbeat true -service-name azumilocal
-heartbeat-interval 30
```
<div align="right">
  
- کامند های کلاینت (ایران)
 <div align="left">
   
```
./tun-client_amd64 -server-addr KHAREJ_IPV4 -server-port 800 -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1400 -verbose true -smux true -tcp-nodelay true -heartbeat true -service-name azumilocal -heartbeat-interval
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
ExecStart=/root/localTUN/tun-client_amd64 -server-addr KHAREJ_IPV4 -server-port 800 -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1400 -verbose true -smux true -tcp-nodelay true -heartbeat true -heartbeat-interval 30 -service-name azumilocal
   

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
./tun-server_amd64 -server-port 800 -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -key azumi -mtu 1480 -verbose true -smux true -tcp-nodelay true -heartbeat true -heartbeat-interval 30  
```
<div align="right">
  
- کامند های کلاینت (ایران)
 <div align="left">
   
```
./tun-client_amd64 -server-addr KHAREJ_IPV6 -server-port 800 -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1400 -verbose true -smux true -tcp-nodelay true -heartbeat true -heartbeat-interval 30 -service-name azumilocal
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
ExecStart=/root/localTUN/tun-client_amd64 -server-addr KHAREJ_IPV6 -server-port 800 -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1400 -verbose true -smux true -tcp-nodelay true -heartbeat true -heartbeat-interval 30 -service-name azumilocal
   

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
./tun-server_amd64 -server-port 800 -server-private 2001:db8::1 -client-private 2001:db8::2 -subnet 64 -device tun2 -key azumi -mtu 1480 -verbose true -smux true -tcp-nodelay true -heartbeat true -heartbeat-interval 30
```
<div align="right">
  
- کامند های کلاینت (ایران)
 <div align="left">
   
```
./tun-client_amd64 -server-addr KHAREJ_IPV6 -server-port 800 -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1400 -verbose true -smux true -tcp-nodelay true -heartbeat true -heartbeat-interval 30 -service-name azumilocal
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
ExecStart=/root/localTUN/tun-client_amd64 -server-addr KHAREJ_IPV6 -server-port 800 -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1400 -verbose true -smux true -tcp-nodelay true -heartbeat true -heartbeat-interval 30 -service-name azumilocal
   

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
./tun-server_amd64 -server-port 800 -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -key azumi -mtu 1480 -verbose true -smux true -tcp-nodelay true -heartbeat true -heartbeat-interval 30    
```
<div align="right">
  
- کامند های کلاینت (خارج)
 <div align="left">
   
```
./tun-client_amd64 -server-addr IRAN_IPV4 -server-port 800 -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1400 -verbose true -smux true -tcp-nodelay true -heartbeat true -heartbeat-interval 30 -service-name azumilocal
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
ExecStart=/root/localTUN/tun-server_amd64 -server-port 800 -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -key azumi -mtu 1480 -verbose true -smux true -tcp-nodelay true -heartbeat true -heartbeat-interval 30
   

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
./tun-server_amd64 -server-port 800 -server-private 2001:db8::1 -client-private 2001:db8::2 -subnet 64 -device tun2 -key azumi -mtu 1480 -verbose true -smux true -tcp-nodelay true -heartbeat true -heartbeat-interval 30
```
<div align="right">
  
- کامند های کلاینت (خارج)
 <div align="left">
   
```
./tun-client_amd64 -server-addr IRAN_IPV4 -server-port 800 -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1400 -verbose true -smux true -tcp-nodelay true -heartbeat true -heartbeat-interval 30 -service-name azumilocal
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
ExecStart=/root/localTUN/tun-client_amd64 -server-addr IRAN_IPV4 -server-port 800 -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1400 -verbose true -smux true -tcp-nodelay true -heartbeat true -heartbeat-interval 30 -service-name azumilocal
   

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
./tun-server_amd64 -server-port 800 -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -key azumi -mtu 1480 -verbose true -smux true -tcp-nodelay true -heartbeat true -heartbeat-interval 30
```
<div align="right">
  
- کامند های کلاینت (خارج)
 <div align="left">
   
```
./tun-client_amd64 -server-addr IRAN_IPV6 -server-port 800 -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1400 -verbose true -smux true -tcp-nodelay true -heartbeat true -heartbeat-interval 30 -service-name azumilocal
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
ExecStart=/root/localTUN/tun-client_amd64 -server-addr IRAN_IPV6 -server-port 800 -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1400 -verbose true -smux true -tcp-nodelay true -heartbeat true -heartbeat-interval 30 -service-name azumilocal
   

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
./tun-server_amd64 -server-port 800 -server-private 2001:db8::1 -client-private 2001:db8::2 -subnet 64 -device tun2 -key azumi -mtu 1480 -verbose true -smux true -tcp-nodelay true -heartbeat true -heartbeat-interval 30
```
<div align="right">
  
- کامند های کلاینت (خارج)
 <div align="left">
   
```
./tun-client_amd64 -server-addr IRAN_IPV6 -server-port 800 -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1400 -verbose true -smux true -tcp-nodelay true -heartbeat true -heartbeat-interval 30 -service-name azumilocal
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
ExecStart=/root/localTUN/tun-client_amd64 -server-addr IRAN_IPV6 -server-port 800 -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1400 -verbose true -smux true -tcp-nodelay true -heartbeat true -heartbeat-interval 30 -service-name azumilocal
   

[Install]
WantedBy=multi-user.target
##### do not copy this ###
chmod u+x /etc/systemd/system/azumilocal.service
systemctl enable /etc/systemd/system/azumilocal.service
systemctl start azumilocal.service
 ```
 </details>
</div>



