# This is an example Docker File
#  Command to build it
# docker built -t <image name > .
FROM osrf/ros:humble-desktop

RUN apt-get update
RUN apt-get install -y git && apt-get install -y python3-pip
RUN echo "source /opt/ros/humble/setup.bash" >> ~/.bashrc
RUN cd '/home'
RUN echo "ALL Done"






