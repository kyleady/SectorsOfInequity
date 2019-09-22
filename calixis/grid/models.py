from django.db import models
from django.forms.models import model_to_dict
import json

class BaseAsset(models.Model):
    class Meta:
        abstract = True

    def __repr__(self):
        return json.dumps(model_to_dict(
            self,
            fields=[field.name for field in self._meta.fields]
        ))

    def __str__(self):
        return self.name

    name = models.CharField(default="-", max_length=25)

class Sector(BaseAsset):
    config = models.ForeignKey(
            'config.Grid',
            on_delete=models.SET_NULL,
            null=True,
            blank=False
        )

class System(models.Model):
    sector = models.ForeignKey(Sector, on_delete=models.CASCADE)
    x = models.PositiveSmallIntegerField()
    y = models.PositiveSmallIntegerField()
    connections = models.ManyToManyField("self")
