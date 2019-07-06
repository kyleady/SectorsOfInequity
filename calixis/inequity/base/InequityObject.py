import pprint

class InequityObject:
    def __init__(self, values):
        self._values = {}

    def _has_key(self, key):
        if key in self._values:
            return True
        else:
            raise Exception("\"" + key + "\" is not a valid key for " + self.__class__.__name__)

    def get(self, key):
        if self._has_key(key):
            return self._values[key]

    def set(self, key, value):
        if self._has_key(key):
            self._values[key] = value
            return self._values[key]

    def __getitem__(self, key):
        return self.get(key)

    def __eq__(self, obj):
        for key in self._values:
            if self[key] != obj[key]:
                return False

        return True


    def __repr__(self):
        return self.__class__.__name__ + "(" + pprint.pformat(self._values) + ")"
