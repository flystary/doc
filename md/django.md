#### 1.django 模型models 常用字段

​    

> 1、models.AutoField　　
>
> - ​    自增列 = int(11)
> - ​    如果没有的话，默认会生成一个名称为 id 的列
> - ​    如果要显式的自定义一个自增列，必须设置primary_key=True。
>
>  
>
> 2、models.CharField　　
>
> - ​    字符串字段
> - 　　必须设置max_length参数
>
>  
>
> 3、models.BooleanField　 
>
> - ​    布尔类型=tinyint(1)
> - 　　不能为空，可添加Blank=True
>
>  
>
> 4、models.ComaSeparatedIntegerField　　
>
> - ​    用逗号分割的数字=varchar
> - 　　继承CharField，所以必须 max_lenght 参数
>
>  
>
> 5、models.DateField
>
> - ​    日期类型 date
> - 　　DateField.auto_now：保存时自动设置该字段为现在日期，最后修改日期
> - ​    DateField.auto_now_add：当该对象第一次被创建是自动设置该字段为现在日期，创建日期。

>  
>
> 6、models.DateTimeField　　
>
> - ​    日期时间类型 datetime
> - 　　同DateField的参数
>
>  
>
> 7、models.Decimal　　
>
> - ​    十进制小数类型 = decimal
> - ​    DecimalField.max_digits：数字中允许的最大位数
> - ​    DecimalField.decimal_places：存储的十进制位数
>
>  
>
> 8、models.EmailField　　
>
> - 　　一个带有检查 Email 合法性的 CharField
>
>  
>
> 9、models.FloatField　　
>
> - ​    浮点类型 = double
>
>  
>
> 10、models.IntegerField　　
>
> - ​    整形
>
>  
>
> 11、models.BigIntegerField　　
>
> - ​    长整形
> - 　　integer_field_ranges = {

> > 　　　　'SmallIntegerField': (-32768, 32767),
>
> > 　　　　'IntegerField': (-2147483648, 2147483647),
>
> > 　　　　'BigIntegerField': (-9223372036854775808, 9223372036854775807),
>
> > 　　　　'PositiveSmallIntegerField': (0, 32767),
>
> > 　　　　'PositiveIntegerField': (0, 2147483647),
>
> > 　　}

>  
>
> 12、models.GenericIPAddressField　　
>
> - ​    一个带有检查 IP地址合法性的 CharField
>
>  
>
> 13、models.NullBooleanField　　
>
> - ​    允许为空的布尔类型
>
>  
>
> 14、models.PositiveIntegerFiel　　
>
> - ​    正整数
>
>  
>
> 15、models.PositiveSmallIntegerField　　
>
> - ​    正smallInteger
>
>  
>
> 16、models.SlugField　　
>
> - ​    减号、下划线、字母、数字
>
>  
>
> 17、models.SmallIntegerField　　
>
> - ​    数字
> - 　　数据库中的字段有：tinyint、smallint、int、bigint
>
>  
>
> 18、models.TextField　　
>
> - ​    大文本。默认对应的form标签是textarea。
>
>  
>
> 19、models.TimeField　　
>
> - ​    时间 HH:MM[:ss[.uuuuuu]]
>
>  
>
> 20、models.URLField　　
>
> - ​    一个带有URL合法性校验的CharField。
>
>  
>
> 21、models.BinaryField　　
>
> - ​    二进制
> - ​    存储二进制数据。不能使用filter函数获得QuerySet。
>
>  
>
> 22、models.ImageField  
>
> - ​    图片
> - ​    ImageField.height_field、ImageField.width_field：如果提供这两个参数，则图片将按提供的高度和宽度规格保存。
> - ​    该字段要求 Python Imaging 库Pillow。
> - ​    会检查上传的对象是否是一个合法图片。
>
>  
>
> 23、models.FileField(upload_to=None[, max_length=100, ** options])
>
> - ​    文件
> - ​    FileField.upload_to：一个用于保存上传文件的本地文件系统路径，该路径由 MEDIA_ROOT 中设置
> - ​    这个字段不能设置primary_key和unique选项.在数据库中存储类型是varchar，默认最大长度为100
>
>  
>
> 24、models.FilePathField(path=None[, math=None, recursive=False, max_length=100, **options])
>
> - ​    FilePathField.path：文件的绝对路径，必填
> - ​    FilePathField.match：用于过滤路径下文件名的正则表达式，该表达式将用在文件名上（不包括路径）。
> - ​    FilePathField.recursive：True 或 False，默认为 False，指定是否应包括所有子目录的路径。
> - ​    例如：FilePathField(path="/home/images", match="foo.*", recursive=True)
>
> ​             将匹配“/home/images/foo.gif”但不匹配“/home/images/foo/bar.gif”   

 

2. #### django 模型models 字段常用参数

 

> 1、null
>
> - ​    如果是True，Django会在数据库中将此字段的值置为NULL，默认值是False
>
>  
>
> 2、blank
>
> - 　　如果为True时django的 Admin 中添加数据时可允许空值，可以不填。如果为False则必须填。默认是False。
> - ​    null纯粹是与数据库有关系的。而blank是与页面必填项验证有关的
>
>  
>
> 3、primary_key = False
>
> - 　  主键，对AutoField设置主键后，就会代替原来的自增 id 列
>
>  
>
> 4、auto_now 和 auto_now_add
>
> - 　　auto_now  自动创建---无论添加或修改，都是当前操作的时间
> - 　　auto_now_add  自动创建---永远是创建时的时间
>
>  
>
> 5、choices
>
> -    一个二维的元组被用作choices，如果这样定义，Django会select box代替普通的文本框，
> -    并且限定choices的值是元组中的值
> -    GENDER_CHOICE = (
> - ​      (u'M', u'Male'),
> - ​      (u'F', u'Female'),
> -    )
> -    gender = models.CharField(max_length=2,choices = GENDER_CHOICE)

>  
>
> 6、max_length
>
> - ​    字段长度
>
>  
>
> 7、default
>
> - ​    默认值

>  
>
> 8、verbose_name　　
>
> - ​    Admin中字段的显示名称，如果不设置该参数时，则与属性名。
>
>  
>
> 9、db_column　　
>
> - ​    数据库中的字段名称
>
>  
>
> 10、unique=True　　
>
> -    不允许重复
>
>  
>
> 11、db_index = True　　
>
> -    数据库索引
>
>  
>
> 12、editable=True　　
>
> -    在Admin里是否可编辑
>
>  
>
> 13、error_messages=None　　
>
> -    错误提示
>
>  
>
> 14、auto_created=False　　
>
> -    自动创建
>
>  
>
> 15、help_text　　
>
> -    在Admin中提示帮助信息
>
>  
>
> 16、validators=[]
>
> - ​     验证器
>
>  
>
> 17、upload-to
>
> - ​    文件上传时的保存上传文件的目录

 

models.py

​    \# -*- coding: utf-8 -*-

> from __future__ import unicode_literals
>
> from django.db import models
>
> 
>
> class UserInfo(models.Model):
>
> userName = models.CharField(max_length=30)  #用户名
>
> passWord = models.CharField(max_length=30)  #密码
>
> gendle = models.BooleanField()  #性别
>
> birthday = models.DateField()  #出生日期
>
>  weigth = models.FloatField()   #体重
>
> heigth = models.IntegerField()  #身高
>
> email = models.EmailField()   #邮箱
>
> host = models.GenericIPAddressField()  #IP地址
>
> introduce = models.TextField()  #个人简介
>
> blog = models.URLField()     #博客地址
>
> photo = models.ImageField()   #照片
>
> CV = models.FilePathField()   #个人简历文件
>
> createDate = models.DateTimeField()   #帐号申请时间

​    执行结果：

​         ![img](https://images2017.cnblogs.com/blog/786577/201710/786577-20171007112806677-1881938928.png)

 

------

3.常见异常处理

  1) 

​        ERRORS：

> ​        web.UserInfo.photo: (fields.E210) Cannot use ImageField because Pillow is not installed.
>
> ​        HINT: Get Pillow at https://pypi.python.org/pypi/Pillow or run command "pip install Pillow".
>
>  
>
> ​    原因：    
>
> 
>
> ​    这是因为使用了ImageField()字段，该字段是直接在数据库中存储图片的，数据库中实际存储时要使用python的Pillow模块对图片进行处理后才能存储进去。因此因需使用pip install Pillow 安装该模块即可解决该报错。

  2）

​       ERRORS：

> ​        在执行python manage.py makemigrations 时需要手动选择处理方法：
>
> ​        You are trying to add a non-nullable field 'CV' to userinfo without a default; we can't do that (the database           needs something to populate existing rows).

> Please select a fix:
>
>  1) Provide a one-off default now (will be set on all existing rows with a null value for this column)
>
>  2) Quit, and let me add a default in models.py
>
> Select an option: 1
>
> Please enter the default value now, as valid Python
>
> The datetime and django.utils.timezone modules are available, so you can do e.g. timezone.now
>
> Type 'exit' to exit this prompt
>
> > timezone.now

>    原因：    
>
> ​       这是因为 UserInfo数据表中已经有了"userName"和"passWord" 这两个字段，当在新增字段时就会出现这种Warning。是由于django框架在生成SQL语句时发现数据表不 为空，担心新增不为Null的字段添加到该表中时，表中以前的数据行在填充这些字段时需要填充的值不明确，所以才会询问用户处理方式。
>
> ​      选1，则会在已存在行中添加null,选2，则会要求在models.py中添加默认值。
>
> ​      在models.py中设置默认值的方法：
>
> ​              host = models.GenericIPAddressField(default = '127.0.0.1' 

  3)  执行python makemigrations正常，但是执行python migrate 报错，之后再执行无法生效的处理办法