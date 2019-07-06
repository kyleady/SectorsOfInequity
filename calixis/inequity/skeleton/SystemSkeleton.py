from inequity.base.InequityObject import InequityObject
from inequity.skeleton.CoordsSkeleton import CoordsSkeleton
from inequity.skeleton.RouteSkeleton import RouteSkeleton

class SystemSkeleton(InequityObject):
    def __init__(self, coords):
        self._values = {
            "blob_label": None,
            "coords": coords,
            "routes": []
        }

    @property
    def coords(self):
        return self["coords"]

    @property
    def blob_label(self):
        return self["blob_label"]

    def add_route(self, route):
        for pre_existing in self["routes"]:
            if route == pre_existing:
                return

        self["routes"].append(route)

    def connect(self, endpoint_system):
        if endpoint_system is None:
            return
        route = RouteSkeleton(self, endpoint_system)
        self.add_route(route)
        endpoint_system.add_route(route.reverse())

    def __add__(self, coords):
        return CoordsSkeleton(
            self.coords.x + coords.x,
            self.coords.y + coords.y
        )

    def get_blob_size(self, current_label):
        if self.blob_label is not None:
            return 0
        self.set("blob_label", current_label)
        systems_in_blob = 1
        for route in self["routes"]:
            systems_in_blob += route.end_system.get_blob_size(current_label)
        return systems_in_blob
