![client1](https://github.com/user-attachments/assets/4695fd97-d937-4223-a622-431afc36e116)![R (2)](https://github.com/Azumi67/PrivateIP-Tunnel/assets/119934376/a064577c-9302-4f43-b3bf-3d4f84245a6f)
نام پروژه : لوکال تانل به صورت دایرکت یا ریورس - Tun interface و بر روی پورت TCP
---------------------------------------------------------------

**این تانل برای استفاده شخصی و گیم انلاین است. اموزش استفاده از اسکریپت قرار داده شد.**

![check](https://github.com/Azumi67/PrivateIP-Tunnel/assets/119934376/13de8d36-dcfe-498b-9d99-440049c0cf14)
**امکانات**
- امکان لوکال تانل بین سرور و کلاینت به صورت دایرکت یا ریورس ( برای سرور های محدود)
- استفاده از ایپی پرایوت های ساخته شده در Tun interface برای تانل اصلی یا پورت فوروارد
- امکان اتصال بین سرور و کلاینت بر روی پورت TCP
- امکان اتصال سرور و کلاینت با پابلیک ایپی 4 یا native
- امکان انتخاب ایپی پرایوت انتخابی خودتان هم به صورت پرایوت ایپی 4 یا پرایوت ایپی 6
- دارای worker بر اساس cpu cores یا انتخابی
- امکان انتخاب subnet mask برای پرایوت ایپی های ساخته
- امکان وارد کردن mtu به صورت manual
- دارای ping interval و استفاده از Bin bash برای ریست سرویس ها
- دارای encryption های IPSEC (در اسکریپت)
- دارای tcp nodelay و tcp keepalive
- دارای پرایوت و پابلیک key برای ارتباط بین سرور و کلاینت 
- دارای verbose برای نمایش لاگ (خطا)
- بعدا codereedsolomon و tcp window هم اضافه میشود

-----------------------


  <div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/3cfd920d-30da-4085-8234-1eec16a67460" alt="Image"> آپدیت</strong></summary>
  
------------------------------------ 

- تغییراتی در authentication method انجام شد و از pub key & priv key استفاده خواهد شد
- لاگ های تانل در مسیر etc/server.log/ یا etc/client.log/ ذخیره میشود
- گزینه worker اضافه شد. دیفالت بر اساس cpu cores و انتخابی
- مورد Challenge n Response Authentication به همراه unique nonce و Sha 256 به همراه expiry time حذف شد
- این مورد xtaci/smux دوباره با retry logic حذف شد
- این مورد tcpnodelay و tcp keepalive اضافه شد

  </details>
</div>


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
**آموزش استفاده از اسکریپت**

 <div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image"> </strong>روش دایرکت یا ریورس با پابلیک ایپی 4 یا 6</summary>
  
![green-dot-clipart-3](https://github.com/Azumi67/6TO4-PrivateIP/assets/119934376/902a2efa-f48f-4048-bc2a-5be12143bef3) **سرور خارج** 

 <p align="right">
  <img src="https://github.com/user-attachments/assets/3433bdf1-fb7a-4c93-8e81-624e39c5230e" alt="Image" />
</p>

- نخست یک پرایوت کی generate میشود که در کلاینت استفاده میشود.
- اگر کانفیگ دایرکت بود سرور خارج و کلاینت ایران است و اگر کانفیگ ریورس بود، سرور ایران و کلاینت خارج است

 <p align="right">
  <img src="https://github.com/user-attachments/assets/c09368e3-ce40-4131-81c0-ed22b9c6376e" alt="Image" />
</p>

- پورت سرور را وارد میکنم
- برای پرایوت ایپی تانل از ایپی ورژن 4 استفاده میکنم و برای سرور و کلاینت ایپی پرایوت مورد نظر را انتخاب میکنم
- برای ایپی پرایوت ورژن 4 از ساب نت 24 و برای ایپی پرایوت ورژن 6 از ساب نت 64 استفاده میکنم
- نام tun را ازومی قرار میدهم و مقدار mtu را 1380 وارد میکنم
- نمیخواهم لاگ ها در مسیر etc ذخیره شود و برای همین N را وارد میکنم
- چون میخواهم برای گیم استفاده کنم، گزینه tcpnodelay را فعال میکنم

 <p align="right">
  <img src="https://github.com/user-attachments/assets/dd90c3d5-3686-4988-9ed9-336a08944fd6" alt="Image" />
</p>

- گزینه tcp keepalive را فعال میکنم و مقدار 1m به این معنی که هر 1 دقیقه Keepalive انجام میشود را فعال میکنم
- گزینه worker را ویرایش میکنم و گزینه default بر اساس cpu cores را انتخاب میکنم
- همچنین socket را فعال میکنم و مقدار socket buff را همان مقدار پیش فرض میذارم
- سپس Ipsec را فعال میکنم و secret key را وارد میکنم
- ریست تایمر را برای ipsec فعال میکنم و مقدار انتخابی را وارد میکنم

 <p align="right">
  <img src="https://github.com/user-attachments/assets/8301a932-377c-43f9-8880-e582badc5b18" alt="Image" />
</p>

- ریست تایمر tun را فعال میکنم و عدد انتخابی را وارد میکنم . مثلا من هر 12 ساعت یا 24 ساعت ریست میکنم
- پس از فعال کردن ریست تایمر هم در سرور خارج و هم کلاینت برای sync شدن تایم ها بهتر است یک بار هر در سرویس مربوطه هم زمان در سرور و کلاینت ریست شوند
- سرور خارج : systemctl restart kazumi-local
- کلاینت ایران : systemctl restart iazumi-local

----------------------

![green-dot-clipart-3](https://github.com/Azumi67/6TO4-PrivateIP/assets/119934376/49000de2-53b6-4c5c-888d-f1f397d77b92)**کلاینت ایران**

<p align="right">
  <img src="https://github.com/user-attachments/assets/e7e83dc8-5eb7-4d63-b746-734afbd19cf9" alt="Image" />
</p>

- پورت سرور و ایپی پابلیک 4 سرور را وارد میکنم. اگر ایپی پابلیک 6 را انتخاب کرده بودید ایپی پابلیک 6 سرور را وارد میکنید
- پرایوت کی که در سرور خارج generate شده بود را اینجا وارد میکنم و سپس enter را وارد میکنم. دقت نمایید تمام محتوا کپی شود

<p align="right">
  <img src="https://github.com/user-attachments/assets/be94d8a2-f83f-44ff-9d09-f3faf0072974" alt="Image" />
</p>

- مانند سرور خارج پرایوت ایپی های مربوطه را انتخاب میکنم
- چون پرایوت ایپی 4 است پس ساب نت 24 را انتخاب میکنم
- ذخیره لاگ را غیرفعال میکنم
- گزینه tcp keepalive را فعال و بر روی 1m میذارم
- گزینه tcpnodelay هم فعال میکنم
- گزینه worker را ویرایش میکنم و default را انتخاب میکنم
- سوکت هم مانند سرورفعال میکنم و مقدار پیش فرض را enter میکنم

<p align="right">
  <img src="https://github.com/user-attachments/assets/3e01c62e-015a-4d42-adb5-8cb64f19b84e" alt="Image" />
</p>

- گزینه ipsec را فعال میکنم و secret key که در سرور وارد کردم را در کلاینت هم وارد میکنم
- ریست تایمر ipsec و tun هم فعال میکنم و مقادیر مشابه سرور را وارد میکنم
- اگر پینگ بین دو پرایوت ایپی برقرار نشد خود تانل پس از چند ثانیه connection را فعال میکند
- پس از انجام کانفیگ برای sync کردن ریست تایمر ها یک بار در سرور و کلاینت، سرویس ها را ریست کنید
- سرور خارج : systemctl restart kazumi-local
- کلاینت ایران : systemctl restart iazumi-local

**ویرایش کانفیگ**

<p align="right">
  <img src="https://github.com/user-attachments/assets/701ead57-6ad9-478e-b3ba-756bd2c3fba1" alt="Image" />
</p>

- در اینجا میتوانم هر مقداری را تغییر دهم.
- پابلیک ایپی گزینه دوم تنها برای نشان دادن انلاین بودن key میباشد و نیازی به تغییر ندارد
- گزینه 3 نشان دادن پرایوت کی است و برای زمانی استفاده میشود که میخواهیم از این کلید در کلاینت دیگری استفاده نماییم. به این صورت که کلید نمایشی را در کلاینت مربوطه paste میکنیم.
- هر مقدار اعم از پرایوت ایپی یا پورت را عوض کردید در کلاینت هم عوض کنید و گزینه save را بزنید

------------------

  </details>
</div>
 <div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image"> </strong>روش دایرکت با vxlan پابلیک ایپی 4 یا 6</summary>
  
![green-dot-clipart-3](https://github.com/Azumi67/6TO4-PrivateIP/assets/119934376/902a2efa-f48f-4048-bc2a-5be12143bef3) **سرور خارج** 

 <p align="right">
  <img src="https://github.com/user-attachments/assets/3433bdf1-fb7a-4c93-8e81-624e39c5230e" alt="Image" />
</p>

- نخست یک پرایوت کی generate میشود که در کلاینت استفاده میشود.
- اگر کانفیگ دایرکت بود سرور خارج و کلاینت ایران است و اگر کانفیگ ریورس بود، سرور ایران و کلاینت خارج است

 <p align="right">
  <img src="https://github.com/user-attachments/assets/c09368e3-ce40-4131-81c0-ed22b9c6376e" alt="Image" />
</p>

- پورت سرور را وارد میکنم
- برای پرایوت ایپی تانل از ایپی ورژن 4 استفاده میکنم و برای سرور و کلاینت ایپی پرایوت مورد نظر را انتخاب میکنم
- برای ایپی پرایوت ورژن 4 از ساب نت 24 و برای ایپی پرایوت ورژن 6 از ساب نت 64 استفاده میکنم
- نام tun را ازومی قرار میدهم و مقدار mtu را 1380 وارد میکنم
- نمیخواهم لاگ ها در مسیر etc ذخیره شود و برای همین N را وارد میکنم
- چون میخواهم برای گیم استفاده کنم، گزینه tcpnodelay را فعال میکنم
- گزینه tcp keepalive را فعال میکنم و مقدار 1m به این معنی که هر 1 دقیقه Keepalive انجام میشود را فعال میکنم
- گزینه worker را ویرایش میکنم و گزینه default بر اساس cpu cores را انتخاب میکنم
- همچنین socket را فعال میکنم و مقدار socket buff را همان مقدار پیش فرض میذارم

 <p align="right">
  <img src="https://github.com/user-attachments/assets/85e8d1c4-ff5b-424d-9dbb-ce5fbb9bd46c" alt="Image" />
</p>


- سپس Ipsec را فعال میکنم و secret key را وارد میکنم
- ریست تایمر را برای ipsec و tun فعال میکنم و مقدار انتخابی را وارد میکنم . مثلا من برای ریست تایمر مقدار 12 یا بیشتر را قرار میدهم


 <p align="right">
  <img src="https://github.com/user-attachments/assets/62584ca5-23a4-4901-b9f8-2b97e62ad31d" alt="Image" />
</p>

- سپس از من سوال میشود که کانفیگ vxlan یا geneve را میخواهم که من vxlan را انتخاب میکنم
- سپس از من سوال میشود که ایا میخواهم که خود اسکریپت نام اینترفیس را پیدا کند که گزینه y را وارد میکنم
- مقدار mtu هم اتوماتیک حساب میشود
- پس از فعال کردن ریست تایمر هم در سرور خارج و هم کلاینت برای sync شدن تایم ها بهتر است یک بار هر در سرویس مربوطه هم زمان در سرور و کلاینت ریست شوند
- سرور خارج : systemctl restart kazumi-local
- کلاینت ایران : systemctl restart iazumi-local

----------------------

![green-dot-clipart-3](https://github.com/Azumi67/6TO4-PrivateIP/assets/119934376/49000de2-53b6-4c5c-888d-f1f397d77b92)**کلاینت ایران**

<p align="right">
  <img src="https://github.com/user-attachments/assets/26f0348f-92a7-4773-b83d-03e473197a31" alt="Image" />
</p>

- پورت سرور و ایپی پابلیک 4 سرور را وارد میکنم. اگر ایپی پابلیک 6 را انتخاب کرده بودید ایپی پابلیک 6 سرور را وارد میکنید
- پرایوت کی که در سرور خارج generate شده بود را اینجا وارد میکنم و سپس enter را وارد میکنم. دقت نمایید تمام محتوا کپی شود

<p align="right">
  <img src="https://github.com/user-attachments/assets/956216f3-e3bb-47f6-922d-c26ac756c0d8" alt="Image" />
</p>

- مانند سرور خارج پرایوت ایپی های مربوطه را انتخاب میکنم
- چون پرایوت ایپی 4 است پس ساب نت 24 را انتخاب میکنم
- ذخیره لاگ را غیرفعال میکنم
- گزینه tcp keepalive را فعال و بر روی 1m میذارم
- گزینه tcpnodelay هم فعال میکنم
- گزینه worker را ویرایش میکنم و default را انتخاب میکنم
- سوکت هم مانند سرورفعال میکنم و مقدار پیش فرض را enter میکنم
- گزینه ipsec را فعال میکنم و secret key که در سرور وارد کردم را در کلاینت هم وارد میکنم

<p align="right">
  <img src="https://github.com/user-attachments/assets/3e01c62e-015a-4d42-adb5-8cb64f19b84e" alt="Image" />
</p>


- ریست تایمر ipsec و tun هم فعال میکنم و مقادیر مشابه سرور را وارد میکنم
- اگر پینگ بین دو پرایوت ایپی برقرار نشد خود تانل پس از چند ثانیه connection را فعال میکند
- کانفیگ vxlan را انتخاب میکنم
- سپس از من سوال میشود که ایا میخواهم که خود اسکریپت نام اینترفیس را پیدا کند که گزینه y را وارد میکنم
- مقدار mtu هم اتوماتیک حساب میشود
- پس از انجام کانفیگ برای sync کردن ریست تایمر ها یک بار در سرور و کلاینت، سرویس ها را ریست کنید
- سرور خارج : systemctl restart kazumi-local
- کلاینت ایران : systemctl restart iazumi-local

**ویرایش کانفیگ**

<p align="right">
  <img src="https://github.com/user-attachments/assets/7b93beab-5687-4ea4-a311-6733fec83274" alt="Image" />
</p>

- در اینجا میتوانم هر مقداری را تغییر دهم.
- پابلیک ایپی گزینه دوم تنها برای نشان دادن انلاین بودن key میباشد و نیازی به تغییر ندارد
- گزینه 3 نشان دادن پرایوت کی است و برای زمانی استفاده میشود که میخواهیم از این کلید در کلاینت دیگری استفاده نماییم. به این صورت که کلید نمایشی را در کلاینت مربوطه paste میکنیم.
- هر مقدار اعم از پرایوت ایپی یا پورت را عوض کردید در کلاینت هم عوض کنید و گزینه save را بزنید
- به طور مثال من پرایوت ایپی tun را عوض کردم. چون کانفیگ vxlan است باید ایپی پرایوت نخست با ایپی ورژن 4 شروع شود اما پرایوت ایپی خود vxlan میتواند ایپی ورژن 4 یا 6 باشد
- برای تغییر ایپی پرایوت vxlan هم مانند اسکرین شات عمل کنید و اگر ایپی پرایوت 6 میخواهید گزینه دوم را انتخاب و ایپی را بدون ساب نت وارد نمایید.
- در کلاینت هم تغییرات را انجام دهید. دقت نمایید اگر به طور مثال ایپی پرایوت در سرور db3::1 است باید در کلاینت db3::2 باشد

------------------

  </details>
</div>
<div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image"> </strong>روش ریورس با geneve پابلیک ایپی 4 یا 6</summary>
  
![green-dot-clipart-3](https://github.com/Azumi67/6TO4-PrivateIP/assets/119934376/902a2efa-f48f-4048-bc2a-5be12143bef3) **سرور ایران** 

 <p align="right">
  <img src="https://github.com/user-attachments/assets/3433bdf1-fb7a-4c93-8e81-624e39c5230e" alt="Image" />
</p>

- نخست یک پرایوت کی generate میشود که در کلاینت استفاده میشود.
- اگر کانفیگ دایرکت بود سرور خارج و کلاینت ایران است و اگر کانفیگ ریورس بود، سرور ایران و کلاینت خارج است

 <p align="right">
  <img src="https://github.com/user-attachments/assets/a63888cd-034c-4115-8ec6-e5d23e0417d1" alt="Image" />
</p>

- پورت سرور را وارد میکنم
- برای پرایوت ایپی تانل از ایپی ورژن 4 استفاده میکنم و برای سرور و کلاینت ایپی پرایوت مورد نظر را انتخاب میکنم
- برای ایپی پرایوت ورژن 4 از ساب نت 24 و برای ایپی پرایوت ورژن 6 از ساب نت 64 استفاده میکنم
- نام tun را ازومی قرار میدهم و مقدار mtu را 1380 وارد میکنم
- نمیخواهم لاگ ها در مسیر etc ذخیره شود و برای همین N را وارد میکنم
- چون میخواهم برای گیم استفاده کنم، گزینه tcpnodelay را فعال میکنم
- گزینه tcp keepalive را فعال میکنم و مقدار 1m به این معنی که هر 1 دقیقه Keepalive انجام میشود را فعال میکنم
- گزینه worker را ویرایش میکنم و گزینه default بر اساس cpu cores را انتخاب میکنم
- همچنین socket را فعال میکنم و مقدار socket buff را همان مقدار پیش فرض میذارم

 <p align="right">
  <img src="https://github.com/user-attachments/assets/6238a098-91ef-4d39-b4a3-9bdfbf5485fc" alt="Image" />
</p>


- سپس Ipsec را فعال میکنم و secret key را وارد میکنم
- ریست تایمر را برای ipsec و tun فعال میکنم و مقدار انتخابی را وارد میکنم . مثلا من برای ریست تایمر مقدار 12 یا بیشتر را قرار میدهم
- سپس از من سوال میشود که کانفیگ vxlan یا geneve را میخواهم که من geneve را انتخاب میکنم
- مقدار mtu هم اتوماتیک حساب میشود
- پس از فعال کردن ریست تایمر هم در سرور خارج و هم کلاینت برای sync شدن تایم ها بهتر است یک بار هر در سرویس مربوطه هم زمان در سرور و کلاینت ریست شوند
- کلاینت خارج : systemctl restart krazumi-local
- سرور ایران : systemctl restart irazumi-local

----------------------

![green-dot-clipart-3](https://github.com/Azumi67/6TO4-PrivateIP/assets/119934376/49000de2-53b6-4c5c-888d-f1f397d77b92)**کلاینت خارج**


<p align="right">
  <img src="https://github.com/user-attachments/assets/1e583d49-c9a5-43c2-9e51-0e714391c95a" alt="Image" />
</p>

- پورت سرور و ایپی پابلیک 4 سرور را وارد میکنم. اگر ایپی پابلیک 6 را انتخاب کرده بودید ایپی پابلیک 6 سرور را وارد میکنید
- پرایوت کی که در سرور ایران generate شده بود را اینجا وارد میکنم و سپس enter را وارد میکنم. دقت نمایید تمام محتوا کپی شود

<p align="right">
  <img src="https://github.com/user-attachments/assets/4baa3b64-e7d6-4d70-b2fd-4e2ae11bddfd" alt="Image" />
</p>

- مانند سرور ایران پرایوت ایپی های مربوطه را انتخاب میکنم
- چون پرایوت ایپی 4 است پس ساب نت 24 را انتخاب میکنم
- ذخیره لاگ را غیرفعال میکنم
- گزینه tcp keepalive را فعال و بر روی 1m میذارم
- گزینه tcpnodelay هم فعال میکنم
- گزینه worker را ویرایش میکنم و default را انتخاب میکنم
- سوکت هم مانند سرورفعال میکنم و مقدار پیش فرض را enter میکنم
- گزینه ipsec را فعال میکنم و secret key که در سرور وارد کردم را در کلاینت هم وارد میکنم

<p align="right">
  <img src="https://github.com/user-attachments/assets/49088dbd-e151-4640-bdb8-40d9548f898e" alt="Image" />
</p>


- ریست تایمر ipsec و tun هم فعال میکنم و مقادیر مشابه سرور را وارد میکنم
- اگر پینگ بین دو پرایوت ایپی برقرار نشد خود تانل پس از چند ثانیه connection را فعال میکند
- کانفیگ geneve را انتخاب میکنم
- مقدار mtu هم اتوماتیک حساب میشود
- پس از انجام کانفیگ برای sync کردن ریست تایمر ها یک بار در سرور و کلاینت، سرویس ها را ریست کنید
- کلاینت خارج : systemctl restart krazumi-local
- سرور ایران : systemctl restart irazumi-local

**ویرایش کانفیگ**

<p align="right">
  <img src="https://github.com/user-attachments/assets/c01eeff7-3238-45f9-8355-8165d45c822f" alt="Image" />
</p>

- در اینجا میتوانم هر مقداری را تغییر دهم.
- پابلیک ایپی گزینه دوم تنها برای نشان دادن انلاین بودن key میباشد و نیازی به تغییر ندارد
- گزینه 3 نشان دادن پرایوت کی است و برای زمانی استفاده میشود که میخواهیم از این کلید در کلاینت دیگری استفاده نماییم. به این صورت که کلید نمایشی را در کلاینت مربوطه paste میکنیم.
- هر مقدار اعم از پرایوت ایپی یا پورت را عوض کردید در کلاینت هم عوض کنید و گزینه save را بزنید
- به طور مثال من پرایوت ایپی tun را عوض کردم. چون کانفیگ geneve است باید ایپی پرایوت نخست با ایپی ورژن 4 شروع شود اما پرایوت ایپی خود geneve میتواند ایپی ورژن 4 یا 6 باشد
- برای تغییر ایپی پرایوت geneve هم مانند اسکرین شات عمل کنید و اگر ایپی پرایوت 6 میخواهید گزینه دوم را انتخاب و ایپی را بدون ساب نت وارد نمایید.
- در کلاینت هم تغییرات را انجام دهید. دقت نمایید اگر به طور مثال ایپی پرایوت در سرور db3::1 است باید در کلاینت db3::2 باشد

------------------

  </details>
</div>
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
  wget https://github.com/Azumi67/LocalTun_TCP/releases/download/v1.7/amd64.zip
  unzip amd64.zip -d /root/localTUN
  cd localTUN
  chmod +x tun-server 
  chmod +x tun-client
Key creation :
openssl genrsa -out private_key.pem 2048
openssl rsa -pubout -in private_key.pem -out public_key.pem
 ```
 </details>
</div>
<div align="right">
  <details>
    <summary><strong><img src="https://github.com/Azumi67/Rathole_reverseTunnel/assets/119934376/fcbbdc62-2de5-48aa-bbdd-e323e96a62b5" alt="Image">دایرکت لوکال تانل پرایوت ایپی 4 - public ipv4 </strong></summary>

  - کامند های سرور (خارج)
  
 <div align="left">
   
```
./tun-server -server-port 800 -pub-key=/root/keys/public_key.pem -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -mtu 1380 -tcpnodelay -keepalive 1m -worker default
   
```
<div align="right">
  
- کامند های کلاینت (ایران)

 <div align="left">
   
```
./tun-client -server-addr KHAREJ_IPV4 -server-port 800 -priv-key=/root/private_key.pem -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -mtu 1280 -tcpnodelay -keepalive 1m -worker default
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
ExecStart=/root/localTUN/tun-server -server-port 800 -pub-key=/root/keys/public_key.pem -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -mtu 1480 -tcpnodelay -keepalive 1m -worker default
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
 <div align="left">
   
```
./tun-server -server-port 800 -pub-key=/root/keys/public_key.pem -server-private 2001:db8::1 -client-private 2001:db8::2 -subnet 64 -device tun2 -mtu 1380 -tcpnodelay -keepalive 1m -worker default

```
<div align="right">
  
- کامند های کلاینت (ایران)
 <div align="left">
   
```
./tun-client -server-addr KHAREJ_IPV4 -server-port 800 -priv-key=/root/private_key.pem -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -mtu 1280 -tcpnodelay -keepalive 1m -worker default
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
ExecStart=/root/localTUN/tun-client -server-addr KHAREJ_IPV4 -server-port 800 -priv-key=/root/private_key.pem -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -mtu 1280 -tcpnodelay -keepalive 1m -worker default
   

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
 <div align="left">
   
```
./tun-server_amd64 -server-port 800 -pub-key=/root/keys/public_key.pem -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -mtu 1480 -tcpnodelay -keepalive 1m -worker default
```
<div align="right">
  
- کامند های کلاینت (ایران)
 <div align="left">
   
```
./tun-client -server-addr KHAREJ_IPV6 -server-port 800 -priv-key=/root/private_key.pem -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -mtu 1280 -tcpnodelay -keepalive 1m -worker default
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
ExecStart=/root/localTUN/tun-client -server-addr KHAREJ_IPV6 -server-port 800 -priv-key=/root/private_key.pem -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -mtu 1380 -tcpnodelay -keepalive 1m -worker default
   

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
./tun-server -server-port 800 -pub-key=/root/keys/public_key.pem -server-private 2001:db8::1 -client-private 2001:db8::2 -subnet 64 -device tun2 -mtu 1380 -tcpnodelay -keepalive 1m -worker default
```
<div align="right">
  
- کامند های کلاینت (ایران)
 <div align="left">
   
```
./tun-client -server-addr [KHAREJ_IPV6] -server-port 800 -priv-key=/root/private_key.pem -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -mtu 1380 -tcpnodelay -keepalive 1m -worker default
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
ExecStart=/root/localTUN/tun-client -server-addr [KHAREJ_IPV6] -server-port 800 -priv-key=/root/private_key.pem -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -mtu 1280 -tcpnodelay -keepalive 1m -worker default
   

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

 <div align="left">
   
```
./tun-server -server-port 800 -pub-key=/root/keys/public_key.pem -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -mtu 1380 -tcpnodelay -keepalive 1m -worker default
```
<div align="right">
  
- کامند های کلاینت (خارج)
 <div align="left">
   
```
./tun-client -server-addr IRAN_IPV4 -server-port 800 -priv-key=/root/private_key.pem -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -mtu 1280 -tcpnodelay -keepalive 1m -worker default
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
ExecStart=/root/localTUN/tun-server -server-port 800 -pub-key=/root/keys/public_key.pem -server-private 30.0.0.1 -client-private 30.0.0.2 -subnet 24 -device tun2 -mtu 1380 -tcpnodelay -keepalive 1m -worker default
   

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
 <div align="left">
   
```
./tun-server -server-port 800 -pub-key=/root/keys/public_key.pem -server-private 2001:db8::1 -client-private 2001:db8::2 -subnet 64 -device tun2 -mtu 1380 -tcpnodelay -keepalive 1m -worker default
```
<div align="right">
  
- کامند های کلاینت (خارج)
 <div align="left">
   
```
./tun-client -server-addr IRAN_IPV4 -server-port 800 -priv-key=/root/private_key.pem -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -mtu 1280 -tcpnodelay -keepalive 1m -worker default
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
ExecStart=/root/localTUN/tun-client -server-addr IRAN_IPV4 -server-port 800 -priv-key=/root/private_key.pem -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -mtu 1280 -tcpnodelay -keepalive 1m -worker default
   

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
./tun-server -server-port 800 -pub-key=/root/keys/public_key.pem -server-private 30.0.0.1 -client-private 30.0.0.2 -pub-key=/root/keys/public_key.pem -subnet 24 -device tun2 -mtu 1380 -tcpnodelay -keepalive 1m -worker default
```
<div align="right">
  
- کامند های کلاینت (خارج)
 <div align="left">
   
```
./tun-client -server-addr [IRAN_IPV6] -server-port 800 -priv-key=/root/private_key.pem -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -mtu 1280 -tcpnodelay -keepalive 1m -worker default
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
ExecStart=/root/localTUN/tun-client -server-addr [IRAN_IPV6] -server-port 800 -priv-key=/root/private_key.pem -client-private 30.0.0.2 -server-private 30.0.0.1 -subnet 24 -device tun2 -mtu 1280 -tcpnodelay -keepalive 1m -worker default
   

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
./tun-server -server-port 800 -pub-key=/root/keys/public_key.pem -server-private 2001:db8::1 -client-private 2001:db8::2 -subnet 64 -device tun2 -mtu 1480 -tcpnodelay -keepalive 1m -worker default
```
<div align="right">
  
- کامند های کلاینت (خارج)
 <div align="left">
   
```
./tun-client -server-addr [IRAN_IPV6] -server-port 800 -priv-key=/root/private_key.pem -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -mtu 1280 -tcpnodelay -keepalive 1m -worker default
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
ExecStart=/root/localTUN/tun-client -server-addr [IRAN_IPV6] -server-port 800 -priv-key=/root/private_key.pem -client-private 2001:db8::2 -server-private 2001:db8::1 -subnet 64 -device tun2 -mtu 1280 -tcpnodelay -keepalive 1m -worker default
   

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



