### 树莓派人脸识别
1、识别熟悉的人
2、语音播报

### 环境准备

```
# Ubuntu
sudo apt-get install libdlib-dev libblas-dev liblapack-dev libjpeg-turbo8-dev
# Debian
sudo apt-get install libdlib-dev libblas-dev liblapack-dev libjpeg62-turbo-dev

# 声音
sudo apt-get install libasound2-dev
```

### 设置声音输出方式
```
sudo raspi-config
设置使用HDMI音频输出

测试工具
sudo apt install mplayer mpg123 -y 

测试声音
sudo mplayer xx.mp3

```

### 设置帧率
```
https://www.raspberrypi.org/forums/viewtopic.php?f=43&t=268139&p=1627824#p1627824
sudo touch /etc/modprobe.d/bcm2835-v4l2.conf
add 
options bcm2835-v4l2 max_video_width=1920 max_video_height=1088
```

### 其余优化帧率方法
```
https://www.e-consystems.com/Articles/Camera/accessing_cameras_in_opencv_with_high_performance.asp
https://github.com/econsystems/opencv_v4l2
```


	
### FAQ
树莓派遇到编译问题
“{aka ‘__gnu_cxx::__normal_iterator<dlib::vector<double, 2>*, std::vector<dlib::vector<double, 2> > >’} changed in GCC 7.1”
，修改
/home/pi/gopro/pkg/mod/github.com/!kagami/go-face@v0.0.0-20200825065730-3dd2d74dccfb/face.go

添加忽略警告标识: -Wno-psabi
// #cgo CXXFLAGS: -std=c++1z -Wall -O3 -DNDEBUG -march=native 
为
// #cgo CXXFLAGS: -std=c++1z -Wall -O3 -DNDEBUG -march=native -Wno-psabi


### 参考文献：
http://dlib.net/ 
https://hackernoon.com/face-recognition-with-go-676a555b8a7e
https://studygolang.com/articles/20172


