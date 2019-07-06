from random import random
import math

from inequity.base.InequityObject import InequityObject
from inequity.skeleton.SystemSkeleton import SystemSkeleton
from inequity.skeleton.CoordsSkeleton import CoordsSkeleton

class SectorSkeleton(InequityObject):
    def __init__(self, skeleton_config):
        self._values = {
            "grid": [],
            "config": skeleton_config
        }

    def get_system(self, coords):
        min = CoordsSkeleton(0, 0)
        max = CoordsSkeleton(self["config"]["x"] - 1, self["config"]["y"] - 1)
        if coords.outside(min, max):
            return
        else:
            return self["grid"][coords.x][coords.y]

    def set_system(self, coords, system):
        self["grid"][coords.x][coords.y] = system
        if system is not None:
            system["coords"] = coords

    def systems(self):
        for row in self["grid"]:
            for system in row:
                if system is not None:
                    yield system

    def nearby_systems(self, system, reach):
        for r in range(1, reach + 1):
            for i_r in range(-r, r+1):
                for j_r in range(-r, r+1):
                    dr = CoordsSkeleton(i_r, j_r)
                    nearby_system = self.get_system(system + dr)
                    if nearby_system is not None and nearby_system is not system:
                        yield r, nearby_system

    def coords(self):
        for i in range(self["config"]["x"]):
            for j in range(self["config"]["y"]):
                yield CoordsSkeleton(i, j)

    def randomize(self):
        self.set("grid", [])
        self._populate_grid()
        self._randomize_routes()
        blob_sizes = self._label_blobs()
        self._reduce_to_largest_blob(blob_sizes)

    def debug_display(self, newline='\n'):
        output = ""
        systems = 0
        for coords in self.coords():
            if coords.y == 0:
                output += newline

            system = self["grid"][coords.x][coords.y]
            if system is not None:
                systems += 1
                connections = len(system["routes"])
                output += str(connections)
                if connections < 10:
                    output += " "
            else:
                output += "- "
            output += " "

        output += newline
        output += "Systems: " + str(systems)

        return output

    def _populate_grid(self):
        for coords in self.coords():
            if coords.y == 0:
                self["grid"].append([])

            if random() <= self["config"]["frequency_rate"]:
                system = SystemSkeleton(coords)
            else:
                system = None

            self["grid"][coords.x].append(system)

    def _randomize_routes(self):
        for system in self.systems():
            for range, endpoint_system in self.nearby_systems(system, self["config"]["reach"]):
                range_multiplier = math.pow(self["config"]["range_multiplier"], range-1)
                if random() <= self["config"]["connection_rate"] * range_multiplier:
                    system.connect(endpoint_system)

    def _label_blobs(self):
        blob_label = 0
        blob_sizes = []
        for system in self.systems():
            if system is not None:
                blob_size = system.get_blob_size(blob_label)
                if blob_size > 0:
                    blob_sizes.append(blob_size)
                    blob_label += 1

        return blob_sizes

    def _reduce_to_largest_blob(self, blob_sizes):
        largest_size = 0
        largest_label = None
        for blob_label in range(len(blob_sizes)):
            blob_size = blob_sizes[blob_label]
            if blob_size > largest_size:
                largest_size = blob_size
                largest_label = blob_label

        for system in self.systems():
            if system["blob_label"] != blob_label:
                self.set_system(system["coords"], None)
