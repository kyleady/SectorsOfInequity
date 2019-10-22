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

class GridSector(BaseAsset):
     pass

class GridSystem(BaseAsset):
    sector = models.ForeignKey(GridSector, on_delete=models.CASCADE)
    x = models.PositiveSmallIntegerField()
    y = models.PositiveSmallIntegerField()

class GridRoute(BaseAsset):
    start = models.ForeignKey(GridSystem, on_delete=models.CASCADE, related_name='start')
    end = models.ForeignKey(GridSystem, on_delete=models.CASCADE, related_name='end')
