from inequity.base.InequityObject import InequityObject

class RouteSkeleton(InequityObject):
    def __init__(self, start_system, end_system):
        self._values = {
            "start": start_system.coords,
            "end": end_system.coords
        }

        self._hidden = {
            "start_system": start_system,
            "end_system": end_system
        }

    @property
    def start_system(self):
        return self._hidden["start_system"]

    @property
    def end_system(self):
        return self._hidden["end_system"]

    def reverse(self):
        return RouteSkeleton(self._hidden["end_system"], self._hidden["start_system"])
