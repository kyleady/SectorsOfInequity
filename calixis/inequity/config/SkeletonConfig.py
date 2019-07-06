from inequity.base.InequityObject import InequityObject

class SkeletonConfig(InequityObject):
    def __init__(self):
        self._values = {
            "frequency_rate": 0.5,
            "reach": 5,
            "connection_rate": 0.4,
            "range_multiplier": 0.25,
            "x": 20,
            "y": 20
        }
