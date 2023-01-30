# XshellCrack
XshellCrack is a secondary development based on SharpXDecrypt, rewritten with go language, increased the registry query Settings, more convenient for xshell local password cracking
# introduce

### Use
```
Usage:
  root SshCrack [flags]

Flags:
  -P, --Path string   eg -P C:\xxxx\xxx.xsh
  -h, --help          help for SshCrack
```
If the xshell path is not configured, the system automatically searches the xshell local cache file in the registry, as shown in the following figure

![图片](https://user-images.githubusercontent.com/113832601/215474069-d22e6af1-8c6e-4f2c-b1ad-fe81a3f1c12a.png)


If you specify a path, the decryption will follow the path you specified, as shown in the figure below
![图片](https://user-images.githubusercontent.com/113832601/215474787-67729908-8ebc-49b3-ae9a-38be928ff301.png)


### Kill free effect
Because there is no windows api call involved, the effect is still OK, and the use of go language construction to remove the symbol table, and use the upx shell for compression, so the volume is also very impressive

![图片](https://user-images.githubusercontent.com/113832601/215475502-7ce075a2-3bc8-4d84-ba9c-25e80fdecb94.png)


#reference
```
https://github.com/JDArmy/SharpXDecrypt
```
