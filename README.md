# A Simple File Server using Iris

Simple and unstable file server web app written with <a href="https://github.com/kataras/iris">Iris</a>. Just for learning. Recommended to use in private network.

Reference:

<a href="https://www.cnblogs.com/pu369/p/10950746.html">www.cnblogs.com/pu369/p/10950746.html</a>

<a href="https://github.com/kataras/iris/tree/master/_examples">github.com/kataras/iris/tree/master/_examples</a>

<br>

Usage:
>$ go run ./main.go

or
>$ ./main

<br>

To specify arguments, you can add:
>--rootdir="~/SharedFiles" --port=3000

or edit "config" file as text. Key and value are separated by a space.
>RootDir /home/alexander/SharedFiles
>
>Port 3000