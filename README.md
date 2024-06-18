# LocalTun_TCP
This implements a TCP tun between a client &amp; server using a TUN interface. The client and server communicate by encapsulating network packets within TCP connections. Encapsulation &amp; decapsulates is Simultaneous.

- از طریق tcp دو سرور به هم وصل میشوند و از طریق اینترفیس tun و پرایوت ایپی به هم دیگه متصل خواهند بود. encapsulation & decapsulation هم زمان انجام میشود.
- هدف نوشتن این برنامه برای این بوده است که از طریق پورت tcp و tun interface،‌ دو سرور به هم متصل شوند و از پرایوت آی‌پی های آنها برای تانل استقاده کرد .
- میتوان سرور اصلی ایران باشد و کلاینت خارج باشد و بدین صورت میتوان تانل لوکالی بر روی سرور های خارج محدود هم ایجاد کرد.
- یا میتوانید به صورت دایرکت یعنی سرور اصلی خارج و کلاینت ایران، این تانل را اجرا کنید
- هم میشه با ایپی 4 سرور و هم ایپی 6 سرور وصل شد .
- میتوانید تست کنید و بهم اطلاع بدید اما برای استفاده باید صبر کنید تا داخل اسکریپت برای ترکیب با ipsec آورده شود.
- کلید authentication اضافه شد که در کلاینت خوانده میشود و در سرور احراز هویت میشود
