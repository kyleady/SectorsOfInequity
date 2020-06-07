from django.db import models

from .config import Config_Grid, Config_System, Config_Star_Cluster
from .perterbation import Perterbation
from .inspiration import Inspiration, Inspiration_Nested

class BaseWeighted(models.Model):
    class Meta:
        abstract = True

    weight = models.SmallIntegerField(blank=True, default=1)

    def __str__(self):
        return "({weight}) {value_name}".format(weight=self.weight, value_name=self.value.name)

    def __repr__(self):
        return json.dumps(model_to_dict(
            self,
            fields=[field.name for field in self._meta.fields]
        ))

class Weighted_Inspiration(BaseWeighted):
    value = models.ForeignKey(Inspiration, on_delete=models.CASCADE)
    systems = models.ManyToManyField(Config_System)
    star_clusters = models.ManyToManyField(Config_Star_Cluster)
    nested_inspirations = models.ManyToManyField(Inspiration_Nested)

class Weighted_Perterbation(BaseWeighted):
    value = models.ForeignKey(Perterbation, on_delete=models.CASCADE)
    parent = models.ForeignKey(Config_Grid, on_delete=models.CASCADE)
