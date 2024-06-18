# LocalTun_TCP
This implements a TCP tun between a client &amp; server using a TUN interface. The client and server communicate by encapsulating network packets within TCP connections. Encapsulation &amp; decapsulates is Simultaneous.

- فعلا فقط private ip ها تست شده و تستی روی ترافیک انجام نشده است. این گام هی بعدی این پروژه میباشد و برای همین در حال حاضر قابل استفاده نمیباشد
- از طریق tcp دو سرور به هم وصل میشوند و از طریق اینترفیس tun و پرایوت ایپی به هم دیگه متصل خواهند بود. encapsulation & decapsulation هم زمان انجام میشود.
- هدف نوشتن این برنامه برای این بوده است که از طریق پورت tcp و tun interface،‌ دو سرور به هم متصل شوند و از پرایوت آی‌پی های آنها برای تانل استقاده کرد .
- میتوان سرور اصلی ایران باشد و کلاینت خارج باشد و بدین صورت میتوان تانل لوکالی بر روی سرور های خارج محدود هم ایجاد کرد.
- یا میتوانید به صورت دایرکت یعنی سرور اصلی خارج و کلاینت ایران، این تانل را اجرا کنید
- هم میشه با ایپی 4 سرور و هم ایپی 6 سرور وصل شد .
- پورت تنها برای ارتباط بین سرور و کلاینت میباشد و شما تنها باید از پرایوت ایپی ها برای تانل اصلی استفاده نمایید.
- اول دستورات سرور را اجرا کنید و سپس دستورات کلاینت . میتوانید هم به صورت دایرکت یا ریورس انجام دهید. یعنی سرور اصلی خارج و کلاینت ایران و یا سرور اصلی ایران و کلاینت خارج باشد
- میتوانید تست کنید و بهم اطلاع بدید اما برای استفاده باید صبر کنید تا داخل اسکریپت برای ترکیب با ipsec آورده شود.
- کلید authentication اضافه شد که در کلاینت خوانده میشود و در سرور احراز هویت میشود

**usage**

```
  apt update -y
  apt install wget -y
  apt install unzip -y
  wget https://github.com/Azumi67/LocalTun_TCP/releases/download/v1.1/tun-server.zip
  wget https://github.com/Azumi67/LocalTun_TCP/releases/download/v1.1/tun-client.zip
  unzip tun-server.zip -d /root/localTUN
  unzip tun-client.zip -d /root/localTUN
  cd localTUN
  chmod +x tun-server-amd64   << for amd64
  chmod +x tun-client-amd64   << for amd64
 ```
**Server : KHAREJ  , Client: IRAN - [PUBLIC IPV4] - DIRECT**

SERVER & Client IPV4 [ Private IPV4]:
 - Server[kharej] command : ./tun-server-amd64 -server-port 800 -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -key azumi -mtu 1480
 - Client[iran] command : ./tun-client-amd64 -server-addr KHAREJ_IPV4 -server-port 800 -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1480
 - 
SERVER & Client IPV4 [ Private IPV6]:
 - Server command : ./tun-server-amd64 -server-port 800 -server-private 2001:db8::1 -client-private 2001:db8::2 -subnet 64 -device tun2 -key azumi -mtu 1480
 - Client command : ./tun-client-amd64 -server-addr KHAREJ_IPV4 -server-port 800 -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1480
--------------
**Server : KHAREJ  , Client: IRAN [PUBLIC IPV6] - DIRECT**
SERVER & Client IPV6 [ Private IPV4]:
 - Server[kharej] command : ./tun-server-amd64 -server-port 800 -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -key azumi -mtu 1480
 - Client[iran] command : ./tun-client-amd64 -server-addr KHAREJ_IPV6 -server-port 800 -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1480

SERVER & Client IPV6 [ Private IPV6]:
 - Server command : ./tun-server-amd64 -server-port 800 -server-private 2001:db8::1 -client-private 2001:db8::2 -subnet 64 -device tun2 -key azumi -mtu 1480
 - Client command : ./tun-client-amd64 -server-addr KHAREJ_IPV6 -server-port 800 -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1480

-----------------
**Server : IRAN  , Client: KHAREJ - [PUBLIC IPV4] - REVERSE**

SERVER & Client IPV4 [ Private IPV4]:
 - Server[iran] command : ./tun-server-amd64 -server-port 800 -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -key azumi -mtu 1480
 - Client[kharej] command : ./tun-client-amd64 -server-addr IRAN_IPV4 -server-port 800 -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1480
 - 
SERVER & Client IPV4 [ Private IPV6]:
 - Server command : ./tun-server-amd64 -server-port 800 -server-private 2001:db8::1 -client-private 2001:db8::2 -subnet 64 -device tun2 -key azumi -mtu 1480
 - Client command : ./tun-client-amd64 -server-addr IRAN_IPV4 -server-port 800 -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1480

------------
**Server : IRAN  , Client: KHAREJ - [PUBLIC IPV6] - REVERSE**

SERVER & Client IPV6 [ Private IPV4]:
 - Server[iran] command : ./tun-server-amd64 -server-port 800 -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -key azumi -mtu 1480
 - Client[kharej] command : ./tun-client-amd64 -server-addr IRAN_IPV6 -server-port 800 -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -key azumi -mtu 1400

SERVER & Client IPV6 [ Private IPV6]:
 - Server command : ./tun-server-amd64 -server-port 800 -server-private 2001:db8::1 -client-private 2001:db8::2 -subnet 64 -device tun2 -key azumi -mtu 1480
 - Client command : ./tun-client-amd64 -server-addr IRAN_IPV6 -server-port 800 -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -key azumi -mtu 1400

