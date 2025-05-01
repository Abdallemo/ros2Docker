#!/usr/bin/env python3
import rclpy
from rclpy.node import Node
from geometry_msgs.msg import Twist

class drawCircleNode(Node):
    
    def __init__(self):
        super().__init__("draw_circle")
        self.cmd_vel_pub_ = self.create_publisher(Twist, "/turtle/cmd_vel", 10)
        self.get_logger().info("Draw circle node has been created")
    
    def



def main():
    rclpy.init()
    rclpy.shutdown
    
