#!/usr/bin/env python3
import rclpy
from rclpy.node import Node
from geometry_msgs.msg import Twist

class CircleDrawer(Node):
    def _init_(self):
        super()._init_("circle_drawer")
        self.publisher_ = self.create_publisher(Twist, "/turtle1/cmd_vel", 10)
        self.timer = self.create_timer(0.1, self.timer_callback)  # 10 Hz

    def timer_callback(self):
        msg = Twist()
        msg.linear.x = 2.0    # Forward speed
        msg.angular.z = 1.0   # Rotational speed to create a circle
        self.publisher_.publish(msg)
        self.get_logger().info("Publishing circle command...")

def main():
    rclpy.init()
    node = CircleDrawer()
    try:
        rclpy.spin(node)
    except KeyboardInterrupt:
        pass
    finally:
        node.destroy_node()
        rclpy.shutdown()

if _name_ == "_main_":
    main()