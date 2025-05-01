#!/usr/bin/env python3
import rclpy
import rclpy.exceptions
from rclpy.node import Node
class myNode(Node):

    def __init__(self):
        super().__init__("first_node")
        self.counter_=0
        self.create_timer(1.0, self.timer_callback)

    def timer_callback(self):
        self.get_logger().info("Hello" + str(self.counter_))
        self.counter_ +=1



def main():
    try:
        rclpy.init()
        node = myNode()
        rclpy.spin(node)
        rclpy.shutdown()
    except:
        KeyboardInterrupt
        print('exited')
    



if __name__ == '__main__':
    main()
