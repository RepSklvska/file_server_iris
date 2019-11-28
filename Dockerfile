FROM i386/debian:buster
MAINTAINER AlexanderSokolov<gunpowderfans@mail.ru>

EXPOSE 3000

RUN apt-get update &&\
	apt-get -y install \
	wget \
	unzip &&\
	cd /usr/local/share/ &&\
	wget https://github.com/RepSklvska/file_server_iris/archive/master.zip &&\
	unzip master.zip &&\
	rm -r master.zip &&\
	mv file_server_iris-master file_server_iris &&\
	cd file_server_iris &&\
	rm -rf `ls * | egrep -v '(main|config|public|views)'` &&\
	rm main.go main.exe &&\
	apt autoremove wget unzip -y &&\
	apt autoclean -y &&\
	apt clean -y &&\
	rm -rf /var/lib/apt/lists/* /tmp/* /vat/tmp/*

CMD ["sh","-c","/usr/local/share/file_server_iris/main -rootdir=/usr/Downloads"]