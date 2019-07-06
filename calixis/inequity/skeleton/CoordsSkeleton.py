from inequity.base.InequityObject import InequityObject

class CoordsSkeleton(InequityObject):
    def __init__(self, x, y):
        self._values = {
            "x": x,
            "y": y
        }

    @property
    def x(self):
        return self["x"]

    @property
    def y(self):
        return self["y"]

    def outside(self, min, max):
         return (
            self.x < min.x
            or
            self.y < min.y
            or
            self.x > max.x
            or
            self.y > max.y
         )
