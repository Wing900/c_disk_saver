# c_disk_saver
一个通过创建软链接从而清理C盘的脚本实现

# 使用c_disk_saver转移C盘内存教程

## 原理

最近发现了一个清理C盘的办法,在很多博客都有写,其原理是通过创建软链接把物理位置搬走,但保留虚拟位置

这个过程可以手动用终端完成,无需任何软件和脚本

1. 移动物理位置:(以谷歌数据文件夹为例)

```
robocopy "C:\Users\Bin\AppData\Local\Google" "D:\NewAppdata\Google" /E /MOVE /XJ
```

2.创建软链接(其中NewAppdata是我保存物理位置的地方)

```
mklink /J "C:\Users\Bin\AppData\Local\Google" "D:\NewAppdata\Google"
```

---

完成这个过程之后,以后电脑访问`C:\Users\Bin\AppData\Local\Google`就会因为它创建过软连接,而不直接打开它,而是连接到`D:\NewAppdata\Google`,从而不影响源文件的使用

以上是原理部分,我为了更好地完成这个操作,写了一个exe脚本便于使用,已经[开源](https://github.com/Wing900/c_disk_saver/blob/main)

## 步骤 1

先确定我们的搬家对象,我一般搬家的都是`C:\Users\Bin\AppData\`下的local文件夹,具体怎么确定,你进行以下步骤

- 下载Windirstat软件,分析磁盘

- 在`C:\Users\Bin\AppData\Local`下发现占用极大的文件,可能有很多
- 这里我以`GitHubDesktop`为例

<img width="2560" height="1600" alt="image" src="https://github.com/user-attachments/assets/8c31e670-d3cf-44dd-a8e5-5d808a1c4b6b" />


- 找到了之后选中-右键-复制路径地址,这个地址就是我们的**源文件夹**

## 步骤2

**创建搬家地点--即目标文件夹**

- 在你想迁移过去的位置比如D盘,创建文件夹
- 为了方便整理,我创建的是`D:\CAppdata\Local`

这样对应了原来C盘的`local`文件夹

然后复制这个地址,作为我们的**目标文件夹**

> 这里注意,不需要和源文件夹同名,比如弄个GitHubDesktop,不需要,因为程序做的事情是复制粘贴,我们只是创建一个父目录

## 步骤3

打开c_disk_saver

- 允许管理员权限

- 阅读警告

如果没有问题就输入刚刚的两个地址即可


<img width="1132" height="637" alt="image" src="https://github.com/user-attachments/assets/a1a66bc9-c403-4efa-b6c4-debdacc7737a" />


## 步骤4

等待移动,没有特殊情况就会移动完毕了,可以打开相关软件进行测试

一般来说,对于一些软件,会自动退出登录,但是数据还是在的

<img width="1255" height="723" alt="image" src="https://github.com/user-attachments/assets/25f46d4f-5c74-47fb-864b-d91bec71395d" />


<img width="1252" height="325" alt="image" src="https://github.com/user-attachments/assets/5bca95b7-7017-496c-9733-4f614a696fd8" />


但是有的文件需要更深层的文件权限,可能会导致转移失败
