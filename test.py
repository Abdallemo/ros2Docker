#!/usr/bin/env python3
import rclpy
from rclpy.node import Node
from std_msgs.msg import String  # Import message type

class MyTalkerNode(Node):
    def __init__(self):
        super().__init__("talker_node")
        self.publisher_ = self.create_publisher(String, "chatter", 10)  # Topic name: 'chatter'
        self.counter_ = 0
        self.create_timer(1.0, self.timer_callback)  # Call every 1 second

    def timer_callback(self):
        msg = String()
        msg.data = "Hello " + str(self.counter_)
        self.publisher_.publish(msg)
        self.get_logger().info(f"Published: {msg.data}")
        self.counter_ += 1

def main():
    rclpy.init()
    node = MyTalkerNode()
    try:
        rclpy.spin(node)
    except KeyboardInterrupt:
        print("Node interrupted")
    finally:
        node.destroy_node()
        rclpy.shutdown()

if __name__ == "__main__":
    main()
