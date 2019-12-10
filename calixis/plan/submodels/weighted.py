from django.db import models

from .config import Config_Region, Config_Grid, Config_System
from .inspiration import Inspiration_System_Feature

class BaseWeighted(models.Model):
    class Meta:
        abstract = True

    weight = models.SmallIntegerField()
    value = None
    parent = None

    def __str__(self):
        return "({weight}) {value_name}".format(weight=self.weight, value_name=self.value.name)

    def __repr__(self):
        return json.dumps(model_to_dict(
            self,
            fields=[field.name for field in self._meta.fields]
        ))

class Weighted_Config_Region(BaseWeighted):
    value = models.ForeignKey(Config_Region, on_delete=models.CASCADE)
    parent = models.ForeignKey(Config_Grid, on_delete=models.CASCADE)


class Weighted_Inspiration_System(BaseWeighted):
    value = models.ForeignKey(Inspiration_System_Feature, on_delete=models.CASCADE)
    parent = models.ForeignKey(Config_System, on_delete=models.CASCADE)
