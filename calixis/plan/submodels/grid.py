from django.db import models
from django.forms.models import model_to_dict
import json

class Grid_Sector(models.Model):
    def __repr__(self):
        return json.dumps(model_to_dict(
            self,
            fields=[field.name for field in self._meta.fields]
        ))

    def __str__(self):
        return self.name
    name = models.CharField(default="-", max_length=100)

class Grid_System(models.Model):
    sector = models.ForeignKey(Grid_Sector, on_delete=models.CASCADE)
    x = models.PositiveSmallIntegerField()
    y = models.PositiveSmallIntegerField()
    region = models.ForeignKey('Perterbation', on_delete=models.CASCADE)

class Grid_Route(models.Model):
    start = models.ForeignKey(Grid_System, on_delete=models.CASCADE, related_name='start')
    end = models.ForeignKey(Grid_System, on_delete=models.CASCADE, related_name='end')
